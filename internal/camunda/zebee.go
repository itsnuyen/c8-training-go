package camunda

import (
	"fmt"
	"os"

	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

type Client struct {
	Client *zbc.Client
}

type ClientInterface interface {
	InitZeebeClient() (*Client, error)
}

func (c *Client) InitZeebeClient(clientId string, clientSecret, clientUrl string) (*Client, error) {
	var config zbc.ClientConfig
	if os.Getenv("env") != "prod" {
		fmt.Println("Using local dev")
		config = zbc.ClientConfig{UsePlaintextConnection: true, GatewayAddress: "localhost:26500"}
	} else {
		credentials, err := zbc.NewOAuthCredentialsProvider(&zbc.OAuthProviderConfig{
			ClientID:               clientId,
			ClientSecret:           clientSecret,
			AuthorizationServerURL: "https://login.cloud.camunda.io/oauth/token",
			Audience:               "zeebe.camunda.io",
		})

		if err != nil {
			panic(err)
		}
		config = zbc.ClientConfig{
			GatewayAddress:      clientUrl,
			CredentialsProvider: credentials,
		}
	}
	client, err := zbc.NewClient(&config)
	if err != nil {
		panic(err)
	}
	c.Client = &client
	return c, nil
}

func (c *Client) CloseZeebeClient() {
	client := *c.Client
	client.Close()
}
