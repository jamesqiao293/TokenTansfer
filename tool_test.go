package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestSignTokenTx(t *testing.T) {
	//ethNet := "https://rinkeby.infura.io"
	ethNet := "https://mainnet.infura.io"
	//privateKey, toAddress, value, tokenAddress, ethNet, (gas, gaslimit, nonce)
	tx, err := SignTokenTx("a8c0a89236eb41cf3d419677bb7f9b3f9cd8ca93d908cf9dbb077aee13a10eb1", "0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d", "10", "0xfdac63e4aed8ab64ee6a58b4985363f703e1bdc8", ethNet,  "0.01"/*, "21000", "0"*/)
	if err != nil {
		t.Error("SignTx failed!", err)
	}
	fmt.Println(tx)
}

func TestTokenTransfer(t *testing.T) {
	filename := "/Users/apple/go/src/TokenTransfer/token.txt"
	TokenTransfer(filename, "a8c0a89236eb41cf3d419677bb7f9b3f9cd8ca93d908cf9dbb077aee13a10eb1", "0x28b149020d2152179873ec60bed6bf7cd705775d")
}

func TestToLower(t *testing.T)  {
	addr := strings.ToLower("0x39465d31336Ed20ACa2952B2B99C2C6Ab8579973")
	fmt.Println(addr)
}

func TestCompare(t *testing.T) {
	Compare()
}
