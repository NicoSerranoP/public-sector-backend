# Architecture

A Go (Gin) API backend service to support an on-chain DAO. The DAO voters are NFT holders where 1 token = 1 vote (tokens can be delegated). The NFT minting authority is a EOA (Externally Owned Account) managed by the public sector controlling entity. The entity can mint and burn NFTs depending on off-chain registration that require citizens to present their identification card and information. The minted NFT can contain a ilustrative image and metadata to be uniquely identified. All this process is performed by the API backend service

## Functionalities

### Health endpoint

Notify if the system is healthy and operational (including database access)

### Register voter endpoint

Receive citizen information (contact information, identification card, and other relevant details), perform validation and check existence in the database. If not, create a new voter record and mint a corresponding NFT.

### Read voter information endpoint

Retrieve voter information by Ethereum address or identity document number, including their registration details and associated NFT metadata.

### Cast a vote endpoint

Receive a vote from a voter through the official UI, validate it (validate signature, check sender address is authorized, check receiver is part of the DAO framework, check the voting is not blacklisted by misbehavior), and record it on-chain using the OpenZeppelin relayer

### Create a proposal endpoint

Receive a proposal submission from a voter through the official UI, validate it (validate signature, check sender address is authorized, check the proposal meets the DAO framework requirements), and record it on-chain using the OpenZeppelin relayer.

### Read blacklist endpoint

Return the list of blacklisted voters in this API backend service, including their Ethereum addresses and the reasons for blacklisting.
