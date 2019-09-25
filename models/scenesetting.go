package models

import (
	"github.com/astaxie/beego/orm"
)

type SceneSetting struct {
	Id           int64 `orm:"column(id);auto"`
	TouchHigh    int   `orm:"column(touch_high);"`
	TouchLow     int   `orm:"column(touch_low);"`
	MonsterLarge int   `orm:"column(monster_large);"`
	MonsterSmall int   `orm:"column(monster_small);"`
}

func (t *SceneSetting) TableName() string {
	return "scenesetting"
}

func init() {
	orm.RegisterModel(new(SceneSetting))
}

func GetSceneSettingByUUID(s *SceneSetting) error {
	return orm.NewOrm().Read(s, "uuid")
}

func InsertSceneSetting(s *SceneSetting) (int64, error) {
	return orm.NewOrm().Insert(s)
}

func UpdateSceneSetting(s *SceneSetting, cols ...string) (int64, error) {
	return orm.NewOrm().Update(s, cols...)
}

func GetSceneSettingById(s *SceneSetting) error {
	return orm.NewOrm().Read(s)
}
