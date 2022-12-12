package core

// #cgo CFLAGS: -I../../../include
// #cgo LDFLAGS: -L../../../build -L../../../build/local/lib -L../../../build/trezor-crypto -lTrustWalletCore -lwallet_core_rs -lprotobuf -lTrezorCrypto -lstdc++ -lm
// #include <TrustWalletCore/TWCoinType.h>
// #include <TrustWalletCore/TWAnySigner.h>
import "C"

import (
	"fmt"
	"math"
	"tw/protos/bitcoin"
	"tw/protos/common"
	"tw/protos/ethereum"
	"math/big"
	"tw/types"
	"encoding/hex"
	"google.golang.org/protobuf/proto"
)

func CreateSignedTx(inputData proto.Message, ct CoinType, outputData proto.Message) error {
	ibytes, _ := proto.Marshal(inputData)
	idata := types.TWDataCreateWithGoBytes(ibytes)
	defer C.TWDataDelete(idata)

	odata := C.TWAnySignerSign(idata, C.enum_TWCoinType(ct))
	defer C.TWDataDelete(odata)

	err := proto.Unmarshal(types.TWDataGoBytes(odata), outputData)
	if err != nil {
		return err
	}
	return nil
}

func CreateEthTransaction(receiver string, amount float64, ew *Wallet) string {
	priKeyByte, _ := hex.DecodeString(ew.PriKey)

	input := ethereum.SigningInput{
		ChainId:    big.NewInt(4).Bytes(), // mainnet: 1, rinkeby: 4 https://chainlist.org/
		Nonce:      big.NewInt(0).Bytes(), // get nonce from network
		TxMode:     ethereum.TransactionMode_Legacy,
		GasPrice:   big.NewInt(100000000000).Bytes(), // 100 gwei
		GasLimit:   big.NewInt(21000).Bytes(),
		ToAddress:  receiver,
		PrivateKey: priKeyByte,
		Transaction: &ethereum.Transaction{
			TransactionOneof: &ethereum.Transaction_Transfer_{
				Transfer: &ethereum.Transaction_Transfer{
					// amount should be in wei unit, eth * (10^decimals) = wei
					Amount: big.NewInt(int64(
						amount * math.Pow10(ew.CoinType.Decimals()),
					)).Bytes(),
					Data: []byte{},
				},
			},
		},
	}

	var output ethereum.SigningOutput
	err := CreateSignedTx(&input, ew.CoinType, &output)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(output.GetEncoded())
}

func CreateBtcTransaction(bw *Wallet) string {
	lockScript := BitcoinScriptLockScriptForAddress(bw.Address, bw.CoinType)
	fmt.Println("\nBitcoin address lock script:")
	fmt.Println("\t", hex.EncodeToString(lockScript))

	utxoHash, _ := hex.DecodeString("fff7f7881a8099afa6940d42d1e7f6362bec38171ea3edf433541db4e4ad969f")

	utxo := bitcoin.UnspentTransaction{
		OutPoint: &bitcoin.OutPoint{
			Hash:     utxoHash,
			Index:    0,
			Sequence: 4294967295,
		},
		Amount: 625000000,
		Script: lockScript,
	}

	priKeyByte, _ := hex.DecodeString(bw.PriKey)

	input := bitcoin.SigningInput{
		HashType:      uint32(BitcoinSigHashTypeAll),
		Amount:        1000000,
		ByteFee:       1,
		ToAddress:     "1Bp9U1ogV3A14FMvKbRJms7ctyso4Z4Tcx",
		ChangeAddress: "1FQc5LdgGHMHEN9nwkjmz6tWkxhPpxBvBU",
		PrivateKey:    [][]byte{priKeyByte},
		Utxo:          []*bitcoin.UnspentTransaction{&utxo},
		CoinType:      uint32(CoinTypeBitcoin),
	}

	var output bitcoin.SigningOutput
	err := CreateSignedTx(&input, bw.CoinType, &output)
	if err != nil {
		panic(err)
	}
	if output.GetError() != common.SigningError_OK {
		panic(output.GetError().String())
	}
	return hex.EncodeToString(output.GetEncoded())
}

func PrintWallet(w *Wallet) {
	fmt.Printf("%s wallet: \n", w.CoinType.GetName())
	fmt.Printf("\t address: %s \n", w.Address)
	fmt.Printf("\t pri key: %s \n", w.PriKey)
	fmt.Printf("\t pub key: %s \n", w.PubKey)
	fmt.Println("")
}
