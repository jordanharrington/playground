package data

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	now := time.Now().Unix()
	byteData := []byte(data)

	timestamp := []byte(strconv.FormatInt(now, 10))
	headers := bytes.Join([][]byte{prevBlockHash, byteData, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	return &Block{now, byteData, prevBlockHash, hash[:]}
}
