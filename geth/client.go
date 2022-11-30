package geth

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const goerli_endpoint = "https://goerli.infura.io/v3/ad0b336d7f2f4082b5a624e50d27df5c"

func InitializeEthClient(ctx context.Context, log *zerolog.Logger) (*ethclient.Client, error) {
	ethClient, err := ethclient.Dial(goerli_endpoint)
	if err != nil {
		return nil, err
	}

	chainId, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Connected to Goerli - Chain ID: %v", chainId)
	return ethClient, nil
}

type contextKey struct{}

var clientKey = &contextKey{}

func MW(ctx context.Context, client *ethclient.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, clientKey, client)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func From(ctx context.Context) *ethclient.Client {
	service := ctx.Value(clientKey)
	client, ok := service.(*ethclient.Client)
	if client == nil || !ok {
		return nil
	}
	return client
}
