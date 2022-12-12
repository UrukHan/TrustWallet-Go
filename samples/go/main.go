package main

import (
	"fmt"
	"tw/core"
	"github.com/gin-gonic/gin"
	"net/http"
	"tw/src"
)


type rec struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

func transact(context *gin.Context) {
	var newReq rec
	if err := context.BindJSON(&newReq); err != nil {
		return
	}
	configuration := src.GetConfig()

	ew, err := core.CreateWalletWithMnemonic(configuration.SEED, core.CoinTypeEthereum)
        if err != nil {
                panic(err)
        }
	ethTxn := core.CreateEthTransaction(newReq.Address, newReq.Amount, ew)
	fmt.Println("Ethereum signed tx:")
	fmt.Println("\t", ethTxn)
	fmt.Printf("%v\n", newReq.Address)
	fmt.Printf("%v\n", newReq.Amount)
        core.PrintWallet(ew)
}

func helthCheck(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, map[string]interface{}{"code": 200, "status": "OK"})
}

func main() {
	configuration := src.GetConfig()

	fmt.Println("==> calling wallet core from go")
	fmt.Println("==> mnemonic is valid: ", core.IsMnemonicValid(configuration.SEED))

	// bitcoin wallet
	bw, err := core.CreateWalletWithMnemonic(configuration.SEED, core.CoinTypeBitcoin)
	if err != nil {
		panic(err)
	}
	core.PrintWallet(bw)

	// ethereum wallet
	ew, err := core.CreateWalletWithMnemonic(configuration.SEED, core.CoinTypeEthereum)
	if err != nil {
		panic(err)
	}
	core.PrintWallet(ew)

	// tron wallet
	tw, err := core.CreateWalletWithMnemonic(configuration.SEED, core.CoinTypeTron)
	if err != nil {
		panic(err)
	}
	core.PrintWallet(tw)

	// Bitcion transaction
	btcTxn := core.CreateBtcTransaction(bw)
	fmt.Println("\nBitcoin signed tx:")
	fmt.Println("\t", btcTxn)

	fmt.Printf("Listening on PORT %s\n", configuration.PORT)
	router := gin.Default()
	router.GET("/health", helthCheck)
	router.POST("/transact", transact)
	router.Run(configuration.HOST + ":" + configuration.PORT)
	fmt.Println("(+_+)")
}
