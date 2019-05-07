go实现的以太坊代币转账程序
## 说明
- 将要转账的地址放到token.txt中
```
0x301edd7f9565e56b7db0702b8fd39f288f5cabaa
0x881914d8758616999f8a3f5cb9ec72404da21f0b
0x47a58b32ab181eebc74e48672fff568412f91b71
```
- 才main函数中写入私钥，代币地址等
```
err := TokenTransfer("filepath/token.txt", "prviKey", "tokenAddress")
```
- 在Tokentransfer()中修改转账的数量gasPrice

```200```为转账数量，```0.05```为gasprice，字符串类型，其中转账数量的单位为
``` 10^18Wei```
```
gas := "0.05"
tx, err := SignTokenTx(privKey, str, "200", tokenAddress, ethNet, gas)
```
