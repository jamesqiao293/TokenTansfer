package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//var w sync.WaitGroup

func main() {
	err := TokenTransfer("/Users/apple/go/src/TokenTransfer/token.txt", "a8c0a89236eb41cf3d419677bb7f9b3f9cd8ca93d908cf9dbb077aee13a10eb1", "0xfdac63e4aed8ab64ee6a58b4985363f703e1bdc8")
	if err != nil {
		log.Println(err)
	}
}

func TokenTransfer(filename, privKey, tokenAddress string) error {
	ethNet := "https://mainnet.infura.io"
	tokenInfo, err := ReadFile(filename)
	if err != nil {
		return errors.New("read from the file failed")
	}
	slice := *tokenInfo
	for i := 0; i < len(slice) - 1; i++ {
		str := strings.TrimRight(slice[i], " ")
		str = strings.ToLower(str)
		isok := CheckAddress(str)
		if str != "" && isok{
			time.Sleep(5 * time.Second)
			//privateKey, toAddress, value, tokenAddress, ethNet
			//w.Add(1)
			//go func() {
				gas := "0.05"
				tx, err := SignTokenTx(privKey, str, "200", tokenAddress, ethNet, gas)
				if err != nil {
					fmt.Println("********************************************************************")
					fmt.Println("["+str+"]","[failed]")
					fmt.Println("["+tx+"]")
				}else{
					fmt.Println("********************************************************************")
					fmt.Println("["+str+"]","[success]")
					fmt.Println("["+tx+"]")
				}
				//w.Done()
			//}()
		}else{
			return errors.New("address invalid")
		}
	}
	//w.Wait()
	return nil
}

func ReadFile(filename string) (*[]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("open the file failed!")
		return nil, err
	}
	defer file.Close()
	buf := make([]byte, 4096*1000)
	_, err = file.Read(buf)
	if err != nil {
		log.Println("read the file failed!")
		return nil, err
	}
	slice := strings.Split(string(buf), "\n")

	return &slice, nil
}

func Compare() error{
	s1, err := ReadFile("/Users/apple/go/src/TokenTransfer/address.txt")
	if err != nil {
		log.Println("read the file address.txt failed!")
		return err
	}
	s2, err := ReadFile("/Users/apple/go/src/TokenTransfer/tokenaddress2.txt")
	if err != nil {
		log.Println("read the file tokenaddress.txt failed!")
		return err
	}
	slice1 := *s1
	slice2 := *s2
	for i := 0; i < len(slice1) - 1; i++ {
		slice1[i] = strings.ToLower(slice1[i])
		slice1[i] = strings.TrimLeft(slice1[i], " ")
		slice1[i] = strings.TrimRight(slice1[i], " ")
	}
	for i := 0; i < len(slice2) - 1; i++ {
		slice2[i] = strings.ToLower(slice2[i])
		slice2[i] = strings.TrimLeft(slice2[i], " ")
		slice2[i] = strings.TrimRight(slice2[i], " ")
	}

	addressMap1 := make(map[string]string, 2000)
	for i := 0; i< len(slice2) - 1; i++ {
		addressMap1[slice2[i]] = slice2[i]
	}
	addressMap2 := make(map[string]string, 2000)
	for i := 0; i< len(slice1) - 1; i++ {
		addressMap2[slice1[i]] = slice1[i]
	}
	for i := 0; i<len(slice1) - 1; i++ {
		if _, ok := addressMap1[slice1[i]]; !ok {
			fmt.Println(slice1[i])
		}
	}
	fmt.Println("s1:", len(slice1), "s2:", len(slice2), len(addressMap1), len(addressMap2))
	return nil
}

func SignTokenTx(privateKey, toAddress, value, tokenAddress, ethNet , gas/*, gaslimit, nonce*/ string) (string, error) {
	client, err := ethclient.Dial(ethNet)
	if err != nil {
		log.Println("set ethclient failed!", err)
		return "", err
	}
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Println("hextoECDSA privateKey failed!", err)
		return "", err
	}

	pubKey := privKey.Public()
	publicKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey", err)
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nn, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
		return "", err
	}

	//单位转换ether=>Wei
	valuef, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Println("is not a number")
	}
	valueWei, isOk := new(big.Int).SetString(fmt.Sprintf("%.0f", valuef*1000000000000000000), 10)
	if !isOk {
		log.Println("float to bigInt failed!")
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("set the gasPrice failed!", err)
		return "", err
	}

	toAddr := common.HexToAddress(toAddress)
	tokenAddr := common.HexToAddress(tokenAddress)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	paddedToAddress := common.LeftPadBytes(toAddr.Bytes(), 32)

	paddedValue := common.LeftPadBytes(valueWei.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedToAddress...)
	data = append(data, paddedValue...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddr,
		Data: data,
	})
	if err != nil {
		log.Println("set the gasLimit failed!", err)
		return "", err
	}

	gasF, err := strconv.ParseFloat(gas, 64)
	if err != nil {
		return "", err
	}
	gasPrice, isSuccess := new(big.Int).SetString(fmt.Sprintf("%.0f", gasF*1000000000000000), 10)
	if !isSuccess {
		return "", err
	}
	gasLimitInt := new(big.Int).SetInt64(int64(gasLimit))
	gasPrice = gasPrice.Div(gasPrice, gasLimitInt)

	ethValue := big.NewInt(0)

	//需要覆盖前面的交易时
	//nn = uint64(44)

	gasLimit = uint64(90000)
	tx := types.NewTransaction(nn, tokenAddr, ethValue, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("sign the tx failed!", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)

	if err != nil {
		log.Println("sign the tx failed!", err)
	}

	/*var buff bytes.Buffer
	if err := signedTx.EncodeRLP(&buff); err != nil {
		return "", err
	}
	sTx := fmt.Sprintf("0x%x", buff.Bytes())*/

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("client send tx failed!", err)
		return "", err
	}
	//log.Printf("tx sent: %s\n", signedTx.Hash().Hex())

	return signedTx.Hash().Hex(), nil
	//return sTx, nil
}

func CheckAddress(address string) bool {
	b, err := regexp.MatchString("^(0x)[0-9a-f]{40}$", address)
	if err != nil {
		return false
	}
	if b {
		return true
	}
	return false
}