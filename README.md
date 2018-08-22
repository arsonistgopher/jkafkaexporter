# jKafkaexporter

[![GoDoc](https://godoc.org/github.com/arsonistgopher/jkafkaexporter?status.svg)](https://godoc.org/github.com/arsonistgopher/go-netconf/jkafkaexporter)
[![Build Status](https://travis-ci.org/arsonistgopher/jkafkaexporter.svg?branch=master)](https://travis-ci.org/arsonistgopher/jkafkaexporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/arsonistgopher/jkafkaexporter)](https://goreportcard.com/report/github.com/arsonistgopher/jkafkaexporter)

This is an application for Junos that drives NETCONF in order to retrieve data that is transformed into JSON and then placed on to a Kafka bus. Each Kafka topic is derived from the collector name.

> **Note:** The idea here is that it is `easy` to create a new collector and re-build the binary.

## Release

Currently beta version. Seems to behave itself!

Kafka options are fairly limited, but the idea here is if we start with the basics and add features as they're required, the project will receive enhancements naturally.

## Features
* Simple user land application 

## Install
* Requires Go 1.4 or later
* `go get github.com/arsonistgopher/jkafkaexporter`

## Documentation
You can view full API documentation at GoDoc: http://godoc.org/github.com/arsonistgopher/jkafkaexporter.

## Usage

`./jkafkaexporter -identity vmx01 -kafkaperiod 1 -kafkahost localhost:9092 -password Passw0rd -username bob -sshport 22 -target 10.42.0.132`

If you want to run this in the background or as a service, my advice here is to put it in a Docker container (see the Dockerfile). 

Also, this application can also work with ssh-keys using the '-sshkey` switch which the argument is the fully qualified path to the ssh-private key. Ensure the public key has been installed on to the Junos device in question.

One might wonder how the Kafka topics are created. Simple. They're automatically created from the collector name located in the [main.go](https://github.com/arsonistgopher/jkafkaexporter/blob/master/main.go#L95) file. For commit [a4977bc4dbef370b52f20a1d51e4dc94d50fd952](https://github.com/arsonistgopher/jkafkaexporter/commit/a4977bc4dbef370b52f20a1d51e4dc94d50fd952), the Kafka topics will be:

* alarm
* interfaces
* routing-engine
* environment
* route
* bgp
* interfacediagnostics

Within each topic, a record discriminator will be the identity which is passed in to the client as a command line argument (the example uses 'vmx01'). The identity appears as "Node:" within each JSON blob of data published to each Kafka topic.

## Custom Collectors

The collector pattern allows for easy custom creation of a collector and insertion into the application. Simply copy an existing collector, modify code and rename files to suite. Then, include the import in `main.go` at both the package import section and relevant code section. See examples below:

```go
    import (
    ...
	// Add new collectors here
    "github.com/arsonistgopher/jkafkaexporter/collectors/something"
    )

    ...
    
    // And also add new collectors here...
	...
	c.Add("something", something.NewCollector())
```

In the last line of config, the "something" is also the Kafka topic.

## Docker

In order to build the Docker image, follow the steps below. Feel free to modify as you see fit.

```bash
docker build -f Dockerfile . -t arsonistgopher/jkafkaexporter
docker run -d --name jkafkaexporter1 arsonistgopher/jkafkaexporter
```

You can also check Docker's logs
```bash
docker logs jkafkaexporter1
```

## Tutorial

WIP

## Attribution

Some code patterns were borrowed from [Daniel Czerwonk's Junos Prometheus exporter](https://github.com/czerwonk/junos_exporter) and as such his name is correctly present in the license. 

## Contributing

Welcoming contributions with arms open. Steps to contribute:

1.  Fork this repo using GitHub in to your own account
2.  Create changes on master and create tests ideally
3.  Document changes in comments and/or commit statements
4.  Add your name to the license. If the year exists in the license, append your name to the last person
5.  Create pull request
6.  Maintainer (currently just me) will comment and/or request changes
7.  Feel good and enjoy the kudos tokens (note> this is not a type of crypto-currency)

## Plans

I intend to create a library of collectors. Right now these are somewhat limited. Please create a collector using any of the existing ones as a base pattern and push back here!

## License

MIT

## Disclaimer

I will support this application best effort through GitHub issues and pull-requests. This application is not tied to Juniper Networks in anyway and as such, all risk is your own.