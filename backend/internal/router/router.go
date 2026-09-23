package router

import (
	"net/http"

	"inktype-backend/internal/handler"
	"inktype-backend/internal/middleware"
	"inktype-backend/internal/server"
	"inktype-backend/internal/service"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func NewRouter(s *server.Server, h *handler.Handlers, services *service.Services) *echo.Echo {
	middlewares := middleware.NewMiddlewares(s)

	router := echo.New()

	router.HTTPErrorHandler = middlewares.Global.GlobalErrorHandler

	// global middlewares
	router.Use(
		echoMiddleware.RateLimiterWithConfig(echoMiddleware.RateLimiterConfig{
			Store: echoMiddleware.NewRateLimiterMemoryStore(rate.Limit(20)),
			DenyHandler: func(c echo.Context, identifier string, err error) error {
				// Record rate limit hit metrics
				if rateLimitMiddleware := middlewares.RateLimit; rateLimitMiddleware != nil {
					rateLimitMiddleware.RecordRateLimitHit(c.Path())
				}

				s.Logger.Warn().
					Str("request_id", middleware.GetRequestID(c)).
					Str("identifier", identifier).
					Str("path", c.Path()).
					Str("method", c.Request().Method).
					Str("ip", c.RealIP()).
					Msg("rate limit exceeded")

				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			},
		}),
		middlewares.Global.CORS(),
		middlewares.Global.Secure(),
		middleware.RequestID(),
		middlewares.Tracing.NewRelicMiddleware(),
		middlewares.Tracing.EnhanceTracing(),
		middlewares.ContextEnhancer.EnhanceContext(),
		middlewares.Global.RequestLogger(),
		middlewares.Global.Recover(),
	)

	// register system routes
	registerSystemRoutes(router, h)

	// register versioned routes
	v1 := router.Group("/api/v1")

	// public routes
	v1.PUT("/local-upload", h.Document.LocalUpload)

	// auth required routes
	v1.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("Authorization") == "" {
				token := c.QueryParam("token")
				if token != "" {
					c.Request().Header.Set("Authorization", "Bearer "+token)
				}
			}
			return next(c)
		}
	})
	v1.Use(middlewares.Auth.RequireAuth)
	{
		v1.GET("/me", h.Document.Me)
		v1.POST("/documents", h.Document.Create)
		v1.GET("/documents", h.Document.List)
		v1.GET("/documents/:id", h.Document.Get)
		v1.GET("/documents/:id/download", h.Document.Download)
		v1.GET("/documents/:id/ws", h.Document.WS)
		v1.DELETE("/documents/:id", h.Document.Delete)
		v1.POST("/documents/:id/convert", h.Document.Convert)
	}

	return router
}
