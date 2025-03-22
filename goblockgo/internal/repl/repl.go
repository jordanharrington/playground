package repl

import (
	"fmt"
	"github.com/jordanharrington/playground/goblockgo/internal/service"
	"github.com/jordanharrington/playground/goreplgo/repl"
)

func Eval(service service.GoBlockGo) {
	r := repl.NewREPL()

	commands := []*repl.Command{
		{
			Name:        "exit",
			Description: "Exit the REPL",
			Args:        nil,
			Handler: func(args map[string]interface{}) (bool, error) {
				fmt.Println("Exiting REPL.")
				return true, nil
			},
		},
		{
			Name:        "create",
			Description: "Create a new blockchain",
			Args: []repl.Arg{
				{Name: "name", ArgType: repl.StringArg, Required: true},
			},
			Handler: func(args map[string]interface{}) (bool, error) {
				name := args["name"].(string)
				b, err := service.CreateBlockchain(name)
				if err != nil {
					return false, fmt.Errorf("error creating blockchain: %v", err)
				}
				fmt.Printf("Blockchain created with ID: %s\n", b.ID)
				return false, nil
			},
		},
		{
			Name:        "add",
			Description: "Add block to a blockchain",
			Args: []repl.Arg{
				{Name: "blockchain_id", ArgType: repl.StringArg, Required: true},
				{Name: "data", ArgType: repl.StringArg, Required: true},
			},
			Handler: func(args map[string]interface{}) (bool, error) {
				id := args["blockchain_id"].(string)
				data := args["data"].(string)
				if _, err := service.AddBlock(id, data); err != nil {
					return false, fmt.Errorf("error adding block: %v", err)
				}
				fmt.Println("Block added.")
				return false, nil
			},
		},
		{
			Name:        "describe",
			Description: "Describe a blockchain's blocks",
			Args: []repl.Arg{
				{Name: "blockchain_id", ArgType: repl.StringArg, Required: true},
			},
			Handler: func(args map[string]interface{}) (bool, error) {
				id := args["blockchain_id"].(string)
				bc, err := service.GetBlockchain(id)
				if err != nil {
					return false, fmt.Errorf("error describing blocks: %v", err)
				}
				fmt.Println("Blocks in Blockchain:")
				for _, block := range bc.Blocks {
					fmt.Printf("- Timestamp: %d, Data: %s, PrevHash: %x, Hash: %x\n",
						block.Timestamp, string(block.Data), block.PrevBlockHash, block.Hash)
				}
				return false, nil
			},
		},
		{
			Name:        "get",
			Description: "Get a block by hash",
			Args: []repl.Arg{
				{Name: "blockchain_id", ArgType: repl.StringArg, Required: true},
				{Name: "block_hash", ArgType: repl.StringArg, Required: true},
			},
			Handler: func(args map[string]interface{}) (bool, error) {
				bcID := args["blockchain_id"].(string)
				hash := args["block_hash"].(string)
				block, err := service.GetBlock(bcID, hash)
				if err != nil {
					return false, fmt.Errorf("error finding block: %v", err)
				}
				fmt.Println("Block:")
				fmt.Printf("- Timestamp: %d, Data: %s, PrevHash: %x, Hash: %x\n",
					block.Timestamp, string(block.Data), block.PrevBlockHash, block.Hash)
				return false, nil
			},
		},
		{
			Name:        "ls",
			Description: "List all blockchains",
			Args:        nil,
			Handler: func(args map[string]interface{}) (bool, error) {
				blockchains, err := service.ListBlockchains()
				if err != nil {
					return false, fmt.Errorf("error listing blockchains: %v", err)
				}
				fmt.Println("Blockchains:")
				for _, bc := range blockchains {
					fmt.Printf("- %s\n", bc.ID)
				}
				return false, nil
			},
		},
		{
			Name:        "help",
			Description: "Show available commands",
			Args:        nil,
			Handler:     r.Help(),
		},
	}

	for _, cmd := range commands {
		r.Add(cmd)
	}

	r.Run()
}
