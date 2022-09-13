# Q: How do gas fees work on Mantle?

Fees on Mantle chains are collected on L2 in the chains' native currency (ETH on Mantle One and Nova).

A transaction fee is comprised of both an L1 and an L2 component:

The L1 component is meant to compensate the Sequencer for the cost of posting transactions on L1 (but no more). (See [L1 Pricing](../mtos/l1-pricing.md).)

The L2 component covers the cost of operating the L2 chain; it uses Geth for gas calculation and thus behaves nearly identically to L1 Ethereum (See [Gas](../mtos/gas.md)).

L2 Gas price adjusts responsively to chain congestion, ala EIP 1559.

Calling an Mantle node's `eth_estimateGas` RPC returns a value sufficient to cover both the L1 and L2 components of the fee for the current gas price; this is the value that, e.g., will appear in users' wallets. (See [2-D fees](https://medium.com/offchainlabs/understanding-mantle-2-dimensional-fees-fd1d582596c9) for more.)
