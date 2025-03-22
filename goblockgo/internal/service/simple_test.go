package service

import (
	"github.com/jordanharrington/playground/goblockgo/internal/data"
	"log"
	"os"
	"testing"
)

func TestLoadAndSave(t *testing.T) {
	// Setup temporary file
	tmpFile, err := os.CreateTemp("", "blockchain_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			log.Fatalf("Failed to remove temp file: %v", err)
		}
	})

	// Set env var for persistence path
	t.Setenv("GBG_SIMPLE_PATH", tmpFile.Name())

	// Create a service and populate it
	svc := &simpleBlockchainService{
		blockchains: make(map[string]*data.Blockchain),
	}
	svc.blockchains["demo"] = &data.Blockchain{
		ID: "demo",
	}

	// Save to disk
	save(svc)

	// Now create a new instance and load from disk
	loadedSvc := &simpleBlockchainService{
		blockchains: make(map[string]*data.Blockchain),
	}
	load(loadedSvc)

	// Assert it loaded correctly
	if _, ok := loadedSvc.blockchains["demo"]; !ok {
		t.Fatalf("Expected blockchain 'demo' to be loaded")
	}
}
