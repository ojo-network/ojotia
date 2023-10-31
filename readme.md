# CI Tool To Post Data from Ojo -> Celestia

This is a tool that anyone can use to post the current price data from the Ojo blockchain to Celestia. It was inspired by [@dferrersan](https://x.com/dferrersan/status/1719427046663725381?s=20) on twitter / x.

## Installation

```
git clone https://github.com/ojo-network/ojotia
cd ojotia
make install
```

## Usage

This assumes you have a celestia light node installed and have access to an ojo GRPC endpoint.

* [Celestia light node docs](https://docs.celestia.org/developers/node-tutorial#instantiate-a-celestia-light-node)
* [Ojo Agamotto Network](https://agamotto.ojo.network/)

To submit data, run:

```
ojotia [auth-token] [celestia-rpc-addr] [ojo-grpc-addr]
```

It should return something similar to:

```
Succesfully submitted blob to Celestia
Height:  12345
Commitment string:  d4bfd96812d50d52f3966fd5e847596ed3176ea1d7958cfx1981804688d89ec8
```

To query that data, run:

```
ojotia query [auth-token] [celestia-rpc-addr] [commitment] [height]
```

Which should return:

```
Celestia queried successfully!
Ojo Price Data:
{7.815472986524036485ATOM}
```
