# Gidari CLI

[![PkgGoDev](https://img.shields.io/badge/go.dev-docs-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/alpstable/gidari-cli)
![Build Status](https://github.com/alpstable/gidari-cli/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpstable/gidari-cli)](https://goreportcard.com/report/github.com/alpstable/gidari-cli)
[![Discord](https://img.shields.io/discord/987810353767403550)](https://discord.gg/3jGYQz74s7)

Gidari CLI is a "web-to-storage" tool for querying web APIs and persisting the resulting data onto local storage. A configuration file is used to define how this querying and storing should occur. Once you have a configuration file, you can initiate this transport using the command `gidari --config <configuration.yml>`. See [here](https://youtu.be/NgeOJ50IWhY) for a quick demonstration.

## Installation

```
go install github.com/alpstable/gidari-cli/cmd/gidari@latest
```

For information on using the Go library, see [here](https://github.com/alpstable/gidari).

## Usage

Using Gidari in command mode is a two step process:

1. Create a configuraiton file to instruct the binary on how to make the RESful HTTP requests and where to store the data
2. Run `gidari --config your_configuration.yml`

The `configuration.yml` file is used to define a set of rules for making RESTful HTTP requests and where to store the data. See [here](https://github.com/alpstable/gidari/tree/main/e2e/testdata/upsert) for example configurations.

### Configurations

|Key                             |Required|Type  |Description                                                                                                                                                                       |
|--------------------------------|--------|------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|url                             |T       |string|The API base URL                                                                                                                                                                  |
|authentication                  |F       |map   |Data required for authenticating the web API HTTP Requests                                                                                                                        |
|authentication.apiKey.passphrase|T       |string|                                                                                                                                                                                  |
|authentication.apiKey.Key       |T       |string|                                                                                                                                                                                  |
|authentication.apiKey.Secret    |T       |string|                                                                                                                                                                                  |
|authentication.auth2.Bearer     |T       |string|                                                                                                                                                                                  |
|storage.ConnectionString        |T       |String|String used to connect to a storage device                                                                                                                                        |
|storage.Database                |F       |String|Name of the database to run operations against. This is an optional field and will not be needed for every storage device. This field is currently needed for MongoDB.            |
|rateLimit                       |T       |map   |Data required for limiting the number of requests per second, avoiding 429 errors                                                                                                 |
|rateLimit.burst                 |T       |uint  |Number of requests that can be made per second                                                                                                                                    |
|rateLimit.period                |T       |uint  |Period for the rateLimit.burst                                                                                                                                                    |
|truncate                        |F       |bool  |Truncate all tables in the databse before performing upserts                                                                                                                      |
|requests                        |F       |list  |List of requests to receive data from the web API for upserting into storage                                                                                                      |
|request.endpoint                |T       |string|Endpoint for making the RESTful API request                                                                                                                                       |
|request.table                   |F       |string|Name of the table in the storage for upserting data. This field defaults to the last string in the endpoint path                                                                  |
|request.clobColumn              |F       |string|Name of the column where data will be stored if the response data is not a valid JSON. Logs a warning if invalid JSON is received and this field is not set, and no data is saved.|
|request.timseries               |F       |map   |Data required for upserting timeseries data, which are batched and can be resource intensive                                                                                      |
|request.timeseries.startName    |T       |string|Name of the query/path parameter for the "start" datetime of the timeseries                                                                                                       |
|request.timeseries.endName      |T       |string|Name of the query/path parameter for the "end" datetime of the timeseries                                                                                                         |
|request.timeseries.period       |T       |uint  |How often (in seconds) to build a new datetime range to batch.                                                                                                                    |
|request.timeseries.layout       |T       |string|The layout for how to build a datetime to query over (e.g. RFC3339 would be "2006-01-02T15:04:05Z07:00")                                                                          |
|request.query                   |N       |map   |A hash of data that holds the query parameters for a request                                                                                                                      |
