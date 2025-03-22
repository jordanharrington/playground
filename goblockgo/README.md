# 🧱 GoBlockGo

**GoBlockGo** is a lightweight, modular blockchain service written in Go. It provides both a REPL (read-eval-print loop) and HTTP API for interacting with in-memory blockchains. This project is ideal for experimentation, education, or rapid prototyping with blockchain concepts.

---

## 🚀 Features

- 🛠 Create and manage multiple blockchains
- 📦 Add and query blocks
- 🔍 Get blockchains or individual blocks by ID/hash
- 🧵 Built-in REPL for interactive local control

---

## 💻 REPL Commands

Launch the REPL by launching the service in interactive mode.

The REPL supports the following commands:

| Command       | Description                             | Usage Example                                     |
|---------------|-----------------------------------------|---------------------------------------------------|
| `exit`        | Exit the REPL                           | `exit`                                            |
| `create`      | Create a new blockchain                 | `create mychain`                                  |
| `add`         | Add a block to a blockchain             | `add mychain "some data"`                         |
| `describe`    | List all blocks in a blockchain         | `describe mychain`                                |
| `get`         | Retrieve a block by hash                | `get mychain abcd1234...`                         |
| `ls`          | List all blockchain IDs                 | `ls`                                              |
| `help`        | Show help with available commands       | `help`                                            |

All commands are validated, typed, and provide structured feedback.

---

## 🧠 Service Interface

The REPL relies on the following interface, which must be implemented by your blockchain service:

```go
type GoBlockGo interface {
	CreateBlockchain(name string) (*data.Blockchain, error)
	AddBlock(id string, data string) (*data.Block, error)
	GetBlockchain(id string) (*data.Blockchain, error)
	GetBlock(id string, hash string) (*data.Block, error)
	ListBlockchains() ([]*data.Blockchain, error)
	ListBlocks(id string) ([]*data.Block, error)
	DeleteBlockchain(id string) error
}