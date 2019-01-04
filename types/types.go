package types

const (
	PAY_METAMASK = 1
	PAY_ERCODE   = 2

	STATUS_WAITING = 0
	STATUS_SUCCESS = 1
	STATUS_FAILED  = -1
)

type ReqAddMoney struct {
	Type     int    `json:"type"`
	UserName string `json:"username"`
	UserId   int64  `json:"userid"`
	Amount   string `json:"amount"`
	Hash     string `json:"hash"`
	Address  string `json:"address"`
	Time     int64  `json:"time"`
}
