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
	t.Log(GetTx(7080164, 7080174, "0x7c3a1c1a59fd3e62c99f6578c5d17bdfa5e9618b"))
}
