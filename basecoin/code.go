package basecoin

const (
	TypeOK                       uint32 = 0
	TypeAddressInvalid           uint32 = 1
	TypeJsonEncodingError        uint32 = 2
	TypeExecuteError             uint32 = 3
	CodeTypeValidateSendTxError  uint32 = 4
	CodeTypeValidateReceiveError uint32 = 5
	CodeAccountNotFoundError     uint32 = 6
)

const (
	OperatorAdd uint32 = 0
	OperatorSub uint32 = 1
)
