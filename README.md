straightforward
===

Minimal SOCKS proxy over SSH.

## Build
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fatpons%2Fstraightforward.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fatpons%2Fstraightforward?ref=badge_shield)

```shell
$ bazel build cmd/straightforward
```

## Usage

```
+-------------------+   +-------------------+   +---------------+
|                   |   |                   |   |               |
|   Local Machine   +---> Remote SSH Server +--->  Target Host  |
|                   |   |                   |   |               |
+----------+--------+   +-------------------+   +---------------+
           |
  SOCKS5   |
           |
+----------+--------+
|                   |
|  Any Application  |
|                   |
+-------------------+
```
```shell
$ bazel run cmd/straighforward -- -i <SSH Private Key> -u <SSH User> -h <Destination SSH Server>:<Port> -p <Listen SOCKS Proxy Port>
```

### over SOCKS server (ocproxy mode)

You can use SOCKS5 proxy behind connection for SSH server (like ocproxy).

```
 +-------------------+   +------------------+   +-------------------+   +---------------+
 |                   |   |                  |   |                   |   |               |
 |   Local Machine   +--->SOCKS5 Local Proxy+---+ Remote SSH Server +--->  Target Host  |
 |                   |   |                  |   |                   |   |               |
 +----------+--------+   +------------------+   +-------------------+   +---------------+
            |
   SOCKS5   |
            |
 +----------+--------+
 |                   |
 |  Any Application  |
 |                   |
 +-------------------+
```
```shell
$ bazel run cmd/straighforward -- -i <SSH Private Key> -u <SSH User> -h <Destination SSH Server>:<Port> -p <Listen SOCKS Proxy Port> -ocproxy <Local SOCKS Server>:<Port>
```



## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fatpons%2Fstraightforward.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fatpons%2Fstraightforward?ref=badge_large)