package automan

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/hongyuefan/tmpserver/models"
)

type AutoMan struct {
	intervel  int64
	mapCity   map[string]string
	chanClose chan bool
}

func NewAutoMan(intervel int64) *AutoMan {
	return &AutoMan{
		intervel:  intervel,
		chanClose: make(chan bool, 0),
		mapCity:   make(map[string]string, 0),
	}
}

func (s *AutoMan) OnStart() {

	if s.intervel < 300 {
		s.intervel = 367
	}

	ticker := time.NewTicker(time.Second * time.Duration(s.intervel))

	for {
		select {
		case <-ticker.C:
			s.handler()
		case <-s.chanClose:
			ticker.Stop()
			return
		}
	}
}

func (s *AutoMan) handler() {

	man, err := s.selectMan()
	if err != nil {
		fmt.Println("selectMan error:", err.Error())
		return
	}

	shopList, err := s.selectShop()
	if err != nil {
		fmt.Println("selectShop error:", err.Error())
		return
	}

	code, err := s.getShopCode(shopList.ID)
	if err != nil {
		fmt.Println("getShopCode error:", err.Error())
		return
	}

	_, err = models.AddMgoRecord(&models.MgoRecord{
		Code:        "A" + fmt.Sprintf("%v", time.Now().UnixNano()/100),
		UserName:    man.UserName,
		Uphoto:      man.Img,
		UID:         man.UID,
		CodeTmp:     0,
		ShopID:      shopList.ID,
		ShopName:    shopList.Title,
		ShopQiShu:   shopList.QiShu,
		GoNumber:    1,
		GouCode:     code,
		MoneyCount:  1,
		Status:      "已付款,未发货,未完成",
		PayType:     "账户",
		Company:     " ",
		ComCode:     " ",
		Address:     " ",
		Phone:       " ",
		ConfirmAddr: 0,
		Time:        fmt.Sprintf("%v.%v", time.Now().Unix(), int(GetRand(100, 999))),
	})
	if err != nil {
		fmt.Println("addMgoRecord error:", err.Error())
		return
	}

	if err := models.UpdateMember(&models.Member{
		UID:   man.UID,
		Money: man.Money - 1,
	}, "money"); err != nil {
		fmt.Println("updateMember error:", err.Error())
		return
	}

	if err := s.updateShop(shopList.ID, shopList.CanY+1, shopList.SanY-1); err != nil {
		fmt.Println("updateShop error:", err.Error())
		return
	}

}

func (s *AutoMan) OnClose() {
	select {
	case <-s.chanClose:
	default:
		close(s.chanClose)
	}
}

func (s *AutoMan) updateShop(shopId, canyu, shenyu int64) error {

	return models.UpdateShopList(&models.ShopList{ID: shopId, CanY: canyu, SanY: shenyu}, "canyurenshu", "shenyurenshu")
}

func (s *AutoMan) selectMan() (models.Member, error) {

	query := make(map[string]string, 0)

	query["level"] = "3"

	mls, err := models.GetMembers(query, []string{}, []string{"money"}, []string{"desc"}, 0, 1)
	if err != nil {
		return models.Member{}, err
	}

	if len(mls) <= 0 {
		return models.Member{}, fmt.Errorf("no member")
	}

	return mls[0].(models.Member), nil
}

func (s *AutoMan) selectShop() (models.ShopList, error) {

	var index int32

	mapShopList := make(map[int32]models.ShopList, 0)

	query := make(map[string]string, 0)

	mls, err := models.GetShopLists(query, []string{}, []string{"shenyurenshu"}, []string{"desc"}, 0, 100)
	if err != nil {
		return models.ShopList{}, err
	}

	for _, ml := range mls {
		if ml.(models.ShopList).SanY > 5 {
			mapShopList[index] = ml.(models.ShopList)
			index++
		}
	}

	if len(mapShopList) <= 0 {
		return models.ShopList{}, fmt.Errorf("no data")
	}

	return mapShopList[int32(GetRand(0, float64(len(mapShopList))))], nil
}

func (s *AutoMan) getShopCode(shopId int64) (string, error) {

	query := make(map[string]string, 0)

	query["s_id"] = fmt.Sprintf("%v", shopId)

	mls, err := models.GetShopCodes(query, []string{}, []string{"s_len"}, []string{"asc"}, 0, 10)
	if err != nil {
		return "", err
	}

	for _, ml := range mls {
		if ml.(models.ShopCode).SLen > 0 {
			code, index, newCodes := parseStr(ml.(models.ShopCode).SCodes)
			if code == "" {
				continue
			}
			models.UpdateShopCode(&models.ShopCode{
				ID:     ml.(models.ShopCode).ID,
				SLen:   index,
				SCodes: newCodes,
			}, "s_len", "s_codes")
			return code, nil
		}
	}

	return "", fmt.Errorf("getshopcode error")
}

func parseData(str string) (array []string) {

	elems := strings.Split(str, ";")

	for _, elem := range elems {

		if len(elem) <= 0 {
			continue
		}

		es := strings.Split(elem, ":")

		array = append(array, es[1])
	}

	return
}

func parseIndex(str string) int64 {

	sIndex := strings.Split(str, ":")

	if len(sIndex) != 2 {
		return 0
	}

	index, _ := strconv.ParseInt(sIndex[1], 10, 64)

	return index
}

func parseStr(str string) (code string, index int64, scode string) {

	var strArray []string

	strs := strings.Split(str, "{")

	for _, st := range strs {

		st = st[:len(st)-1]

		strArray = append(strArray, st)
	}

	if len(strArray) < 2 {
		return "", -1, ""
	}

	index = parseIndex(strArray[0]) - 1

	datas := parseData(strArray[1])

	code = datas[len(datas)-1]

	scode = "a:" + fmt.Sprintf("%v", index) + ":{"

	for i := 0; i < len(datas)-2; i++ {
		scode = scode + "i:" + datas[i] + ";"
	}

	scode = scode + "}"

	return
}

func GetRand(min float64, max float64) (result float64) {
	source := rand.NewSource(time.Now().UnixNano())
	nRand := rand.New(source)
	return nRand.Float64()*(max-min) + min
}

//func initMapCity() map[string][]string {

//	mapCity := make(map[string]([]string), 0)

//	mapCity = {"北京":{"北京"},"天津":{"天津"},"河北":{"石家庄","唐山","秦皇岛","邯郸","邢台","保定","张家口","承德"},"山东":{"滨州","德州","枣庄","青岛","济南","烟台","威海","淄博","潍坊","东营"},"山西":{"长治","晋城","朔州","忻州","吕梁","晋中","临汾"}}
//}
