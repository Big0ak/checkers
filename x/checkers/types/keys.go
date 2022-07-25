package types

const (
	// ModuleName defines the module name
	ModuleName = "checkers"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_checkers"

	// this line is used by starport scaffolding # ibc/keys/name
)

// this line is used by starport scaffolding # ibc/keys/port

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	NextGameKey = "NextGame-value-"
)

// события созданые игой
const (
	StoredGameKey          = "StoredGame-value-"
	StoredGameEventKey     = "NewGameCreated" // указывает, на какой ключ обращать внимание
	StoredGameEventCreator = "Creator"
	StoredGameEventIndex   = "Index" // какая игра актуальна
	StoredGameEventRed     = "Red"   // актуально ли для этого игрока
	StoredGameEventBlack   = "Black" // актуально ли для этого игрока
)

//события созданные игроком
const (
	PlayMoveEventKey       = "MovePlayed"
	PlayMoveEventCreator   = "Creator"
	PlayMoveEventIdValue   = "IdValue"
	PlayMoveEventCapturedX = "CapturedX"
	PlayMoveEventCapturedY = "CapturedY"
	PlayMoveEventWinner    = "Winner"
)
