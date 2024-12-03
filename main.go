package main

import (
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	enforcer, err := casbin.NewEnforcer("model.conf", "policy.csv")
	if err != nil {
		log.Fatalf("NewEnforcer: %s", err)
	}

	e.Use(authenticateMiddleware(enforcer))

	e.GET("/news", func(c echo.Context) error {
		jsonMap := map[string]string{
			"article": "hoge",
		}
		return c.JSON(http.StatusOK, jsonMap)
	})

	e.POST("/news", func(c echo.Context) error {
		return c.JSON(http.StatusNoContent, "")
	})

	e.Logger.Fatal(e.Start("0.0.0.0:3000"))
}

func authenticateMiddleware(enforcer *casbin.Enforcer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			sub := e.Request().Header.Get("user")
			method := e.Request().Method
			path := e.Request().URL.Path

			ok, err := enforcer.Enforce(sub, path, method)
			if err != nil {
				log.Fatalf("enforce error: %s", err)
			}
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "permission denined")
			}

			return next(e)
		}
	}
}
