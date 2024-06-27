package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
	"github.com/xander42280/mpc/mpc_trade_test/core/sdk"
	"github.com/xander42280/mpc/mpc_trade_test/log"

	"github.com/xander42280/mpc/mpc_trade_test/common"
)

var (
	cmd      = flag.String("cmd", "check_node_status", "list_vaults|check_node_status|list_address|sign!list_wallet|create_wallet")
	private  = flag.String("private", common.FakePrivateKey, "Your private key.")
	url      = flag.String("url", common.BaseUrl, "Base url to Sinohope WaaS service.")
	valultId = flag.String("valult_id", "560095170453829", "valult_id")
	walletId = flag.String("wallet_id", "561456649624389", "wallet_id")
	from     = flag.String("from", "", "from")
	to       = flag.String("to", "", "to")
	address  = flag.String("address", "muKRKTpmp7tts9Ztfw8huUrzgWMr7BMDzP", "address")
)

func main() {
	flag.Parse()

	if err := check(); err != nil {
		panic(err)
	}

	a := log.App{
		Name: "sinohope-golang-sdk",
	}
	l := log.Log{
		Stdout: struct {
			Enable bool `toml:"enable"`
			Level  int  `toml:"level"`
		}{
			Enable: true,
			Level:  5,
		},
		File: struct {
			Enable bool   `toml:"enable"`
			Level  int    `toml:"level"`
			Path   string `toml:"path"`
			MaxAge int    `toml:"max-age"`
		}{
			Enable: true,
			Level:  5,
			Path:   "./logs/sinohope-golang-sdk.log",
			MaxAge: 7,
		},
	}
	log.SetLogDetailsByConfig(a, l)

	if *cmd == "check_node_status" {
		checkMPCNodeStatus()
	} else if *cmd == "list_address" {
		listAddress()
	} else if *cmd == "sign" {
		signRawData()
	} else if *cmd == "list_wallet" {
		listWallet()
	} else if *cmd == "create_wallet" {
		createWallet()
	} else if *cmd == "list_vaults" {
		getVaults()
	} else if *cmd == "list_supported_chains" {
		listSupportedChains()
	} else if *cmd == "get_supported_coins" {
		getSupportedCoins()
	} else if *cmd == "create_address" {
		createAddress()
	} else if *cmd == "gen_advance_address" {
		genAdvanceAddress()
	} else if *cmd == "transfer" {
		transfer()
	} else if *cmd == "list_transfer" {
		listTransfer(*address)
	}
}

func genRequestId() string {
	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return strconv.FormatInt(time.Now().UnixMicro(), 10)
	}
	return requestId
}

func check() error {
	if *private == "" {
		return fmt.Errorf("private key can not be empty")
	}
	if *url == "" {
		return fmt.Errorf("base url can not be empty")
	}
	return nil
}

func checkMPCNodeStatus() {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("create mpc node sdk failed, %v", err)
		return
	}
	request := &common.WaasMpcNodeExecRecordParam{
		BusinessExecType:   1,
		BusinessExecStatus: 10,
		SinoId:             "fake-sino-id",
		PageIndex:          0,
		PageSize:           40,
	}
	var result *common.WaaSMPCNodeRequestRes
	if result, err = client.MPCNode.ListMPCRequests(request); err != nil {
		logrus.Errorf("list mpc requests failed, %v", err)
	} else {
		// TODO: do something with result
		fmt.Printf("-----------> [%v]", result.List)
	}
	var status *common.WaaSMpcNodeStatusDTOData
	if status, err = client.MPCNode.Status(); err != nil {
		logrus.Errorf("get mpc node status failed, %v", err)
	} else {
		logrus.Infof("get mpc node status success, %v", status)
	}
}

func listSupportedChains() {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("list wallet failed, %v", err)
		return
	}
	var result []*common.WaasChainData
	if result, err = client.Common.GetSupportedChains(); err != nil {
		logrus.Errorf("list vault requests failed, %v", err)
	} else {
		fmt.Printf("list chain-----------> [%v]\n", result)
		for _, item := range result {
			fmt.Printf("chain: [%+v]\n", *item)
		}
	}
}

func getSupportedCoins() {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("list wallet failed, %v", err)
		return
	}
	request := &common.WaasChainParam{
		ChainSymbol: "BTC_TEST",
	}
	var result []*common.WaaSCoinDTOData
	if result, err = client.Common.GetSupportedCoins(request); err != nil {
		logrus.Errorf("get supported coins requests failed, %v", err)
	} else {
		fmt.Printf("get supported coins-----------> [%v]\n", result)
		for _, item := range result {
			fmt.Printf("coin: [%+v]\n", *item)
		}
	}
}

func getVaults() {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("list wallet failed, %v", err)
		return
	}
	var result []*common.WaaSVaultInfoData
	if result, err = client.Common.GetVaults(); err != nil {
		logrus.Errorf("list vault requests failed, %v", err)
	} else {
		fmt.Printf("list valult-----------> [%v]\n", result)
		for _, item := range result {
			fmt.Printf("valult: [%+v]\n", *item)
		}
	}
}

func listWallet() {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("list wallet failed, %v", err)
		return
	}
	request := &common.WaaSListWalletsParam{
		VaultId:   *valultId,
		PageIndex: 0,
		PageSize:  40,
	}
	var result *common.WaaSListWalletsResult
	if result, err = client.AccountAndAddress.ListWallets(request); err != nil {
		logrus.Errorf("list wallet requests failed, %v", err)
	} else {
		fmt.Printf("list wallet-----------> [%v]\n", *result)
		for _, item := range result.List {
			fmt.Printf("list wallet-----------> [%+v]\n", *item)
		}
	}
}

func createWallet() {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("create mpc node sdk failed, %v", err)
		return
	}
	walletInfo := common.WaaSCreateWalletInfo{
		WalletName:      "wall_1",
		AdvancedEnabled: 1,
	}
	request := &common.WaaSCreateBatchWalletParam{
		VaultId:     *valultId,
		RequestId:   genRequestId(),
		WalletInfos: []common.WaaSCreateWalletInfo{walletInfo},
	}
	var result []*common.WaaSWalletInfoData
	if result, err = client.AccountAndAddress.CreateWallets(request); err != nil {
		logrus.Errorf("create wallet requests failed, %v", err)
	} else {
		// TODO: do something with result
		fmt.Printf("create wallet-----------> [%v]\n", result)
		for i := 0; i < len(result); i++ {
			fmt.Printf("create wallet-----------> [%v]\n", result)
		}
	}
}

func createAddress() {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("create mpc node sdk failed, %v", err)
		return
	}
	request := &common.WaaSGenerateChainAddressParam{
		VaultId:     *valultId,
		RequestId:   genRequestId(),
		WalletId:    *walletId,
		Count:       1,
		ChainSymbol: "BTC_TEST",
	}
	var result []*common.WaaSAddressInfoData
	if result, err = client.AccountAndAddress.GenerateChainAddresses(request); err != nil {
		logrus.Errorf("create address requests failed, %v", err)
	} else {
		// TODO: do something with result
		fmt.Printf("create address-----------> [%+v]\n", result)
		for i := 0; i < len(result); i++ {
			fmt.Printf("create address-----------> [%+v]\n", *result[i])
		}
	}
}

func listAddress() {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("create mpc node sdk failed, %v", err)
		return
	}
	request := &common.WaaSListAddressesParam{
		VaultId:     *valultId,
		WalletId:    *walletId,
		ChainSymbol: "BTC_TEST",
		PageIndex:   0,
		PageSize:    40,
	}
	var result *common.WaaSListAddressesResult
	if result, err = client.AccountAndAddress.ListAddress(request); err != nil {
		logrus.Errorf("list address requests failed, %v", err)
	} else {
		fmt.Printf("-----------> [%+v]\n", *result)
		for _, item := range result.List {
			fmt.Printf("-----------> [%+v]\n", *item)
			listBalance(item.Address)
		}
	}
}

func listBalance(address string) {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("create mpc node sdk failed, %v", err)
		return
	}
	request := &common.WaaSGetAddressBalanceParam{
		AssetId: "BTC_BTC_TEST",
		Address: address,
	}
	var result *common.WaaSGetWalletBalanceDTOData
	if result, err = client.AccountAndAddress.GetAddressBalance(request); err != nil {
		logrus.Errorf("list balance requests failed, %v", err)
	} else {
		fmt.Printf("-----------> [%+v]\n", *result)
	}
}

func signRawData() {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("sign failed, %v", err)
		return
	}
	request := &common.WaaSSignRawDataParam{
		VaultId:   *valultId,
		RequestId: genRequestId(),
		WalletId:  *walletId,
		HdPath:    "m/1/1/1/1",
		RawData:   "0x4dac0911bbb5f363e04c425d84a84a98355285fb359ca212701528bf9f4164d4",
	}
	var result *common.WaaSSignRawDataRes
	if result, err = client.Advance.SignRawData(request); err != nil {
		logrus.Errorf("sign failed, %v", err)
	} else {
		fmt.Printf("-----------> [%+v]", *result)
	}
}

func genAdvanceAddress() {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("sign failed, %v", err)
		return
	}
	request := &common.WaaSAddressPathParam{
		VaultId:       *valultId,
		WalletId:      *walletId,
		Index:         1,
		AlgorithmType: 0,
		CoinType:      1, // 1 testnet
	}
	var result *common.WaaSAddressInfoData
	if result, err = client.Advance.GenAddressByPath(request); err != nil {
		logrus.Errorf("gen advance address failed, %v", err)
	} else {
		fmt.Printf("gen advance address-----------> [%+v]", *result)
	}
}

func transfer() {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("sign failed, %v", err)
		return
	}

	requestGas := &common.WalletTransactionFeeWAASParam{
		OperationType: "TRANSFER",
		ChainSymbol:   "BTC_TEST",
		AssetId:       "BTC_BTC_TEST",
		From:          *from,
		To:            *to,
		Amount:        "0.00001",
	}
	var resultGas *common.WalletTransactionFeeWAASResponse
	if resultGas, err = client.Transaction.Fee(requestGas); err != nil {
		logrus.Errorf("transfer failed, %v", err)
	} else {
		fmt.Printf("transfer fee-----------> [%+v]", *resultGas)
	}

	request := &common.WalletTransactionSendWAASParam{
		VaultId:     *valultId,
		WalletId:    *walletId,
		RequestId:   genRequestId(),
		ChainSymbol: "BTC_TEST",
		AssetId:     "BTC_BTC_TEST",
		From:        *from,
		To:          *to,
		ToTag:       "test",
		Amount:      "0.00001",
	}
	var result *common.CreateSettlementTxResData
	if result, err = client.Transaction.CreateTransfer(request); err != nil {
		logrus.Errorf("transfer failed, %v", err)
	} else {
		fmt.Printf("transfer-----------> [%+v]", *result)
	}
}

func listTransfer(address string) {
	logrus.
		WithField("private", private).
		Infof("after generate ECDSA keypair")
	client, err := sdk.NewApiClient(*url, *private)
	if err != nil {
		logrus.Errorf("create mpc node sdk failed, %v", err)
		return
	}
	request := &common.WalletTransactionQueryWAASParam{
		AssetId:     "BTC_BTC_TEST",
		ChainSymbol: "BTC_TEST",
		Address:     address,
	}
	var result *common.TransferHistoryWAASDTO
	if result, err = client.Transaction.ListTransactions(request); err != nil {
		logrus.Errorf("list balance requests failed, %v", err)
	} else {
		fmt.Printf("-----------> [%+v]\n", *result)
		for _, item := range result.List {
			fmt.Printf("-----------> [%+v]\n", item)
		}
	}
}
