// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package mttest

import (
	"bytes"
	"context"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/mantlenetworkio/mantle/blsSignatures"
	"github.com/mantlenetworkio/mantle/cmd/genericconf"
	"github.com/mantlenetworkio/mantle/das"
	"github.com/mantlenetworkio/mantle/mtos"
	"github.com/mantlenetworkio/mantle/mtos/util"
	"github.com/mantlenetworkio/mantle/mtstate"
	"github.com/mantlenetworkio/mantle/mtutil"
	"github.com/mantlenetworkio/mantle/util/headerreader"
	"github.com/mantlenetworkio/mantle/util/mtmath"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"

	"github.com/mantlenetworkio/mantle/mtnode"
	_ "github.com/mantlenetworkio/mantle/nodeInterface"
	"github.com/mantlenetworkio/mantle/solgen/go/bridgegen"
	"github.com/mantlenetworkio/mantle/solgen/go/mocksgen"
	"github.com/mantlenetworkio/mantle/solgen/go/precompilesgen"
	"github.com/mantlenetworkio/mantle/statetransfer"
	"github.com/mantlenetworkio/mantle/util/testhelpers"
	"github.com/mantlenetworkio/mantle/validator"
)

type info = *BlockchainTestInfo
type client = mtutil.L1Interface

func SendWaitTestTransactions(t *testing.T, ctx context.Context, client client, txs []*types.Transaction) {
	t.Helper()
	for _, tx := range txs {
		Require(t, client.SendTransaction(ctx, tx))
	}
	for _, tx := range txs {
		_, err := EnsureTxSucceeded(ctx, client, tx)
		Require(t, err)
	}
}

func TransferBalance(
	t *testing.T, from, to string, amount *big.Int, l2info info, client client, ctx context.Context,
) (*types.Transaction, *types.Receipt) {
	return TransferBalanceTo(t, from, l2info.GetAddress(to), amount, l2info, client, ctx)
}

func TransferBalanceTo(
	t *testing.T, from string, to common.Address, amount *big.Int, l2info info, client client, ctx context.Context,
) (*types.Transaction, *types.Receipt) {
	tx := l2info.PrepareTxTo(from, &to, l2info.TransferGas, amount, nil)
	err := client.SendTransaction(ctx, tx)
	Require(t, err)
	res, err := EnsureTxSucceeded(ctx, client, tx)
	Require(t, err)
	return tx, res
}

func SendSignedTxViaL1(
	t *testing.T,
	ctx context.Context,
	l1info *BlockchainTestInfo,
	l1client mtutil.L1Interface,
	l2client mtutil.L1Interface,
	delayedTx *types.Transaction,
) *types.Receipt {
	delayedInboxContract, err := bridgegen.NewInbox(l1info.GetAddress("Inbox"), l1client)
	Require(t, err)
	usertxopts := l1info.GetDefaultTransactOpts("User", ctx)

	txbytes, err := delayedTx.MarshalBinary()
	Require(t, err)
	txwrapped := append([]byte{mtos.L2MessageKind_SignedTx}, txbytes...)
	l1tx, err := delayedInboxContract.SendL2Message(&usertxopts, txwrapped)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, l1client, l1tx)
	Require(t, err)

	// sending l1 messages creates l1 blocks.. make enough to get that delayed inbox message in
	for i := 0; i < 30; i++ {
		SendWaitTestTransactions(t, ctx, l1client, []*types.Transaction{
			l1info.PrepareTx("Faucet", "Faucet", 30000, big.NewInt(1e12), nil),
		})
	}
	receipt, err := EnsureTxSucceeded(ctx, l2client, delayedTx)
	Require(t, err)
	return receipt
}

func SendUnsignedTxViaL1(
	t *testing.T,
	ctx context.Context,
	l1info *BlockchainTestInfo,
	l1client mtutil.L1Interface,
	l2client mtutil.L1Interface,
	templateTx *types.Transaction,
) *types.Receipt {
	delayedInboxContract, err := bridgegen.NewInbox(l1info.GetAddress("Inbox"), l1client)
	Require(t, err)

	usertxopts := l1info.GetDefaultTransactOpts("User", ctx)
	remapped := util.RemapL1Address(usertxopts.From)
	nonce, err := l2client.NonceAt(ctx, remapped, nil)
	Require(t, err)

	unsignedTx := types.NewTx(&types.MantleUnsignedTx{
		ChainId:   templateTx.ChainId(),
		From:      remapped,
		Nonce:     nonce,
		GasFeeCap: templateTx.GasFeeCap(),
		Gas:       templateTx.Gas(),
		To:        templateTx.To(),
		Value:     templateTx.Value(),
		Data:      templateTx.Data(),
	})

	l1tx, err := delayedInboxContract.SendUnsignedTransaction(
		&usertxopts,
		mtmath.UintToBig(unsignedTx.Gas()),
		unsignedTx.GasFeeCap(),
		mtmath.UintToBig(unsignedTx.Nonce()),
		*unsignedTx.To(),
		unsignedTx.Value(),
		unsignedTx.Data(),
	)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, l1client, l1tx)
	Require(t, err)

	// sending l1 messages creates l1 blocks.. make enough to get that delayed inbox message in
	for i := 0; i < 30; i++ {
		SendWaitTestTransactions(t, ctx, l1client, []*types.Transaction{
			l1info.PrepareTx("Faucet", "Faucet", 30000, big.NewInt(1e12), nil),
		})
	}
	receipt, err := EnsureTxSucceeded(ctx, l2client, unsignedTx)
	Require(t, err)
	return receipt
}

func GetBaseFee(t *testing.T, client client, ctx context.Context) *big.Int {
	header, err := client.HeaderByNumber(ctx, nil)
	Require(t, err)
	return header.BaseFee
}

type lifecycle struct {
	start func() error
	stop  func() error
}

func (l *lifecycle) Start() error {
	if l.start != nil {
		return l.start()
	}
	return nil
}

func (l *lifecycle) Stop() error {
	if l.start != nil {
		return l.stop()
	}
	return nil
}

func createTestL1BlockChain(t *testing.T, l1info info) (info, *ethclient.Client, *eth.Ethereum, *node.Node) {
	if l1info == nil {
		l1info = NewL1TestInfo(t)
	}
	l1info.GenerateAccount("Faucet")

	chainConfig := params.MantleDevTestChainConfig()
	chainConfig.MantleChainParams = params.MantleChainParams{}

	stackConf := node.DefaultConfig
	stackConf.HTTPPort = 0
	stackConf.WSPort = 0
	stackConf.UseLightweightKDF = true
	stackConf.P2P.ListenAddr = ""
	stackConf.P2P.NoDial = true
	stackConf.P2P.NoDiscovery = true
	stackConf.P2P.NAT = nil
	var err error
	stackConf.DataDir = t.TempDir()
	stack, err := node.New(&stackConf)
	Require(t, err)

	nodeConf := ethconfig.Defaults
	nodeConf.NetworkId = chainConfig.ChainID.Uint64()
	l1Genesys := core.DeveloperGenesisBlock(0, 15_000_000, l1info.GetAddress("Faucet"))
	infoGenesys := l1info.GetGenesysAlloc()
	for acct, info := range infoGenesys {
		l1Genesys.Alloc[acct] = info
	}
	l1Genesys.BaseFee = big.NewInt(50 * params.GWei)
	nodeConf.Genesis = l1Genesys
	nodeConf.Miner.Etherbase = l1info.GetAddress("Faucet")

	l1backend, err := eth.New(stack, &nodeConf)
	Require(t, err)
	tempKeyStore := keystore.NewPlaintextKeyStore(t.TempDir())
	faucetAccount, err := tempKeyStore.ImportECDSA(l1info.Accounts["Faucet"].PrivateKey, "passphrase")
	Require(t, err)
	Require(t, tempKeyStore.Unlock(faucetAccount, "passphrase"))
	l1backend.AccountManager().AddBackend(tempKeyStore)
	l1backend.SetEtherbase(l1info.GetAddress("Faucet"))

	stack.RegisterLifecycle(&lifecycle{stop: func() error {
		l1backend.StopMining()
		return nil
	}})

	Require(t, stack.Start())
	Require(t, l1backend.StartMining(1))

	rpcClient, err := stack.Attach()
	Require(t, err)

	l1Client := ethclient.NewClient(rpcClient)

	return l1info, l1Client, l1backend, stack
}

func DeployOnTestL1(
	t *testing.T, ctx context.Context, l1info info, l1client client, chainId *big.Int,
) *mtnode.RollupAddresses {
	l1info.GenerateAccount("RollupOwner")
	l1info.GenerateAccount("Sequencer")
	l1info.GenerateAccount("User")

	SendWaitTestTransactions(t, ctx, l1client, []*types.Transaction{
		l1info.PrepareTx("Faucet", "RollupOwner", 30000, big.NewInt(9223372036854775807), nil),
		l1info.PrepareTx("Faucet", "Sequencer", 30000, big.NewInt(9223372036854775807), nil),
		l1info.PrepareTx("Faucet", "User", 30000, big.NewInt(9223372036854775807), nil)})

	l1TransactionOpts := l1info.GetDefaultTransactOpts("RollupOwner", ctx)
	config := mtnode.GenerateRollupConfig(false, common.Hash{}, l1info.GetAddress("RollupOwner"), chainId, common.Address{})
	addresses, err := mtnode.DeployOnL1(
		ctx,
		l1client,
		&l1TransactionOpts,
		l1info.GetAddress("Sequencer"),
		0,
		headerreader.TestConfig,
		validator.DefaultMantleMachineConfig,
		config,
	)
	Require(t, err)
	l1info.SetContract("Bridge", addresses.Bridge)
	l1info.SetContract("SequencerInbox", addresses.SequencerInbox)
	l1info.SetContract("Inbox", addresses.Inbox)
	return addresses
}

func createL2BlockChain(
	t *testing.T, l2info *BlockchainTestInfo, dataDir string, chainConfig *params.ChainConfig,
) (*BlockchainTestInfo, *node.Node, ethdb.Database, ethdb.Database, *core.BlockChain) {
	if l2info == nil {
		l2info = NewmttestInfo(t, chainConfig.ChainID)
	}
	stack, err := mtnode.CreateDefaultStackForTest(dataDir)
	Require(t, err)
	chainDb, err := stack.OpenDatabase("chaindb", 0, 0, "", false)
	Require(t, err)
	arbDb, err := stack.OpenDatabase("arbdb", 0, 0, "", false)
	Require(t, err)

	initReader := statetransfer.NewMemoryInitDataReader(&l2info.ArbInitData)
	blockchain, err := mtnode.WriteOrTestBlockChain(chainDb, nil, initReader, chainConfig, mtnode.ConfigDefaultL2Test(), 0)
	Require(t, err)

	return l2info, stack, chainDb, arbDb, blockchain
}

func ClientForStack(t *testing.T, backend *node.Node) *ethclient.Client {
	rpcClient, err := backend.Attach()
	Require(t, err)
	return ethclient.NewClient(rpcClient)
}

// Create and deploy L1 and mtnode for L2
func createTestNodeOnL1(
	t *testing.T,
	ctx context.Context,
	isSequencer bool,
) (
	l2info info, node *mtnode.Node, l2client *ethclient.Client, l2stack *node.Node, l1info info,
	l1backend *eth.Ethereum, l1client *ethclient.Client, l1stack *node.Node,
) {
	conf := mtnode.ConfigDefaultL1Test()
	return createTestNodeOnL1WithConfig(t, ctx, isSequencer, conf, params.MantleDevTestChainConfig())
}

func createTestNodeOnL1WithConfig(
	t *testing.T,
	ctx context.Context,
	isSequencer bool,
	nodeConfig *mtnode.Config,
	chainConfig *params.ChainConfig,
) (
	l2info info, currentNode *mtnode.Node, l2client *ethclient.Client, l2stack *node.Node, l1info info,
	l1backend *eth.Ethereum, l1client *ethclient.Client, l1stack *node.Node,
) {
	fatalErrChan := make(chan error, 10)
	l1info, l1client, l1backend, l1stack = createTestL1BlockChain(t, nil)
	var l2chainDb ethdb.Database
	var l2arbDb ethdb.Database
	var l2blockchain *core.BlockChain
	l2info, l2stack, l2chainDb, l2arbDb, l2blockchain = createL2BlockChain(t, nil, "", chainConfig)
	addresses := DeployOnTestL1(t, ctx, l1info, l1client, chainConfig.ChainID)
	var sequencerTxOptsPtr *bind.TransactOpts
	if isSequencer {
		sequencerTxOpts := l1info.GetDefaultTransactOpts("Sequencer", ctx)
		sequencerTxOptsPtr = &sequencerTxOpts
	}

	if !isSequencer {
		nodeConfig.BatchPoster.Enable = false
		nodeConfig.Sequencer.Enable = false
		nodeConfig.DelayedSequencer.Enable = false
	}

	var err error
	currentNode, err = mtnode.CreateNode(
		ctx, l2stack, l2chainDb, l2arbDb, nodeConfig, l2blockchain, l1client,
		addresses, sequencerTxOptsPtr, nil, fatalErrChan,
	)
	Require(t, err)

	Require(t, l2stack.Start())

	l2client = ClientForStack(t, l2stack)

	StartWatchChanErr(t, ctx, fatalErrChan, l2stack)

	return
}

// L2 -Only. Enough for tests that needs no interface to L1
// Requires precompiles.AllowDebugPrecompiles = true
func CreateTestL2(t *testing.T, ctx context.Context) (*BlockchainTestInfo, *mtnode.Node, *ethclient.Client, *node.Node) {
	return CreateTestL2WithConfig(t, ctx, nil, mtnode.ConfigDefaultL2Test(), true)
}

func CreateTestL2WithConfig(
	t *testing.T, ctx context.Context, l2Info *BlockchainTestInfo, nodeConfig *mtnode.Config, takeOwnership bool,
) (*BlockchainTestInfo, *mtnode.Node, *ethclient.Client, *node.Node) {
	feedErrChan := make(chan error, 10)
	l2info, stack, chainDb, arbDb, blockchain := createL2BlockChain(t, l2Info, "", params.MantleDevTestChainConfig())
	currentNode, err := mtnode.CreateNode(ctx, stack, chainDb, arbDb, nodeConfig, blockchain, nil, nil, nil, nil, feedErrChan)
	Require(t, err)

	// Give the node an init message
	err = currentNode.TxStreamer.AddFakeInitMessage()
	Require(t, err)

	Require(t, stack.Start())
	client := ClientForStack(t, stack)

	if takeOwnership {
		debugAuth := l2info.GetDefaultTransactOpts("Owner", ctx)

		// make auth a chain owner
		arbdebug, err := precompilesgen.NewMtDebug(common.HexToAddress("0xff"), client)
		Require(t, err, "failed to deploy MtDebug")

		tx, err := arbdebug.BecomeChainOwner(&debugAuth)
		Require(t, err, "failed to deploy MtDebug")

		_, err = EnsureTxSucceeded(ctx, client, tx)
		Require(t, err)
	}

	StartWatchChanErr(t, ctx, feedErrChan, stack)

	return l2info, currentNode, client, stack
}

func StartWatchChanErr(t *testing.T, ctx context.Context, feedErrChan chan error, stack *node.Node) {
	go func() {
		select {
		case <-ctx.Done():
			return
		case err := <-feedErrChan:
			t.Errorf("error occurred: %v", err)
			if stack != nil {
				err = stack.Close()
				if err != nil {
					t.Errorf("error closing stack: %v", err)
				}
			}
		}
	}()
}

func Require(t *testing.T, err error, text ...interface{}) {
	t.Helper()
	testhelpers.RequireImpl(t, err, text...)
}

func Fail(t *testing.T, printables ...interface{}) {
	t.Helper()
	testhelpers.FailImpl(t, printables...)
}

func Create2ndNode(
	t *testing.T,
	ctx context.Context,
	first *mtnode.Node,
	l1stack *node.Node,
	l2InitData *statetransfer.MtosInitializationInfo,
	dasConfig *das.DataAvailabilityConfig,
) (*ethclient.Client, *mtnode.Node, *node.Node) {
	nodeConf := mtnode.ConfigDefaultL1NonSequencerTest()
	if dasConfig == nil {
		nodeConf.DataAvailability.Enable = false
	} else {
		nodeConf.DataAvailability = *dasConfig
	}
	return Create2ndNodeWithConfig(t, ctx, first, l1stack, l2InitData, nodeConf)
}

func Create2ndNodeWithConfig(
	t *testing.T,
	ctx context.Context,
	first *mtnode.Node,
	l1stack *node.Node,
	l2InitData *statetransfer.MtosInitializationInfo,
	nodeConfig *mtnode.Config,
) (*ethclient.Client, *mtnode.Node, *node.Node) {
	feedErrChan := make(chan error, 10)
	l1rpcClient, err := l1stack.Attach()
	if err != nil {
		Fail(t, err)
	}
	l1client := ethclient.NewClient(l1rpcClient)
	l2stack, err := mtnode.CreateDefaultStackForTest("")
	Require(t, err)

	l2chainDb, err := l2stack.OpenDatabase("chaindb", 0, 0, "", false)
	Require(t, err)
	l2arbDb, err := l2stack.OpenDatabase("arbdb", 0, 0, "", false)
	Require(t, err)
	initReader := statetransfer.NewMemoryInitDataReader(l2InitData)

	l2blockchain, err := mtnode.WriteOrTestBlockChain(l2chainDb, nil, initReader, first.MtInterface.BlockChain().Config(), mtnode.ConfigDefaultL2Test(), 0)
	Require(t, err)

	currentNode, err := mtnode.CreateNode(ctx, l2stack, l2chainDb, l2arbDb, nodeConfig, l2blockchain, l1client, first.DeployInfo, nil, nil, feedErrChan)
	Require(t, err)

	err = l2stack.Start()
	Require(t, err)
	l2client := ClientForStack(t, l2stack)

	StartWatchChanErr(t, ctx, feedErrChan, l1stack)

	return l2client, currentNode, l2stack
}

func GetBalance(t *testing.T, ctx context.Context, client *ethclient.Client, account common.Address) *big.Int {
	t.Helper()
	balance, err := client.BalanceAt(ctx, account, nil)
	Require(t, err, "could not get balance")
	return balance
}

func requireClose(t *testing.T, s *node.Node, text ...interface{}) {
	t.Helper()
	Require(t, s.Close(), text...)
}

func authorizeDASKeyset(
	t *testing.T,
	ctx context.Context,
	dasSignerKey *blsSignatures.PublicKey,
	l1info info,
	l1client mtutil.L1Interface,
) {
	if dasSignerKey == nil {
		return
	}
	keyset := &mtstate.DataAvailabilityKeyset{
		AssumedHonest: 1,
		PubKeys:       []blsSignatures.PublicKey{*dasSignerKey},
	}
	wr := bytes.NewBuffer([]byte{})
	err := keyset.Serialize(wr)
	Require(t, err, "unable to serialize DAS keyset")
	keysetBytes := wr.Bytes()
	sequencerInbox, err := bridgegen.NewSequencerInbox(l1info.Accounts["SequencerInbox"].Address, l1client)
	Require(t, err, "unable to create sequencer inbox")
	trOps := l1info.GetDefaultTransactOpts("RollupOwner", ctx)
	tx, err := sequencerInbox.SetValidKeyset(&trOps, keysetBytes)
	Require(t, err, "unable to set valid keyset")
	_, err = EnsureTxSucceeded(ctx, l1client, tx)
	Require(t, err, "unable to ensure transaction success for setting valid keyset")
}

func setupConfigWithDAS(
	t *testing.T, ctx context.Context, dasModeString string,
) (*params.ChainConfig, *mtnode.Config, *das.LifecycleManager, string, *blsSignatures.PublicKey) {
	l1NodeConfigA := mtnode.ConfigDefaultL1Test()
	chainConfig := params.MantleDevTestChainConfig()
	var dbPath string
	var err error

	enableFileStorage, enableDbStorage, enableDas := false, false, true
	switch dasModeString {
	case "db":
		enableDbStorage = true
		chainConfig = params.MantleDevTestDASChainConfig()
	case "files":
		enableFileStorage = true
		chainConfig = params.MantleDevTestDASChainConfig()
	case "onchain":
		enableDas = false
	default:
		Fail(t, "unknown storage type")
	}
	dbPath = t.TempDir()
	dasSignerKey, _, err := das.GenerateAndStoreKeys(dbPath)
	Require(t, err)

	dasConfig := &das.DataAvailabilityConfig{
		Enable: enableDas,
		KeyConfig: das.KeyConfig{
			KeyDir: dbPath,
		},
		LocalFileStorageConfig: das.LocalFileStorageConfig{
			Enable:  enableFileStorage,
			DataDir: dbPath,
		},
		LocalDBStorageConfig: das.LocalDBStorageConfig{
			Enable:  enableDbStorage,
			DataDir: dbPath,
		},
		RequestTimeout:           5 * time.Second,
		L1NodeURL:                "none",
		SequencerInboxAddress:    "none",
		PanicOnError:             true,
		DisableSignatureChecking: true,
	}

	l1NodeConfigA.DataAvailability = das.DefaultDataAvailabilityConfig
	var lifecycleManager *das.LifecycleManager
	if dasModeString != "onchain" {
		var dasServerStack das.DataAvailabilityService
		dasServerStack, lifecycleManager, err = mtnode.SetUpDataAvailability(ctx, dasConfig, nil, nil)

		Require(t, err)
		rpcLis, err := net.Listen("tcp", "localhost:0")
		Require(t, err)
		restLis, err := net.Listen("tcp", "localhost:0")
		Require(t, err)
		_, err = das.StartDASRPCServerOnListener(ctx, rpcLis, genericconf.HTTPServerTimeoutConfigDefault, dasServerStack)
		Require(t, err)
		_, err = das.NewRestfulDasServerOnListener(restLis, genericconf.HTTPServerTimeoutConfigDefault, dasServerStack)
		Require(t, err)

		beConfigA := das.BackendConfig{
			URL:                 "http://" + rpcLis.Addr().String(),
			PubKeyBase64Encoded: blsPubToBase64(dasSignerKey),
			SignerMask:          1,
		}
		l1NodeConfigA.DataAvailability.AggregatorConfig = aggConfigForBackend(t, beConfigA)
		l1NodeConfigA.DataAvailability.Enable = true
		l1NodeConfigA.DataAvailability.RestfulClientAggregatorConfig = das.DefaultRestfulClientAggregatorConfig
		l1NodeConfigA.DataAvailability.RestfulClientAggregatorConfig.Enable = true
		l1NodeConfigA.DataAvailability.RestfulClientAggregatorConfig.Urls = []string{"http://" + restLis.Addr().String()}
		l1NodeConfigA.DataAvailability.L1NodeURL = "none"
	}

	return chainConfig, l1NodeConfigA, lifecycleManager, dbPath, dasSignerKey
}

func getDeadlineTimeout(t *testing.T, defaultTimeout time.Duration) time.Duration {
	testDeadLine, deadlineExist := t.Deadline()
	var timeout time.Duration
	if deadlineExist {
		timeout = time.Until(testDeadLine) - (time.Second * 10)
		if timeout > time.Second*10 {
			timeout = timeout - (time.Second * 10)
		}
	} else {
		timeout = defaultTimeout
	}

	return timeout
}

func deploySimple(
	t *testing.T, ctx context.Context, auth bind.TransactOpts, client *ethclient.Client,
) (common.Address, *mocksgen.Simple) {
	addr, tx, simple, err := mocksgen.DeploySimple(&auth, client)
	Require(t, err, "could not deploy Simple.sol contract")
	_, err = EnsureTxSucceeded(ctx, client, tx)
	Require(t, err)
	return addr, simple
}
