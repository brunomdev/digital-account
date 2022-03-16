package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nobuyo/nrfiber"
)

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func FiberMiddleware(a *fiber.App, newRelic *newrelic.Application, appDebug bool) {
	a.Use(
		recover.New(),
		requestid.New(),
		cors.New(),
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),
		nrfiber.New(nrfiber.Config{
			NewRelicApp: newRelic,
		}),
		NewLog(appDebug),
	)
}
