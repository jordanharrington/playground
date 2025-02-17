package service

import "github.com/jordanharrington/playground/goblockgo/internal/data"

type GoBlockGo interface {
	CreateBlockchain(name string) (*data.Blockchain, error)
	AddBlock(id string, data string) (*data.Block, error)
	GetBlockchain(id string) (*data.Blockchain, error)
	GetBlock(id string, hash string) (*data.Block, error)
	ListBlockchains() ([]*data.Blockchain, error)
	ListBlocks(id string) ([]*data.Block, error)
	DeleteBlockchain(id string) error
}
