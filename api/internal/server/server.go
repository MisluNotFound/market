package server

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	resourcemanager "github.com/mislu/market-api/internal/core/resource_manager"
	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/server/controllers"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/utils/log"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
}

func newServer(logger *zap.Logger) *Server {
	engine := gin.New()

	engine.Use(func(ctx *gin.Context) {
		if ctx.Writer.Status() == http.StatusNotFound {
			return
		}

		now := time.Now()
		defer func() {
			var (
				abortErr error
				apiError exceptions.APIError
				httpCode int
				msg      string
			)

			if err := recover(); err != nil {
				stackInfo := string(debug.Stack())
				logger.Error("got panic", zap.String("panic", fmt.Sprintf("%+v", err)), zap.String("stack", stackInfo))

				ctx.AbortWithStatus(http.StatusInternalServerError)
			}

			if ctx.IsAborted() {
				for i := range ctx.Errors {
					multierr.AppendInto(&abortErr, ctx.Errors[i].Err)
				}

				if apiError = controllers.GetAbortError(ctx); apiError != nil {
					multierr.AppendInto(&abortErr, apiError)
					httpCode = apiError.ToResponse().Code
					msg = apiError.ToResponse().Msg
				} else {
					httpCode = http.StatusInternalServerError
					msg = "Internal server error"
				}

				ctx.JSON(httpCode, apiError.ToResponse())
			} else {
				resp := controllers.GetPayLoad(ctx)
				ctx.JSON(http.StatusOK, resp)
			}

			decodedUrl, _ := url.QueryUnescape(ctx.Request.RequestURI)
			cost := time.Since(now).Seconds()

			logger.Error(
				"result",
				zap.Any("method", ctx.Request.Method),
				zap.Any("path", decodedUrl),
				zap.Any("code", httpCode),
				zap.String("msg", msg),
				zap.Any("cost", cost),
				zap.Error(abortErr),
			)
		}()

		ctx.Next()
	})

	return &Server{
		engine: engine,
	}
}

func Run() {
	// init logger
	logger := log.NewLogger()

	// init db
	db.Init(logger)

	// init resource manager
	resourcemanager.InitGlobalResourceManager()

	// init gin
	server := newServer(logger)

	// init routers
	server.serve()

	select{}
}
