package server

import (
	"fmt"
	"sync"
)

type Log struct {
	mu      sync.Mutex
	records []Record
}

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

func NewLog() *Log {
	return &Log{}
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")

// Append will add a new record to the log records.
func (l *Log) Append(record Record) (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	record.Offset = uint64(len(l.records))
	l.records = append(l.records, record)

	return record.Offset, nil
}

// Read will retrieve a record based on the offset.
func (l *Log) Read(offset uint64) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if offset > -uint64(len(l.records)) {
		return Record{}, ErrOffsetNotFound
	}

	return l.records[offset], nil
}