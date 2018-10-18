package mintcoin

const (
	TypeOK                       uint32 = 0
	TypeAddressInvalid           uint32 = 1
	TypeJsonEncodingError        uint32 = 2
	TypeExecuteError             uint32 = 3
	TypeVerifyError              uint32 = 4
	CodeTypeValidateSendTxError  uint32 = 5
	CodeTypeValidateReceiveError uint32 = 6
	CodeAccountNotFoundError     uint32 = 7
)

const (
	OperatorAdd uint32 = 0
	OperatorSub uint32 = 1
)
