package backserver

import (
	"time"

	"github.com/hongyuefan/tmpserver/ethscan"
	"github.com/hongyuefan/tmpserver/models"
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
}

func NewServer(founder string, interval, judge int64) *Server {
	return &Server{
		limit:        200,
		waitingDatas: make(map[int64]*models.AddMoneyRecord, 0),
		founder:      founder,
		chanClose:    make(chan bool, 0),
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
				if _, ok := s.waitingDatas[ml.(models.AddMoneyRecord).ID]; !ok {
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
		ID:      record.ID,
		UID:     record.UID,
		Address: record.Address,
		Hash:    record.Hash,
		Money:   record.Money,
		Type:    record.Type,
		Status:  record.Status,
		Time:    record.Time,
	}
}
func (s *Server) changeRecord(id int64, status int, money, hash string) error {
	return models.UpdateRecord(&models.AddMoneyRecord{ID: id, Status: status}, "status", "money", "hash")
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
					if status, err := s.checkTxByHash(data.Hash); err != nil {
						switch status {
						case 1:
							models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_SUCCESS}, "status")
							break
						case 2:
							if data.CheckedBlock+s.judge < s.curBlockNumber {
								models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_FAILED}, "status")
								delete(s.waitingDatas, data.ID)
							} else {
								data.CheckedBlock = s.curBlockNumber
							}
							break
						case -1:
							models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_FAILED}, "status")
							delete(s.waitingDatas, data.ID)

						}
					}
				}
				if data.Type == types.PAY_ERCODE {
					if status, hash, money, err := s.checkTx(s.preBlockNumber, s.curBlockNumber, data.Address); err != nil {
						switch status {
						case 1:
							models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Hash: hash, Money: money, Status: types.STATUS_SUCCESS}, "hash", "money", "status")
							break
						case 2:
							if data.CheckedBlock+s.judge < s.curBlockNumber {
								models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_FAILED}, "status")
								delete(s.waitingDatas, data.ID)
							} else {
								data.CheckedBlock = s.curBlockNumber
							}
							break
						case -1:
							models.UpdateRecord(&models.AddMoneyRecord{ID: data.ID, Status: types.STATUS_FAILED}, "status")
							delete(s.waitingDatas, data.ID)
						}
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
