package models

import (
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	if err := orm.RegisterDataBase("default", "mysql", "root:30bb9f7ed8@tcp(39.105.5.188:3306)/cyyungoucms"); err != nil {
		panic(err)
	}
}

func TestGetShopLists(t *testing.T) {

	query := make(map[string]string)

	query["s_id"] = "5"
	query["s_cid"] = "2"

	ml, err := GetShopCodes(query, []string{}, []string{"id"}, []string{"asc"}, 0, 100)

	if err != nil {
		panic(err)
	}

	t.Log(ml)
}
