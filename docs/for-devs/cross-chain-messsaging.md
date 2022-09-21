# Cross Chain Messaging

The Mantle protocol and related tooling makes it easy for developers to build cross-chain applications; i.e., applications that involve sending messages from Ethereum to an Mantle chain, and/or from an Mantle chain to Ethereum.

## Ethereum to Mantle Messaging 

Arbitrary L1 to L2 contract calls can be created via the `Inbox`'s `createRetryableTicket` method; upon publishing the L1 transaction, the L2 side will typically get included within minutes. Happily / commonly, the L2 execution will automatically succeed, but if reverts, and it can be rexecuted via a call to the [`redeem`](../mtos/precompiles.md#MtRetryableTx) method of the [`MtRetryableTx`](../mtos/precompiles.md#MtRetryableTx) precompile.

For details and protocol specification, see [L1 to L2 Messages](../mtos/l1-to-l2-messaging.md).

For an example of retryable tickets in action, see the [Greeter](https://github.com/mantlenetwork/mantle-tutorials/tree/master/packages/greeter) tutorial, which uses the [Mantle SDK](./sdk). 


## Mantle to Ethereum Messaging

Similarly, L2 contracts can send Arbitrary messages for execution on L1. These are initiated via calls to the [`MtSys`](../mtos/precompiles.md#MtSys) precompile contract's `sendTxToL1` method. Upon confirmation (about 1 week later), they can executed by retrieving the relevant data via a call to `NodeInterface` contract's `constructOutboxProof` method, and then executing them via the  `Outbox`'s `executeTransaction` method. 

For details and protocol specification, see [L2 to L1 Messages](../mtos/l2-to-l1-messaging.md).

For a demo, see the [Outbox Tutorial](https://github.com/mantlenetwork/mantle-tutorials/tree/master/packages/outbox-execute). 
