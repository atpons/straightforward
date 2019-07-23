straightforward
===

Minimal SOCKS proxy over SSH.

## Build
```shell
$ bazel build cmd/straightforward
```

## Usage
```shell
$ bazel run cmd/straighforward -- -i <SSH Private Key> -u <SSH User> -h <Destination SSH Server>:<Port> -p <Listen SOCKS Proxy Port>
```