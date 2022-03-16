package middleware

import (
	"fmt"
	"github.com/brunomdev/digital-account/infra/log"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// NewLog create middleware to add requestId to log context and logging the request
func NewLog(appDebug bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(log.LoggerKeyType,
			log.WithContext(c.Context()).With(
				zap.String("X-Request-ID", c.GetRespHeader(fiber.HeaderXRequestID)),
			),
		)

		if appDebug {
			log.Info(
				c.Context(),
				fmt.Sprintf("%s %s %d", c.Method(), c.Path(), c.Response().StatusCode()),
			)
		}

		return errors.Wrap(c.Next(), "NewLog")
	}
}
