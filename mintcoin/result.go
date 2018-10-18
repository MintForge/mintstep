package mintcoin

type Result struct {
	Msg  string `json:"msg"`
	Code uint32 `json:"code"`
}

func (res Result) IsErr() bool {
	return res.Code != TypeOK
}
