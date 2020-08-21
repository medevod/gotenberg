package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/middleware"
	"strings"
	"time"

	"github.com/gofiber/fiber"
	flag "github.com/spf13/pflag"
	"github.com/thecodingmachine/gotenberg/v7/pkg/core"
)

func init() {
	core.MustRegisterModule(API{})
}

type API struct {
	port           int
	bodyLimit      int
	readTimeout    float64
	processTimeout float64
	writeTimeout   float64
	rootPath       string

	routers []Router
	srv     *fiber.App
}

func (API) Descriptor() core.ModuleDescriptor {
	return core.ModuleDescriptor{
		ID: "API",
		FlagSet: func() *flag.FlagSet {
			fs := flag.NewFlagSet("api", flag.ExitOnError)
			fs.Int("api-port", 3000, "Set the port on which the API should listen")
			fs.Int("api-body-limit", 16, "Set the maximum allowed size in MB for a request body")
			fs.Float64("api-read-timeout", 10, "Set the maximum duration in seconds allowed to read the full request, including body")
			fs.Float64("api-process-timeout", 10, "Set the maximum duration in seconds allowed to process a request")
			fs.Float64("api-write-timeout", 10, "Set the maximum duration in seconds before timing out writes of the response")
			fs.String("api-root-path", "/", "Set the root path of the API - useful for service discovery via URL paths")

			return fs
		}(),
		New: func() core.Module { return new(API) },
	}
}

func (a *API) Provision(ctx *core.Context) error {
	flags := ctx.ParsedFlags()

	a.port = flags.MustInt("api-port")
	a.bodyLimit = flags.MustInt("api-body-limit")
	a.readTimeout = flags.MustFloat("api-read-timeout")
	a.processTimeout = flags.MustFloat("api-process-timeout")
	a.writeTimeout = flags.MustFloat("api-write-timeout")
	a.rootPath = strings.ToLower(flags.MustString("api-root-path"))

	// As the root path must begin and end with a slash,
	// we correct its value if necessary.
	if !strings.HasPrefix(a.rootPath, "/") {
		a.rootPath = fmt.Sprintf("/%s", a.rootPath)
	}

	if !strings.HasSuffix(a.rootPath, "/") {
		a.rootPath = fmt.Sprintf("%s/", a.rootPath)
	}

	routers, err := ctx.Modules(new(Router))
	if err != nil {
		return err
	}

	for _, router := range routers {
		a.routers = append(a.routers, router.(Router))
	}

	return nil
}

func (a *API) Validate() error {
	var errs core.ErrorArray

	// TODO: migrate common validations to a dedicated package.
	if a.port < 1 || a.port > 65535 {
		errs = append(errs,
			errors.New("port must be more than 1 and less than 65535"),
		)
	}

	if a.bodyLimit <= 0 {
		errs = append(errs,
			errors.New("body limit must be more than 0"),
		)
	}

	if a.readTimeout <= 0 {
		errs = append(errs,
			errors.New("read timeout must be more than 0"),
		)
	}

	if a.processTimeout <= 0 {
		errs = append(errs,
			errors.New("process timeout must be more than 0"),
		)
	}

	if a.writeTimeout <= 0 {
		errs = append(errs,
			errors.New("write timeout must be more than 0"),
		)
	}

	if err := a.validateRouters(); err != nil {
		errs = append(errs, err)
	}

	// TODO required?
	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (a *API) validateRouters() error {
	routesMap := make(map[string]*Route)

	for _, router := range a.routers {
		routes := router.Routes()

		for _, route := range routes {
			if route.Path == "" {
				return errors.New("route with empty path cannot be registered")
			}

			// Fiber panics if a route is using an
			// invalid method. Therefore, we check
			// first that the given method is valid.
			validMethods := []string{
				fiber.MethodConnect,
				fiber.MethodDelete,
				fiber.MethodGet,
				fiber.MethodHead,
				fiber.MethodOptions,
				fiber.MethodPatch,
				fiber.MethodPost,
				fiber.MethodPut,
				fiber.MethodTrace,
			}

			isValidMethod := false
			for _, method := range validMethods {
				if method == route.Method {
					isValidMethod = true
					break
				}
			}

			if !isValidMethod {
				return fmt.Errorf("method '%s' from route '%s' is invalid", route.Method, route.Path)
			}

			if _, ok := routesMap[route.Path]; ok {
				return fmt.Errorf("route '%s' is already registered", route.Path)
			}

			routesMap[route.Path] = route
		}
	}

	return nil
}

func (a *API) Start() error {
	// TODO: try prefork setting.
	a.srv = fiber.New(&fiber.Settings{
		DisableStartupMessage: false, // TODO: remove.
		//DisableStartupMessage: true,
		BodyLimit:    a.bodyLimit * 1024 * 1024, // TODO: bytes?
		ReadTimeout:  time.Duration(a.readTimeout*1000) * time.Millisecond,
		WriteTimeout: time.Duration(a.writeTimeout*1000) * time.Millisecond,
	})

	// Add routes from other modules.
	for _, router := range a.routers {
		for _, route := range router.Routes() {
			a.srv.Add(
				route.Method,
				fmt.Sprintf("%s%s", a.rootPath, route.Path),
				route.Handlers...,
			)
		}
	}

	// TODO custom middleware.
	a.srv.Use(middleware.Logger())

	core.Log().Error("error message")
	core.Log().Warn("warn message")
	core.Log().Info("info message")
	core.Log().Debug("debug message")

	return a.srv.Listen(a.port)
}

func (a *API) Stop() error {
	return a.srv.Shutdown()
}

type Router interface {
	Routes() []*Route
}

type Route struct {
	Path     string
	Method   string
	Handlers []fiber.Handler
}

// Interface guards.
var (
	_ core.Module      = (*API)(nil)
	_ core.Provisioner = (*API)(nil)
	_ core.Validator   = (*API)(nil)
	_ core.App         = (*API)(nil)
)
