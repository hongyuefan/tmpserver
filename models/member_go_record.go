package models

import (
	"errors"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type MgoRecord struct {
	ID         int64   `orm:"column(id)"`
	Code       string  `orm:"column(code);size(20)"`
	UserName   string  `orm:"column(username);size(30)"`
	Uphoto     string  `orm:"column(uphoto);size(255)"`
	UID        int64   `orm:"column(uid);"`
	ShopID     int64   `orm:"column(shopid);"`
	ShopName   string  `orm:"column(shopname);size(255)"`
	ShopQiShu  int32   `orm:"column(shopqishu);"`
	GoNumber   int32   `orm:"column(gonumber);"`
	GouCode    string  `orm:"column(goucode);size(2048)"`
	MoneyCount float64 `orm:"column(moneycount);"`
	Status     string  `orm:"column(status);size(32)"`
	Time       string  `orm:"column(time);size(21)"`
}

func (t *MgoRecord) TableName() string {
	return "go_member_go_record"
}

func init() {
	orm.RegisterModel(new(MgoRecord))
}

func AddMgoRecord(m *MgoRecord) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func UpdateMgoRecord(m *MgoRecord, cols ...string) (err error) {
	o := orm.NewOrm()
	v := MgoRecord{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		if _, err = o.Update(m, cols...); err != nil {
			return err
		}
	}
	return
}

func GetMgoRecords(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (ml []interface{}, err error) {

	o := orm.NewOrm()
	qs := o.QueryTable(new(MgoRecord))
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

	var l []MgoRecord
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