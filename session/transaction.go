package session

type TxnContext struct {
	active bool

	queuedCommands []QueuedCommand

	// Transaction error flag; set if an invalid command was issued after MULTI. Causes EXEC to abort
	isDirty bool
}

type QueuedCommand struct {
	name string
	args []string
}

func NewTxnContext() *TxnContext {
	return &TxnContext{
		active:         false,
		queuedCommands: make([]QueuedCommand, 0),
		isDirty:        false,
	}
}

func (t *TxnContext) BeginTransaction() {
	t.active = true
}

func (t *TxnContext) QueueCommand(name string, args []string) {
	// FIFO
	t.queuedCommands = append(t.queuedCommands, QueuedCommand{name, args})
}

func (t *TxnContext) InTransaction() bool {
	return t.active
}

func (t *TxnContext) MarkDirty() {
	t.isDirty = true
}

func (t *TxnContext) EndTransaction() {
	t.active = false
}
