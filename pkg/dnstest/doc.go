// SPDX-License-Identifier: GPL-3.0-or-later

/*
Package dnstest contains helpers for writing tests for DNS clients.

This package provides stdlib-independent helpers for testing various
kinds of DNS clients. For now, there is support for DNS over TCP,
TLS, and HTTPS.

The overall intention is to support writing tests against servers that
are created and managed by this package. While this package does not
aim to make its servers mimic production servers, the intention is to
provide enough functionality for testing.

This package's API design is inspired by net/http/httptest. Like
net/http/httptest, we panic when we cannot create a testing server
because, in a test, such a failure should be loud and obvious.
*/
package dnstest
