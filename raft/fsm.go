package raft

import (
	"encoding/json"
	"fmt"
	"github.com/Open-Twin/alexandria/dns/metadata"
	"github.com/hashicorp/raft"
	"io"
	"log"
)

/*type fsm struct {
	mutex sync.Mutex
	stateValue int
}*/
type fsm struct{
	Repo *metadata.InMemoryStorageRepository
}
//var repo metadata.InMemoryStorageRepository
type event struct {
	Service string
	Ip      string
	Type    string
	Key     string
	Value   string
}
// Apply log is invoked once a log entry is committed.
// It returns a value which will be made available in the
// ApplyFuture returned by Raft.Apply method if that
// method was called on the same Raft Node as the FSM.
func (fsm *fsm) Apply(logEntry *raft.Log) interface{} {
	var e event
	if err := json.Unmarshal(logEntry.Data, &e); err != nil {
		panic("Failed unmarshaling Raft log entry. This is a bug.")
	}

	switch e.Type {
	case "store":
		err := fsm.Repo.Create(e.Service,e.Ip,e.Key,e.Value)
		if err != nil{
			log.Print("store error: "+err.Error())
		}
		return nil
	case "update":
		err := fsm.Repo.Update(e.Service,e.Ip,e.Key,e.Value)
		if err != nil{
			log.Print("update error: "+err.Error())
		}
		return nil
	case "delete":
		err := fsm.Repo.Delete(e.Service,e.Ip,e.Key)
		if err != nil{
			log.Print("delete error: "+err.Error())
		}
		return nil
	default:
		panic(fmt.Sprintf("Unrecognized event type in Raft log entry: %s. This is a bug.", e.Type))
	}
}
// Snapshot is used to support log compaction. This call should
// return an FSMSnapshot which can be used to save a point-in-time
// snapshot of the FSM. Apply and Snapshot are not called in multiple
// threads, but Apply will be called concurrently with Persist. This means
// the FSM should be implemented in a fashion that allows for concurrent
// updates while a snapshot is happening.
func (fsm *fsm) Snapshot() (raft.FSMSnapshot, error) {
	/*fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	return &fsmSnapshot{StateValue: fsm.stateValue}, nil*/
	return nil,nil
}

// Restore is used to restore an FSM from a snapshot. It is not called
// concurrently with any other command. The FSM must discard all previous
// state.
func (fsm *fsm) Restore(serialized io.ReadCloser) error {
	/*var snapshot fsmSnapshot
	if err := json.NewDecoder(serialized).Decode(&snapshot); err != nil {
		return err
	}

	fsm.stateValue = snapshot.StateValue
	return nil*/
	return nil
}