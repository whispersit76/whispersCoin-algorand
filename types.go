package whispers76Coin

// Initialize coin ledger.
type Initialize struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Supply designates how much the coin creator will start with.
	Supply uint64 `codec:"supply"`
}

// Transfer represents an coin transfer.
type Transfer struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Destination designates who is sending examplecoin.
	Source string `codec:"source"`

	// Destination designates who is receiving coin.
	Destination string `codec:"destination"`

	// Amount designates how much coin is being transferred.
	Amount uint64 `codec:"Amount"`
}

// NoteFieldType indicates a type of coin message encoded into a
// transaction's Note field.
type NoteFieldType string

const (
	// NoteInitialize indicates an Initialize message.
	NoteInitialize NoteFieldType = "i"

	// NoteTransfer indicates a Transfer message.
	NoteTransfer NoteFieldType = "t"
)

// NoteField is the struct that represents an coin message.
type NoteField struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Type indicates which type of a message this is
	Type NoteFieldType `codec:"type"`

	// Initialize, for NoteInitialize type
	Initialize Initialize `codec:"i"`

	// Transfer, for NoteTransfer type
	Transfer Transfer `codec:"t"`
}
