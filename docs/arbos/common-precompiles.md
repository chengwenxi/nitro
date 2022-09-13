# Overview
MtOS provides L2-specific precompiles with methods smart contracts can call the same way they can solidity functions. This reference details those we expect users to most frequently use. For an exhaustive reference including those we don't expect most users to ever call, please refer to the [Full Precompiles documentation](precompiles.md).

From the perspective of user applications, precompiles live as contracts at the following addresses. Click on any to jump to their section.

| Precompile                                 | Address &nbsp; | Purpose                             |
| :----------------------------------------- | :------------- | :---------------------------------- |
| [`MtAggregator`](#MtAggregator)          | `0x6d`         | Configuring transaction aggregation |
| [`MtGasInfo`](#MtGasInfo)                | `0x6c`         | Info about gas pricing              |
| [`MtRetryableTx`](#MtRetryableTx) &nbsp; | `0x6e`         | Managing retryables                 |
| [`MtSys`](#MtSys)                        | `0x64`         | System-level functionality          |

[MtAggregator_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtAddressTable.go
[MtGasInfo_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtGasInfo.go
[MtRetryableTx_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtRetryableTx.go
[MtSys_link]: https://github.com/mantlenetworkio/mantle/blob/master/precompiles/MtSys.go

# [MtAggregator][MtAggregator_link]
Provides aggregators and their users methods for configuring how they participate in L1 aggregation. Mantle One's default aggregator is the Sequencer, which a user will prefer unless `SetPreferredAggregator` is invoked to change it.

| Methods                                                        |                                                         |
| :------------------------------------------------------------- | :------------------------------------------------------ |
| [![](e.png)][As0] [`GetPreferredAggregator`][A0]`(account)`    | Gets an account's preferred aggregator                  |
| [![](e.png)][As1] [`SetPreferredAggregator`][A1]`(aggregator)` | Sets the caller's preferred aggregator to that provided |
| [![](e.png)][As2] [`GetDefaultAggregator`][A2]`()`             | Gets the chain's default aggregator                     |

[A0]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtAggregator.go#L22
[A1]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtAggregator.go#L39
[A2]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtAggregator.go#L48

[As0]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtAggregator.sol#L28
[As1]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtAggregator.sol#L32
[As2]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtAggregator.sol#L35


# [MtGasInfo][MtGasInfo_link]
Provides insight into the cost of using the chain. These methods have been adjusted to account for Nitro's heavy use of calldata compression. Of note to end-users, we no longer make a distinction between non-zero and zero-valued calldata bytes.

| Methods                                                |                                                                   |
| :----------------------------------------------------- | :---------------------------------------------------------------- |
| [![](e.png)][GIs1] [`GetPricesInWei`][GI1]`()`         | Get prices in wei when using the caller's preferred aggregator    |
| [![](e.png)][GIs3] [`GetPricesInArbGas`][GI3]`()`      | Get prices in ArbGas when using the caller's preferred aggregator |
| [![](e.png)][GIs4] [`GetGasAccountingParams`][GI4]`()` | Get the chain speed limit, pool size, and tx gas limit            |
| [![](e.png)][GIs11] [`GetL1BaseFeeEstimate`][GI11]`()` | Get MtOS's estimate of the L1 basefee in wei                     |

[GI1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L63
[GI3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L99
[GI4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L111
[GI11]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/precompiles/MtGasInfo.go#L150

[GIs1]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L58
[GIs3]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L83
[GIs4]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L94
[GIs11]: https://github.com/mantlenetworkio/mantle/blob/3f504c57fba8ddf0759b7a55b4108e0bf5a078b3/solgen/src/precompiles/MtGasInfo.sol#L122

# [MtRetryableTx][MtRetryableTx_link]
Provides methods for managing retryables. The model has been adjusted for Nitro, most notably in terms of how retry transactions are scheduled. For more information on retryables, please see [the retryable documentation](mtos.md#Retryables).


| Methods                                                     |                                                                                    | Nitro changes          |
| :---------------------------------------------------------- | :--------------------------------------------------------------------------------- | :--------------------- |
| [![](e.png)][RTs0] [`Cancel`][RT0]`(ticket)`                | Cancel the ticket and refund its callvalue to its beneficiary                      |                        |
| [![](e.png)][RTs1] [`GetBeneficiary`][RT1]`(ticket)` &nbsp; | Gets the beneficiary of the ticket                                                 |                        |
| [![](e.png)][RTs3] [`GetTimeout`][RT3]`(ticket)`            | Gets the timestamp for when ticket will expire                                     |                        |
| [![](e.png)][RTs4] [`Keepalive`][RT4]`(ticket)`             | Adds one lifetime period to the ticket's expiry                                    | Doesn't add callvalue  |
| [![](e.png)][RTs5] [`Redeem`][RT5]`(ticket)`                | Schedule an attempt to redeem the retryable, donating all of the call's gas &nbsp; | Happens in a future tx |

[RT0]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtRetryableTx.go#L184
[RT1]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtRetryableTx.go#L171
[RT3]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtRetryableTx.go#L115
[RT4]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtRetryableTx.go#L132
[RT5]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtRetryableTx.go#L36

[RTs0]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtRetryableTx.sol#L70
[RTs1]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtRetryableTx.sol#L63
[RTs3]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtRetryableTx.sol#L45
[RTs4]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtRetryableTx.sol#L55
[RTs5]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtRetryableTx.sol#L32


# [MtSys][MtSys_link]
Provides system-level functionality for interacting with L1 and understanding the call stack.

| Methods                                                            |                                                             |
| :----------------------------------------------------------------- | :---------------------------------------------------------- |
| [![](e.png)][Ss0] [`ArbBlockNumber`][S0]`()`                       | Gets the current L2 block number                            |
| [![](e.png)][Ss1] [`ArbBlockHash`][S1]`()`                         | Gets the L2 block hash, if the block is sufficiently recent |
| [![](e.png)][Ss5] [`IsTopLevelCall`][S5]`()`                       | Checks if the call is top-level                             |
| [![](e.png)][Ss9] [`SendTxToL1`][S9]`(destination, calldataForL1)` | Sends a transaction to L1, adding it to the outbox          |
| [![](e.png)][Ss11] [`WithdrawEth`][S11]`(destination)`             | Send paid eth to the destination on L1                      |

[S0]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtSys.go#L30
[S1]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtSys.go#L35
[S5]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtSys.go#L66
[S9]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtSys.go#L98
[S11]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/precompiles/MtSys.go#L187

[Ss0]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtSys.sol#L31
[Ss1]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtSys.sol#L37
[Ss5]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtSys.sol#L61
[Ss9]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtSys.sol#L100
[Ss11]: https://github.com/mantlenetworkio/mantle/blob/704e82bb38ae3ccd70c35e31934c7b45f6c25561/solgen/src/precompiles/MtSys.sol#L92
