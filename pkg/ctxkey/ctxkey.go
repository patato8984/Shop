package ctxkey

type key int
type role int
type txKey int

const (
	UserIDKey      key   = iota
	Role           role  = iota
	TransactionKey txKey = iota
)
