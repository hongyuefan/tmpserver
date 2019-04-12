package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type Picture struct {
	Id       int64  `orm:"column(id);auto"`
	OpenId   int64  `orm:"column(open_id);size(128)"`
	FilePath string `orm:"column(file_path);size(260)"`
	FileId   string `orm:"column(file_id);size(128)"`
	Result   string `orm:"column(result);size(2048);"`
	Time     int64  `orm:"column(time);"`
}

func (t *Picture) TableName() string {
	return "face_picture"
}

func init() {
	orm.RegisterModel(new(Picture))
}

func AddPicture(m *Picture) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func UpdatePictureByFid(m *Picture, cols ...string) (err error) {
	o := orm.NewOrm()
	v := Picture{FileId: m.FileId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, cols...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func GetPictureBy(m *Picture, cols ...string) error {
	return orm.NewOrm().Read(m, cols...)
}
