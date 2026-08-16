package exit

type Code int

const (
	Success         Code = 0
	Internal        Code = 1
	InvalidInput    Code = 2
	PreDispatch     Code = 3
	Rejected        Code = 4
	Uncertain       Code = 5
	SafetyConflict  Code = 6
	ReceiptDelivery Code = 7
)
