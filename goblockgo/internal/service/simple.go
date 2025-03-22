package service

import (
	"encoding/json"
	"fmt"
	"github.com/jordanharrington/playground/goblockgo/internal/data"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type simpleBlockchainService struct {
	blockchains map[string]*data.Blockchain
	lock        sync.RWMutex
}

func NewSimpleBlockchainService() GoBlockGo {
	s := &simpleBlockchainService{
		blockchains: make(map[string]*data.Blockchain),
	}

	load(s)
	setupShutdownHook(s)

	return s
}

func load(s *simpleBlockchainService) {
	s.lock.Lock()
	defer s.lock.Unlock()

	filename := os.Getenv("GBG_SIMPLE_PATH")
	if filename == "" {
		log.Println("GBG_SIMPLE_PATH is not set, could not load blockchains")
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error closing file:", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&s.blockchains); err != nil {
		log.Println("Error unmarshaling blockchains:", err)
	}
}

func save(s *simpleBlockchainService) {
	s.lock.Lock()
	defer s.lock.Unlock()

	value := os.Getenv("GBG_SIMPLE_PATH")
	if value == "" {
		log.Println("Could not find GBG_SIMPLE_PATH")
		return
	}

	bytes, err := json.Marshal(s.blockchains)
	if err != nil {
		log.Println("Error marshaling blockchains:", err)
		return
	}

	writeFile(value, bytes)
}

func writeFile(filename string, data []byte) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error closing file:", err)
		}
	}(file)

	if _, err := file.Write(data); err != nil {
		log.Println("Error writing to file:", err)
	}
}

func setupShutdownHook(s *simpleBlockchainService) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Received shutdown signal, saving blockchains...")
		save(s)
		os.Exit(0)
	}()
}

func (s *simpleBlockchainService) CreateBlockchain(id string) (*data.Blockchain, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, exists := s.blockchains[id]; exists {
		return nil, fmt.Errorf("blockchain already exists")
	}

	bc := data.NewBlockchain(id)
	s.blockchains[id] = bc
	return bc, nil
}

func (s *simpleBlockchainService) AddBlock(id string, data string) (*data.Block, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	bc, exists := s.blockchains[id]
	if !exists {
		return nil, fmt.Errorf("blockchain not found")
	}

	newBlock := bc.AddBlock(data)
	return newBlock, nil
}

func (s *simpleBlockchainService) GetBlockchain(id string) (*data.Blockchain, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	bc, exists := s.blockchains[id]
	if !exists {
		return nil, fmt.Errorf("blockchain not found")
	}
	return bc, nil
}

func (s *simpleBlockchainService) GetBlock(id string, hash string) (*data.Block, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	bc, exists := s.blockchains[id]
	if !exists {
		return nil, fmt.Errorf("blockchain not found")
	}

	for _, block := range bc.Blocks {
		if fmt.Sprintf("%x", block.Hash) == hash {
			return block, nil
		}
	}

	return nil, fmt.Errorf("block not found")
}

func (s *simpleBlockchainService) ListBlockchains() ([]*data.Blockchain, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var bcs []*data.Blockchain
	for _, bc := range s.blockchains {
		bcs = append(bcs, bc)
	}
	return bcs, nil
}

func (s *simpleBlockchainService) ListBlocks(id string) ([]*data.Block, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	bc, exists := s.blockchains[id]
	if !exists {
		return nil, fmt.Errorf("blockchain not found")
	}

	return bc.Blocks, nil
}

func (s *simpleBlockchainService) DeleteBlockchain(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, exists := s.blockchains[id]; !exists {
		return fmt.Errorf("blockchain not found")
	}

	delete(s.blockchains, id)
	return nil
}
