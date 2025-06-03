package controllers

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/core/recommend"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/utils/lib"
)

const (
	_ctx_user_id = "userID"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")
		if len(tokenHeader) == 0 {
			AbortWithError(c, exceptions.NewGenericError(http.StatusUnauthorized, "No token provided", errors.New("no token provided")))
			c.Abort()
			return
		}

		token := tokenHeader[len("Bearer "):]
		if len(token) == 0 {
			AbortWithError(c, exceptions.NewGenericError(http.StatusUnauthorized, "No token provided", errors.New("no token provided")))
			c.Abort()
			return
		}

		claims, err := lib.VerifyToken(token, true)
		if err != nil {
			AbortWithError(c, exceptions.NewGenericError(http.StatusUnauthorized, "Invalid token", errors.New("invalid token")))
			c.Abort()
			return
		}

		c.Set(_ctx_user_id, claims.UserID)
		c.Next()
	}
}

var productRegex = regexp.MustCompile(`^/api/product/([^/]+)/([^/]+)/?$`)
var orderRegex = regexp.MustCompile(`^/api/order/([^/]+)/([^/]+)/?$`)

func ScrapMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.RequestURI
		userID, exists := GetContextUserID(c)
		if !exists {
			c.Next()
			return
		}

		if len(userID) == 0 {
			c.Next()
			return
		}

		switch c.Request.Method {
		case http.MethodGet:
			if productRegex.MatchString(uri) {
				feedBacks := []recommend.Feedback{{
					UserId:       userID,
					ItemId:       productRegex.FindStringSubmatch(uri)[2],
					FeedbackType: "view",
				}}

				recommend.GlobalWorker.InsertFeedback(context.Background(), feedBacks)
			}
		case http.MethodPost:
			if productRegex.MatchString(uri) {
				feedBacks := []recommend.Feedback{{
					UserId:       userID,
					ItemId:       productRegex.FindStringSubmatch(uri)[2],
					FeedbackType: "like",
				}}

				recommend.GlobalWorker.InsertFeedback(context.Background(), feedBacks)
			} else if orderRegex.MatchString(uri) {
				feedBacks := []recommend.Feedback{{
					UserId:       userID,
					ItemId:       orderRegex.FindStringSubmatch(uri)[2],
					FeedbackType: "purchase",
				}}

				recommend.GlobalWorker.InsertFeedback(context.Background(), feedBacks)
			}
		}

		c.Next()
	}
}

func GetContextUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(_ctx_user_id)
	if !exists {
		return "", false
	}
	return userID.(string), true
}
