package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type Member struct {
	Id       int64  `orm:"column(id);auto"`
	OpenId   string `orm:"column(openid);size(128)"`
	AppId    string `orm:"column(appid);size(128)"`
	Nick     string `orm:"column(nick);size(128)"`
	Gender   int    `json:"gender"`
	Province string `json:"province"`
	Count    int64  `orm:"column(counts);"`
	Time     int64  `orm:"column(time);"`
}

func (t *Member) TableName() string {
	return "face_user"
}

func init() {
	orm.RegisterModel(new(Member))
}

func AddMember(m *Member) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func UpdateMemberByOpenId(m *Member, cols ...string) (err error) {
	o := orm.NewOrm()
	v := Member{OpenId: m.OpenId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, cols...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func GetMember(member *Member, cols ...string) error {
	return orm.NewOrm().Read(member, cols...)
}
