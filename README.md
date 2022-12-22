# Cryptographic Commitment Schemes

This repo aims to provide a simple and easy to use implementation of cryptographic commitment schemes. The goal is to provide a simple and easy to use interface for the user to create commitments and verify them. The user should not have to worry about the underlying cryptographic primitives.

Primitives for cryptographic commitment schemes.
- operations of group, field and polynomial.
  - import group/mod from "github.com/drand/kyber/group/mod"
  - import bn256 from "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
  - polynomial.go 
- Hash commitment
  - hash_commitment.go
- Polynomial Commitment
  - kzg.go ([KZG commitment](https://cacr.uwaterloo.ca/techreports/2010/cacr2010-10.pdf))