package session

type TxnContext struct {
	active bool
}

func NewTxnContext() *TxnContext {
	return &TxnContext{
		active: false,
	}
}

func (t *TxnContext) BeginTransaction() {
	t.active = true
}

func (t *TxnContext) InTransaction() bool {
	return t.active
}

func (t *TxnContext) EndTransaction() {
	t.active = false
}
