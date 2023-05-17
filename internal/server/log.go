package server

import (
	"fmt"
	"sync"
)

// Logs are a data structure for an append only sequence of records,
// ordered by time, this is a simple commit log with slice.

// Define a single record
type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

// Now define the structure log that will keep
// a slice of Records. A log is defined over a
// Log and each each log is concurrency protected.
type Log struct {
	mu      sync.Mutex
	records []Record
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")

// This function will initiate the log structure itself.
// This fuction returns a pointer to log's address.
func NewLog() *Log {
	return &Log{}
}

// This function will first calculate the offset
// that is basically the size of the current record
// and append the incoming record to slice of records
// in the log. Will take a Record object as parameter,
// will return Offset value, and an error.
func (r *Log) Append(record Record) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	record.Offset = uint64(len(r.records))
	r.records = append(r.records, record)
	return record.Offset, nil
}

// This function will first check if we the provided offset
// is within the range of the offset limits of the records that are
// currently present.
// If the offset is within the range, then we will return the record.
// If the offset is not within the range, then we will return an error.
func (r *Log) Read(offset uint64) (Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if offset >= uint64(len(r.records)) {
		return Record{}, ErrOffsetNotFound
	}

	return r.records[offset], nil
}
