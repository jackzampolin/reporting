# `reporting`

This is a tool for reporting on validator income for stargate cosmos based chains. It is a cobra CLI command and can be run the following way:

## Installing:

1. clone the repo: `git clone https://github.com/jackzampolin/reporting && cd reporting`
2. copy `.reporting.yaml` to `~/.reporting.yaml` and modify (see suggestion below)

```bash
$ go run main.go
reporting for stargate enabled cosmos chains

Usage:
  reporting [command]

Available Commands:
  day-blocks       get a list of blocks that happened close to midnight local time from start block to now
  help             Help about any command
  validator-report outputs a csv of the data required for validator income reporting

Flags:
      --config string   config file (default is $HOME/.reporting.yaml)
  -h, --help            help for reporting
  -t, --toggle          Help message for toggle

Use "reporting [command] --help" for more information about a command.
```

Currently there is a large amount of hardcoding in this. The configuration for each individual chain needs to be put in the `~/.reporting.yaml` file and read in prior to commands
This interface will be similar to the go relayer and will just add a chain-id to each command. Arguements also need to be added for address. 


```yaml
# ~/.reporting.yaml
networks:
  akashnet-2:
    chain-id: "akashnet-2"
    archive: "https://akash.technofractal.com:443"
    prefix: "akash"
    token: "uakt"
    coin-gecko-id: "akash-network"
  cosmoshub-4:
    chain-id: "cosmoshub-4"
    archive: "https://hub.technofractal.com:443"
    prefix: "cosmos"
    token: "uatom"
    coin-gecko-id: "atom"
  osmosis-1:
    chain-id: "osmosis-1"
    archive: "https://osmosis.technofractal.com:443"
    prefix: "osmo"
    token: "uosmo"
    coin-gecko-id: "osmosis" 
  sentinelhub-2:
    chain-id: "sentinelhub-2"
    archive: "https://dvpn.technofractal.com:443"
    prefix: "sent"
    token: "udvpn"
    coin-gecko-id: "sentinel"
```