// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package mttest

import (
	"bytes"
	"context"
	"io"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/mantlenetworkio/mantle/mtcompress"
	"github.com/mantlenetworkio/mantle/mtnode"
	"github.com/mantlenetworkio/mantle/mtos"
	"github.com/mantlenetworkio/mantle/mtstate"
	"github.com/mantlenetworkio/mantle/mtutil"
	"github.com/mantlenetworkio/mantle/solgen/go/challengegen"
	"github.com/mantlenetworkio/mantle/solgen/go/mocksgen"
	"github.com/mantlenetworkio/mantle/solgen/go/ospgen"
	"github.com/mantlenetworkio/mantle/validator"
)

func DeployOneStepProofEntry(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client) common.Address {
	osp0, _, _, err := ospgen.DeployOneStepProver0(auth, client)
	if err != nil {
		Fail(t, err)
	}
	ospMem, _, _, err := ospgen.DeployOneStepProverMemory(auth, client)
	if err != nil {
		Fail(t, err)
	}
	ospMath, _, _, err := ospgen.DeployOneStepProverMath(auth, client)
	if err != nil {
		Fail(t, err)
	}
	ospHostIo, _, _, err := ospgen.DeployOneStepProverHostIo(auth, client)
	if err != nil {
		Fail(t, err)
	}
	ospEntry, tx, _, err := ospgen.DeployOneStepProofEntry(auth, client, osp0, ospMem, ospMath, ospHostIo)
	if err != nil {
		Fail(t, err)
	}
	_, err = EnsureTxSucceeded(ctx, client, tx)
	if err != nil {
		Fail(t, err)
	}
	return ospEntry
}

func CreateChallenge(
	t *testing.T,
	ctx context.Context,
	auth *bind.TransactOpts,
	client *ethclient.Client,
	ospEntry common.Address,
	sequencerInbox common.Address,
	delayedBridge common.Address,
	wasmModuleRoot common.Hash,
	startGlobalState validator.GoGlobalState,
	endGlobalState validator.GoGlobalState,
	numBlocks uint64,
	asserter common.Address,
	challenger common.Address,
) (*mocksgen.MockResultReceiver, common.Address) {
	challengeManagerLogic, tx, _, err := challengegen.DeployChallengeManager(auth, client)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, client, tx)
	Require(t, err)
	challengeManagerAddr, tx, _, err := mocksgen.DeploySimpleProxy(auth, client, challengeManagerLogic)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, client, tx)
	Require(t, err)
	challengeManager, err := challengegen.NewChallengeManager(challengeManagerAddr, client)
	Require(t, err)

	resultReceiverAddr, _, resultReceiver, err := mocksgen.DeployMockResultReceiver(auth, client, challengeManagerAddr)
	Require(t, err)
	tx, err = challengeManager.Initialize(auth, resultReceiverAddr, sequencerInbox, delayedBridge, ospEntry)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, client, tx)
	Require(t, err)
	tx, err = resultReceiver.CreateChallenge(
		auth,
		wasmModuleRoot,
		[2]uint8{
			validator.StatusFinished,
			validator.StatusFinished,
		},
		[2]mocksgen.GlobalState{
			{
				Bytes32Vals: [2][32]byte{startGlobalState.BlockHash, startGlobalState.SendRoot},
				U64Vals:     [2]uint64{startGlobalState.Batch, startGlobalState.PosInBatch},
			},
			{
				Bytes32Vals: [2][32]byte{endGlobalState.BlockHash, endGlobalState.SendRoot},
				U64Vals:     [2]uint64{endGlobalState.Batch, endGlobalState.PosInBatch},
			},
		},
		numBlocks,
		asserter,
		challenger,
		big.NewInt(100000),
		big.NewInt(100000),
	)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, client, tx)
	Require(t, err)
	return resultReceiver, challengeManagerAddr
}

func writeTxToBatch(writer io.Writer, tx *types.Transaction) error {
	txData, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	var segment []byte
	segment = append(segment, mtstate.BatchSegmentKindL2Message)
	segment = append(segment, mtos.L2MessageKind_SignedTx)
	segment = append(segment, txData...)
	err = rlp.Encode(writer, segment)
	return err
}

func makeBatch(t *testing.T, l2Node *mtnode.Node, l2Info *BlockchainTestInfo, backend *ethclient.Client, sequencer *bind.TransactOpts, seqInbox *mocksgen.SequencerInboxStub, seqInboxAddr common.Address, isChallenger bool) {
	ctx := context.Background()

	batchBuffer := bytes.NewBuffer([]byte{})
	for i := int64(0); i < 10; i++ {
		value := i
		if i == 5 && isChallenger {
			value++
		}
		err := writeTxToBatch(batchBuffer, l2Info.PrepareTx("Owner", "Destination", 1000000, big.NewInt(value), []byte{}))
		Require(t, err)
	}
	compressed, err := mtcompress.CompressWell(batchBuffer.Bytes())
	Require(t, err)
	message := append([]byte{0}, compressed...)

	tx, err := seqInbox.AddSequencerL2BatchFromOrigin0(sequencer, big.NewInt(0), message, big.NewInt(0), common.Address{}, big.NewInt(0), big.NewInt(1))
	Require(t, err)
	receipt, err := EnsureTxSucceeded(ctx, backend, tx)
	Require(t, err)

	nodeSeqInbox, err := mtnode.NewSequencerInbox(backend, seqInboxAddr, 0)
	Require(t, err)
	batches, err := nodeSeqInbox.LookupBatchesInRange(ctx, receipt.BlockNumber, receipt.BlockNumber)
	Require(t, err)
	if len(batches) == 0 {
		Fail(t, "batch not found after AddSequencerL2BatchFromOrigin")
	}
	err = l2Node.InboxTracker.AddSequencerBatches(ctx, backend, batches)
	Require(t, err)
	_, err = l2Node.InboxTracker.GetBatchMetadata(0)
	Require(t, err, "failed to get batch metadata after adding batch:")
}

func confirmLatestBlock(ctx context.Context, t *testing.T, l1Info *BlockchainTestInfo, backend mtutil.L1Interface) {
	for i := 0; i < 12; i++ {
		SendWaitTestTransactions(t, ctx, backend, []*types.Transaction{
			l1Info.PrepareTx("Faucet", "Faucet", 30000, big.NewInt(1e12), nil),
		})
	}
}

func RunChallengeTest(t *testing.T, asserterIsCorrect bool) {
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.LvlInfo)
	log.Root().SetHandler(glogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initialBalance := new(big.Int).Lsh(big.NewInt(1), 200)
	l1Info := NewL1TestInfo(t)
	l1Info.GenerateGenesysAccount("deployer", initialBalance)
	l1Info.GenerateGenesysAccount("asserter", initialBalance)
	l1Info.GenerateGenesysAccount("challenger", initialBalance)
	l1Info.GenerateGenesysAccount("sequencer", initialBalance)

	chainConfig := params.MantleDevTestChainConfig()
	l1Info, l1Backend, _, _ := createTestL1BlockChain(t, l1Info)
	conf := mtnode.ConfigDefaultL1Test()
	conf.BlockValidator.Enable = false
	conf.BatchPoster.Enable = false
	conf.InboxReader.CheckDelay = time.Second
	fatalErrChan := make(chan error, 10)
	rollupAddresses := DeployOnTestL1(t, ctx, l1Info, l1Backend, chainConfig.ChainID)

	deployerTxOpts := l1Info.GetDefaultTransactOpts("deployer", ctx)
	sequencerTxOpts := l1Info.GetDefaultTransactOpts("sequencer", ctx)
	asserterTxOpts := l1Info.GetDefaultTransactOpts("asserter", ctx)
	challengerTxOpts := l1Info.GetDefaultTransactOpts("challenger", ctx)
	delayedBridge, tx, _, err := mocksgen.DeployBridgeStub(&deployerTxOpts, l1Backend)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, l1Backend, tx)
	Require(t, err)

	timeBounds := mocksgen.ISequencerInboxMaxTimeVariation{
		DelayBlocks:   big.NewInt(10000),
		FutureBlocks:  big.NewInt(10000),
		DelaySeconds:  big.NewInt(10000),
		FutureSeconds: big.NewInt(10000),
	}
	asserterSeqInboxAddr, tx, asserterSeqInbox, err := mocksgen.DeploySequencerInboxStub(
		&deployerTxOpts,
		l1Backend,
		delayedBridge,
		l1Info.GetAddress("sequencer"),
		timeBounds,
	)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, l1Backend, tx)
	Require(t, err)
	tx, err = asserterSeqInbox.AddInitMessage(&deployerTxOpts)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, l1Backend, tx)
	Require(t, err)
	challengerSeqInboxAddr, tx, challengerSeqInbox, err := mocksgen.DeploySequencerInboxStub(
		&deployerTxOpts,
		l1Backend,
		delayedBridge,
		l1Info.GetAddress("sequencer"),
		timeBounds,
	)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, l1Backend, tx)
	Require(t, err)
	tx, err = challengerSeqInbox.AddInitMessage(&deployerTxOpts)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, l1Backend, tx)
	Require(t, err)

	asserterL2Info, asserterL2Stack, asserterL2ChainDb, asserterL2MtDb, asserterL2Blockchain := createL2BlockChain(t, nil, "", chainConfig)
	rollupAddresses.SequencerInbox = asserterSeqInboxAddr
	asserterL2, err := mtnode.CreateNode(ctx, asserterL2Stack, asserterL2ChainDb, asserterL2MtDb, conf, asserterL2Blockchain, l1Backend, rollupAddresses, nil, nil, fatalErrChan)
	Require(t, err)
	err = asserterL2Stack.Start()
	Require(t, err)

	challengerL2Info, challengerL2Stack, challengerL2ChainDb, challengerL2MtDb, challengerL2Blockchain := createL2BlockChain(t, nil, "", chainConfig)
	rollupAddresses.SequencerInbox = challengerSeqInboxAddr
	challengerL2, err := mtnode.CreateNode(ctx, challengerL2Stack, challengerL2ChainDb, challengerL2MtDb, conf, challengerL2Blockchain, l1Backend, rollupAddresses, nil, nil, fatalErrChan)
	Require(t, err)
	err = challengerL2Stack.Start()
	Require(t, err)

	asserterL2Info.GenerateAccount("Destination")
	challengerL2Info.SetFullAccountInfo("Destination", asserterL2Info.GetInfoWithPrivKey("Destination"))
	makeBatch(t, asserterL2, asserterL2Info, l1Backend, &sequencerTxOpts, asserterSeqInbox, asserterSeqInboxAddr, false)
	makeBatch(t, challengerL2, challengerL2Info, l1Backend, &sequencerTxOpts, challengerSeqInbox, challengerSeqInboxAddr, true)

	trueSeqInboxAddr := challengerSeqInboxAddr
	expectedWinner := l1Info.GetAddress("challenger")
	if asserterIsCorrect {
		trueSeqInboxAddr = asserterSeqInboxAddr
		expectedWinner = l1Info.GetAddress("asserter")
	}
	ospEntry := DeployOneStepProofEntry(t, ctx, &deployerTxOpts, l1Backend)

	wasmModuleRoot, err := validator.DefaultMantleMachineConfig.ReadLatestWasmModuleRoot()
	if err != nil {
		Fail(t, err)
	}

	asserterGenesis := asserterL2.MtInterface.BlockChain().Genesis()
	challengerGenesis := challengerL2.MtInterface.BlockChain().Genesis()
	if asserterGenesis.Hash() != challengerGenesis.Hash() {
		Fail(t, "asserter and challenger have different genesis hashes")
	}
	asserterLatestBlock := asserterL2.MtInterface.BlockChain().CurrentBlock()
	challengerLatestBlock := challengerL2.MtInterface.BlockChain().CurrentBlock()
	if asserterLatestBlock.Hash() == challengerLatestBlock.Hash() {
		Fail(t, "asserter and challenger have the same end block")
	}

	asserterStartGlobalState := validator.GoGlobalState{
		BlockHash:  asserterGenesis.Hash(),
		Batch:      1,
		PosInBatch: 0,
	}
	asserterEndGlobalState := validator.GoGlobalState{
		BlockHash:  asserterLatestBlock.Hash(),
		Batch:      2,
		PosInBatch: 0,
	}
	numBlocks := asserterLatestBlock.NumberU64() - asserterGenesis.NumberU64()

	resultReceiver, challengeManagerAddr := CreateChallenge(
		t,
		ctx,
		&deployerTxOpts,
		l1Backend,
		ospEntry,
		trueSeqInboxAddr,
		delayedBridge,
		wasmModuleRoot,
		asserterStartGlobalState,
		asserterEndGlobalState,
		numBlocks,
		l1Info.GetAddress("asserter"),
		l1Info.GetAddress("challenger"),
	)

	confirmLatestBlock(ctx, t, l1Info, l1Backend)
	machineLoader := validator.NewMantleMachineLoader(validator.DefaultMantleMachineConfig, fatalErrChan)
	asserterManager, err := validator.NewChallengeManager(ctx, l1Backend, &asserterTxOpts, asserterTxOpts.From, challengeManagerAddr, 1, asserterL2Blockchain, nil, asserterL2.InboxReader, asserterL2.InboxTracker, asserterL2.TxStreamer, machineLoader, 0, 4, 0)
	if err != nil {
		Fail(t, err)
	}

	challengerManager, err := validator.NewChallengeManager(ctx, l1Backend, &challengerTxOpts, challengerTxOpts.From, challengeManagerAddr, 1, challengerL2Blockchain, nil, challengerL2.InboxReader, challengerL2.InboxTracker, challengerL2.TxStreamer, machineLoader, 0, 4, 0)
	if err != nil {
		Fail(t, err)
	}

	for i := 0; i < 100; i++ {
		var tx *types.Transaction
		var currentCorrect bool
		// Gas cost is slightly reduced if done in the same timestamp or block as previous call.
		// This might make gas estimation undersestimate next move.
		// Invoke a new L1 block, with a new timestamp, before estimating.
		time.Sleep(time.Second)
		SendWaitTestTransactions(t, ctx, l1Backend, []*types.Transaction{
			l1Info.PrepareTx("Faucet", "User", 30000, big.NewInt(1e12), nil),
		})

		if i%2 == 0 {
			currentCorrect = !asserterIsCorrect
			tx, err = challengerManager.Act(ctx)
		} else {
			currentCorrect = asserterIsCorrect
			tx, err = asserterManager.Act(ctx)
		}
		if err != nil {
			if !currentCorrect && (strings.Contains(err.Error(), "lost challenge") ||
				strings.Contains(err.Error(), "SAME_OSP_END") ||
				strings.Contains(err.Error(), "BAD_SEQINBOX_MESSAGE")) {
				t.Log("challenge completed! asserter hit expected error:", err)
				return
			}
			Fail(t, "challenge step", i, "hit error:", err)
		}
		if tx == nil {
			Fail(t, "no move")
		}
		_, err = EnsureTxSucceeded(ctx, l1Backend, tx)
		if err != nil {
			if !currentCorrect && strings.Contains(err.Error(), "BAD_SEQINBOX_MESSAGE") {
				t.Log("challenge complete! Tx failed as expected:", err)
				return
			}
			Fail(t, err)
		}

		confirmLatestBlock(ctx, t, l1Info, l1Backend)

		winner, err := resultReceiver.Winner(&bind.CallOpts{})
		if err != nil {
			Fail(t, err)
		}
		if winner == (common.Address{}) {
			continue
		}
		if winner != expectedWinner {
			Fail(t, "wrong party won challenge")
		}
	}

	Fail(t, "challenge timed out without winner")
}
