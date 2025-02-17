package data

type Blockchain struct {
	ID     string
	Blocks []*Block
}

func NewBlockchain(id string) *Blockchain {
	return &Blockchain{
		ID:     id,
		Blocks: []*Block{NewBlock("In the beginning....", []byte{})},
	}
}

func (bc *Blockchain) AddBlock(data string) *Block {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
	return newBlock
}
