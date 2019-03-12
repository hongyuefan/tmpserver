package ethscan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/hongyuefan/tmpserver/tools"
)

type RspHashByBlockNumber struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      int64       `json:"id"`
	Result  BlockDetail `json:"result"`
}

type BlockDetail struct {
	Hash      string `json:"hash"`
	Number    string `json:"number"`
	TimeStamp string `json:"timestamp"`
}

type RspBlockNumber struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int64  `json:"id"`
	Result  string `json:"result"`
}

type RspBlockTxs struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Transaction `json:"result"`
}

type Transaction struct {
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	IsError     string `json:"isError"`
	ErrCode     string `json:"errCode"`
	Hash        string `json:"hash"`
}

func GetBlockByBlockNumber(blockNumber int64) (blockHash, hexBlockNumber, timeStamp string, error error) {

	var rspHashByBlockNumber RspHashByBlockNumber

	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=%#x&boolean=true&apikey=FIE4TCCJBDQ48VMA99H8F5X7FVQJ98JISA", blockNumber)

	rsp, err := tools.Get(url)
	if err != nil {
		return "", "", "", err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", "", "", err
	}

	if err := json.Unmarshal(body, &rspHashByBlockNumber); err != nil {
		return "", "", "", err
	}

	return rspHashByBlockNumber.Result.Hash, rspHashByBlockNumber.Result.Number, rspHashByBlockNumber.Result.TimeStamp, nil

}

func GetLastBlock() (int64, error) {

	var rspBlockNumber RspBlockNumber

	rsp, err := tools.Get("https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey=FIE4TCCJBDQ48VMA99H8F5X7FVQJ98JISA")
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(body, &rspBlockNumber); err != nil {
		return 0, err
	}
	return strconv.ParseInt(rspBlockNumber.Result, 0, 64)
}

func GetTxByHash(hash string) (*RspBlockTxs, error) {

	var rspBlockTxs RspBlockTxs

	url := "https://api.etherscan.io/api?module=account&action=txlistinternal&txhash=" + hash + "&apikey=FIE4TCCJBDQ48VMA99H8F5X7FVQJ98JISA"

	rsp, err := tools.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &rspBlockTxs); err != nil {
		return nil, err
	}
	return &rspBlockTxs, nil
}

func GetTx(start, end int64, address string) (*RspBlockTxs, error) {

	var rspBlockTxs RspBlockTxs

	url := "http://api.etherscan.io/api?module=account&action=txlist&address=" + address + "&startblock=" + fmt.Sprintf("%v", start) + "&endblock=" + fmt.Sprintf("%v", end) + "&sort=desc" + "&apikey=FIE4TCCJBDQ48VMA99H8F5X7FVQJ98JISA"

	rsp, err := tools.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &rspBlockTxs); err != nil {
		return nil, err
	}
	return &rspBlockTxs, nil
}
