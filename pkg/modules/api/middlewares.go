package api

import (
	"github.com/gofiber/fiber"
)

func MultipartFormDataEnforcer(ctx *fiber.Ctx) {
	if ctx.Accepts(fiber.MIMEMultipartForm) != "" {
		ctx.Next()
	}

	ctx.Next(fiber.ErrUnsupportedMediaType)
}
