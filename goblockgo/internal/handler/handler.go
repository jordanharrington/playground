package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jordanharrington/playground/goblockgo/internal/service"
	"net/http"
)

type CreateBlockchainRequest struct {
	Name string `json:"name"`
}

type CreateBlockchainResponse struct {
	ID string `json:"id"`
}

type AddBlockRequest struct {
	Data string `json:"data"`
}

type BlockResponse struct {
	Timestamp     int64  `json:"timestamp"`
	Data          string `json:"data"`
	PrevBlockHash string `json:"prevBlockHash"`
	Hash          string `json:"hash"`
}

type ListBlocksResponse struct {
	Blocks []BlockResponse `json:"blocks"`
}

type DeleteBlockchainResponse struct {
	Status string `json:"status"`
}

// s is the singleton instance of the GoBlockGo service
var s service.GoBlockGo

func Route(goBlockGo service.GoBlockGo) http.Handler {
	s = goBlockGo

	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	r.Route("/v1", func(r chi.Router) {
		// Create a new blockchain
		r.Post("/", createBlockchain)
		// Append a block to an existing blockchain
		r.Post("/{id}", addBlock)
		// List all blockchains
		r.Get("/", listBlockchains)
		// List all Blocks in the blockchain
		r.Get("/{id}", listBlocks)
		// Get a block in a blockchain
		r.Get("/{id}/{hash}", getBlock)
		// Delete a blockchain
		r.Delete("/{id}", deleteBlockchain)
	})

	return r
}

func createBlockchain(w http.ResponseWriter, r *http.Request) {
	var req CreateBlockchainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bc, err := s.CreateBlockchain(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreateBlockchainResponse{ID: bc.ID}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func addBlock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req AddBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	block, err := s.AddBlock(id, req.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := BlockResponse{
		Timestamp:     block.Timestamp,
		Data:          string(block.Data),
		PrevBlockHash: fmt.Sprintf("%x", block.PrevBlockHash),
		Hash:          fmt.Sprintf("%x", block.Hash),
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func listBlockchains(w http.ResponseWriter, _ *http.Request) {
	blockchains, err := s.ListBlockchains()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ids []string
	for _, bc := range blockchains {
		ids = append(ids, bc.ID)
	}

	err = json.NewEncoder(w).Encode(ids)
	if err != nil {
		return
	}
}

func listBlocks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	blocks, err := s.ListBlocks(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var resp ListBlocksResponse
	for _, block := range blocks {
		resp.Blocks = append(resp.Blocks, BlockResponse{
			Timestamp:     block.Timestamp,
			Data:          string(block.Data),
			PrevBlockHash: fmt.Sprintf("%x", block.PrevBlockHash),
			Hash:          fmt.Sprintf("%x", block.Hash),
		})
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func getBlock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	hash := chi.URLParam(r, "hash")

	block, err := s.GetBlock(id, hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := BlockResponse{
		Timestamp:     block.Timestamp,
		Data:          string(block.Data),
		PrevBlockHash: fmt.Sprintf("%x", block.PrevBlockHash),
		Hash:          fmt.Sprintf("%x", block.Hash),
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func deleteBlockchain(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := s.DeleteBlockchain(id)
	if err != nil {
		// Assume the error is due to the blockchain not being found
		http.Error(w, "Blockchain not found", http.StatusNotFound)
		return
	}

	response := DeleteBlockchainResponse{
		Status: "Deleted",
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
