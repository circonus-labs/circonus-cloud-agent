module github.com/circonus-labs/circonus-cloud-agent

go 1.12

require (
	cloud.google.com/go v0.34.0
	contrib.go.opencensus.io/exporter/ocagent v0.4.7 // indirect
	github.com/Azure/azure-sdk-for-go v26.4.0+incompatible
	github.com/Azure/go-autorest v11.5.2+incompatible
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf
	github.com/aws/aws-sdk-go v1.18.4
	github.com/circonus-labs/go-apiclient v0.6.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/golang/protobuf v1.3.0
	github.com/googleapis/gax-go v2.0.2+incompatible // indirect
	github.com/pelletier/go-toml v1.4.0
	github.com/pkg/errors v0.8.1
	github.com/rs/zerolog v1.13.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2
	golang.org/x/oauth2 v0.0.0-20190226205417-e64efc72b421
	golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys v0.0.0-20190318195719-6c81ef8f67ca
	google.golang.org/api v0.2.0
	google.golang.org/genproto v0.0.0-20190307195333-5fe7a883aa19
	gopkg.in/yaml.v2 v2.2.2
)
