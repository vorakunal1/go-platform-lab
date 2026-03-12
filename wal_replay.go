package main

import (
    "encoding/binary"
    "encoding/json"
    "path/filepath"
    "fmt"
    "hash/crc32"
    "io"
    "os"
)

const HeaderSize = 20

type Replayer struct {
    file *os.File
    lastAppliedLSN uint64
}

type ReplayResult struct {
    LastValidLSN uint64
    LastValidOffset int64
    Error error
}

func (r *Replayer) Recover (applyFc func(lsn uint64, data []byte) error) ReplayResult {
    var currentOffset int64 = 0
    var lastValidLSN uint64 = r.lastAppliedLSN

    for {
        headerBuf := make([]byte, HeaderSize)
        _, err := io.ReadFull(r.file, headerBuf)

        if err == io.EOF {
            break
        }
        if err != nil {
            return ReplayResult{lastValidLSN, currentOffset, err}
        }

        lsn := binary.BigEndian.Uint64(headerBuf[0:8])
        dataLen := binary.BigEndian.Uint32(headerBuf[12:16])
        expectedCRC := binary.BigEndian.Uint32(headerBuf[16:20])

        dataBuf := make([]byte, dataLen)
        _, err = io.ReadFull(r.file, dataBuf)

        if err != nil {
            return ReplayResult{lastValidLSN, currentOffset, err}
        }

        if crc32.ChecksumIEEE(dataBuf) != expectedCRC {
            return ReplayResult{lastValidLSN, currentOffset, fmt.Errorf("checksum mismatch")}
        }

        if lsn > r.lastAppliedLSN {
            if err := applyFc(lsn, dataBuf); err != nil {
                return ReplayResult{lastValidLSN, currentOffset, err}
            }
            lastValidLSN = lsn
        }

        currentOffset += int64(HeaderSize + dataLen)

    }

    return ReplayResult{lastValidLSN, currentOffset, nil}
}


func SaveCheckpoint(path string, lsn uint64) error {
    cp := Checkpoint{lastAppliedLSN: lsn}
    data, err := json.Marshal(cp)
    if err != nil {
        return err
    }

    tmpPath := path + ".tmp"

    if err := os.WriteFile(tmpPath, data, 0644); err != nil {
        return err
    }

    f, err := os.Open(tmpPath)
    if err != nil {
        return err
    }

    if err := f.Sync(); err != nil {
        f.Close()
        return err
    }
    f.Close()

    if err := os.Rename(tmpPath, path); err != nil {
        return err
    }

    return nil
}

func loadCheckpoint(path string) uint64 {
    data, err := os.ReadFile(path)

    if err != nil {
        return 0
    }

    var cp Checkpoint
    json.Unmarshal(data, &cp)
    return cp.lastAppliedLSN
}


func main() {
    walPath := "simulation.wal"

    // f, _ := os.OpenFile(walPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

    // for i := uint64(1); i<=13; i++ {
    //     writeRecord(f, i, []byte(fmt.Sprintf("Data_%d", i)))
    // }

    fmt.Println("Starting Recovery....")
    rf, _ := os.Open(walPath)
    replayer := &Replayer{file: rf, lastAppliedLSN: 0}

    result := replayer.Recover(func(lsn uint64, data []byte) error {
        fmt.Println("Replayed LSN %d: %s", lsn, string(data))
        return nil
    })

    if result.Error != nil {
        fmt.Println("Stopped at LSn %d due to crash. Repairing...", result.LastValidLSN+1)

        os.Truncate(walPath, result.LastValidOffset)
    }

    fmt.Println("system ready")


}


func writeRecord(f *os.File, lsn uint64, data []byte) {
    h := make([]byte, HeaderSize)
    binary.BigEndian.PutUint64(h[0:8], lsn)
    binary.BigEndian.PutUint32(h[12:16], uint32(len(data)))
    binary.BigEndian.PutUint32(h[16:20], crc32.ChecksumIEEE(data))
    f.Write(h)
    f.Write(data)
}
