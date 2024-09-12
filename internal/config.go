package internal

import "github.com/spf13/viper"

var (
	CamundaClientId     string
	CamundaClientSecret string
	CamundaGatewayURL   string
)

// InitViper initializes the viper configuration
func InitViper() {
	CamundaClientId = "client-id"
	CamundaClientSecret = "client"
	CamundaGatewayURL = "gateway-url"
	viper.SetConfigFile("./config/config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	CamundaClientSecret = viper.GetString("camunda.clientSecret")
	CamundaClientId = viper.GetString("camunda.clientId")
	CamundaGatewayURL = viper.GetString("camunda.gatewayAddress")
}
