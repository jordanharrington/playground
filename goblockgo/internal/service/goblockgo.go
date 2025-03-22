package service

import "github.com/jordanharrington/playground/goblockgo/internal/data"

// GoBlockGo defines the core service interface for interacting with blockchains.
// It provides methods for creating, appending to, retrieving, and deleting blockchains and their blocks.
type GoBlockGo interface {
	// CreateBlockchain creates a new blockchain with the given name.
	// Returns the created blockchain or an error if creation fails.
	CreateBlockchain(name string) (*data.Blockchain, error)

	// AddBlock adds a new block with the given data to the blockchain identified by id.
	// Returns the newly added block or an error if the operation fails.
	AddBlock(id string, data string) (*data.Block, error)

	// GetBlockchain retrieves the full blockchain identified by id.
	// Returns the blockchain or an error if it doesn't exist.
	GetBlockchain(id string) (*data.Blockchain, error)

	// GetBlock retrieves a single block by its hash from the blockchain identified by id.
	// Returns the block or an error if not found.
	GetBlock(id string, hash string) (*data.Block, error)

	// ListBlockchains returns a list of all existing blockchains.
	// Each blockchain is returned without its full block history.
	ListBlockchains() ([]*data.Blockchain, error)

	// ListBlocks returns all blocks from the blockchain identified by id.
	// Returns a slice of blocks or an error if the blockchain doesn't exist.
	ListBlocks(id string) ([]*data.Block, error)

	// DeleteBlockchain deletes the blockchain identified by id from the system.
	// Returns an error if deletion fails.
	DeleteBlockchain(id string) error
}
