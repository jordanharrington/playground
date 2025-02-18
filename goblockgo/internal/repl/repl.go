package repl

import (
	"bufio"
	"fmt"
	"github.com/jordanharrington/playground/goblockgo/internal/service"
	"os"
	"strings"
)

func Eval(service service.GoBlockGo) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}
		handleInput(scanner.Text(), service)
	}
}

func handleInput(input string, service service.GoBlockGo) {
	tokens := strings.Fields(input)
	if len(tokens) == 0 {
		return
	}

	command := tokens[0]
	args := tokens[1:]

	switch command {
	case "exit":
		handleExit(args)
	case "create":
		handleCreate(args, service)
	case "add":
		handleAdd(args, service)
	case "describe":
		handleDescribe(args, service)
	case "get":
		handleGet(args, service)
	case "ls":
		handleList(args, service)
	default:
		fmt.Println("Unknown command. Available commands: exit, create, describe, get, add, ls")
	}
}

func handleExit(args []string) {
	if len(args) > 0 {
		fmt.Println("Usage: exit")
		return
	}
	fmt.Println("Exiting REPL.")
	os.Exit(0)
}

func handleCreate(args []string, service service.GoBlockGo) {
	if len(args) != 1 {
		fmt.Println("Usage: create [blockchain_name]")
		return
	}
	b, err := service.CreateBlockchain(args[0])
	if err != nil {
		fmt.Printf("Error creating blockchain: %s\n", err)
		return
	}
	fmt.Printf("Blockchain created with ID: %s\n", b.ID)
}

func handleAdd(args []string, service service.GoBlockGo) {
	if len(args) != 2 {
		fmt.Println("Usage: add [blockchain_id] [data]")
		return
	}
	if _, err := service.AddBlock(args[0], args[1]); err != nil {
		fmt.Printf("Error adding block: %v\n", err)
		return
	}
	fmt.Println("Block added.")
}

func handleDescribe(args []string, service service.GoBlockGo) {
	if len(args) != 1 {
		fmt.Println("Usage: describe [blockchain_id]")
		return
	}
	bc, err := service.GetBlockchain(args[0])
	if err != nil {
		fmt.Printf("Error describing blocks: %s\n", err)
		return
	}
	fmt.Println("Blocks in Blockchain:")
	for _, block := range bc.Blocks {
		fmt.Printf("- Timestamp: %d, Data: %s, PrevHash: %x, Hash: %x\n", block.Timestamp, string(block.Data), block.PrevBlockHash, block.Hash)
	}
}

func handleGet(args []string, service service.GoBlockGo) {
	if len(args) != 2 {
		fmt.Println("Usage: get [blockchain_id] [block_hash]")
		return
	}
	block, err := service.GetBlock(args[0], args[1])
	if err != nil {
		fmt.Printf("Error finding block: %s\n", err)
		return
	}
	fmt.Println("Block:")
	fmt.Printf("- Timestamp: %d, Data: %s, PrevHash: %x, Hash: %x\n", block.Timestamp, string(block.Data), block.PrevBlockHash, block.Hash)
}

func handleList(args []string, service service.GoBlockGo) {
	if len(args) > 0 {
		fmt.Println("Usage: ls")
		return
	}
	blockchains, err := service.ListBlockchains()
	if err != nil {
		fmt.Printf("Error listing blockchains: %s\n", err)
		return
	}
	fmt.Println("Blockchains:")
	for _, bc := range blockchains {
		fmt.Printf("- %s\n", bc.ID)
	}
}
