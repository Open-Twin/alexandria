package raft

import (
	"encoding/json"
	"errors"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/hashicorp/raft"
	"github.com/rs/zerolog/log"
	"io"
	"strconv"
	"sync"
)

type Fsm struct {
	MetadataRepo *storage.InMemoryStorageRepository
	DnsRepo      *storage.StorageRepository
	snapMutex    sync.RWMutex
}

// Apply log is invoked once a log entry is committed.
// It returns a value which will be made available in the
// ApplyFuture returned by Raft.Apply method if that
// method was called on the same Raft node as the FSM.
func (fsm *Fsm) Apply(logEntry *raft.Log) interface{} {
	var m storage.Metadata
	var d storage.Dnsresource
	dnsormeta := struct {
		DnsOrMetadata bool
	}{}
	if err := json.Unmarshal(logEntry.Data, &dnsormeta); err != nil {
		return errors.New("failed unmarshaling dnsormetadata log entry")
	}
	log.Debug().Msg("dnsormetadata: " + strconv.FormatBool(dnsormeta.DnsOrMetadata))
	if dnsormeta.DnsOrMetadata {
		if err := json.Unmarshal(logEntry.Data, &d); err != nil {
			return errors.New("failed unmarshaling Raft log entry")
		}
		err := applyToDnsStore(fsm, d)
		if err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(logEntry.Data, &m); err != nil {
			return errors.New("failed unmarshaling Raft log entry")
		}
		err := applyToMetadataStore(fsm, m)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyToMetadataStore(fsm *Fsm, e storage.Metadata) error {
	switch e.Type {
	case "store":
		err := fsm.MetadataRepo.Create(e.Service, e.Ip, e.Key, e.Value)
		if err != nil {
			log.Error().Msgf("store error: %s", err.Error())
			return err
		}
		return nil
	case "update":
		err := fsm.MetadataRepo.Update(e.Service, e.Ip, e.Key, e.Value)
		if err != nil {
			log.Error().Msgf("update error: %s", err.Error())
			return err
		}
		return nil
	case "delete":
		err := fsm.MetadataRepo.Delete(e.Service, e.Ip, e.Key)
		if err != nil {
			log.Error().Msgf("delete error: %s", err.Error())
			return err
		}
		return nil
	default:
		log.Error().Msgf("Unrecognized metadata event type in Raft log entry: %s", e.Type)
		return errors.New("Unrecognized metadata event type in Raft log entry: " + e.Type)
	}
}

func applyToDnsStore(fsm *Fsm, e storage.Dnsresource) error {
	switch e.RequestType {
	case "store":
		err := fsm.DnsRepo.Create(e.Hostname, e.Ip, e.ResourceRecord)
		if err != nil {
			log.Error().Msgf("store error: %s", err.Error())
			return err
		}
		return nil
	case "update":
		err := fsm.DnsRepo.Update(e.Hostname, e.Ip, e.ResourceRecord)
		if err != nil {
			log.Error().Msgf("update error: %s", err.Error())
			return err
		}
		return nil
	case "delete":
		err := fsm.DnsRepo.Delete(e.Hostname, e.Ip)
		if err != nil {
			log.Error().Msgf("delete error: %s", err.Error())
			return err
		}
		return nil
	default:
		log.Error().Msgf("Unrecognized dns event type in Raft log entry: %s", e.RequestType)
		return errors.New("Unrecognized dns event type in Raft log entry: " + e.RequestType)
		// TODO: Return error
	}
}

// Snapshot is used to support log compaction. This call should
// return an FSMSnapshot which can be used to save a point-in-time
// snapshot of the FSM. Apply and Snapshot are not called in multiple
// threads, but Apply will be called concurrently with Persist. This means
// the FSM should be implemented in a fashion that allows for concurrent
// updates while a snapshot is happening.
func (fsm *Fsm) Snapshot() (raft.FSMSnapshot, error) {
	fsm.snapMutex.Lock()
	defer fsm.snapMutex.Unlock()

	return &fsmSnapshot{
		MetadataRepo: fsm.MetadataRepo,
		DnsRepo:      fsm.DnsRepo,
	}, nil
}

// Restore is used to restore an FSM from a snapshot. It is not called
// concurrently with any other command. The FSM must discard all previous
// state.
func (fsm *Fsm) Restore(serialized io.ReadCloser) error {
	var snapshot fsmSnapshot
	if err := json.NewDecoder(serialized).Decode(&snapshot); err != nil {
		return err
	}
	fsm.snapMutex.Lock()
	defer fsm.snapMutex.Unlock()
	fsm.DnsRepo = snapshot.DnsRepo
	fsm.MetadataRepo = snapshot.MetadataRepo
	return nil
}
