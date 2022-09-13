# Copyright 2021-2022, Mantlenetwork, Inc.
# For license information, see https://github.com/nitro/blob/master/LICENSE

# Docker builds mess up file timestamps. Then again, in docker builds we never
# have to update an existing file. So - for docker, convert all dependencies
# to order-only dependencies (timestamps ignored).
# WARNING: when using this trick, you cannot use the $< automatic variable
ifeq ($(origin NITRO_BUILD_IGNORE_TIMESTAMPS),undefined)
 DEP_PREDICATE:=
 ORDER_ONLY_PREDICATE:=|
else
 DEP_PREDICATE:=|
 ORDER_ONLY_PREDICATE:=
endif


ifneq ($(origin NITRO_VERSION),undefined)
 GOLANG_LDFLAGS += -X github.com/mantlenetworkio/mantle/cmd/util.version=$(NITRO_VERSION)
endif

ifneq ($(origin NITRO_DATETIME),undefined)
 GOLANG_LDFLAGS += -X github.com/mantlenetworkio/mantle/cmd/util.datetime=$(NITRO_DATETIME)
endif

ifneq ($(origin NITRO_MODIFIED),undefined)
 GOLANG_LDFLAGS += -X github.com/mantlenetworkio/mantle/cmd/util.modified=$(NITRO_MODIFIED)
endif

ifneq ($(origin GOLANG_LDFLAGS),undefined)
 GOLANG_PARAMS = -ldflags="$(GOLANG_LDFLAGS)"
endif

precompile_names = AddressTable Aggregator BLS Debug FunctionTable GasInfo Info osTest Owner RetryableTx Statistics Sys
precompiles = $(patsubst %,./solgen/generated/%.go, $(precompile_names))

output_root=target

repo_dirs = mtos mtnode mtstate cmd precompiles solgen system_tests util validator wavmio
go_source = $(wildcard $(patsubst %,%/*.go, $(repo_dirs)) $(patsubst %,%/*/*.go, $(repo_dirs)))

color_pink = "\e[38;5;161;1m"
color_reset = "\e[0;0m"

done = "%bdone!%b\n" $(color_pink) $(color_reset)

replay_deps=mtos wavmio mtstate mtcompress solgen/go/node-interfacegen blsSignatures cmd/replay

replay_wasm=$(output_root)/machines/latest/replay.wasm

mtitrator_generated_header=$(output_root)/include/mtitrator.h
mtitrator_wasm_libs_nogo=$(output_root)/machines/latest/wasi_stub.wasm $(output_root)/machines/latest/host_io.wasm $(output_root)/machines/latest/soft-float.wasm
mtitrator_wasm_libs=$(mtitrator_wasm_libs_nogo) $(patsubst %,$(output_root)/machines/latest/%.wasm, go_stub brotli)
mtitrator_prover_lib=$(output_root)/lib/libprover.a
mtitrator_prover_bin=$(output_root)/bin/prover
mtitrator_jit=$(output_root)/bin/jit

mtitrator_cases=mtitrator/prover/test-cases

mtitrator_tests_wat=$(wildcard $(mtitrator_cases)/*.wat)
mtitrator_tests_rust=$(wildcard $(mtitrator_cases)/rust/src/bin/*.rs)

mtitrator_test_wasms=$(patsubst %.wat,%.wasm, $(mtitrator_tests_wat)) $(patsubst $(mtitrator_cases)/rust/src/bin/%.rs,$(mtitrator_cases)/rust/target/wasm32-wasi/release/%.wasm, $(mtitrator_tests_rust)) $(mtitrator_cases)/go/main

WASI_SYSROOT?=/opt/wasi-sdk/wasi-sysroot

mtitrator_wasm_lib_flags_nogo=$(patsubst %, -l %, $(mtitrator_wasm_libs_nogo))
mtitrator_wasm_lib_flags=$(patsubst %, -l %, $(mtitrator_wasm_libs))

mtitrator_wasm_wasistub_files = $(wildcard mtitrator/wasm-libraries/wasi-stub/src/*/*)
mtitrator_wasm_gostub_files = $(wildcard mtitrator/wasm-libraries/go-stub/src/*/*)
mtitrator_wasm_hostio_files = $(wildcard mtitrator/wasm-libraries/host-io/src/*/*)

# user targets

push: lint test-go .make/fmt
	@printf "%bdone building %s%b\n" $(color_pink) $$(expr $$(echo $? | wc -w) - 1) $(color_reset)
	@printf "%bready for push!%b\n" $(color_pink) $(color_reset)

all: build build-replay-env test-gen-proofs
	@touch .make/all

build: $(patsubst %,$(output_root)/bin/%, nitro deploy relay daserver datool seq-coordinator-invalidate)
	@printf $(done)

build-node-deps: $(go_source) build-prover-header build-prover-lib build-jit .make/solgen .make/cbrotli-lib

test-go-deps: \
	build-replay-env \
	$(patsubst %,$(mtitrator_cases)/%.wasm, global-state read-inboxmsg-10 global-state-wrapper const)

build-prover-header: $(mtitrator_generated_header)

build-prover-lib: $(mtitrator_prover_lib)

build-prover-bin: $(mtitrator_prover_bin)

build-jit: $(mtitrator_jit)

build-replay-env: $(mtitrator_prover_bin) $(mtitrator_wasm_libs) $(replay_wasm) $(output_root)/machines/latest/machine.wavm.br

build-wasm-libs: $(mtitrator_wasm_libs)

build-wasm-bin: $(replay_wasm)

build-solidity: .make/solidity

contracts: .make/solgen
	@printf $(done)

format fmt: .make/fmt
	@printf $(done)

lint: .make/lint
	@printf $(done)

test-go: .make/test-go
	@printf $(done)

test-go-challenge: test-go-deps
	go test -v -timeout 120m ./system_tests/... -run TestFullChallenge -tags fullchallengetest
	@printf $(done)

test-go-redis: test-go-deps
	go test -p 1 -run TestSeqCoordinator -tags redistest ./system_tests/... ./mtnode/...
	@printf $(done)

test-gen-proofs: \
	$(patsubst $(mtitrator_cases)/%.wat,contracts/test/prover/proofs/%.json, $(mtitrator_tests_wat)) \
	$(patsubst $(mtitrator_cases)/rust/src/bin/%.rs,contracts/test/prover/proofs/rust-%.json, $(mtitrator_tests_rust)) \
	contracts/test/prover/proofs/go.json

wasm-ci-build: $(mtitrator_wasm_libs) $(mtitrator_test_wasms)
	@printf $(done)

clean:
	go clean -testcache
	rm -rf $(mtitrator_cases)/rust/target
	rm -f $(mtitrator_cases)/*.wasm $(mtitrator_cases)/go/main
	rm -rf mtitrator/wasm-testsuite/tests
	rm -rf $(output_root)
	rm -f contracts/test/prover/proofs/*.json contracts/test/prover/spec-proofs/*.json
	rm -rf mtitrator/target
	rm -rf mtitrator/wasm-libraries/target
	rm -f mtitrator/wasm-libraries/soft-float/soft-float.wasm
	rm -f mtitrator/wasm-libraries/soft-float/*.o
	rm -f mtitrator/wasm-libraries/soft-float/SoftFloat/build/Wasm-Clang/*.o
	rm -f mtitrator/wasm-libraries/soft-float/SoftFloat/build/Wasm-Clang/*.a
	@rm -rf contracts/build contracts/cache solgen/go/
	@rm -f .make/*

docker:
	docker build -t nitro-node-slim --target nitro-node-slim .
	docker build -t nitro-node --target nitro-node .
	docker build -t nitro-node-dev --target nitro-node-dev .

# regular build rules

$(output_root)/bin/nitro: $(DEP_PREDICATE) build-node-deps
	go build $(GOLANG_PARAMS) -o $@ "$(CURDIR)/cmd/nitro"

$(output_root)/bin/deploy: $(DEP_PREDICATE) build-node-deps
	go build $(GOLANG_PARAMS) -o $@ "$(CURDIR)/cmd/deploy"

$(output_root)/bin/relay: $(DEP_PREDICATE) build-node-deps
	go build $(GOLANG_PARAMS) -o $@ "$(CURDIR)/cmd/relay"

$(output_root)/bin/daserver: $(DEP_PREDICATE) build-node-deps
	go build $(GOLANG_PARAMS) -o $@ "$(CURDIR)/cmd/daserver"

$(output_root)/bin/datool: $(DEP_PREDICATE) build-node-deps
	go build $(GOLANG_PARAMS) -o $@ "$(CURDIR)/cmd/datool"

$(output_root)/bin/seq-coordinator-invalidate: $(DEP_PREDICATE) build-node-deps
	go build $(GOLANG_PARAMS) -o $@ "$(CURDIR)/cmd/seq-coordinator-invalidate"

# recompile wasm, but don't change timestamp unless files differ
$(replay_wasm): $(DEP_PREDICATE) $(go_source) .make/solgen
	mkdir -p `dirname $(replay_wasm)`
	GOOS=js GOARCH=wasm go build -o $(output_root)/tmp/replay.wasm ./cmd/replay/...
	if ! diff -qN $(output_root)/tmp/replay.wasm $@ > /dev/null; then cp $(output_root)/tmp/replay.wasm $@; fi

$(mtitrator_prover_bin): $(DEP_PREDICATE) mtitrator/prover/src/*.rs mtitrator/prover/Cargo.toml
	mkdir -p `dirname $(mtitrator_prover_bin)`
	cargo build --manifest-path mtitrator/Cargo.toml --release --bin prover
	install mtitrator/target/release/prover $@

$(mtitrator_prover_lib): $(DEP_PREDICATE) mtitrator/prover/src/*.rs mtitrator/prover/Cargo.toml
	mkdir -p `dirname $(mtitrator_prover_lib)`
	cargo build --manifest-path mtitrator/Cargo.toml --release --lib -p prover
	install mtitrator/target/release/libprover.a $@

$(mtitrator_jit): $(DEP_PREDICATE) .make/cbrotli-lib mtitrator/jit/src/*.rs mtitrator/jit/*.rs mtitrator/jit/Cargo.toml
	mkdir -p `dirname $(mtitrator_jit)`
	cargo build --manifest-path mtitrator/Cargo.toml --release --bin jit
	install mtitrator/target/release/jit $@

$(mtitrator_cases)/rust/target/wasm32-wasi/release/%.wasm: $(mtitrator_cases)/rust/src/bin/%.rs $(mtitrator_cases)/rust/src/lib.rs
	cargo build --manifest-path $(mtitrator_cases)/rust/Cargo.toml --release --target wasm32-wasi --bin $(patsubst $(mtitrator_cases)/rust/target/wasm32-wasi/release/%.wasm,%, $@)

$(mtitrator_cases)/go/main: $(mtitrator_cases)/go/main.go $(mtitrator_cases)/go/go.mod $(mtitrator_cases)/go/go.sum
	cd $(mtitrator_cases)/go && GOOS=js GOARCH=wasm go build main.go

$(mtitrator_generated_header): $(DEP_PREDICATE) mtitrator/prover/src/lib.rs mtitrator/prover/src/utils.rs
	@echo creating ${PWD}/$(mtitrator_generated_header)
	mkdir -p `dirname $(mtitrator_generated_header)`
	cd mtitrator && cbindgen --config cbindgen.toml --crate prover --output ../$(mtitrator_generated_header)

$(output_root)/machines/latest/wasi_stub.wasm: $(DEP_PREDICATE) $(mtitrator_wasm_wasistub_files)
	mkdir -p $(output_root)/machines/latest
	cargo build --manifest-path mtitrator/wasm-libraries/Cargo.toml --release --target wasm32-unknown-unknown --package wasi-stub
	install mtitrator/wasm-libraries/target/wasm32-unknown-unknown/release/wasi_stub.wasm $@

mtitrator/wasm-libraries/soft-float/SoftFloat/build/Wasm-Clang/softfloat.a: $(DEP_PREDICATE) \
		mtitrator/wasm-libraries/soft-float/SoftFloat/build/Wasm-Clang/Makefile \
		mtitrator/wasm-libraries/soft-float/SoftFloat/build/Wasm-Clang/platform.h \
		mtitrator/wasm-libraries/soft-float/SoftFloat/source/*.c \
		mtitrator/wasm-libraries/soft-float/SoftFloat/source/include/*.h \
		mtitrator/wasm-libraries/soft-float/SoftFloat/source/8086/*.c \
		mtitrator/wasm-libraries/soft-float/SoftFloat/source/8086/*.h
	cd mtitrator/wasm-libraries/soft-float/SoftFloat/build/Wasm-Clang && make $(MAKEFLAGS)

mtitrator/wasm-libraries/soft-float/bindings32.o: $(DEP_PREDICATE) mtitrator/wasm-libraries/soft-float/bindings32.c
	clang mtitrator/wasm-libraries/soft-float/bindings32.c --sysroot $(WASI_SYSROOT) -I mtitrator/wasm-libraries/soft-float/SoftFloat/source/include -target wasm32-wasi -Wconversion -c -o $@

mtitrator/wasm-libraries/soft-float/bindings64.o: $(DEP_PREDICATE) mtitrator/wasm-libraries/soft-float/bindings64.c
	clang mtitrator/wasm-libraries/soft-float/bindings64.c --sysroot $(WASI_SYSROOT) -I mtitrator/wasm-libraries/soft-float/SoftFloat/source/include -target wasm32-wasi -Wconversion -c -o $@

$(output_root)/machines/latest/soft-float.wasm: $(DEP_PREDICATE) \
		mtitrator/wasm-libraries/soft-float/bindings32.o \
		mtitrator/wasm-libraries/soft-float/bindings64.o \
		mtitrator/wasm-libraries/soft-float/SoftFloat/build/Wasm-Clang/softfloat.a
	mkdir -p $(output_root)/machines/latest
	wasm-ld \
		mtitrator/wasm-libraries/soft-float/bindings32.o \
		mtitrator/wasm-libraries/soft-float/bindings64.o \
		mtitrator/wasm-libraries/soft-float/SoftFloat/build/Wasm-Clang/*.o \
		--no-entry -o $@ \
		$(patsubst %,--export wavm__f32_%, abs neg ceil floor trunc nearest sqrt add sub mul div min max) \
		$(patsubst %,--export wavm__f32_%, copysign eq ne lt le gt ge) \
		$(patsubst %,--export wavm__f64_%, abs neg ceil floor trunc nearest sqrt add sub mul div min max) \
		$(patsubst %,--export wavm__f64_%, copysign eq ne lt le gt ge) \
		$(patsubst %,--export wavm__i32_trunc_%,     f32_s f32_u f64_s f64_u) \
		$(patsubst %,--export wavm__i32_trunc_sat_%, f32_s f32_u f64_s f64_u) \
		$(patsubst %,--export wavm__i64_trunc_%,     f32_s f32_u f64_s f64_u) \
		$(patsubst %,--export wavm__i64_trunc_sat_%, f32_s f32_u f64_s f64_u) \
		$(patsubst %,--export wavm__f32_convert_%, i32_s i32_u i64_s i64_u) \
		$(patsubst %,--export wavm__f64_convert_%, i32_s i32_u i64_s i64_u) \
		--export wavm__f32_demote_f64 \
		--export wavm__f64_promote_f32

$(output_root)/machines/latest/go_stub.wasm: $(DEP_PREDICATE) $(wildcard mtitrator/wasm-libraries/go-stub/src/*/*)
	mkdir -p $(output_root)/machines/latest
	cargo build --manifest-path mtitrator/wasm-libraries/Cargo.toml --release --target wasm32-wasi --package go-stub
	install mtitrator/wasm-libraries/target/wasm32-wasi/release/go_stub.wasm $@

$(output_root)/machines/latest/host_io.wasm: $(DEP_PREDICATE) $(wildcard mtitrator/wasm-libraries/host-io/src/*/*)
	mkdir -p $(output_root)/machines/latest
	cargo build --manifest-path mtitrator/wasm-libraries/Cargo.toml --release --target wasm32-wasi --package host-io
	install mtitrator/wasm-libraries/target/wasm32-wasi/release/host_io.wasm $@

$(output_root)/machines/latest/brotli.wasm: $(DEP_PREDICATE) $(wildcard mtitrator/wasm-libraries/brotli/src/*/*) .make/cbrotli-wasm
	mkdir -p $(output_root)/machines/latest
	cargo build --manifest-path mtitrator/wasm-libraries/Cargo.toml --release --target wasm32-wasi --package brotli
	install mtitrator/wasm-libraries/target/wasm32-wasi/release/brotli.wasm $@

$(output_root)/machines/latest/machine.wavm.br: $(DEP_PREDICATE) $(mtitrator_prover_bin) $(mtitrator_wasm_libs) $(replay_wasm)
	$(mtitrator_prover_bin) $(replay_wasm) --generate-binaries $(output_root)/machines/latest -l $(output_root)/machines/latest/soft-float.wasm -l $(output_root)/machines/latest/wasi_stub.wasm -l $(output_root)/machines/latest/go_stub.wasm -l $(output_root)/machines/latest/host_io.wasm -l $(output_root)/machines/latest/brotli.wasm

$(mtitrator_cases)/%.wasm: $(mtitrator_cases)/%.wat
	wat2wasm $< -o $@

contracts/test/prover/proofs/float%.json: $(mtitrator_cases)/float%.wasm $(mtitrator_prover_bin) $(output_root)/machines/latest/soft-float.wasm
	$(mtitrator_prover_bin) $< -l $(output_root)/machines/latest/soft-float.wasm -o $@ -b --allow-hostapi --require-success --always-merkleize

contracts/test/prover/proofs/no-stack-pollution.json: $(mtitrator_cases)/no-stack-pollution.wasm $(mtitrator_prover_bin)
	$(mtitrator_prover_bin) $< -o $@ --allow-hostapi --require-success --always-merkleize

contracts/test/prover/proofs/rust-%.json: $(mtitrator_cases)/rust/target/wasm32-wasi/release/%.wasm $(mtitrator_prover_bin) $(mtitrator_wasm_libs_nogo)
	$(mtitrator_prover_bin) $< $(mtitrator_wasm_lib_flags_nogo) -o $@ -b --allow-hostapi --require-success --inbox-add-stub-headers --inbox $(mtitrator_cases)/rust/data/msg0.bin --inbox $(mtitrator_cases)/rust/data/msg1.bin --delayed-inbox $(mtitrator_cases)/rust/data/msg0.bin --delayed-inbox $(mtitrator_cases)/rust/data/msg1.bin --preimages $(mtitrator_cases)/rust/data/preimages.bin

contracts/test/prover/proofs/go.json: $(mtitrator_cases)/go/main $(mtitrator_prover_bin) $(mtitrator_wasm_libs)
	$(mtitrator_prover_bin) $< $(mtitrator_wasm_lib_flags) -o $@ -i 5000000 --require-success

# avoid testing read-inboxmsg-10 in onestepproofs. It's used for go challenge testing.
contracts/test/prover/proofs/read-inboxmsg-10.json:
	echo "[]" > $@

contracts/test/prover/proofs/%.json: $(mtitrator_cases)/%.wasm $(mtitrator_prover_bin)
	$(mtitrator_prover_bin) $< -o $@ --allow-hostapi --always-merkleize

# strategic rules to minimize dependency building

.make/lint: $(DEP_PREDICATE) build-node-deps $(ORDER_ONLY_PREDICATE) .make
	golangci-lint run --fix
	yarn --cwd contracts solhint
	@touch $@

.make/fmt: $(DEP_PREDICATE) build-node-deps .make/yarndeps $(ORDER_ONLY_PREDICATE) .make
	golangci-lint run --disable-all -E gofmt --fix
	cargo fmt --all --manifest-path mtitrator/Cargo.toml -- --check
	cargo fmt --all --manifest-path mtitrator/wasm-testsuite/Cargo.toml -- --check
	yarn --cwd contracts prettier:solidity
	@touch $@

.make/test-go: $(DEP_PREDICATE) $(go_source) build-node-deps test-go-deps $(ORDER_ONLY_PREDICATE) .make
	gotestsum --format short-verbose
	@touch $@

.make/solgen: $(DEP_PREDICATE) solgen/gen.go .make/solidity $(ORDER_ONLY_PREDICATE) .make
	mkdir -p solgen/go/
	go run solgen/gen.go
	@touch $@

.make/solidity: $(DEP_PREDICATE) contracts/src/*/*.sol .make/yarndeps $(ORDER_ONLY_PREDICATE) .make
	yarn --cwd contracts build
	@touch $@

.make/yarndeps: $(DEP_PREDICATE) contracts/package.json contracts/yarn.lock $(ORDER_ONLY_PREDICATE) .make
	yarn --cwd contracts install
	@touch $@

.make/cbrotli-lib: $(DEP_PREDICATE) $(ORDER_ONLY_PREDICATE) .make
	@printf "%btesting cbrotli local build exists. If this step fails, run ./build-brotli.sh -l%b\n" $(color_pink) $(color_reset)
	test -f target/include/brotli/encode.h
	test -f target/include/brotli/decode.h
	test -f target/lib/libbrotlicommon-static.a
	test -f target/lib/libbrotlienc-static.a
	test -f target/lib/libbrotlidec-static.a
	@touch $@

.make/cbrotli-wasm: $(DEP_PREDICATE) $(ORDER_ONLY_PREDICATE) .make
	@printf "%btesting cbrotli wasm build exists. If this step fails, run ./build-brotli.sh -w%b\n" $(color_pink) $(color_reset)
	test -f target/lib-wasm/libbrotlicommon-static.a
	test -f target/lib-wasm/libbrotlienc-static.a
	test -f target/lib-wasm/libbrotlidec-static.a
	@touch $@

.make:
	mkdir .make


# Makefile settings

always:              # use this to force other rules to always build
.DELETE_ON_ERROR:    # causes a failure to delete its target
.PHONY: push all build build-node-deps test-go-deps build-prover-header build-prover-lib build-prover-bin build-jit build-replay-env build-solidity build-wasm-libs contracts format fmt lint test-go test-gen-proofs push clean docker
