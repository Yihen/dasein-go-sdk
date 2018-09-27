package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	logging "gx/ipfs/QmRb5jh8z2E8hMGN2tkvs1yHynUanqnZ3UeKwgN1i9P1F8/go-log"

	sdk "github.com/daseinio/dasein-go-sdk"
	"github.com/daseinio/dasein-go-sdk/crypto"
)

var smallTxt = "QmZyvDNq1gHEkH5USKLUSunZBBya4qox7w4dWGxnt41Zox"
var bigTxt = "QmW5CME8vkw3ndeuDuf5a5oL9x55yPWfhF4fz4R6XTMTBk"

//var largeTxt = "QmU7QRQpSZhukKsraEaa23Re1AzLqpFvyHPwayseVKTbFp"
var deleteTxt = "QmevhnWdtmz89BMXuuX5pSY2uZtqKLz7frJsrCojT5kmb6"

// var node = "/ip4/127.0.0.1/tcp/4001/ipfs/QmR1AqNQBqAjPeLswq86dkJZ5Y7ACVGoXzz2K8tz6MHyUB"
var node = "/ip4/127.0.0.1/tcp/4001/ipfs/Qmdkh8dBb8p99KGDhazTnNZJpM4hDx95NJtnSLGSKp5tTy"
var log = logging.Logger("test")

var encrypt = false
var encryptPassword = "123456"
var walletPwd = "pwd"
var wallet = "./wallet.dat"
var rpc = "http://127.0.0.1:20336"

func testSendFile(fileName string, rate, times uint64, copyNum int32) {
	client, err := sdk.NewClient("", wallet, walletPwd, rpc)
	if err != nil {
		log.Error(err)
		return
	}
	err = client.SendFile(fileName, rate, times, copyNum, encrypt, encryptPassword)
	if err != nil {
		log.Error(err)
		return
	}
}

func testSendSmallFile() {
	client, err := sdk.NewClient(node, wallet, walletPwd, rpc)
	if err != nil {
		log.Error(err)
		return
	}
	var smallFile = "smallfile"
	smallF, err := os.OpenFile(smallFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Error(err)
		return
	}
	defer smallF.Close()
	smallF.WriteString("hello world2234556666789\n")
	err = client.SendFile(smallFile, 1, 1, 0, encrypt, encryptPassword)
	if err != nil {
		log.Error(err)
		return
	}
	err = os.Remove(smallFile)
	if err != nil {
		log.Error(err)
		return
	}

}

func testSendBigFile() {
	client, err := sdk.NewClient(node, wallet, walletPwd, rpc)
	if err != nil {
		log.Error(err)
		return
	}
	var bigFile = "bigfile"
	bigF, err := os.OpenFile(bigFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Error(err)
		return
	}
	defer bigF.Close()

	for i := 1; i < 80000; i++ {
		bigF.WriteString(fmt.Sprintf("%d\n", i))
	}
	err = client.SendFile(bigFile, 1, 1, 0, encrypt, encryptPassword)
	if err != nil {
		log.Errorf("send file err:%s", err)
		return
	}
	err = os.Remove(bigFile)
	if err != nil {
		log.Error(err)
		return
	}
}

func testGetData(fileHashStr string) {
	r := sdk.NewContractRequest(wallet, walletPwd, rpc)
	nodes, err := r.FindStoreFileNodes(fileHashStr)
	if err != nil {
		log.Error(err)
		return
	}
	randIndx := rand.Intn(len(nodes))
	chosenNode := nodes[randIndx]
	log.Infof("get data from :%s", chosenNode.Addr)
	client, err := sdk.NewClient(chosenNode.Addr, wallet, walletPwd, rpc)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("-----------------------")
	log.Info("Single Block Test")
	_, err = client.GetData(fileHashStr, chosenNode.WalletAddr)
	if err != nil {
		log.Error(err)
	}
	// file, err := os.Create("small")
	// if err != nil {
	// 	log.Error(err)
	// } else {
	// 	log.Infof("GetData %s success", smallTxt)
	// }
	// _, err = file.Write(data)
	// if err != nil {
	// 	log.Error(err)
	// }
	// file.Close()
	if encrypt {
		crypto.AESDecryptFile(fileHashStr, encryptPassword, fmt.Sprintf("%s-decrypted", fileHashStr))
	}
	log.Info("-----------------------")
	// log.Info("Multi Block Test")
	// data, err = client.GetData(bigTxt)
	// if err != nil {
	// 	log.Error(err)
	// }

	// file, err = os.Create("big")
	// if err != nil {
	// 	log.Error(err)
	// }
	// _, err = file.Write(data)
	// if err != nil {
	// 	log.Error(err)
	// }
	// file.Close()

	// log.Info("-----------------------")
	// log.Info("Delete Block Test")
	// err = client.DelData(deleteTxt)
	// if err != nil {
	// 	log.Error(err)
	// } else {
	// 	log.Infof("DelData %s success", deleteTxt)
	// }

	/*
		log.Info("Multi Block Test")
		data, err = client.GetData(largeTxt)
		if err != nil {
			log.Error(err)
		}

		file, err = os.Create("large")
		if err != nil {
			log.Error(err)
		}
		_, err = file.Write(data)
		if err != nil {
			log.Error(err)
		}
		file.Close()
	*/
}

func testDelData(fileHashStr string) {
	r := sdk.NewContractRequest(wallet, walletPwd, rpc)
	nodes, err := r.FindStoreFileNodes(smallTxt)
	if err != nil {
		log.Error(err)
		return
	}
	err = r.DeleteFile(fileHashStr)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("delete file sucess:%s", fileHashStr)
	retry := 0
	for {
		if retry > sdk.MAX_RETRY_REQUEST_TIMES {
			log.Error("delete file failed timeout")
			return
		}
		info, _ := r.GetFileInfo(fileHashStr)
		if info == nil {
			break
		}
		retry++
		time.Sleep(time.Duration(sdk.MAX_REQUEST_TIMEWAIT) * time.Second)
	}

	for _, node := range nodes {
		client, err := sdk.NewClient(node.Addr, wallet, walletPwd, rpc)
		if err != nil {
			log.Error(err)
			continue
		}
		err = client.DelData(fileHashStr)
		if err != nil {
			log.Errorf("delete file:%s failed in node:%s", fileHashStr, node.Addr)
		}
	}
}

func testByFlags() {
	var fileName *string = flag.String("file", "", "Use --file <filesource>")
	var operation *string = flag.String("op", "", "Use --op <send|get|del|testsendsmall|testsendbig>")
	var isEncrypt *bool = flag.Bool("encrypt", false, "Encrypt file")
	var ePwd *string = flag.String("encryptpwd", "", "Encrypt password")
	var wPwd *string = flag.String("walletpwd", "pwd", "Wallet password")
	var rates *int = flag.Int("challengerate", 100, "block count of challenge")
	var times *int = flag.Int("challengetimes", 1, "challenge time")
	var copynum *int = flag.Int("copynum", 0, "backup nodes number for copy")
	flag.Parse()
	if len(*operation) == 0 {
		return
	}
	if len(*ePwd) > 0 {
		encryptPassword = *ePwd
	}
	walletPwd = *wPwd
	encrypt = *isEncrypt
	switch *operation {
	case "send":
		testSendFile(*fileName, uint64(*rates), uint64(*times), int32(*copynum))
	case "get":
		if len(*fileName) <= 0 {
			log.Errorf("option: file is missing")
			return
		}
		testGetData(*fileName)
	case "del":
		if len(*fileName) <= 0 {
			log.Errorf("option: file is missing")
			return
		}
		testDelData(*fileName)
	case "testsendsmall":
		testSendSmallFile()
	case "testsendbig":
		testSendBigFile()
	}
}

func main() {
	logging.SetLogLevel("test", "INFO")
	logging.SetLogLevel("daseingosdk", "INFO")
	logging.SetLogLevel("bitswap", "INFO")
	testByFlags()
	// testDelFileAndGet()
	// testSendSmallFile()
	// testGetData("QmafaFyC4DkLcaPTbhBrBxfUWB3BVWNeAnekzzYhgtkztD")
	// testDelData()·
	// testSendBigFile()
	// testSendSmallFile()
	// testGetNodeList()
	// testGetFileProveDetails()
}
