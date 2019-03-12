package ethscan

import (
	"testing"
)

func TestGetBlockNumber(t *testing.T) {
	t.Log(GetLastBlock())
}

func TestGetTxByHash(t *testing.T) {
	t.Log(GetTxByHash("0x2663ba0db6909921c469c6bb49ccc9bd3313dce7db1326d40507df78f3bbdd67"))
}

func TestGetTx(t *testing.T) {
	t.Log(GetTx(7080311, 7080336, "0xFe8E9198CEb395Bd748Aaff3b6f8d8015E34dC01"))
}

func TestGetBlockByBlockNumber(t *testing.T) {

	number, _ := GetLastBlock()

	t.Log(GetBlockByBlockNumber(number))
}
