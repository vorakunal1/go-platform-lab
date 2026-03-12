package main 

import (
	"hash/crc32"
	"fmt"
	"os"
	"sync"
	"encoding/binary"
)

const HeaderSize = 30

type WAL struct {
	file *os.File
	mu sync.Mutex
}

func OpenWal(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{file: f}, nil
}

func (w *WAL) Append(lsn uint64, opType uint32, data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	dataLen := uint32(len(data))
	checksum := crc32.ChecksumIEEE(data)

	packet := make([]byte, HeaderSize+dataLen)

	binary.BigEndian.PutUint64(packet[0:8], lsn)
	binary.BigEndian.PutUint32(packet[8:12], opType)
	binary.BigEndian.PutUint32(packet[12:16], dataLen)
	binary.BigEndian.PutUint32(packet[16:20], checksum)

	copy(packet[20:], data)

	if _, err := w.file.Write(packet); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	if err := w.file.Sync(); err != nil {
		return fmt.Errorf("fsync failed: %w", err)
	}

	return nil
}


func main() {
	wal, err := OpenWal("simulation.wal")
	if err != nil {
		panic(err)
	}

	defer wal.file.Close()

	var wg sync.WaitGroup
	userCount := 5
	const OpUpdate uint32 = 2

	for i:=1; i<userCount; i++ {
		wg.Add(1)
		go func(userId int) {
			defer wg.Done()
			lsn := uint64(100+userId)

			payload := []byte(fmt.Sprintf("user_%d:SET AGE:62\n", userId))
			fmt.Printf("User %d sending data to WAL.")

			err:=wal.Append(lsn, OpUpdate, payload)

			if err != nil {
				fmt.Printf("User %d, Error:", userId, err)
			} else {
				fmt.Printf("User %d, persisted to disk. LSN %d", userId, lsn)
			}

		}(i)
	}

	wg.Wait()
	fmt.Println("All operation are physically safe.")
}

