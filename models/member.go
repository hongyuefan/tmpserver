package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type Member struct {
	UID        int64   `orm:"column(uid);auto"`
	UserName   string  `orm:"column(username);size(20)"`
	Email      string  `orm:"column(email);size(50)"`
	Mobile     string  `orm:"column(mobile);size(11)"`
	PassWord   string  `orm:"column(password);size(20)"`
	UserIP     string  `orm:"column(user_ip);size(255)"`
	Img        string  `orm:"column(img);size(255)"`
	QianMing   string  `orm:"column(qianming);size(255)"`
	GroupId    int32   `orm:"column(groupid);"`
	AddGroup   string  `orm:"column(addgroup);size(255)"`
	Money      float64 `orm:"column(money)"`
	EmailCode  string  `orm:"column(emailcode);size(21)"`
	MobileCode string  `orm:"column(mobilecode);size(21)"`
	OtherCode  int32   `orm:"column(othercode);"`
	PassCode   string  `orm:"column(passcode);size(21)"`
	RegKey     string  `orm:"column(reg_key);size(100)"`
	Score      int32   `orm:"column(score);"`
	JingYan    int32   `orm:"column(jingyan);"`
	YaoQing    int64   `orm:"column(yaoqing);"`
	Band       string  `orm:"column(band);size(255)"`
	Time       int64   `orm:"column(time);"`
	Headimg    string  `orm:"column(headimg);size(255)"`
	WxId       string  `orm:"column(wxid);size(255)"`
	TypeId     int32   `orm:"column(typeid);"`
	AutoUser   int32   `orm:"column(auto_user);"`
	Level      int32   `orm:"column(level);"`
}

func (t *Member) TableName() string {
	return "go_member"
}

func init() {
	orm.RegisterModel(new(Member))
}

func AddMember(m *Member) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func UpdateMember(m *Member, cols ...string) (err error) {
	o := orm.NewOrm()
	v := Member{UID: m.UID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, cols...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func UpdateMemberByUserName(m *Member, cols ...string) (err error) {
	o := orm.NewOrm()
	v := Member{UserName: m.UserName}
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
