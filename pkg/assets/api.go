package assets

import (
	"net/http"

	"github.com/dn3010/sylo-network-technical-assessment/logger"
	"github.com/gin-gonic/gin"
)

func listTokens(c *gin.Context) {
	ctx := c.Request.Context()
	account := c.Query("account")
	if account == "" {
		logger.From(ctx).Error().Msg("No user account was passed-in the url")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user account was provided"})
		return
	}

	logger.From(ctx).Info().Msgf("User Account: %s", account)

	assets, err := queryTransferEventsLogs(ctx, account, contracts)
	if err != nil {
		logger.From(ctx).Error().Msgf("Problem occurred while querying transfer events logs for account %s", account)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No user account was provided"})
		return
	}

	c.JSON(http.StatusOK, assets)
}

func RegisterAssetsAPI(r *gin.RouterGroup) {
	r.GET("/assets", listTokens)
}
