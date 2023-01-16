# Pedersen Commitment

### Introduction
Implementation of Pedersen Commitment scheme based on Monero example from https://medium.com/coinmonks/zero-knowledge-proofs-um-what-a092f0ee9f28

### How it works?
The program calculates input commitment and inputs of amounts according to transaction amount. Then prepares the output commitment and outputs of trasactions including the fee. When output commitment equals input commitment and the total of outputs trasactions equals trasaction amout the commitment is true.

### Instalation 
1. install Go
2. `git clone https://github.com/elefantaka/pedersen_commitment.git`
3. `go run main.go`


