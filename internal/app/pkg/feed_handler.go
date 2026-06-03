package app

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type FeedResponse struct {
    IDs []uint `json:"ids"`
}

func (a *Application) GetFeed(gCtx *gin.Context) {
    userIDAny, ok := gCtx.Get("userid")
    if !ok {
        gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    userID := userIDAny.(uint)

    ctx := gCtx.Request.Context()

    // если ленты нет — генерируем
    ids, err := a.feedRepo.GetFeedForUser(ctx, userID)
    if err != nil || len(ids) == 0 {
        if err := a.feedRepo.GenerateFeedForUser(ctx, userID); err != nil {
            gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        ids, err = a.feedRepo.GetFeedForUser(ctx, userID)
        if err != nil {
            gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    }

    gCtx.JSON(http.StatusOK, FeedResponse{IDs: ids})
}