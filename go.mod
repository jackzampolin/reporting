module github.com/jackzampolin/reporting

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.42.6
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/tendermint/tendermint v0.34.11
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

replace github.com/gogo/grpc => google.golang.org/grpc v1.33.2

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
