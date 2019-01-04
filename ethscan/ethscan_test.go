package ethscan

import (
	"testing"
)

func TestGetBlockNumber(t *testing.T) {
	t.Log(GetLastBlock())
}

func TestGetTxByHash(t *testing.T) {
	t.Log(GetTxByHash("0x40eb908387324f2b575b4879cd9d7188f69c8fc9d87c901b9e2daaea4b442170"))
}

func TestGetTx(t *testing.T) {
	t.Log(GetTx(54092, 54192, "0x5abfec25f74cd88437631a7731906932776356f9"))
}
