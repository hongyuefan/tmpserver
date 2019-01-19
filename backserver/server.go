package backserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"time"

	"github.com/hongyuefan/tmpserver/ethscan"
	"github.com/hongyuefan/tmpserver/models"
	"github.com/hongyuefan/tmpserver/tools"
	"github.com/hongyuefan/tmpserver/types"
)

type Server struct {
	curBlockNumber int64
	preBlockNumber int64
	limit          int64
	interval       int64
	judge          int64
	waitingDatas   map[int64]*models.AddMoneyRecord
	founder        string
	chanClose      chan bool
	chanGetData    chan bool
}

func NewServer(founder string, interval, judge int64) *Server {
	return &Server{
		limit:        200,
		waitingDatas: make(map[int64]*models.AddMoneyRecord, 0),
		founder:      founder,
		chanClose:    make(chan bool, 0),
		chanGetData:  make(chan bool, 0),
		interval:     interval,
		judge:        judge,
	}
}

func (s *Server) OnClose() {
	close(s.chanClose)
}

func (s *Server) getWaitingDatas() {

	query := make(map[string]string, 0)

	query["status"] = "0"

	var offset int64 = 0

FOR:
	for {
		select {
		case <-s.chanClose:
			return
		default:
			mls, err := models.GetRecords(query, []string{}, []string{"id"}, []string{"asc"}, offset, s.limit)
			if err != nil {
				break FOR
			}
			for _, ml := range mls {

				_, ok := s.waitingDatas[ml.(models.AddMoneyRecord).ID]
				if !ok {
					s.waitingDatas[ml.(models.AddMoneyRecord).ID] = s.newRecord(ml.(models.AddMoneyRecord))
				}
			}
			if len(mls) < 200 {
				break FOR
			}
			offset += s.limit
		}
	}

}

func (s *Server) newRecord(record models.AddMoneyRecord) *models.AddMoneyRecord {
	return &models.AddMoneyRecord{
		ID:           record.ID,
		UID:          record.UID,
		Address:      record.Address,
		Hash:         record.Hash,
		Money:        record.Money,
		Type:         record.Type,
		Status:       record.Status,
		Time:         record.Time,
		CheckedBlock: record.CheckedBlock,
	}
}
func (s *Server) changeRecord(id int64, status int, money, hash string) error {
	return models.UpdateRecord(&models.AddMoneyRecord{ID: id, Status: status}, "status", "money", "hash")
}

type EthPrice struct {
	Price string `json:"price_usd"`
}

type Exchange struct {
	Reason    string   `json:"reason"`
	Result    ExResult `json:"result"`
	ErrorCode int64    `json:"error_code"`
}

type ExResult struct {
	Update string     `json:"update"`
	List   [][]string `json:"list"`
}

func (s *Server) getExchange() (float64, error) {

	rsp, err := tools.Get("http://op.juhe.cn/onebox/exchange/query?key=de0438b501aebc296fa9cdeb0632b1f7")

	if err != nil {
		return 0, err
	}

	byt, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return 0, err
	}
	var exchanges Exchange

	if err := json.Unmarshal(byt, &exchanges); err != nil {
		return 0, err
	}

	for _, result := range exchanges.Result.List {
		if result[0] == "美元" {

			iEx, err := strconv.ParseFloat(result[3], 64)
			if err != nil {
				return 0, err
			}

			return iEx / float64(100), nil
		}
	}

	return 0, fmt.Errorf("cannt get exchange")
}

func (s *Server) getPrice() (float64, error) {

	rsp, err := tools.Get("https://api.coinmarketcap.com/v1/ticker/ethereum/")

	if err != nil {
		return 0, err
	}

	byt, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return 0, err
	}

	var eths []EthPrice

	if err := json.Unmarshal(byt, &eths); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(eths[0].Price, 64)
}

func (s *Server) addMoney(uid int64, eth string) {

	var money float64

	member := &models.Member{UID: uid}

	if err := models.GetMember(member, "uid"); err != nil {
		fmt.Println("addMoney GetMember error uid: ", uid, "eth:", eth, "error:", err.Error())
		return
	}

	feth, err := strconv.ParseFloat(eth, 64)
	if err != nil {
		fmt.Println("addMoney ParseFloat error uid: ", uid, "eth:", eth, "error:", err.Error())
		return
	}

	price, err := s.getPrice()
	if err != nil {
		fmt.Println("addMoney getprice error uid: ", uid, "eth:", eth, "error:", err.Error())
		return
	}
	exch, err := s.getExchange()
	if err != nil {
		exch = 6.7
	}

	money = feth * price * exch

	member.Money += money

	if err := models.UpdateMember(member, "money"); err != nil {
		fmt.Println("addMoney UpdateMember error uid: ", uid, "eth:", eth, "error:", err.Error())

		return
	}

	return
}

func (s *Server) Handler() {

	var err error

	if s.preBlockNumber, err = ethscan.GetLastBlock(); err != nil {
		return
	}

	for {

		time.Sleep(time.Second * time.Duration(s.interval))

		select {
		case <-s.chanClose:
			return
		default:

			if s.curBlockNumber, err = ethscan.GetLastBlock(); err != nil {
				continue
			}

			s.getWaitingDatas()

			for _, data := range s.waitingDatas {

				time.Sleep(time.Second)

				if data.Type == types.PAY_METAMASK {
					if data.Hash == "" {
						models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_FAILED}, "status")
					}

					status, _, _, _ := s.checkTx(data.CheckedBlock, s.curBlockNumber, data.Address)

					switch status {
					case 1:
						models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_SUCCESS}, "status")
						s.addMoney(data.UID, data.Money)
						delete(s.waitingDatas, data.ID)
						break
					case 2:
						if data.CheckedBlock+s.judge < s.curBlockNumber {
							models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_FAILED}, "status")
							delete(s.waitingDatas, data.ID)
						}
						break
					case -1:
						models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_FAILED}, "status")
						delete(s.waitingDatas, data.ID)

					}

				}
				if data.Type == types.PAY_ERCODE {

					status, hash, money, _ := s.checkTx(data.CheckedBlock, s.curBlockNumber, data.Address)

					switch status {
					case 1:
						if err := models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Hash: hash, Money: money, Status: types.STATUS_SUCCESS}, "hash", "money", "status"); err != nil {
							models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_FAILED}, "status")
							break
						}
						s.addMoney(data.UID, data.Money)
						delete(s.waitingDatas, data.ID)
						break
					case 2:
						fmt.Println(data.CheckedBlock, s.judge, s.curBlockNumber)
						if data.CheckedBlock+s.judge < s.curBlockNumber {
							models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_FAILED}, "status")
							delete(s.waitingDatas, data.ID)
						}
						break
					case -1:
						fmt.Println("-1")
						models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_FAILED}, "status")
						delete(s.waitingDatas, data.ID)
					}
				}

			}

			s.preBlockNumber = s.curBlockNumber
		}
	}

}

//1:success
//2:not find
//3:net error
//-1:failed
func (s *Server) checkTxByHash(hash string) (status int, err error) {
	txs, err := ethscan.GetTxByHash(hash)
	if err != nil {
		return 3, err
	}
	if txs.Status == "1" {
		if txs.Result[0].IsError == "1" {
			return -1, nil
		}
		return 1, nil
	}
	if txs.Status == "0" {
		return 2, nil
	}
	return 3, nil
}

func (s *Server) checkTx(start, end int64, address string) (status int, txhash, money string, err error) {

	txs, err := ethscan.GetTx(start, end, address)
	if err != nil {
		return 3, "", "", err
	}

	if txs.Status == "0" {
		return 2, "", "", nil
	}

	for _, tx := range txs.Result {
		if tx.From == address && tx.To == s.founder {
			if tx.IsError == "1" {
				return -1, "", "", nil
			} else {
				return 1, tx.Hash, tx.Value, nil
			}
		}
	}
	return 2, "", "", nil
}
