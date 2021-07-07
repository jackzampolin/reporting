# `reporting`

This is a tool for reporting on validator income for stargate cosmos based chains. It is a cobra CLI command and can be run the following way:

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
