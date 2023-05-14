package configuration

type Config struct {
	ListenOn                 string
	LogLevel                 string
	SQLStoreConnectionString string
	EnableGraphql            bool
	GraphqlWebClientEnable   bool
	CORSAllowOrigins         []string
	KratosAdminUrl           string
}
