package userspace

type ReqUserLogin struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type RspUserLogin struct {
	Token string `json:"token"`
}

type ReqUUID struct {
	UUID string `json:"uuid"`
}

type Ant struct {
	Id     int64   `json:"id"`
	Speed  float32 `json:"speed"`
	Power  float32 `json:"power"`
	Attact float32 `json:"attact"`
	Blood  float32 `json:"blood"`
}

//////// search ant //////////////

type ReqUpdateSearchAnt struct {
	Ant
	ReqUUID
	Gold int64 `json:"gold"`
}
type RspUpdateSearchAnt struct {
	Ant
}
type RspUserSearchScene struct {
	Gold      int64  `json:"gold"`
	TouchHigh int    `json:"touch_high"`
	TouchLow  int    `json:"touch_low"`
	Ants      []*Ant `json:"ants"`
}

///////////////////////////////////////

//////// defen ant////////////////////

type ReqUpdateDefenAnt struct {
	Ant
	ReqUUID
	Gold int64 `json:"gold"`
}
type RspUpdateDefenAnt struct {
	Ant
}
type RspUserDefenScene struct {
	Ants []*Ant `json:"ants"`
}

////////// queen ant ///////////////////
type RspUserQueenAnts struct {
	Ants []*Ant `json:"ants"`
}

////////////////////////////////////////
