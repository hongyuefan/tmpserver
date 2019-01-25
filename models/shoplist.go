package models

import (
	"errors"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type ShopList struct {
	ID        int64   `orm:"column(id);auto"`
	SID       int64   `orm:"column(sid);"`
	Total     int64   `orm:"column(zongrenshu);"`
	Title     string  `orm:"column(title);size(100)"`
	CanY      int64   `orm:"column(canyurenshu);"`
	SanY      int64   `orm:"column(shenyurenshu);"`
	QiShu     int32   `orm:"column(qishu);"`
	MaxQS     int32   `orm:"column(maxqishu);"`
	QUid      int64   `orm:"column(q_uid)"`
	QUser     string  `orm:"column(q_user);size(2048)"`
	QUserCode string  `orm:"column(q_user_code);size(32)"`
	YJiaGe    float64 `orm:"column(yunjiage)"`
}

func (t *ShopList) TableName() string {
	return "go_shoplist"
}

func init() {
	orm.RegisterModel(new(ShopList))
}

func AddShopList(m *ShopList) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func UpdateShopList(m *ShopList, cols ...string) (err error) {
	o := orm.NewOrm()
	v := ShopList{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		if _, err = o.Update(m, cols...); err != nil {
			return err
		}
	}
	return
}

func GetShopLists(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (ml []interface{}, err error) {

	o := orm.NewOrm()
	qs := o.QueryTable(new(ShopList))
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

	var l []ShopList
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
