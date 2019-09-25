package userspace

import (
	"encoding/json"

	"github.com/hongyuefan/tmpserver/models"
)

func UserLogin(reqMsg []byte) (interface{}, error) {
	user := new(ReqUserLogin)
	if err := json.Unmarshal(reqMsg, user); err != nil {
		return nil, err
	}
	return &RspUserLogin{Token: "123"}, nil
}

func GetUserSearchScene(reqMsg []byte) (interface{}, error) {
	uuid := new(ReqUUID)
	if err := json.Unmarshal(reqMsg, uuid); err != nil {
		return nil, err
	}
	mQuery := make(map[string]string, 0)
	mQuery["uuid"] = uuid.UUID
	mQuery["scene"] = "1"

	mls, err := models.GetAnts(mQuery, []string{}, []string{"speed"}, []string{"asc"}, 0, 100)
	if err != nil {
		return nil, err
	}
	rsp := new(RspUserSearchScene)
	for _, ant := range mls {
		sAnt := new(Ant)
		sAnt.Id = ant.(models.Ants).Id
		sAnt.Power = ant.(models.Ants).Power
		sAnt.Speed = ant.(models.Ants).Speed
		rsp.Ants = append(rsp.Ants, sAnt)
	}
	user := &models.Users{UUID: uuid.UUID}
	if err := models.GetUserByUUID(user); err != nil {
		return nil, err
	}
	rsp.Gold = user.Gold
	rsp.TouchHigh = user.TouchHigh
	rsp.TouchLow = user.TouchLow
	return rsp, nil
}

func UpdateUserSearchAnt(reqMsg []byte) (interface{}, error) {
	reqUp := new(ReqUpdateSearchAnt)
	if err := json.Unmarshal(reqMsg, reqUp); err != nil {
		return nil, err
	}
	ant := new(models.Ants)
	ant.Id = reqUp.Id
	ant.Speed = reqUp.Speed
	ant.Power = reqUp.Power
	if ant.Speed > 0 {
		if _, err := models.UpdateAnt(ant, "speed"); err != nil {
			return nil, err
		}
	}
	if ant.Power > 0 {
		if _, err := models.UpdateAnt(ant, "power"); err != nil {
			return nil, err
		}
	}
	user := new(models.Users)
	user.UUID = reqUp.UUID
	if err := models.GetUserByUUID(user); err != nil {
		return nil, err
	}
	user.Gold = reqUp.Gold
	if _, err := models.UpdateUser(user, "golds"); err != nil {
		return nil, err
	}
	rspUp := new(RspUpdateSearchAnt)
	rspUp.Id = reqUp.Id
	rspUp.Speed = reqUp.Speed
	rspUp.Power = reqUp.Power
	return rspUp, nil
}

func GetUserDefenScene(reqMsg []byte) (interface{}, error) {
	uuid := new(ReqUUID)
	if err := json.Unmarshal(reqMsg, uuid); err != nil {
		return nil, err
	}
	mQuery := make(map[string]string, 0)
	mQuery["uuid"] = uuid.UUID
	mQuery["Scene"] = "2"

	mls, err := models.GetAnts(mQuery, []string{}, []string{"attact"}, []string{"asc"}, 0, 100)
	if err != nil {
		return err, nil
	}
	rsp := new(RspUserDefenScene)
	for _, ml := range mls {
		ant := new(Ant)
		ant.Id = ml.(models.Ants).Id
		ant.Attact = ml.(models.Ants).Attact
		ant.Speed = ml.(models.Ants).Speed
		ant.Blood = ml.(models.Ants).Blood
		rsp.Ants = append(rsp.Ants, ant)
	}
	return rsp, nil
}

func UpdateUserDefenAnt(reqMsg []byte) (interface{}, error) {
	reqUp := new(ReqUpdateDefenAnt)
	if err := json.Unmarshal(reqMsg, reqUp); err != nil {
		return nil, err
	}
	ant := new(models.Ants)
	ant.Id = reqUp.Id
	ant.Speed = reqUp.Speed
	ant.Attact = reqUp.Attact
	ant.Blood = reqUp.Blood
	if ant.Speed > 0 {
		if _, err := models.UpdateAnt(ant, "speed"); err != nil {
			return nil, err
		}
	}
	if ant.Attact > 0 {
		if _, err := models.UpdateAnt(ant, "attact"); err != nil {
			return nil, err
		}
	}
	if ant.Blood > 0 {
		if _, err := models.UpdateAnt(ant, "blood"); err != nil {
			return nil, err
		}
	}
	user := new(models.Users)
	user.UUID = reqUp.UUID
	if err := models.GetUserByUUID(user); err != nil {
		return nil, err
	}
	user.Gold = reqUp.Gold
	if _, err := models.UpdateUser(user, "golds"); err != nil {
		return nil, err
	}
	rspUp := new(RspUpdateDefenAnt)
	rspUp.Id = reqUp.Id
	rspUp.Speed = reqUp.Speed
	rspUp.Attact = reqUp.Attact
	rspUp.Blood = reqUp.Blood
	return rspUp, nil
}

func GetUserQueenScene(reqMsg []byte) (interface{}, error) {
	uuid := new(ReqUUID)
	if err := json.Unmarshal(reqMsg, uuid); err != nil {
		return nil, err
	}
	mQuery := make(map[string]string, 0)
	mQuery["uuid"] = uuid.UUID
	mQuery["Scene"] = "3"

	mls, err := models.GetAnts(mQuery, []string{}, []string{"level"}, []string{"asc"}, 0, 100)
	if err != nil {
		return err, nil
	}
	rsp := new(RspUserQueenAnts)
	for _, ant := range mls {
		sAnt := new(Ant)
		sAnt.Id = ant.(models.Ants).Id
		sAnt.Speed = ant.(models.Ants).Speed
		rsp.Ants = append(rsp.Ants, sAnt)
	}
	return rsp, nil
}

func UpdateUserQueenAnt(reqMsg []byte) (interface{}, error) {
	reqUp := new(ReqUpdateDefenAnt)
	if err := json.Unmarshal(reqMsg, reqUp); err != nil {
		return nil, err
	}
	ant := new(models.Ants)
	ant.Id = reqUp.Id
	ant.Speed = reqUp.Speed
	ant.Attact = reqUp.Attact
	ant.Blood = reqUp.Blood
	if ant.Speed > 0 {
		if _, err := models.UpdateAnt(ant, "speed"); err != nil {
			return nil, err
		}
	}
	if ant.Attact > 0 {
		if _, err := models.UpdateAnt(ant, "attact"); err != nil {
			return nil, err
		}
	}
	if ant.Blood > 0 {
		if _, err := models.UpdateAnt(ant, "blood"); err != nil {
			return nil, err
		}
	}
	user := new(models.Users)
	user.UUID = reqUp.UUID
	if err := models.GetUserByUUID(user); err != nil {
		return nil, err
	}
	user.Gold = reqUp.Gold
	if _, err := models.UpdateUser(user, "golds"); err != nil {
		return nil, err
	}
	rspUp := new(RspUpdateDefenAnt)
	rspUp.Id = reqUp.Id
	rspUp.Speed = reqUp.Speed
	rspUp.Attact = reqUp.Attact
	rspUp.Blood = reqUp.Blood
	return rspUp, nil
}
