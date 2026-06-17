package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a Application) GetFeed(c *gin.Context) {
	userIDAny, ok := c.Get("userid")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDAny.(uint)

	ids, err := a.feedRepo.GetFeedForUser(c.Request.Context(), userID)
	if err != nil || len(ids) == 0 {
		if err := a.feedRepo.GenerateFeedForUser(c.Request.Context(), userID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ids, err = a.feedRepo.GetFeedForUser(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"ids": ids})
}