# Curated Go Networking Libs

[![GoDoc](https://pkg.go.dev/badge/github.com/bassosimone/nettree)](https://pkg.go.dev/github.com/bassosimone/nettree) [![Build Status](https://github.com/bassosimone/nettree/actions/workflows/go.yml/badge.svg)](https://github.com/bassosimone/nettree/actions) [![codecov](https://codecov.io/gh/bassosimone/nettree/branch/main/graph/badge.svg)](https://codecov.io/gh/bassosimone/nettree)

This repository contains curated versions of the Go networking
libraries that I maintain for fun. Whereas the actual libraries
live at "HEAD" and have no tags, this repository has tags. In
addition, libraries live inside the same import path instead of
being small, scattered, independent repositories.

## Packages

| Name                         | Purpose                              |
| ---------------------------- | ------------------------------------ |
| [dnscodec][dnscodec]         | Marshal/unmarshal DNS messages       |
| [dnsoverhttps][dnsoverhttps] | DNS over HTTPS transport             |
| [dnstest][dnstest]           | Helpers for writing DNS tests        |
| [httptestx][httptestx]       | Helpers for writing HTTP tests       |
| [iotest][iotest]             | Helpers for `io` tests               |
| [iox][iox]                   | Extensions for the `io` package      |
| [minest][minest]             | Minimal network stack                |
| [netstub][netstub]           | Stub for testing `net` code          |
| [pkitest][pkitest]           | Helpers for testing with certs       |
| [runtimex][runtimex]         | Extensions for the `runtime` package |

[dnscodec]: pkg/dnscodec
[dnsoverhttps]: pkg/dnsoverhttps
[dnstest]: pkg/dnstest
[httptestx]: pkg/httptestx
[iotest]: pkg/iotest
[iox]: pkg/iox
[minest]: pkg/minest
[netstub]: pkg/netstub
[pkitest]: pkg/pkitest
[runtimex]: pkg/runtimex

Each directory inside `pkg` contains a README.md pointing to
the upstream repository it derives from, along with an explanation
of how to use the Go package cache to verify it.

The [remotes](remotes) directory contains the JSONs downloaded
from the Go package cache during the last sync.

## Minimum Go Version

The "oldstable" patch release available during the last sync
as indicated in the [go.mod](go.mod) file.

## Testing

```bash
go test -race ./...
```

## Releases

A minimum of two minor releases (i.e. `v0.X.0`) per year:

1. in April and October to sync up with the Go releases
happening in February and August respectively

2. as needed, for bug fixes

We will do our best to keep backwards compatibility. We will
not break APIs unless there are security or performance reasons.

## Syncing

To sync from my upstream repositories:

```bash
make sync
```

which uses [scripts/sync.bash](scripts/sync.bash).

## Contributing

If you open pull requests in this repository, the most likely
outcome is that I will cherry-pick the commits and apply
them in the upstream repositories with attribution.

## Scope

I maintain these libraries in my spare time, mainly for my own
network measurement fun. Bug reports are welcome. Feature requests
may be declined if they don't fit my needs or would add
maintenance burden. The code is GPL-3.0-or-later, so forking
is always an option.

## License

```
SPDX-License-Identifier: GPL-3.0-or-later
```

## History

This code evolved from [ooni/probe-cli][probe-cli] code
or [rbmk-project/rbmk][rbmk] code through various intermediate
steps, as indicated in detail in each upstream repository and by
detailed `// Adapted from: ...` comments in each file
that was not written from scratch.

[probe-cli]: https://github.com/ooni/probe-cli
[rbmk]: https://github.com/rbmk-project/rbmk
