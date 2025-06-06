package server

import (
	"net/http"
	"net/url"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/core/im"
	"github.com/mislu/market-api/internal/core/mq/rabbit"
	"github.com/mislu/market-api/internal/core/payment"
	"github.com/mislu/market-api/internal/core/recommend"
	resourcemanager "github.com/mislu/market-api/internal/core/resource_manager"
	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/es"
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
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent) // 204
			return
		}

		ctx.Next()
	})

	engine.Use(func(ctx *gin.Context) {
		if ctx.Writer.Status() == http.StatusNotFound {
			return
		}

		now := time.Now()
		decodedUrl, _ := url.QueryUnescape(ctx.Request.RequestURI)

		defer func() {
			var (
				abortErr error
				apiError exceptions.APIError
				httpCode int
				msg      string
			)

			if err := recover(); err != nil {
				stackInfo := string(debug.Stack())
				logger.Error("panic recovered",
					zap.String("path", decodedUrl),
					zap.Any("error", err),
					zap.String("stack", stackInfo))
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
				switch controllers.GetContentType(ctx) {
				case controllers.ResponseTypeFile:
					ctx.File(resp.Data.(string))
				case controllers.ResponseTypeStream:
					// TODO
				case controllers.ResponseTypeJSON:
					ctx.JSON(http.StatusOK, resp)
				}

				httpCode = http.StatusOK
				msg = "success"
			}

			cost := time.Since(now).Seconds()
			logFields := []zap.Field{
				zap.String("method", ctx.Request.Method),
				zap.String("path", decodedUrl),
				zap.Int("code", httpCode),
				zap.String("msg", msg),
				zap.Float64("cost", cost),
			}

			if httpCode >= 400 || ctx.IsAborted() {
				logFields = append(logFields, zap.Error(abortErr))
				logger.Error("request error", logFields...)
			} else {
				logger.Info("request success", logFields...)
			}
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

	dbLogger := log.NewLogger()
	// init db
	db.Init(dbLogger)

	es.Init()

	im.Init()
	// init resource manager
	resourcemanager.InitGlobalResourceManager()

	err := rabbit.InitGlobalRabbitMQ()
	if err != nil {
		panic(err)
	}

	recommend.InitGlobalWorker()
	payment.InitPaymentService()
	// init gin
	server := newServer(logger)

	// init routers
	server.serve()

	select {}
}
