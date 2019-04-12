package types

const (
	PAY_METAMASK = 1
	PAY_ERCODE   = 2

	STATUS_WAITING = 0
	STATUS_SUCCESS = 1
	STATUS_FAILED  = -1
)

type ReqMember struct {
	OpenId string `json:"openid"`
	AppId  string `json:"appid"`
	Count  int64  `json:"count"`
	Time   int64  `json:"time"`
}

const (
	EYE_LENGTH         int64 = 1
	EYE_HEIGHT         int64 = 2
	EYE_ANGLE          int64 = 3
	EYE_TO_EYE         int64 = 4
	EYE_TO_EYEBROW     int64 = 5
	EYEBROW_LENGTH     int64 = 6
	EYEBROW_TO_EYEBROW int64 = 7
	EYEBROW_ANGLE      int64 = 8
	EYEBROW_MAX_RATIO  int64 = 9
	NOSE_WIDTH         int64 = 10
	NOSE_LENGTH        int64 = 11
	NOSE_RATIO         int64 = 12
	PHILTRUM_LENGTH    int64 = 13
	MOUTH_WIDTH        int64 = 14
	MOUTH_THICKNESS    int64 = 15
	MOUTH_LIPS_RATIO   int64 = 16
	CHIN_WIDTH         int64 = 17
)
