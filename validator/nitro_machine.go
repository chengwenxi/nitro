// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package validator

/*
#cgo CFLAGS: -g -Wall -I../target/include/
#include "mtitrator.h"
#include <stdlib.h>
*/
import "C"
import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type MantleMachineConfig struct {
	RootPath             string // a folder with various machines in it
	WavmBinaryPath       string
	UntilHostIoStatePath string
	ProverBinPath        string

	// Used for jit validator only
	JitCranelift bool

	// Used for debugging only
	LibraryPaths []string
}

var DefaultMantleMachineConfig = MantleMachineConfig{
	RootPath:             "./target/machines/",
	WavmBinaryPath:       "machine.wavm.br",
	UntilHostIoStatePath: "until-host-io-state.bin",

	JitCranelift:  DefaultBlockValidatorConfig.JitValidatorCranelift,
	ProverBinPath: "replay.wasm",
	LibraryPaths:  []string{"soft-float.wasm", "wasi_stub.wasm", "go_stub.wasm", "host_io.wasm", "brotli.wasm"},
}

func init() {
	_, thisfile, _, _ := runtime.Caller(0)
	projectDir := filepath.Dir(filepath.Dir(thisfile))
	DefaultMantleMachineConfig.RootPath = filepath.Join(projectDir, "target", "machines")
}

func (c MantleMachineConfig) getMachinePath(moduleRoot common.Hash) string {
	if moduleRoot == (common.Hash{}) {
		return filepath.Join(c.RootPath, "latest")
	} else {
		return filepath.Join(c.RootPath, moduleRoot.String())
	}
}

func (c MantleMachineConfig) ReadLatestWasmModuleRoot() (common.Hash, error) {
	fileToRead := filepath.Join(c.getMachinePath(common.Hash{}), "module-root.txt")
	fileBytes, err := os.ReadFile(fileToRead)
	if err != nil {
		return common.Hash{}, err
	}
	s := strings.TrimSpace(string(fileBytes))
	return common.HexToHash(s), nil
}

type loaderMachineStatus struct {
	machine    *MtitratorMachine
	jitMachine *JitMachine
	chanSignal chan struct{}
	jit        bool
	err        error
}

func (s *loaderMachineStatus) signalReady() {
	close(s.chanSignal)
}

func (s *loaderMachineStatus) createZeroStepMachineInternal(config MantleMachineConfig, moduleRoot common.Hash, realModuleRoot common.Hash) {
	defer s.signalReady()
	binPath := filepath.Join(config.getMachinePath(moduleRoot), config.WavmBinaryPath)
	cBinPath := C.CString(binPath)
	defer C.free(unsafe.Pointer(cBinPath))
	log.Info("creating mantle machine", "binpath", binPath)
	baseMachine := C.mtitrator_load_wavm_binary(cBinPath)
	if baseMachine == nil {
		s.err = errors.New("failed to load base machine")
		return
	}
	mantleMachine := machineFromPointer(baseMachine)
	machineModuleRoot := mantleMachine.GetModuleRoot()
	if machineModuleRoot != realModuleRoot {
		s.err = fmt.Errorf("attempting to load module root %v got machine with module root %v", realModuleRoot, machineModuleRoot)
		return
	}
	s.machine = mantleMachine
	s.machine.Freeze()
}

// We try to store/load state before first host_io to a file.
// We will chicken out of that if something fails, but still try to calculate the machine
func (s *loaderMachineStatus) createHostIoMachineInternal(config MantleMachineConfig, moduleRoot common.Hash, zerostep *MtitratorMachine) {
	defer s.signalReady()
	ctx := context.Background()
	machine := zerostep.Clone()

	statePath := filepath.Join(config.getMachinePath(moduleRoot), config.UntilHostIoStatePath)
	_, err := os.Stat(statePath)
	if err == nil {
		log.Info("found cached machine until host io state", "moduleRoot", moduleRoot)

		err := machine.DeserializeAndReplaceState(statePath)
		if err != nil {
			// Safe as if DeserializeAndReplaceState returns an error it will not have mutated the machine
			log.Warn("failed to load machine until host io state; will reexecute", "err", err)
		} else {
			s.machine = machine
			s.machine.Freeze()
			return
		}
	} else if errors.Is(err, os.ErrNotExist) {
		log.Info("didn't find cached machine until host io state", "path", statePath)
	} else {
		log.Warn("error checking if machine until host io state is cached", "path", statePath, "err", err)
	}

	s.err = machine.StepUntilHostIo(ctx)
	if s.err != nil {
		return
	}

	if machine.IsErrored() {
		s.err = errors.New("machine entered errored state while caching execution up to host io")
		return
	}

	s.machine = machine
	s.machine.Freeze()
}

type mantleMachineRequest struct {
	moduleRoot  common.Hash
	untilHostIo bool
	jit         bool
}

type MantleMachineLoader struct {
	config       MantleMachineConfig
	machinesLock sync.Mutex
	machines     map[mantleMachineRequest]*loaderMachineStatus
	fatalErrChan chan error
}

func NewMantleMachineLoader(config MantleMachineConfig, fatalErrChan chan error) *MantleMachineLoader {
	return &MantleMachineLoader{
		config:       config,
		machines:     make(map[mantleMachineRequest]*loaderMachineStatus),
		fatalErrChan: fatalErrChan,
	}
}

func (s *loaderMachineStatus) waitForMachine(ctx context.Context) (*MtitratorMachine, *JitMachine, error) {
	select {
	case <-s.chanSignal:
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
	if s.err != nil {
		return nil, nil, s.err
	}
	if !s.jit && s.machine == nil {
		return nil, nil, errors.New("mantle machine is nil")
	}
	if s.jit && s.jitMachine == nil {
		return nil, nil, errors.New("jit machine is nil")
	}
	return s.machine, s.jitMachine, nil
}

func (l *MantleMachineLoader) createMachineImpl(
	moduleRoot common.Hash, untilHostIo, jit bool,
) (*loaderMachineStatus, error) {
	machineRequest := mantleMachineRequest{
		moduleRoot:  moduleRoot,
		untilHostIo: untilHostIo,
		jit:         jit,
	}

	config := l.config

	// Fast path: check if we already have the machine
	l.machinesLock.Lock()
	machine, ok := l.machines[machineRequest]
	if ok {
		l.machinesLock.Unlock()
		return machine, nil
	}
	l.machinesLock.Unlock()

	// Attempt to resolve any alias to the module root (due to the latest machine being separate).
	realModuleRoot := moduleRoot
	if moduleRoot == (common.Hash{}) {
		var err error
		realModuleRoot, err = config.ReadLatestWasmModuleRoot()
		if err != nil {
			return nil, err
		}
	} else {
		_, err := os.Stat(filepath.Join(config.getMachinePath(moduleRoot), config.WavmBinaryPath))
		if errors.Is(err, os.ErrNotExist) {
			// Attempt to load the latest module root instead (maybe it's what we're looking for).
			originalErr := err
			realModuleRoot, err = config.ReadLatestWasmModuleRoot()
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					// Be nice and return the original error, as it's clarifies what went wrong.
					return nil, originalErr
				} else {
					return nil, err
				}
			}
			if realModuleRoot == moduleRoot {
				// The latest machine is the requested one! Pretend we're loading the latest machine instead.
				moduleRoot = common.Hash{}
				machineRequest.moduleRoot = common.Hash{}
			} else {
				// The latest machine is different, so return the original error loading this machine.
				return nil, originalErr
			}
		} else if err != nil {
			return nil, err
		}
	}

	l.machinesLock.Lock()
	defer l.machinesLock.Unlock()

	realMachineRequest := mantleMachineRequest{
		moduleRoot:  realModuleRoot,
		untilHostIo: untilHostIo,
		jit:         jit,
	}
	machine, ok = l.machines[machineRequest]
	if !ok && moduleRoot != realModuleRoot {
		machine, ok = l.machines[realMachineRequest]
	}

	if !ok {
		machine = &loaderMachineStatus{
			chanSignal: make(chan struct{}),
			jit:        jit,
		}
		l.machines[machineRequest] = machine
		if moduleRoot != realModuleRoot {
			l.machines[realMachineRequest] = machine
		}

		go func() {
			if jit {
				machine.jitMachine, machine.err = createJitMachine(config, moduleRoot, l.fatalErrChan)
				machine.signalReady()
				return
			}
			if untilHostIo {
				zeroStep, err := l.GetMachine(context.Background(), moduleRoot, false)
				if err != nil {
					machine.err = err
					machine.signalReady()
				} else {
					machine.createHostIoMachineInternal(config, moduleRoot, zeroStep)
				}
			} else {
				machine.createZeroStepMachineInternal(config, moduleRoot, realModuleRoot)
			}
		}()
	}

	return machine, nil
}

// Starts work on creating the machine in a separate goroutine
// Returns immediately. Can be called multiple times.
func (l *MantleMachineLoader) CreateMachine(moduleRoot common.Hash, untilHostIo, jit bool) error {
	_, err := l.createMachineImpl(moduleRoot, untilHostIo, jit)
	return err
}

// Gets machine when one is ready
// Returns with proper error if context aborts
func (l *MantleMachineLoader) GetMachine(
	ctx context.Context, moduleRoot common.Hash, untilHostIo bool,
) (*MtitratorMachine, error) {
	loader, err := l.createMachineImpl(moduleRoot, untilHostIo, false)
	if err != nil {
		return nil, err
	}
	machine, _, err := loader.waitForMachine(ctx)
	return machine, err
}

func (l *MantleMachineLoader) GetJitMachine(
	ctx context.Context, moduleRoot common.Hash, untilHostIo bool,
) (*JitMachine, error) {
	loader, err := l.createMachineImpl(moduleRoot, untilHostIo, true)
	if err != nil {
		return nil, err
	}
	_, machine, err := loader.waitForMachine(ctx)
	return machine, err

}

func (l *MantleMachineLoader) GetConfig() MantleMachineConfig {
	return l.config
}
