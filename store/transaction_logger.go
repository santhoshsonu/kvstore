package store

type EventType byte

const (
	_                  = iota
	EventPut EventType = iota
	EventDelete
)

const LOG_FORMAT = "%d\t%d\t%s\t%q"

type Event struct {
	Sequence  uint64    // a unique record ID
	EventType EventType // action taken
	Key       string    // key affected by this transaction
	Value     string    // value of the transaction
}

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}
