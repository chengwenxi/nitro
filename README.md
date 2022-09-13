# Mantle

Mantle is a fully integrated, complete
layer 2 optimistic rollup system, including fraud proofs, the sequencer, the token bridges, 
advanced calldata compression, and more.

The Mantle stack is built on several innovations. At its core is a new prover, which can do Mantle’s classic 
interactive fraud proofs over WASM code. That means the L2 Mantle engine can be written and compiled using 
standard languages and tools, replacing the custom-designed language and compiler used in previous Mantle
versions. In normal execution, 
validators and nodes run the Nitro engine compiled to native code, switching to WASM if a fraud proof is needed. 
We compile the core of Geth, the EVM engine that practically defines the Ethereum standard, right into Mantle. 
So the previous custom-built EVM emulator is replaced by Geth, the most popular and well-supported Ethereum client.

The last piece of the stack is a slimmed-down version of our MtOS component, rewritten in Go, which provides the 
rest of what’s needed to run an L2 chain: things like cross-chain communication, and a new and improved batching 
and compression system to minimize L1 costs.

Essentially, Nitro runs Geth at layer 2 on top of Ethereum, and can prove fraud over the core engine of Geth 
compiled to WASM.

During the devnet period, we have licensed Mantle under a 
Business Source License, similar to our friends at Uniswap and Aave. 
Before mainnet launch, we will be re-licensing the code in a more 
permissive fashion to ensure that everyone can have full comfort 
using and running nodes on Mantle.
