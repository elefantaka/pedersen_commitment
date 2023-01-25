# Pedersen Commitment

### Introduction
The practical part of the report is a Proof of Concept simulation of a piece of Monero transaction.

### How it works?
The program calculates input commitment and inputs of amounts according to transaction amount. Then prepares the output commitment and outputs of trasactions including the fee. When output commitment equals input commitment and the total of outputs trasactions equals trasaction amout the commitment is true. Our implementation does not include the verification phase of the commitment - we print input and output commitment, and inputs with outputs for the user to show how the commitment works.

### Instalation 
1. install Go
2. `git clone https://github.com/elefantaka/pedersen_commitment.git`
3. `go run main.go`


