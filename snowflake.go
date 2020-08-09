package snowflake

import (
	"errors"
	"sync"
	"time"
)

// snowflake format:
//
// |    symbol  |   timestamp |   workerid  |     seq     |
// |------------|-------------|-------------|-------------|
// |<-- 1bit -->|<-- 41bit -->|<-- 10bit -->|<-- 12bit -->|
//

var (
	beginTs     int64 = 1577808000000          //2020-01-01 00:00:00
	workIDBit   int64 = 10                     //10 bit for workerID
	seqBit      int64 = 12                     //12 bit for sequence per second
	maxWorkerID int64 = (-1 << workIDBit) ^ -1 //max workerID is 1023
	maxSeq      int64 = (-1 << seqBit) ^ -1    //max seq is 2^12 -1

	errInvalidWorkderID = errors.New("Bad workerid, must be between 0 and 1023")
)

// IDGenerator implements Twitter snowflake.
type IDGenerator struct {
	sync.Mutex
	workerID int64
	ts       int64
	seq      int64
}

// NewIDGenerator return instance of IDGenerator
func NewIDGenerator(workerID int64) (*IDGenerator, error) {
	if workerID > maxWorkerID || workerID < 0 {
		return nil, errInvalidWorkderID
	}
	return &IDGenerator{
		workerID: workerID,
	}, nil
}

// GenerateID generates distributed unique ID
func (s *IDGenerator) GenerateID() int64 {

	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixNano() / 1000000

	// generate difference seqs within 1ms
	if s.ts == now {
		s.seq = (s.seq + 1) ^ maxSeq

		// The serial number is used up
		if s.seq == 0 {
			for s.ts <= now {
				s.ts = time.Now().UnixNano() / 1000000
			}
		}

	} else {
		s.seq = 0
	}

	s.ts = now

	// generate id
	id := ((now - beginTs) << (workIDBit + seqBit)) | (s.workerID << seqBit) | s.seq

	return id
}
