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
	waitingDatas   []interface{}
	founder        string
	chanClose      chan bool
}

func NewServer(founder string, interval int64) *Server {
	return &Server{
		limit:        200,
		waitingDatas: make([]interface{}, 0),
		founder:      founder,
		chanClose:    make(chan bool, 0),
		interval:     interval,
	}
}

func (s *Server) OnClose() {
	close(s.chanClose)
}

func (s *Server) getWaitingDatas() {

	query := make(map[string]string, 0)

	query["status"] = types.STATUS_WAITING

	offset := 0

	for {
		select {
		case <-s.chanClose:
			return
		default:
			mls, err := models.GetRecords(query, []string, []string{"id"}, []string{"asc"}, offset, s.limit)
			if err != nil {
				break
			}

			s.waitingDatas = append(s.waitingDatas, mls...)

			if len(mls) < 200 {
				break
			}

			offset += limit
		}
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

	ticker := time.NewTicker(time.Duration(s.interval))

	for {
		select {
		case <-s.chanClose:
			ticker.Stop()
			return
		case <-ticker.C:
			if s.curBlockNumber, err = ethscan.GetLastBlock(); err != nil {
				continue
			}

			s.getWaitingDatas()

			for _, data := range s.waitingDatas {

				if data.(models.AddMoneyRecord).Type == types.PAY_METAMASK {

					txs, err := ethscan.GetTxByHash(data.(models.AddMoneyRecord).Hash)
					if err != nil {

					}
				}
				if data.(models.AddMoneyRecord).Type == types.PAY_ERCODE {

				}
			}
		}
	}

}
