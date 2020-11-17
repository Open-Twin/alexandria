package raft

// Persist should dump all necessary state to the WriteCloser 'sink',
// and call sink.Close() when finished or call sink.Cancel() on error.
func Persist(sink SnapshotSink) error {
}

// Release is invoked when we are finished with the snapshot.
func Release() {
}
