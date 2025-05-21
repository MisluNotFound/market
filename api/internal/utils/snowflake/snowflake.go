package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

const (
	epoch        = int64(1688140800000) // 2023-07-01 00:00:00 UTC
	workerIDBits = uint(5)              // 5位workerID
	dataCenterIDBits = uint(5)          // 5位dataCenterID
	sequenceBits = uint(12)             // 12位序列号

	maxWorkerID = -1 ^ (-1 << workerIDBits)
	maxDataCenterID = -1 ^ (-1 << dataCenterIDBits)
	maxSequence = -1 ^ (-1 << sequenceBits)

	timeShift = workerIDBits + dataCenterIDBits + sequenceBits
	workerShift = dataCenterIDBits + sequenceBits
	dataCenterShift = sequenceBits
)

var IdGenerator *IDGenerator

func init() {
	var err error
	// 使用默认workerID和dataCenterID
	IdGenerator, err = NewIDGenerator(1, 1)
	if err != nil {
		panic(err)
	}
}

type IDGenerator struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	dataCenterID  int64
	sequence      int64
}

func NewIDGenerator(workerID, dataCenterID int64) (*IDGenerator, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("worker ID out of range")
	}
	if dataCenterID < 0 || dataCenterID > maxDataCenterID {
		return nil, errors.New("data center ID out of range")
	}
	return &IDGenerator{
		workerID:     workerID,
		dataCenterID: dataCenterID,
	}, nil
}

func (g *IDGenerator) NextID() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1e6

	if timestamp < g.lastTimestamp {
		return "", errors.New("clock moved backwards")
	}

	if timestamp == g.lastTimestamp {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			for timestamp <= g.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastTimestamp = timestamp

	id := ((timestamp - epoch) << timeShift) |
		(g.workerID << workerShift) |
		(g.dataCenterID << dataCenterShift) |
		g.sequence

	return strconv.FormatInt(id, 10), nil
}