# Overview
MtOS provides L2-specific precompiles with methods smart contracts can call the same way they can solidity functions. This reference exhaustively documents the specific calls MtOS makes available. For more details on the infrastructure that makes this possible, please refer to the [MtOS documentation](mtos.md). For an abbreviated reference on the precompiles we expect users to most often use, please see the [common precompiles documentation](common-precompiles.md).

From the perspective of user applications, precompiles live as contracts at the following addresses. Click on any to jump to their section.

| Precompile                                     | Address &nbsp; | Purpose                                            |
|:-----------------------------------------------|:---------------|:---------------------------------------------------|
| [`MtAddressTable`](#MtAddressTable)          | `0x66`         | Supporting compression of addresses                |
| [`MtAggregator`](#MtAggregator)              | `0x6d`         | Configuring transaction aggregation                |
| [`MtBLS`](#MtBLS)                            | `0x67`         | Managing BLS keys                                  |
| [`MtDebug`](#MtDebug)                        | `0xff`         | Testing tools                                      |
| [`MtFunctionTable`](#MtFunctionTable) &nbsp; | `0x68`         | No longer used                                     |
| [`MtGasInfo`](#MtGasInfo)                    | `0x6c`         | Info about gas pricing                             |
| [`MtInfo`](#MtInfo)                          | `0x65`         | Info about accounts                                |
| [`MtOwner`](#MtOwner)                        | `0x70`         | Chain administration, callable only by chain owner |
| [`MtOwnerPublic`](#MtOwnerPublic)            | `0x6b`         | Info about chain owners                            |
| [`MtosTest`](#MtosTest)                      | `0x69`         | No longer used                                     |
| [`MtRetryableTx`](#MtRetryableTx)            | `0x6e`         | Managing retryables                                |
| [`MtStatistics`](#MtStatistics)              | `0x6f`         | Info about the pre-Mantle state                     |
| [`MtSys`](#MtSys)                            | `0x64`         | System-level functionality                         |

[MtAddressTable_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtAddressTable.go
[MtAggregator_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtAddressTable.go
[MtBLS_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtBLS.go
[MtDebug_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtDebug.go
[MtFunctionTable_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtFunctionTable.go
[MtInfo_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtInfo.go
[MtGasInfo_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtGasInfo.go
[MtosTest_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtosTest.go
[MtOwner_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtOwner.go
[MtOwnerPublic_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtOwnerPublic.go
[MtRetryableTx_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtRetryableTx.go
[MtStatistics_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtStatistics.go
[MtSys_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtSys.go

# [MtAddressTable][MtAddressTable_link]
Provides the ability to create short-hands for commonly used accounts.

| Methods                                                         |                                                                                           |
|:----------------------------------------------------------------|:------------------------------------------------------------------------------------------|
| [![](e.png)][ATs0] [`AddressExists`][AT0]`(address)`            | Checks if an address exists in the table                                                  |
| [![](e.png)][ATs1] [`Compress`][AT1]`(address)`                 | Gets bytes that represent the address                                                     |
| [![](e.png)][ATs2] [`Decompress`][AT2]`(buffer, offset)` &nbsp; | Replaces the compressed bytes at the given offset with those of the corresponding account |
| [![](e.png)][ATs3] [`Lookup`][AT3]`(address)`                   | Gets the index of an address in the table                                                 |
| [![](e.png)][ATs4] [`LookupIndex`][AT4]`(index)`                | Gets the address at an index in the table                                                 |
| [![](e.png)][ATs5] [`Register`][AT5]`(address)`                 | Adds an address to the table, shrinking its compressed representation                     |
| [![](e.png)][ATs6] [`Size`][AT6]`()`                            | Gets the number of addresses in the table                                                 |

[AT0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtAddressTable.go#L18
[AT1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtAddressTable.go#L23
[AT2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtAddressTable.go#L28
[AT3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtAddressTable.go#L41
[AT4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtAddressTable.go#L53
[AT5]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtAddressTable.go#L68
[AT6]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtAddressTable.go#L74

[ATs0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtAddressTable.sol#L31
[ATs1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtAddressTable.sol#L38
[ATs2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtAddressTable.sol#L46
[ATs3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtAddressTable.sol#L55
[ATs4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtAddressTable.sol#L61
[ATs5]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtAddressTable.sol#L68
[ATs6]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtAddressTable.sol#L73


# [MtAggregator][MtAggregator_link]
Provides aggregators and their users methods for configuring how they participate in L1 aggregation. Mantle One's default aggregator is the Sequencer, which a user will prefer unless `SetPreferredAggregator` is invoked to change it.

Compression ratios are measured in basis points. Methods that are checkmarked are access-controlled and will revert if not called by the aggregator, its fee collector, or a chain owner.

| Methods                                                                   |                                                         |    |
|:--------------------------------------------------------------------------|:--------------------------------------------------------|:---|
| [![](e.png)][As0] [`GetPreferredAggregator`][A0]`(account)`               | Gets an account's preferred aggregator                  |    |
| [![](e.png)][As1] [`SetPreferredAggregator`][A1]`(aggregator)`            | Sets the caller's preferred aggregator to that provided |    |
| [![](e.png)][As2] [`GetDefaultAggregator`][A2]`()`                        | Gets the chain's default aggregator                     |    |
| [![](e.png)][As3] [`SetDefaultAggregator`][A3]`(default)`                 | Sets the chain's default aggregator                     | ✔️  |
| [![](e.png)][As4] [`GetCompressionRatio`][A4]`(aggregator)`               | Gets the aggregator's compression ratio                 |    |
| [![](e.png)][As5] [`SetCompressionRatio`][A5]`(aggregator, ratio)`        | Set the aggregator's compression ratio                  | ✔️  |
| [![](e.png)][As6] [`GetFeeCollector`][A6]`(aggregator)`                   | Gets an aggregator's fee collector                      |    |
| [![](e.png)][As7] [`SetFeeCollector`][A7]`(aggregator, collector)` &nbsp; | Sets an aggregator's fee collector                      | ✔️  |

[A0]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/precompiles/MtAggregator.go#L25
[A1]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/precompiles/MtAggregator.go#L42
[A2]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/precompiles/MtAggregator.go#L51
[A3]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/precompiles/MtAggregator.go#L56
[A4]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/precompiles/MtAggregator.go#L73
[A5]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/precompiles/MtAggregator.go#L79
[A6]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/precompiles/MtAggregator.go#L91
[A7]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/precompiles/MtAggregator.go#L96

[As0]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/solgen/src/precompiles/MtAggregator.sol#L28
[As1]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/solgen/src/precompiles/MtAggregator.sol#L32
[As2]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/solgen/src/precompiles/MtAggregator.sol#L35
[As3]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/solgen/src/precompiles/MtAggregator.sol#L40
[As4]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/solgen/src/precompiles/MtAggregator.sol#L45
[As5]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/solgen/src/precompiles/MtAggregator.sol#L51
[As6]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/solgen/src/precompiles/MtAggregator.sol#L56
[As7]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/solgen/src/precompiles/MtAggregator.sol#L62

| Deprecated Methods                                          |              |
|:------------------------------------------------------------|:-------------|
| [![](e.png)][Ads0] [`GetTxBaseFee`][Ad0]`(aggregator)`      | Returns 0    |
| [![](e.png)][Ads1] [`SetTxBaseFee`][Ad1]`(aggregator, fee)` | Does nothing |

[Ad0]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/precompiles/MtAggregator.go#L108
[Ad1]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/precompiles/MtAggregator.go#L114

[Ads0]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/solgen/src/precompiles/MtAggregator.sol#L67
[Ads1]: https://github.com/mantlenetworkio/mantle/blob/ba3a86afb2e7057bdc3cce54b28be4c1c0579180/solgen/src/precompiles/MtAggregator.sol#L75

# [MtBLS][MtBLS_link]
Provides a registry of BLS public keys for accounts.

| Methods                                                             |                                                             |
|:--------------------------------------------------------------------|:------------------------------------------------------------|
| [![](e.png)][Bs0] [`RegisterAltBN128`][B0]`(x0, x1, y0, y1)` &nbsp; | Associate an AltBN128 public key with the caller's address  |
| [![](e.png)][Bs1] [`GetAltBN128`][B1]`(account)`                    | Gets the AltBN128 public key associated with an address     |
| [![](e.png)][Bs2] [`RegisterBLS12381`][B2]`(key)`                   | Associate a BLS 12-381 public key with the caller's address |
| [![](e.png)][Bs3] [`GetBLS12381`][B3]`(account)`                    | Gets the BLS 12-381 public key associated with an address   |

[B0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtBLS.go#L27
[B1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtBLS.go#L32
[B2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtBLS.go#L37
[B3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtBLS.go#L46

[Bs0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtBLS.sol#L44
[Bs1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtBLS.sol#L52
[Bs2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtBLS.sol#L63
[Bs3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtBLS.sol#L66

| Deprecated Methods                                            |                                |
|:--------------------------------------------------------------|:-------------------------------|
| [![](e.png)][Bds0] [`Register`][Bd0]`(x0, x1, y0, y1)` &nbsp; | equivalent to registerAltBN128 |
| [![](e.png)][Bds1] [`GetPublicKey`][Bd1]`(account)`           | equivalent to getAltBN128      |

[Bd0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtBLS.go#L17
[Bd1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtBLS.go#L22

[Bds0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtBLS.sol#L25
[Bds1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtBLS.sol#L33


# [MtDebug][MtDebug_link]
Provides mechanisms useful for testing. The methods of `MtDebug` are only available for chains with the `AllowDebugPrecompiles` chain parameter set. Otherwise, calls to this precompile will revert.

| Methods                                                |                                                    |
|:-------------------------------------------------------|:---------------------------------------------------|
| [![](e.png)][Ds0] [`BecomeChainOwner`][D0]`()`         | Caller becomes a chain owner                       |
| [![](e.png)][Ds1] [`Events`][D1]`(flag, value)` &nbsp; | Emit events with values based on the args provided |

[D0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtDebug.go#L38
[D1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtDebug.go#L19

[Ds0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtDebug.sol#L27
[Ds1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtDebug.sol#L30


| Events                                   |                                            |
|:-----------------------------------------|:-------------------------------------------|
| [![](e.png)][Des0] [`Basic`][De0] &nbsp; | Emitted in `Events` for testing            |
| [![](e.png)][Des1] [`Mixed`][De1]        | Emitted in `Events` for testing            |
| [![](e.png)][Des2] [`Store`][De2]        | Never emitted (used for testing log sizes) |

[De0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtDebug.go#L24
[De1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtDebug.go#L29
[De2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtDebug.go#L13

[Des0]:https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtDebug.sol#L33
[Des1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtDebug.sol#L34
[Des2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtDebug.sol#L41


# [MtFunctionTable][MtFunctionTable_link]
Provided aggregator's the ability to manage function tables, to enable one form of transaction compression. The Mantle aggregator implementation does not use these, so these methods have been stubbed and their effects disabled. They are kept for backwards compatibility.

| Methods                                                  |                                            |
|:---------------------------------------------------------|:-------------------------------------------|
| [![](e.png)][FTs0] [`Get`][FT0]`(address, index)` &nbsp; | Reverts since the table is empty           |
| [![](e.png)][FTs1] [`Size`][FT1]`(address)`              | Returns the empty table's size, which is 0 |
| [![](e.png)][FTs2] [`Upload`][FT2]`(bytes)`              | Does nothing                               |

[FT0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtFunctionTable.go#L30
[FT1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtFunctionTable.go#L25
[FT2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtFunctionTable.go#L20

[FTs0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtFunctionTable.sol#L35
[FTs1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtFunctionTable.sol#L32
[FTs2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtFunctionTable.sol#L29


# [MtGasInfo][MtGasInfo_link]
Provides insight into the cost of using the chain. These methods have been adjusted to account for Mantle's heavy use of calldata compression. Of note to end-users, we no longer make a distinction between non-zero and zero-valued calldata bytes.

| Methods                                                                   |                                                                                                  |
|:--------------------------------------------------------------------------|:-------------------------------------------------------------------------------------------------|
| [![](e.png)][GIs0] [`GetPricesInWeiWithAggregator`][GI0]`(aggregator)`    | Get prices in wei when using the provided aggregator                                             |
| [![](e.png)][GIs1] [`GetPricesInWei`][GI1]`()`                            | Get prices in wei when using the caller's preferred aggregator                                   |
| [![](e.png)][GIs2] [`GetPricesInArbGasWithAggregator`][GI2]`(aggregator)` | Get prices in ArbGas when using the provided aggregator                                          |
| [![](e.png)][GIs3] [`GetPricesInArbGas`][GI3]`()`                         | Get prices in ArbGas when using the caller's preferred aggregator                                |
| [![](e.png)][GIs4] [`GetGasAccountingParams`][GI4]`()`                    | Get the chain speed limit, pool size, and tx gas limit                                           |
| [![](e.png)][GIs5] [`GetMinimumGasPrice`][GI5]`()`                        | Get the minimum gas price needed for a transaction to succeed                                    |
| [![](e.png)][GIs6] [`GetGasPoolSeconds`][GI6]`()`                         | Get the number of seconds worth of the speed limit the gas pool contains                         |
| [![](e.png)][GIs7] [`GetGasPoolTarget`][GI7]`()`                          | Get the target fullness in bips the pricing model will try to keep the pool at                   |
| [![](e.png)][GIs8] [`GetGasPoolWeight`][GI8]`()`                          | Get the extent in bips to which the pricing model favors filling the pool over increasing speeds |
| [![](e.png)][GIs9] [`GetRateEstimate`][GI9]`()`                           | Get MtOS's estimate of the amount of gas being burnt per second                                 |
| [![](e.png)][GIs10] [`GetRateEstimateInertia`][GI10]`()`                  | Get how slowly MtOS updates its estimate the amount of gas being burnt per second               |
| [![](e.png)][GIs11] [`GetL1BaseFeeEstimate`][GI11]`()`                    | Get MtOS's estimate of the L1 basefee in wei                                                    |
| [![](e.png)][GIs12] [`GetL1BaseFeeEstimateInertia`][GI12]`()`             | Get how slowly MtOS updates its estimate of the L1 basefee                                      |
| [![](e.png)][GIs13] [`GetL1GasPriceEstimate`][GI13]`()`                   | Deprecated -- Same as getL1BaseFeeEstimate()                                                     |
| [![](e.png)][GIs14] [`GetCurrentTxL1GasFees`][GI14]`()`                   | Get L1 gas fees paid by the current transaction                                                  |


[GI0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L27
[GI1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L63
[GI2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L75
[GI3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L99
[GI4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L111
[GI5]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L120
[GI6]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L125
[GI7]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L130
[GI8]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L135
[GI9]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L140
[GI10]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L145
[GI11]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L150
[GI12]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L155
[GI13]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L160
[GI14]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L165
[GI15]: https://github.com/mantlenetworkio/mantle/blob/6b6de9068662883518cff13c67b885161763f52c/precompiles/MtGasInfo.go#L170
[GI16]: https://github.com/mantlenetworkio/mantle/blob/6b6de9068662883518cff13c67b885161763f52c/precompiles/MtGasInfo.go#L175
[GI17]: https://github.com/mantlenetworkio/mantle/blob/6b6de9068662883518cff13c67b885161763f52c/precompiles/MtGasInfo.go#L180

[GIs0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L36
[GIs1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L58
[GIs2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L72
[GIs3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L83
[GIs4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L94
[GIs5]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L104
[GIs6]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L107
[GIs7]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L110
[GIs8]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L113
[GIs9]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L116
[GIs10]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L119
[GIs11]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L122
[GIs12]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L125
[GIs13]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L128
[GIs14]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L131
[GIs15]: https://github.com/mantlenetworkio/mantle/blob/6b6de9068662883518cff13c67b885161763f52c/contracts/src/precompiles/MtGasInfo.sol#L123
[GIs16]: https://github.com/mantlenetworkio/mantle/blob/6b6de9068662883518cff13c67b885161763f52c/contracts/src/precompiles/MtGasInfo.sol#L126
[GIs17]: https://github.com/mantlenetworkio/mantle/blob/6b6de9068662883518cff13c67b885161763f52c/contracts/src/precompiles/MtGasInfo.sol#L129


# [MtInfo][MtInfo_link]
Provides the ability to lookup basic info about accounts and contracts.

| Methods                                                |                                      |
|:-------------------------------------------------------|:-------------------------------------|
| [![](e.png)][Is0] [`GetBalance`][I0]`(account)` &nbsp; | Retrieves an account's balance       |
| [![](e.png)][Is1] [`GetCode`][I1]`(account)`           | Retrieves a contract's deployed code |

[I0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtInfo.go#L18
[I1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtInfo.go#L26

[Is0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtInfo.sol#L25
[Is1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtInfo.sol#L28


# [MtosTest][MtosTest_link]
Provides a method of burning arbitrary amounts of gas, which exists for historical reasons. In Classic, `MtosTest` had additional methods only the zero address could call. These have been removed since users don't use them and calls to missing methods revert.

| Methods                                               |                                                     | Mantle changes |
|:------------------------------------------------------|:----------------------------------------------------|---------------|
| [![](e.png)][Ts0] [`BurnArbGas`][T0]`(amount)` &nbsp; | unproductively burns the amount of L2 ArbGas &nbsp; | Now pure      |

[T0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtosTest.go#L17

[Ts0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtosTest.sol#L27


# [MtOwner][MtOwner_link]
Provides owners with tools for managing the rollup. Calls by non-owners will always revert.

Most of Mantle Classic's owner methods have been removed since they no longer make sense in Mantle:

- What were once chain parameters are now parts of MtOS's state, and those that remain are set at genesis. 
- MtOS upgrades happen with the rest of the system rather than being independent
- Exemptions to address aliasing are no longer offered. Exemptions were intended to support backward compatibility for contracts deployed before aliasing was introduced, but no exemptions were ever requested.

| Methods                                                          |                                                                                                  |
|:-----------------------------------------------------------------|:-------------------------------------------------------------------------------------------------|
| [![](e.png)][Os0] [`AddChainOwner`][O0]`(account)`               | Add account as a chain owner                                                                     |
| [![](e.png)][Os1] [`RemoveChainOwner`][O1]`(account)`            | Remove account from the list of chain owners                                                     |
| [![](e.png)][Os2] [`IsChainOwner`][O2]`(account)`                | See if account is a chain owner                                                                  |
| [![](e.png)][Os3] [`GetAllChainOwners`][O3]`()`                  | Retrieves the list of chain owners                                                               |
| [![](e.png)][Os4] [`SetL1BaseFeeEstimate`][O4]`(price)`          | Set the L1 basefee estimate directly, bypassing the autoregression                               |
| [![](e.png)][Os5] [`SetL1BaseFeeEstimateInertia`][O5]`(inertia)` | Set how slowly MtOS updates its estimate of the L1 basefee                                      |
| [![](e.png)][Os6] [`SetL2GasPrice`][O6]`(price)`                 | Set the L2 gas price directly, bypassing the pool calculus                                       |
| [![](e.png)][Os7] [`SetMinimumGasPrice`][O7]`(price)`            | Set the minimum gas price needed for a transaction to succeed                                    |
| [![](e.png)][Os8] [`SetSpeedLimit`][O8]`(limit)`                 | Set the computational speed limit for the chain                                                  |
| [![](e.png)][Os9] [`SetGasPoolSeconds`][O9]`(seconds)`           | Set the number of seconds worth of the speed limit the gas pool contains                         |
| [![](e.png)][Os10] [`SetGasPoolTarget`][O10]`(target)`           | Set the target fullness in bips the pricing model will try to keep the pool at                   |
| [![](e.png)][Os11] [`SetGasPoolWeight`][O11]`(weight)`           | Set the extent in bips to which the pricing model favors filling the pool over increasing speeds |
| [![](e.png)][Os12] [`SetRateEstimateInertia`][O12]`(inertia)`    | Set how slowly MtOS updates its estimate the amount of gas being burnt per second               |
| [![](e.png)][Os13] [`SetMaxTxGasLimit`][O13]`(limit)`            | Set the maximum size a tx (and block) can be                                                     |
| [![](e.png)][Os14] [`GetNetworkFeeAccount`][O14]`()`             | Get the network fee collector                                                                    |
| [![](e.png)][Os15] [`SetNetworkFeeAccount`][O15]`(account)`      | Set the network fee collector                                                                    |

[O0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L24
[O1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L29
[O2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L38
[O3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L43
[O4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L48
[O5]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L53
[O6]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L58
[O7]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L63
[O8]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L68
[O9]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L73
[O10]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L78
[O11]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L83
[O12]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L88
[O13]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L93
[O14]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L98
[O15]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwner.go#L103

[Os0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L30
[Os1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L33
[Os2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L36
[Os3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L39
[Os4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L42
[Os5]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L45
[Os6]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L48
[Os7]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L51
[Os8]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L54
[Os9]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L57
[Os10]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L60
[Os11]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L63
[Os12]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L66
[Os13]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L69
[Os14]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L72
[Os15]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L75

| Events                                       |                                                           |
|:---------------------------------------------|:----------------------------------------------------------|
| [![](e.png)][Oes0] [`OwnerActs`][Oe0] &nbsp; | Emitted when a successful call is made to this precompile |

[Oe0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/wrapper.go#L105

[Oes0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwner.sol#L78


# [MtOwnerPublic][MtOwnerPublic_link]
Provides non-owners with info about the current chain owners.

| Methods                                                    |                                 |
|:-----------------------------------------------------------|:--------------------------------|
| [![](e.png)][OPs0] [`IsChainOwner`][OP0]`(account)` &nbsp; | See if account is a chain owner |
| [![](e.png)][OPs1] [`GetAllChainOwners`][OP1]`()`          | Gets the list of chain owners   |
| [![](e.png)][OPs2] [`GetNetworkFeeAccount`][OP2]`()`       | Gets the network fee collector  |

[OP0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwnerPublic.go#L24
[OP1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwnerPublic.go#L19
[OP2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtOwnerPublic.go#L29

[OPs0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwnerPublic.sol#L25
[OPs1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwnerPublic.sol#L28
[OPs2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtOwnerPublic.sol#L31


# [MtRetryableTx][MtRetryableTx_link]
Provides methods for managing retryables. The model has been adjusted for Mantle, most notably in terms of how retry transactions are scheduled. For more information on retryables, please see [the retryable documentation](mtos.md#Retryables).


| Methods                                                     |                                                                                    | Mantle changes          |
|:------------------------------------------------------------|:-----------------------------------------------------------------------------------|:-----------------------|
| [![](e.png)][RTs0] [`Cancel`][RT0]`(ticket)`                | Cancel the ticket and refund its callvalue to its beneficiary                      |                        |
| [![](e.png)][RTs1] [`GetBeneficiary`][RT1]`(ticket)` &nbsp; | Gets the beneficiary of the ticket                                                 |                        |
| [![](e.png)][RTs2] [`GetLifetime`][RT2]`()`                 | Gets the default lifetime period a retryable has at creation                       | Reverts when not found |
| [![](e.png)][RTs3] [`GetTimeout`][RT3]`(ticket)`            | Gets the timestamp for when ticket will expire                                     |                        |
| [![](e.png)][RTs4] [`Keepalive`][RT4]`(ticket)`             | Adds one lifetime period to the ticket's expiry                                    | Doesn't add callvalue  |
| [![](e.png)][RTs5] [`Redeem`][RT5]`(ticket)`                | Schedule an attempt to redeem the retryable, donating all of the call's gas &nbsp; | Happens in a future tx |

[RT0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtRetryableTx.go#L184
[RT1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtRetryableTx.go#L171
[RT2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtRetryableTx.go#L110
[RT3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtRetryableTx.go#L115
[RT4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtRetryableTx.go#L132
[RT5]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtRetryableTx.go#L36

[RTs0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtRetryableTx.sol#L70
[RTs1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtRetryableTx.sol#L63
[RTs2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtRetryableTx.sol#L38
[RTs3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtRetryableTx.sol#L45
[RTs4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtRetryableTx.sol#L55
[RTs5]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtRetryableTx.sol#L32

| Events                                                |                                                    | Mantle Changes                       |
|:------------------------------------------------------|:---------------------------------------------------|:------------------------------------|
| [![](e.png)][RTes0] [`TicketCreated`][RTe0]           | Emitted when creating a retryable                  |                                     |
| [![](e.png)][RTes1] [`LifetimeExtended`][RTe1] &nbsp; | Emitted when extending a retryable's expiry &nbsp; |                                     |
| [![](e.png)][RTes2] [`RedeemScheduled`][RTe2]         | Emitted when scheduling a retryable                | Replaces [Redeemed][old_event_link] |
| [![](e.png)][RTes3] [`Canceled`][RTe3]                | Emitted when cancelling a retryable                |                                     |

[RTe0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/mtos/tx_processor.go#L143
[RTe1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtRetryableTx.go#L163
[RTe2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/mtos/tx_processor.go#L186
[RTe3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtRetryableTx.go#L209

[RTes0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtRetryableTx.sol#L72
[RTes1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtRetryableTx.sol#L73
[RTes2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtRetryableTx.sol#L74
[RTes3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtRetryableTx.sol#L81

[old_event_link]: https://github.com/mantlenetwork/mt-os/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/arb_os/arbretryable.mini#L90

# [MtStatistics][MtStatistics_link]
Provides statistics about the chain as of just before the Mantle upgrade. In Mantle Classic, this was how a user would get info such as the total number of accounts, but there are better ways to get that info in Mantle.

| Methods                                         |                                                                                         |
|:------------------------------------------------|:----------------------------------------------------------------------------------------|
| [![](e.png)][STs0] [`GetStats`][ST0]`()` &nbsp; | Returns the current block number and some statistics about the rollup's pre-Mantle state |

[ST0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtStatistics.go#L19

[STs0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtStatistics.sol#L32


# [MtSys][MtSys_link]
Provides system-level functionality for interacting with L1 and understanding the call stack.

| Methods                                                                                  |                                                                                                              | Mantle changes     |
|:-----------------------------------------------------------------------------------------|:-------------------------------------------------------------------------------------------------------------|:------------------|
| [![](e.png)][Ss0] [`ArbBlockNumber`][S0]`()`                                             | Gets the current L2 block number                                                                             |                   |
| [![](e.png)][Ss1] [`ArbBlockHash`][S1]`(blocknum)`                                       | Gets the L2 block hash at blocknum, if blocknum is sufficiently recent                                       |                   |
| [![](e.png)][Ss2] [`ArbChainID`][S2]`()`                                                 | Gets the chain's ChainID                                                                                     |                   |
| [![](e.png)][Ss3] [`MtOSVersion`][S3]`()`                                               | Gets the current MtOS version                                                                               | Now view          |
| [![](e.png)][Ss4] [`GetStorageGasAvailable`][S4]`()`                                     | Returns 0 since Mantle has no concept of storage gas                                                          | Now always 0      |
| [![](e.png)][Ss5] [`IsTopLevelCall`][S5]`()`                                             | Checks if the caller is top-level (i.e. if the caller was called directly by an EOA or an L1 contract)       |                   |
| [![](e.png)][Ss6] [`MapL1SenderContractAddressToL2Alias`][S6]`(contract, unused)` &nbsp; | Gets contract's L2 alias                                                                                     | 2nd arg is unused |
| [![](e.png)][Ss7] [`WasMyCallersAddressAliased`][S7]`()`                                 | Checks if the caller's caller was aliased                                                                    |                   |
| [![](e.png)][Ss8] [`MyCallersAddressWithoutAliasing`][S8]`()`                            | Gets the caller's caller without any potential address aliasing                                              | New outbox scheme |
| [![](e.png)][Ss9] [`SendTxToL1`][S9]`(destination, calldataForL1)`                       | Sends a transaction to L1, adding it to the outbox; callvalue is sent to L1 attached to the sent transaction | New outbox scheme |
| [![](e.png)][Ss10] [`SendMerkleTreeState`][S10]`()`                                      | Gets the root, size, and partials of the outbox Merkle tree state                                            | New outbox scheme |
| [![](e.png)][Ss11] [`WithdrawEth`][S11]`(destination)`                                   | Send callvalue to the destination address on L1                                                              |                   |

[S0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L30
[S1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L35
[S2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L50
[S3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L55
[S4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L61
[S5]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L66
[S6]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L71
[S7]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L76
[S8]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L82
[S9]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L98
[S10]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L171
[S11]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L187

[Ss0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L31
[Ss1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L37
[Ss2]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L43
[Ss3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L49
[Ss4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L55
[Ss5]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L61
[Ss6]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L69
[Ss7]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L78
[Ss8]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L84
[Ss9]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L100
[Ss10]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L111
[Ss11]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L92


| Events                                               |                                                                 |
|:-----------------------------------------------------|:----------------------------------------------------------------|
| [![](e.png)][Ses0] [`L2ToL1Transaction`][Se0] &nbsp; | Logs a send tx from L2 to L1, including data for outbox proving |
| [![](e.png)][Ses1] [`SendMerkleUpdate`][Se1]         | Logs a new merkle branch needed for constructing outbox proofs  |

[Se0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L152
[Se1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtSys.go#L138

[Ses0]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L124
[Ses1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtSys.sol#L143

| Removed Methods                                                   |                                                                   |
|:------------------------------------------------------------------|:------------------------------------------------------------------|
| [![](e.png)][Srs0] [`GetStorageAt`][Sr0]`(account, index)` &nbsp; | Mantle doesn't need this introspection, and users couldn't call it |
| [![](e.png)][Srs1] [`GetTransactionCount`][Sr1]`(account)`        | Mantle doesn't need this introspection, and users couldn't call it |

[Sr0]: https://github.com/mantlenetwork/mt-os/blob/89e36db597c4857a4dac3efd7cc01b13c7845cc0/arb_os/arbsys.mini#L335
[Sr1]: https://github.com/mantlenetwork/mt-os/blob/89e36db597c4857a4dac3efd7cc01b13c7845cc0/arb_os/arbsys.mini#L315

[Srs0]: https://github.com/mantlenetwork/mt-os/blob/89e36db597c4857a4dac3efd7cc01b13c7845cc0/contracts/mtos/builtin/MtSys.sol#L51
[Srs1]: https://github.com/mantlenetwork/mt-os/blob/89e36db597c4857a4dac3efd7cc01b13c7845cc0/contracts/mtos/builtin/MtSys.sol#L42
