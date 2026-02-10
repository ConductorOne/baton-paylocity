![Baton Logo](./docs/images/baton-logo.png)

# `baton-paylocity` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-paylocity.svg)](https://pkg.go.dev/github.com/conductorone/baton-paylocity) ![verify](https://github.com/conductorone/baton-paylocity/actions/workflows/verify.yaml/badge.svg)

`baton-paylocity` is a connector for [Paylocity](https://www.paylocity.com/) built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Prerequisites
To use this connector, you will need the following information from your Paylocity account:

* Your Paylocity **Client ID**, specified with the `--paylocity-client-id` flag.
* Your Paylocity **Client Secret**, specified with the `--paylocity-client-secret` flag.
* The **Base URL** of your Paylocity API environment, indicated by the `--paylocity-base-url` flag.
* The **Company ID** for the specific entity you want to sync, specified with the `--paylocity-company-id` flag.

For example, you would connect like this:

    baton-paylocity \
      --paylocity-client-id your_client_id \
      --paylocity-client-secret your_client_secret \
      --paylocity-base-url "https://dc1demogwext.paylocity.com" \
      --paylocity-company-id your_company_id

## Connector capabilities

* Syncs employees as **Users**.
* Syncs position codes as **Roles**.
* Creates grants assigning each user to their corresponding position role.

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-paylocity
baton-paylocity
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_DOMAIN_URL=domain_url -e BATON_API_KEY=apiKey -e BATON_USERNAME=username ghcr.io/conductorone/baton-paylocity:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-paylocity/cmd/baton-paylocity@main

baton-paylocity

baton resources
```

# Data Model

`baton-paylocity` will pull down information about the following resources:
- Users (from Employees)
- Roles (from Position Codes)

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-paylocity` Command Line Usage

```
baton-paylocity

Usage:
  baton-paylocity [flags]
  baton-paylocity [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --paylocity-base-url string      The base API Gateway URL for the Paylocity environment. e.g., Testing: https://dc1demogwext.paylocity.com or Production: https://dc1prodgwext.paylocity.com
      --paylocity-client-id string     Client ID for the Paylocity API
      --paylocity-client-secret string Client Secret for the Paylocity API
      --paylocity-company-id string    The ID of the specific company in Paylocity to fetch data from.
      --client-id string               The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string           The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                    The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                           help for baton-paylocity
      --log-format string              The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string               The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -p, --provisioning                   This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --ticketing                      This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                        version for baton-paylocity

Use "baton-paylocity [command] --help" for more information about a command.
```
