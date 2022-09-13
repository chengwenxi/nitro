# Public Chains

The following is a comprehensive list of all of the currently live Mantle chains:

| Name                         | RPC Url(s)                                                                                                                         | ID     | Native Currency | Explorer(s)                                                          | Underlying L1 | Current Tech Stack  | Sequencer Feed                         | Nitro Seed Database URLs                 |
| ---------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- | ------ | --------------- | -------------------------------------------------------------------- | ------------- | ------------------- | -------------------------------------- | ---------------------------------------- |
| Mantle One                 | `https://arb1.mantle.io/rpc` `https://mantle-mainnet.infura.io/v3/YOUR-PROJECT-ID` `https://arb-mainnet.g.alchemy.com/v2/-KEY` | 42161  | ETH             | `https://arbiscan.io/` `https://explorer.mantle.io/`               | Ethereum      | Nitro Rollup (8/31) | `wss://arb1.mantle.io/feed`          |  `snapshot.mantle.io/mainnet/nitro.tar`                        |
| Mantle Nova                | `https://nova.mantle.io/rpc`                                                                                                     | 42170  | ETH             | `https://nova-explorer.mantle.io/`                                 | Ethereum      | Nitro AnyTrust      | `wss://nova.mantle.io/feed`          | N/A                                      |
| RinkArby^                    | `https://rinkeby.mantle.io/rpc`                                                                                                  | 421611 | RinkebyETH      | `https://testnet.arbiscan.io` `https://rinkeby-explorer.mantle.io` | Rinkeby       | Nitro Rollup        | `wss://rinkeby.mantle.io/feed`       | `snapshot.mantle.io/rinkeby/nitro.tar` |
| Nitro Goerli Rollup Testnet^ | `https://goerli-rollup.mantle.io/rpc`                                                                                            | 421613 | GoerliETH       | `https://goerli-rollup-explorer.mantle.io`                         | Goerli        | Nitro Rollup        | `wss://goerli-rollup.mantle.io/feed` | N/A                                      |

^ Testnet

All chains use [bridge.mantle.io/](https://bridge.mantle.io/) for bridging assets and [retryable-dashboard.mantle.io](https://retryable-dashboard.mantle.io/) for executing [retryable tickets](l1-to-l2-messagaing) if needed.

For a list of useful contract addresses, see [here](useful-addresses).

### Mantle Chains Summary

**Mantle One**: Mantle One is the flagship Mantle mainnet chain; it is an Optimistic Rollup chain running on top of Ethereum Mainnet, and is open to all users. In an upgrade on 8/31, the Mantle One chain is/was upgraded to use the [Nitro](https://medium.com/offchainlabs/its-nitro-time-86944693bf29) tech stack, maintaining the same state.
Users can now use [Alchemy](https://alchemy.com/?a=mantle-docs), [Infura](https://infura.io/), [QuickNode](https://www.quicknode.com), [Moralis](https://moralis.io/), [Ankr](https://www.ankr.com/), [BlockVision](https://blockvision.org/), and [GetBlock](https://getblock.io/) to interact with the Mantle One. See [node providers](node-providers) for more.

**Mantle Nova**: Mantle Nova is the first mainnet [AnyTrust](inside-anytrust) chain. The following are the members of the initial data availability committee (DAC):
- Consensys
- FTX
- Google Cloud
- Mantlenetwork
- P2P
- Quicknode
- Reddit

Users can now use [QuickNode](https://www.quicknode.com) to interact with the Mantle Nova chain. For a full guide of how to set up an Mantle node on QuickNode, see the QuickNode's Mantle RPC documentation.

**RinkArby**: RinkArby is the longest running Mantle testnet. It previously ran on the classic stack, but at block 7/28/2022 it was migrated use the Nitro stack! Rinkarby will be deprecated [when Rinkeby itself gets deprecated](https://blog.ethereum.org/2022/06/21/testnet-deprecation/); plan accordingly!
Users can now use [Alchemy](https://alchemy.com/?a=mantle-docs), [Infura](https://infura.io/), [QuickNode](https://www.quicknode.com), [Moralis](https://moralis.io/), [Ankr](https://www.ankr.com/), [BlockVision](https://blockvision.org/), and [GetBlock](https://getblock.io/) to interact with the Mantle One. See [node providers](node-providers) for the full guide.

**Nitro Goerli Rollup Testnet**: This testnet (421613) uses the Nitro rollup tech stack; it is expected to be the primary, stable Mantle testnet moving forward.
Users can now use [Alchemy](https://alchemy.com/?a=mantle-docs), [Infura](https://infura.io/), and [QuickNode](https://www.quicknode.com) to interact with the Mantle One. See [node providers](./node-running/node-providers.md) for more.
