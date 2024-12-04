package main

import (
	"context"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/dgraph-io/badger/v4"
)

type App struct {
	db           *badger.DB
	onGoingBlock *badger.Txn
}

var _ abcitypes.Application = (*App)(nil)

func NewApp(db *badger.DB) *App {
	return &App{db: db}
}

// Input: a transaction
// Output: a response code, 0 means the transaction format is valid
//
// This method is called whenever a user sends a new transaction to the node.
// If the transaction is invalid, it will not be gossiped to other nodes, and
// will never be included in a proposal.
func (app *App) CheckTx(_ context.Context, check *abcitypes.RequestCheckTx) (*abcitypes.ResponseCheckTx, error) {
	// CometBFT doesn't know what's inside a transaction, it just passes
	// it as a byte slice.

	// Unmarshal the bytes into our custom Tx type.
	tx, err := UnmarshalTx(check.Tx)
	if err != nil {
		// If the transaction doesn't even follow the right format, is invalid.
		return &abcitypes.ResponseCheckTx{
			Code: 112233,
			Log:  err.Error(),
		}, nil
	}

	// Here we can perform additional checks other than just the format. For
	// example we could enforce that "key" is at least 3 characters long.
	if err := tx.Valid(); err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: 112233,
			Log:  err.Error(),
		}, nil
	}

	return &abcitypes.ResponseCheckTx{
		Code: abcitypes.CodeTypeOK,
	}, nil
}

// Input: a block (list of transactions)
// Output: a list of responses, one for each transaction
//
// The node should not persist any change yet, but it should be prepared to do
// so when Commit() is called.
func (app *App) FinalizeBlock(_ context.Context, req *abcitypes.RequestFinalizeBlock) (*abcitypes.ResponseFinalizeBlock, error) {
	var results = make([]*abcitypes.ExecTxResult, len(req.Txs))

	app.onGoingBlock = app.db.NewTransaction(true)
	for i, txBytes := range req.Txs {
		// CometBFT doesn't know what's inside a transaction, it just passes
		// it as a byte slice.

		// Unmarshal the bytes into our custom Tx type.
		tx, err := UnmarshalTx(txBytes)
		if err != nil {
			return nil, err
		}

		// Check if the transaction is valid.
		if err := tx.Valid(); err != nil {
			results[i] = &abcitypes.ExecTxResult{
				Code: 112233,
				Log:  err.Error(),
			}
			continue
		}

		log.Printf("Adding key %s with value %s", tx.Key, tx.Val)

		// Write the new data to the database (we're only writing it to the
		// badger's transaction, it's not persisted yet).
		if err := app.onGoingBlock.Set(tx.Key, tx.Val); err != nil {
			log.Panicf("Error writing to database, unable to execute tx: %v", err)
		}

		log.Printf("Successfully added key %s with value %s", tx.Key, tx.Val)

		results[i] = &abcitypes.ExecTxResult{}
	}

	return &abcitypes.ResponseFinalizeBlock{
		TxResults: results,
	}, nil
}

// Commit is called when a new block has been committed, included in the chain.
// The application must now persist all the changes to the disk.
func (app App) Commit(_ context.Context, commit *abcitypes.RequestCommit) (*abcitypes.ResponseCommit, error) {
	err := app.onGoingBlock.Commit()
	return &abcitypes.ResponseCommit{}, err
}

// Query is called when a user wants to query the application state. This is a
// read-only operation that is not part of the consensus protocol.
func (app *App) Query(_ context.Context, req *abcitypes.RequestQuery) (*abcitypes.ResponseQuery, error) {
	resp := abcitypes.ResponseQuery{Key: req.Data}

	dbErr := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(req.Data)
		if err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
			resp.Log = "key does not exist"
			return nil
		}

		return item.Value(func(val []byte) error {
			resp.Log = "exists"
			resp.Value = val
			return nil
		})
	})
	if dbErr != nil {
		log.Panicf("Error reading database, unable to execute query: %v", dbErr)
	}

	return &resp, nil
}

// The following methods are left to the bare minimum to just make the app work
// They allow the app to tune a bit more the behaviour of the node and the
// consensus protocol itself.

func (app *App) Info(_ context.Context, info *abcitypes.RequestInfo) (*abcitypes.ResponseInfo, error) {
	return &abcitypes.ResponseInfo{}, nil
}

func (app *App) InitChain(_ context.Context, chain *abcitypes.RequestInitChain) (*abcitypes.ResponseInitChain, error) {
	return &abcitypes.ResponseInitChain{}, nil
}

func (app *App) PrepareProposal(_ context.Context, proposal *abcitypes.RequestPrepareProposal) (*abcitypes.ResponsePrepareProposal, error) {
	return &abcitypes.ResponsePrepareProposal{Txs: proposal.Txs}, nil
}

func (app *App) ProcessProposal(_ context.Context, proposal *abcitypes.RequestProcessProposal) (*abcitypes.ResponseProcessProposal, error) {
	return &abcitypes.ResponseProcessProposal{Status: abcitypes.ResponseProcessProposal_ACCEPT}, nil
}

func (app *App) ListSnapshots(_ context.Context, snapshots *abcitypes.RequestListSnapshots) (*abcitypes.ResponseListSnapshots, error) {
	return &abcitypes.ResponseListSnapshots{}, nil
}

func (app *App) OfferSnapshot(_ context.Context, snapshot *abcitypes.RequestOfferSnapshot) (*abcitypes.ResponseOfferSnapshot, error) {
	return &abcitypes.ResponseOfferSnapshot{}, nil
}

func (app *App) LoadSnapshotChunk(_ context.Context, chunk *abcitypes.RequestLoadSnapshotChunk) (*abcitypes.ResponseLoadSnapshotChunk, error) {
	return &abcitypes.ResponseLoadSnapshotChunk{}, nil
}

func (app *App) ApplySnapshotChunk(_ context.Context, chunk *abcitypes.RequestApplySnapshotChunk) (*abcitypes.ResponseApplySnapshotChunk, error) {
	return &abcitypes.ResponseApplySnapshotChunk{Result: abcitypes.ResponseApplySnapshotChunk_ACCEPT}, nil
}

func (app App) ExtendVote(_ context.Context, extend *abcitypes.RequestExtendVote) (*abcitypes.ResponseExtendVote, error) {
	return &abcitypes.ResponseExtendVote{}, nil
}

func (app *App) VerifyVoteExtension(_ context.Context, verify *abcitypes.RequestVerifyVoteExtension) (*abcitypes.ResponseVerifyVoteExtension, error) {
	return &abcitypes.ResponseVerifyVoteExtension{}, nil
}
