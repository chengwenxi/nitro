# Mantle/Ethereum Differences

Mantle is designed to be as compatible and consistent with Ethereum as possible, from its high-level RPCs to its low-level bytecode and everything in between. Dapp developers with experience building on Ethereum will likely find that little-to-no new L2-specific knowledge is required to build on Mantle.

This document presents an overview of some of the minor differences, perks, and gotchas that devs are advised to be aware of.

## Transactions / Blocks

##### Blocks and Time

Time in L2 is tricky; the timing assumptions one is used to making about L1 blocks don't exactly carry over into the timing of Mantle blocks. For details, see [Block Numbers and Time](./time.md).

##### Block hashes and randomness

Mantle's L2 block hashes should not be relied on as a secure source of randomness (see ['blockhash(x);](./solidity-support.md))

##### L1 Fees

The L2 fees an Mantle transaction pays essentially work identically to gas fees on Ethereum. Mantle transactions must also, however, pay an L1-fee component to cover the cost of their calldata. (See [L1 pricing](./mtos/l1-pricing.md).)

##### Tx Receipts

Mantle transaction receipts include two additional fields:

1. `l1BlockNumber`: The l1 block number that would be used [for block.number calls](time).
1. `gasUsedForL1`: Amount of gas spent on l1 computation in units of l2 gas.

## L1 to L2 Messages

Mantle chains support arbitrary L1 to L2 message passing; developers using this functionality should familiarize themselves with how they work (see [L1 to L2 Messaging](./mtos/l1-to-l2-messaging.md)). Of particular note:

- The result of a successful initial/"auto"-execution of an L1 to L2 message will be an unsigned L2 tx receipt.
- The `msg.sender` of the L2 side of an L1 to L2 message will be not the initiating L1 address, but rather its address alias.
- Using the special `ethDeposit` method will _not_ result in an L2 contract's fallback function getting triggered.

Etc.

## Precompiles

Mantle chains include a number of special precompiles not present on Ethereum; see [Common Precompiles](./mtos/common-precompiles.md) / [All Precompiles](./mtos/precompiles.md).

Of particular note is the [MtAddressTable](./mtos/precompiles.md#MtAddressTable), which allows contracts to map addresses to integers, saving calldata / fees for addresses expected to be reused as parameters; see [Arb Address Table tutorial](https://github.com/mantlenetworkio/mantle-tutorials/tree/master/packages/address-table) for example usage.

## Solidity

You can deploy Solidity contracts onto Mantle just like you do Ethereum; there are only a few minor differences in behavior. See [Solidity Support](./solidity-support.md) for details.
