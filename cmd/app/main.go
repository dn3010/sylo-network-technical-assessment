package main

import (
	"context"
	"net/http"
	"os"

	"github.com/dn3010/sylo-network-technical-assessment/geth"
	"github.com/dn3010/sylo-network-technical-assessment/logger"
	"github.com/dn3010/sylo-network-technical-assessment/pkg/assets"
	"github.com/gin-gonic/gin"
)

func main() {
	mainCtx := context.Background()
	zeroLogger := logger.InitializeLogger()

	ethClient, err := geth.InitializeEthClient(mainCtx, zeroLogger)
	if err != nil {
		zeroLogger.Error().Err(err).Msg("Unable to connect to Goerli network")
		os.Exit(1)
	}

	api := gin.New()
	api.GET("/state", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	unauthedRoutes := api.Group("/")
	//middleware to set logger into request context
	unauthedRoutes.Use(logger.LoggerMW(mainCtx, zeroLogger))
	//middleware to set geth client into request context
	unauthedRoutes.Use(geth.MW(mainCtx, ethClient))
	//assets specific routes
	assets.RegisterAssetsAPI(unauthedRoutes)

	api.Run("0.0.0.0:3333") // listen and serve on 0.0.0.0:3333
}
