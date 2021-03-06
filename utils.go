package whispers76coin

import (
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
)

// BytesBase64 is a base64-encoded binary blob (i.e., []byte), for
// use with text encodings like JSON.
type BytesBase64 []byte

// BuildInitializeNote takes in the desired supply and produces a blob for your note field
func BuildInitializeNote(supply uint64) (initializeBlob BytesBase64) {
	initializeBlob = BytesBase64(msgpack.Encode(NoteField{
		Type: NoteInitialize,
		Initialize: Initialize{
			Supply: supply,
		},
	}))
	return
}

// BuildTransferNote takes in the desired recipient as well as amount to send, and produces a blob for your note field
func BuildTransferNote(amount uint64, from, to string) (transferBlob BytesBase64) {
	transferBlob = BytesBase64(msgpack.Encode(NoteField{
		Type: NoteTransfer,
		Transfer: Transfer{
			Amount:      amount,
			Source:      from,
			Destination: to,
		},
	}))
	return
}

// ProcessInitialize accepts the current ledger, the initialize message, and the wrapping txn.
// it updates the ledger, and returns an error if something went wrong.
func ProcessInitialize(curState map[string]uint64, initialize Initialize, wrappingTxn models.Transaction) (map[string]uint64, error) {
	if len(curState) != 0 {
		return curState, fmt.Errorf("attempted to process an initialize message against a ledger that was already initialized")
	}
	curState[wrappingTxn.From] = initialize.Supply
	return curState, nil
}

// ProcessTransfer accepts the current ledger, the transfer message, and the wrapping txn.
// it updates the ledger, and returns an error if something went wrong.
func ProcessTransfer(curState map[string]uint64, transfer Transfer, wrappingTxn models.Transaction) (map[string]uint64, error) {
	if transfer.Source != wrappingTxn.From {
		return curState, fmt.Errorf("transaction submitted by %s tries to spend %s's whispers76coin", wrappingTxn.Payment.To, transfer.Source)
	}
	senderBalance, exists := curState[transfer.Source]
	if !exists {
		return curState, fmt.Errorf("sender %v does not exist in the ledger", transfer.Source)
	}
	if transfer.Amount > senderBalance {
		return curState, fmt.Errorf("sender %v is trying to spend %d whispers76coin, greater than balance %d", transfer.Source, transfer.Amount, senderBalance)
	}
	curState[transfer.Source] = senderBalance - transfer.Amount
	receiverBalance := curState[transfer.Destination]
	curState[transfer.Destination] = receiverBalance + transfer.Amount

	return curState, nil
}

// ReadTransferNote takes in the desired blob notefield and produces amount, to, from
func ReadTransferNote(transferBlob BytesBase64, transfer Transfer) error {
	err := msgpack.Decode(transferBlob, transfer)

	if err != nil {
		return err
	}
	return nil
}
