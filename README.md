# `go-open-service-broker-client`

[![Build Status](https://travis-ci.org/pmorie/go-open-service-broker-client.svg?branch=master)](https://travis-ci.org/pmorie/go-open-service-broker-client)

A golang client for service brokers implementing the
[Open Service Broker API](https://github.com/openservicebrokerapi/servicebroker).

## Goals

Overall, to make an excellent golang client for the Open Service Broker API.
Specifically:

- Support new auth modes in a backward-compatible manner
- Support alpha features in the Open Service Broker API in a clear manner
- Allow advanced configuration of TLS configuration to a broker
- Provide a fake client suitable for unit-type testing

## Non-goals

This project does not aim to provide:

- A fake _service broker_
- A conformance suite for service brokers
- Any 'custom' API features with are not either in a released version of the
  Open Service Broker API spec or accepted into the spec but not yet released

## Current status

This repository is a pet project of mine.  During the day, I work on the
[Kubernetes](https://github.com/kubernetes/kubernetes)
[Service Catalog](https://github.com/pmorie/go-open-service-broker-client) as well as
the Open Service Broker API specification.  My overall goal is for this repo to
contain a great golang client for the API, as well as a fake that facilitates
unit-style testing.

I'm currently just sketching around in my free time, but would love to
eventually donate this repo to an organization if people decide that it is
useful.