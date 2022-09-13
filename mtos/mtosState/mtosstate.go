// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package mtosState

import (
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"

	"github.com/mantlenetworkio/mantle/mtos/blockhash"
	"github.com/mantlenetworkio/mantle/mtos/l2pricing"

	"github.com/mantlenetworkio/mantle/mtos/addressSet"
	"github.com/mantlenetworkio/mantle/mtos/burn"

	"github.com/mantlenetworkio/mantle/mtos/addressTable"
	"github.com/mantlenetworkio/mantle/mtos/l1pricing"
	"github.com/mantlenetworkio/mantle/mtos/merkleAccumulator"
	"github.com/mantlenetworkio/mantle/mtos/retryables"
	"github.com/mantlenetworkio/mantle/mtos/storage"
	"github.com/mantlenetworkio/mantle/mtos/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

// MtosState contains MtOS-related state. It is backed by MtOS's storage in the persistent stateDB.
// Modifications to the MtosState are written through to the underlying StateDB so that the StateDB always
// has the definitive state, stored persistently. (Note that some tests use memory-backed StateDB's that aren't
// persisted beyond the end of the test.)

type MtosState struct {
	mtosVersion       uint64                      // version of the MtOS storage format and semantics
	upgradeVersion    storage.StorageBackedUint64 // version we're planning to upgrade to, or 0 if not planning to upgrade
	upgradeTimestamp  storage.StorageBackedUint64 // when to do the planned upgrade
	networkFeeAccount storage.StorageBackedAddress
	l1PricingState    *l1pricing.L1PricingState
	l2PricingState    *l2pricing.L2PricingState
	retryableState    *retryables.RetryableState
	addressTable      *addressTable.AddressTable
	chainOwners       *addressSet.AddressSet
	sendMerkle        *merkleAccumulator.MerkleAccumulator
	blockhashes       *blockhash.Blockhashes
	chainId           storage.StorageBackedBigInt
	genesisBlockNum   storage.StorageBackedUint64
	infraFeeAccount   storage.StorageBackedAddress
	backingStorage    *storage.Storage
	Burner            burn.Burner
}

var ErrUninitializedMtOS = errors.New("MtOS uninitialized")
var ErrAlreadyInitialized = errors.New("MtOS is already initialized")

func OpenMtosState(stateDB vm.StateDB, burner burn.Burner) (*MtosState, error) {
	backingStorage := storage.NewGeth(stateDB, burner)
	mtosVersion, err := backingStorage.GetUint64ByUint64(uint64(versionOffset))
	if err != nil {
		return nil, err
	}
	if mtosVersion == 0 {
		return nil, ErrUninitializedMtOS
	}
	return &MtosState{
		mtosVersion,
		backingStorage.OpenStorageBackedUint64(uint64(upgradeVersionOffset)),
		backingStorage.OpenStorageBackedUint64(uint64(upgradeTimestampOffset)),
		backingStorage.OpenStorageBackedAddress(uint64(networkFeeAccountOffset)),
		l1pricing.OpenL1PricingState(backingStorage.OpenSubStorage(l1PricingSubspace)),
		l2pricing.OpenL2PricingState(backingStorage.OpenSubStorage(l2PricingSubspace)),
		retryables.OpenRetryableState(backingStorage.OpenSubStorage(retryablesSubspace), stateDB),
		addressTable.Open(backingStorage.OpenSubStorage(addressTableSubspace)),
		addressSet.OpenAddressSet(backingStorage.OpenSubStorage(chainOwnerSubspace)),
		merkleAccumulator.OpenMerkleAccumulator(backingStorage.OpenSubStorage(sendMerkleSubspace)),
		blockhash.OpenBlockhashes(backingStorage.OpenSubStorage(blockhashesSubspace)),
		backingStorage.OpenStorageBackedBigInt(uint64(chainIdOffset)),
		backingStorage.OpenStorageBackedUint64(uint64(genesisBlockNumOffset)),
		backingStorage.OpenStorageBackedAddress(uint64(infraFeeAccountOffset)),
		backingStorage,
		burner,
	}, nil
}

func OpenSystemMtosState(stateDB vm.StateDB, tracingInfo *util.TracingInfo, readOnly bool) (*MtosState, error) {
	burner := burn.NewSystemBurner(tracingInfo, readOnly)
	state, err := OpenMtosState(stateDB, burner)
	burner.Restrict(err)
	return state, err
}

func OpenSystemMtosStateOrPanic(stateDB vm.StateDB, tracingInfo *util.TracingInfo, readOnly bool) *MtosState {
	state, err := OpenSystemMtosState(stateDB, tracingInfo, readOnly)
	if err != nil {
		panic(err)
	}
	return state
}

// Create and initialize a memory-backed MtOS state (for testing only)
func NewMtosMemoryBackedMtOSState() (*MtosState, *state.StateDB) {
	raw := rawdb.NewMemoryDatabase()
	db := state.NewDatabase(raw)
	statedb, err := state.New(common.Hash{}, db, nil)
	if err != nil {
		log.Fatal("failed to init empty statedb", err)
	}
	burner := burn.NewSystemBurner(nil, false)
	state, err := InitializeMtosState(statedb, burner, params.MantleDevTestChainConfig())
	if err != nil {
		log.Fatal("failed to open the MtOS state", err)
	}
	return state, statedb
}

// Get the MtOS version
func MtOSVersion(stateDB vm.StateDB) uint64 {
	backingStorage := storage.NewGeth(stateDB, burn.NewSystemBurner(nil, false))
	mtosVersion, err := backingStorage.GetUint64ByUint64(uint64(versionOffset))
	if err != nil {
		log.Fatal("faled to get the MtOS version", err)
	}
	return mtosVersion
}

type MtosStateOffset uint64

const (
	versionOffset MtosStateOffset = iota
	upgradeVersionOffset
	upgradeTimestampOffset
	networkFeeAccountOffset
	chainIdOffset
	genesisBlockNumOffset
	infraFeeAccountOffset
)

type MtosStateSubspaceID []byte

var (
	l1PricingSubspace    MtosStateSubspaceID = []byte{0}
	l2PricingSubspace    MtosStateSubspaceID = []byte{1}
	retryablesSubspace   MtosStateSubspaceID = []byte{2}
	addressTableSubspace MtosStateSubspaceID = []byte{3}
	chainOwnerSubspace   MtosStateSubspaceID = []byte{4}
	sendMerkleSubspace   MtosStateSubspaceID = []byte{5}
	blockhashesSubspace  MtosStateSubspaceID = []byte{6}
)

// Returns a list of precompiles that only appear in Mantle chains (i.e. MtOS precompiles) at the genesis block
func getMantleOnlyPrecompiles(chainConfig *params.ChainConfig) []common.Address {
	rules := chainConfig.Rules(big.NewInt(0), false)
	arbPrecompiles := vm.ActivePrecompiles(rules)
	rules.IsMantle = false
	ethPrecompiles := vm.ActivePrecompiles(rules)

	ethPrecompilesSet := make(map[common.Address]bool)
	for _, addr := range ethPrecompiles {
		ethPrecompilesSet[addr] = true
	}

	var arbOnlyPrecompiles []common.Address
	for _, addr := range arbPrecompiles {
		if !ethPrecompilesSet[addr] {
			arbOnlyPrecompiles = append(arbOnlyPrecompiles, addr)
		}
	}
	return arbOnlyPrecompiles
}

// During early development we sometimes change the storage format of version 1, for convenience. But as soon as we
// start running long-lived chains, every change to the storage format will require defining a new version and
// providing upgrade code.

func InitializeMtosState(stateDB vm.StateDB, burner burn.Burner, chainConfig *params.ChainConfig) (*MtosState, error) {
	sto := storage.NewGeth(stateDB, burner)
	mtosVersion, err := sto.GetUint64ByUint64(uint64(versionOffset))
	if err != nil {
		return nil, err
	}
	if mtosVersion != 0 {
		return nil, ErrAlreadyInitialized
	}

	desiredMtosVersion := chainConfig.MantleChainParams.InitialMtOSVersion
	if desiredMtosVersion == 0 {
		return nil, errors.New("cannot initialize to MtOS version 0")
	}

	// Solidity requires call targets have code, but precompiles don't.
	// To work around this, we give precompiles fake code.
	for _, precompile := range getMantleOnlyPrecompiles(chainConfig) {
		stateDB.SetCode(precompile, []byte{byte(vm.INVALID)})
	}

	// may be the zero address
	initialChainOwner := chainConfig.MantleChainParams.InitialChainOwner

	_ = sto.SetUint64ByUint64(uint64(versionOffset), 1) // initialize to version 1; upgrade at end of this func if needed
	_ = sto.SetUint64ByUint64(uint64(upgradeVersionOffset), 0)
	_ = sto.SetUint64ByUint64(uint64(upgradeTimestampOffset), 0)
	if desiredMtosVersion >= 2 {
		_ = sto.SetByUint64(uint64(networkFeeAccountOffset), util.AddressToHash(initialChainOwner))
	} else {
		_ = sto.SetByUint64(uint64(networkFeeAccountOffset), common.Hash{}) // the 0 address until an owner sets it
	}
	_ = sto.SetByUint64(uint64(chainIdOffset), common.BigToHash(chainConfig.ChainID))
	_ = sto.SetUint64ByUint64(uint64(genesisBlockNumOffset), chainConfig.MantleChainParams.GenesisBlockNum)

	initialRewardsRecipient := l1pricing.BatchPosterAddress
	if desiredMtosVersion >= 2 {
		initialRewardsRecipient = initialChainOwner
	}
	_ = l1pricing.InitializeL1PricingState(sto.OpenSubStorage(l1PricingSubspace), initialRewardsRecipient)
	_ = l2pricing.InitializeL2PricingState(sto.OpenSubStorage(l2PricingSubspace))
	_ = retryables.InitializeRetryableState(sto.OpenSubStorage(retryablesSubspace))
	addressTable.Initialize(sto.OpenSubStorage(addressTableSubspace))
	merkleAccumulator.InitializeMerkleAccumulator(sto.OpenSubStorage(sendMerkleSubspace))
	blockhash.InitializeBlockhashes(sto.OpenSubStorage(blockhashesSubspace))

	ownersStorage := sto.OpenSubStorage(chainOwnerSubspace)
	_ = addressSet.Initialize(ownersStorage)
	_ = addressSet.OpenAddressSet(ownersStorage).Add(initialChainOwner)

	aState, err := OpenMtosState(stateDB, burner)
	if err != nil {
		return nil, err
	}
	if desiredMtosVersion > 1 {
		aState.UpgradeMtosVersion(desiredMtosVersion, true)
	}
	return aState, err
}

func (state *MtosState) UpgradeMtosVersionIfNecessary(currentTimestamp uint64) {
	upgradeTo, err := state.upgradeVersion.Get()
	state.Restrict(err)
	flagday, _ := state.upgradeTimestamp.Get()
	if state.mtosVersion < upgradeTo && currentTimestamp >= flagday {
		state.UpgradeMtosVersion(upgradeTo, false)
	}
}

func (state *MtosState) UpgradeMtosVersion(upgradeTo uint64, firstTime bool) {
	for state.mtosVersion < upgradeTo {
		ensure := func(err error) {
			if err != nil {
				message := fmt.Sprintf(
					"Failed to upgrade MtOS version %v to version %v: %v",
					state.mtosVersion, state.mtosVersion+1, err,
				)
				panic(message)
			}
		}

		switch state.mtosVersion {
		case 1:
			ensure(state.l1PricingState.SetLastSurplus(common.Big0))
		case 2:
			ensure(state.l1PricingState.SetPerBatchGasCost(0))
			ensure(state.l1PricingState.SetAmortizedCostCapBips(math.MaxUint64))
		case 3:
			// no state changes needed
		case 4:
			// no state changes needed
		case 5:
			// no state changes needed
		default:
			panic("Unable to perform requested MtOS upgrade")
		}
		state.mtosVersion++
	}

	if firstTime && upgradeTo >= 6 {
		state.Restrict(state.l1PricingState.SetPerBatchGasCost(l1pricing.InitialPerBatchGasCostV6))
		state.Restrict(state.l1PricingState.SetEquilibrationUnits(l1pricing.InitialEquilibrationUnitsV6))
		state.Restrict(state.l2PricingState.SetSpeedLimitPerSecond(l2pricing.InitialSpeedLimitPerSecondV6))
		state.Restrict(state.l2PricingState.SetMaxPerBlockGasLimit(l2pricing.InitialPerBlockGasLimitV6))
	}

	state.Restrict(state.backingStorage.SetUint64ByUint64(uint64(versionOffset), state.mtosVersion))
}

func (state *MtosState) ScheduleMtOSUpgrade(newVersion uint64, timestamp uint64) error {
	err := state.upgradeVersion.Set(newVersion)
	if err != nil {
		return err
	}
	return state.upgradeTimestamp.Set(timestamp)
}

func (state *MtosState) BackingStorage() *storage.Storage {
	return state.backingStorage
}

func (state *MtosState) Restrict(err error) {
	state.Burner.Restrict(err)
}

func (state *MtosState) FormatVersion() uint64 {
	return state.mtosVersion
}

func (state *MtosState) SetFormatVersion(val uint64) {
	state.mtosVersion = val
	state.Restrict(state.backingStorage.SetUint64ByUint64(uint64(versionOffset), val))
}

func (state *MtosState) RetryableState() *retryables.RetryableState {
	return state.retryableState
}

func (state *MtosState) L1PricingState() *l1pricing.L1PricingState {
	return state.l1PricingState
}

func (state *MtosState) L2PricingState() *l2pricing.L2PricingState {
	return state.l2PricingState
}

func (state *MtosState) AddressTable() *addressTable.AddressTable {
	return state.addressTable
}

func (state *MtosState) ChainOwners() *addressSet.AddressSet {
	return state.chainOwners
}

func (state *MtosState) SendMerkleAccumulator() *merkleAccumulator.MerkleAccumulator {
	if state.sendMerkle == nil {
		state.sendMerkle = merkleAccumulator.OpenMerkleAccumulator(state.backingStorage.OpenSubStorage(sendMerkleSubspace))
	}
	return state.sendMerkle
}

func (state *MtosState) Blockhashes() *blockhash.Blockhashes {
	return state.blockhashes
}

func (state *MtosState) NetworkFeeAccount() (common.Address, error) {
	return state.networkFeeAccount.Get()
}

func (state *MtosState) SetNetworkFeeAccount(account common.Address) error {
	return state.networkFeeAccount.Set(account)
}

func (state *MtosState) InfraFeeAccount() (common.Address, error) {
	return state.infraFeeAccount.Get()
}

func (state *MtosState) SetInfraFeeAccount(account common.Address) error {
	return state.infraFeeAccount.Set(account)
}

func (state *MtosState) Keccak(data ...[]byte) ([]byte, error) {
	return state.backingStorage.Keccak(data...)
}

func (state *MtosState) KeccakHash(data ...[]byte) (common.Hash, error) {
	return state.backingStorage.KeccakHash(data...)
}

func (state *MtosState) ChainId() (*big.Int, error) {
	return state.chainId.Get()
}

func (state *MtosState) GenesisBlockNum() (uint64, error) {
	return state.genesisBlockNum.Get()
}
