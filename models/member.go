package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

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

func GetMembers(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Member))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Member
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}
