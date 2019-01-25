package automan

import (
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func TestPares(t *testing.T) {

	if err := orm.RegisterDataBase("default", "mysql", "root:30bb9f7ed8@tcp(39.105.5.188:3306)/cyyungoucms"); err != nil {
		panic(err)
	}

	auto := NewAutoMan(10)

	auto.OnStart()
}
