package api

import (
	"fmt"
	"net/http"

	"github.com/devanand100/gym/api/auth"
	"github.com/devanand100/gym/internal/util"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(server *APIService, next echo.HandlerFunc, secret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path

		if server.defaultAuthSkipper(c) {
			return next(c)
		}

		if util.HasPrefixes(path, "/api/user") {
			return next(c)
		}

		cookie, err := c.Cookie(auth.CookieName)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing access token")
		}
		claims, err := auth.VerifyToken(cookie.Value)
		fmt.Println("claims.....", claims)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token expired")
		}
		Id := claims["Id"].(string)
		c.Set("UserId", Id)

		return next(c)
	}
}

func (*APIService) defaultAuthSkipper(c echo.Context) bool {
	path := c.Path()
	return util.HasPrefixes(path, "/api/v1/me")
}
