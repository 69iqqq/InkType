package middleware

import (
	"encoding/json"
	"net/http"

	"inktype-backend/internal/errs"
	"inktype-backend/internal/server"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AuthMiddleware struct {
	server *server.Server
}

func NewAuthMiddleware(s *server.Server) *AuthMiddleware {
	return &AuthMiddleware{
		server: s,
	}
}

func (auth *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.WrapMiddleware(
		clerkhttp.WithHeaderAuthorization(
			clerkhttp.AuthorizationFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				response := map[string]interface{}{
					"code":    "UNAUTHORIZED",
					"message": "Unauthorized",
					"status":  401,
				}

				if err := json.NewEncoder(w).Encode(response); err != nil {
					auth.server.Logger.Error().Err(err).Msg("failed to write JSON response")
				}
				auth.server.Logger.Warn().Msg("unauthorized request rejected")
			}))))(func(c echo.Context) error {
		claims, ok := clerk.SessionClaimsFromContext(c.Request().Context())

		if !ok {
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		c.Set("user_id", claims.Subject)

		if claims.ActiveOrganizationRole != "" {
			c.Set("user_role", claims.ActiveOrganizationRole)
		}

		// Enrich the existing context logger
		log := GetLogger(c)
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("user_id", claims.Subject)
		})

		return next(c)
	})
}
