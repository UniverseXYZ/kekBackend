package barn

const (
	DepositEvent                = "Deposit"
	WithdrawEvent               = "Withdraw"
	LockEvent                   = "Lock"
	DelegateEvent               = "Delegate"
	DelegatePowerIncreasedEvent = "DelegatePowerIncreased"
	DelegatePowerDecreasedEvent = "DelegatePowerDecreased"
)

const (
	DEPOSIT = iota
	WITHDRAW
)

const (
	DELEGATE_START = iota
	DELEGATE_STOP
)

const (
	DELEGATE_INCREASE = iota
	DELEGATE_DECREASE
)

const ZeroAddress = "0x0000000000000000000000000000000000000000"
