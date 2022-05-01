# “Word of Wisdom” TCP server protected from DOS with Proof of Work algorithm

## Description
"Word of Wisdom" TCP server which protected from Denial of Service attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work) 
algorithm via Challenge-Response protocol. Used Client Puzzle Protocol (https://en.wikipedia.org/wiki/Client_Puzzle_Protocol) as a POW algorithm. 
After the POW verification server sends quotes from the wise words collection.

## Getting started
### Requirements
+ [Go 1.17+](https://go.dev/)
+ [Docker](https://docs.docker.com/)
+ [Docker Compose](https://docs.docker.com/compose/install/)
+ [Make](https://www.gnu.org/software/make/)

### Install
```
make install
```

### Run tests
```
make test
```

### Run TCP server and client
```
make start
```

## More about the Proof of Work
Proof of work (PoW) is a form of cryptographic proof in which one party (the client) proves to others 
(the server) that a certain amount of a specific computational effort has been expended, which is generally 
a hash computation by the given criteria. POW is mainly used in blockchain, protection from DOS attacks, 
protection from spam abuses. There are multiple POW algorithms, among of which:
+ Client Puzzle Protocol (https://en.wikipedia.org/wiki/Client_Puzzle_Protocol)
+ HashCash (https://en.wikipedia.org/wiki/Hashcash)
+ Merkle Tree (https://en.wikipedia.org/wiki/Merkle_tree)
+ Equihash (https://en.wikipedia.org/wiki/Equihash)

## More about the Client Puzzle Protocol and why it is chosen
Client Puzzle Protocol (CPP) serves a cryptographic puzzle to clients that must be solved before their request is served.
In CPP two levels of hashing are used. First, the server-side secret and timestamp are hashed (puzzle hash). 
Then this hash is hashed again (target hash). The puzzle consists of target hash and most of the bits of puzzle hash.
Puzzle is solved by brute forcing missing bits of the hash.

CPP's main use case is protection against DOS attacks which makes it perfect to use here. While solving it requires some 
resources from client, the verification of the puzzle solution is very fast, so the protocol imposes negligible impact 
on the server.

## To Do
+ The "Wisdom Book" should be stored in storage instead of in the slice.
+ It would be great if the CPP puzzle strength is adjusted by the server connections amount, so more loaded the server
more difficult is the puzzle.

Inspired by: https://github.com/nightlord189/tcp-pow-go
