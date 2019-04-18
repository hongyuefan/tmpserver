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

type RspDetect struct {
	IsSuccess bool        `json:"isSuccess"`
	Result    interface{} `json:"result"`
}

type Detect struct {
	Age        int64         `json:"age"`
	Gender     string        `json:"gender"`
	Race       string        `json:"race"`
	Emotion    string        `json:"emotion"`
	Expression string        `json:"expression"`
	IsGlass    bool          `json:"glass"`
	IsAngle    bool          `json:"angle"`
	Beauty     float64       `json:"beauty"`
	FaceType   string        `json:"face_type"`
	Descrips   []Description `json:"descriptions"`
}

type Description struct {
	Title   string `json:"title"`
	Content string `json:"content"`
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
	EYEBROW_HEIGHT     int64 = 18
	MOUTH_ANGLE        int64 = 19
	EYEBROW_ANGLE_MID  int64 = 20
	MOUTH_LIPS_EQUAL   int64 = 30
	FACE_ANGRY         int64 = 31
	FACE_SHAP          int64 = 32
)
