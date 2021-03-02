package raft

import (
	"encoding/json"
	"github.com/Open-Twin/alexandria/dns"
	"github.com/Open-Twin/alexandria/storage"
	"github.com/hashicorp/raft"
	"io"
	"log"
	"strconv"
)


type Fsm struct{
	MetadataRepo *storage.InMemoryStorageRepository
	DnsRepo *storage.StorageRepository
}
type metadata struct {
	Dnsormetadata bool
	Service string
	Ip      string
	Type    string
	Key     string
	Value   string
}
type dnsresource struct {
	Dnsormetadata bool
	Hostname string
	Ip string
	RequestType string
	ResourceRecord dns.DNSResourceRecord
}

// Apply log is invoked once a log entry is committed.
// It returns a value which will be made available in the
// ApplyFuture returned by Raft.Apply method if that
// method was called on the same Raft node as the FSM.
func (fsm *Fsm) Apply(logEntry *raft.Log) interface{} {
	var m metadata
	var d dnsresource

	dnsormeta := struct{
		DnsOrMetadata bool
	}{}
	if err := json.Unmarshal(logEntry.Data, &dnsormeta); err != nil {
		panic("Failed unmarshaling dnsormetadata log entry.")
	}
	log.Println("dnsormetadata: "+strconv.FormatBool(dnsormeta.DnsOrMetadata))
	if dnsormeta.DnsOrMetadata {
		if err := json.Unmarshal(logEntry.Data, &d); err != nil {
			panic("Failed unmarshaling Raft log entry.")
		}
		err := applyToDnsStore(fsm, d)
		if err != nil{
			return err
		}
	}else{
		if err := json.Unmarshal(logEntry.Data, &m); err != nil {
			panic("Failed unmarshaling Raft log entry.")
		}
		err := applyToMetadataStore(fsm, m)
		if err != nil{
			return err
		}
	}

	return nil
}

func applyToMetadataStore(fsm *Fsm, e metadata) error{
	switch e.Type {
	case "store":
		err := fsm.MetadataRepo.Create(e.Service,e.Ip,e.Key,e.Value)
		if err != nil{
			log.Print("store error: "+err.Error())
		}
		return nil
	case "update":
		err := fsm.MetadataRepo.Update(e.Service,e.Ip,e.Key,e.Value)
		if err != nil{
			log.Print("update error: "+err.Error())
		}
		return nil
	case "delete":
		err := fsm.MetadataRepo.Delete(e.Service,e.Ip,e.Key)
		if err != nil{
			log.Print("delete error: "+err.Error())
		}
		return nil
	default:
		log.Printf("Unrecognized event type in Raft log entry: %s.", e.Type)
	}
	return nil
}

func applyToDnsStore(fsm *Fsm, e dnsresource) error {
	switch e.RequestType {
	case "store":
		err := fsm.DnsRepo.Create(e.Hostname, e.ResourceRecord)
		if err != nil{
			log.Print("store error: "+err.Error())
		}
		return nil
	case "update":
		err := fsm.DnsRepo.Update(e.Hostname,e.ResourceRecord)
		if err != nil{
			log.Print("update error: "+err.Error())
		}
		return nil
	case "delete":
		err := fsm.DnsRepo.Delete(e.Hostname)
		if err != nil{
			log.Print("delete error: "+err.Error())
		}
		return nil
	default:
		log.Printf("Unrecognized event type in Raft log entry: %s.", e.RequestType)
	}
	return nil
}

// Snapshot is used to support log compaction. This call should
// return an FSMSnapshot which can be used to save a point-in-time
// snapshot of the FSM. Apply and Snapshot are not called in multiple
// threads, but Apply will be called concurrently with Persist. This means
// the FSM should be implemented in a fashion that allows for concurrent
// updates while a snapshot is happening.
func (fsm *Fsm) Snapshot() (raft.FSMSnapshot, error) {
	/*fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	return &fsmSnapshot{StateValue: fsm.stateValue}, nil*/
	return nil,nil
}

// Restore is used to restore an FSM from a snapshot. It is not called
// concurrently with any other command. The FSM must discard all previous
// state.
func (fsm *Fsm) Restore(serialized io.ReadCloser) error {
	/*var snapshot fsmSnapshot
	if err := json.NewDecoder(serialized).Decode(&snapshot); err != nil {
		return err
	}

	fsm.stateValue = snapshot.StateValue
	return nil*/
	return nil
}