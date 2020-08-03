package api

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber"
	"github.com/gotenberg/gotenberg/v7/pkg/core"
	flag "github.com/spf13/pflag"
)

func init() {
	core.MustRegisterModule(API{})
}

type API struct {
	Port           int
	BodyLimit      int
	ReadTimeout    float64
	ProcessTimeout float64
	WriteTimeout   float64
	RootPath       string

	srv *fiber.App
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
			fs.String("api-root-path", "/", "Set the root path of the API (useful for service discovery via URL paths)")

			return fs
		}(),
		New: func() core.Module { return new(API) },
	}
}

func (a *API) Provision(flags core.ParsedFlags) error {
	a.Port = flags.MustInt("api-port")
	a.BodyLimit = flags.MustInt("api-body-limit")
	a.ReadTimeout = flags.MustFloat("api-read-timeout")
	a.ProcessTimeout = flags.MustFloat("api-process-timeout")
	a.WriteTimeout = flags.MustFloat("api-write-timeout")
	a.RootPath = strings.ToLower(flags.MustString("api-root-path"))

	// As the root path must begin and end with a slash,
	// we correct its value if necessary.
	if !strings.HasPrefix(a.RootPath, "/") {
		a.RootPath = fmt.Sprintf("/%s", a.RootPath)
	}

	if !strings.HasSuffix(a.RootPath, "/") {
		a.RootPath = fmt.Sprintf("%s/", a.RootPath)
	}

	return nil
}

func (a *API) Validate() error {
	var errs core.ErrorArray

	// TODO: migrate common validations to a dedicated package.
	if a.Port < 1 || a.Port > 65535 {
		errs = append(errs,
			errors.New("port must be more than 1 and less than 65535"),
		)
	}

	if a.BodyLimit <= 0 {
		errs = append(errs,
			errors.New("body limit must be more than 0"),
		)
	}

	if a.ReadTimeout <= 0 {
		errs = append(errs,
			errors.New("read timeout must be more than 0"),
		)
	}

	if a.ProcessTimeout <= 0 {
		errs = append(errs,
			errors.New("process timeout must be more than 0"),
		)
	}

	if a.WriteTimeout <= 0 {
		errs = append(errs,
			errors.New("write timeout must be more than 0"),
		)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (a *API) Start() error {
	// TODO: try prefork setting.
	a.srv = fiber.New(&fiber.Settings{
		DisableStartupMessage: false, // TODO: remove.
		//DisableStartupMessage: true,
		BodyLimit:    a.BodyLimit * 1024 * 1024, // TODO: bytes?
		ReadTimeout:  time.Duration(a.ReadTimeout*1000) * time.Millisecond,
		WriteTimeout: time.Duration(a.WriteTimeout*1000) * time.Millisecond,
	})

	// Add routes from other modules.
	routesMu.Lock()
	defer routesMu.Unlock()

	for _, route := range routes {
		a.srv.Add(
			route.Method,
			fmt.Sprintf("%s%s", a.RootPath, route.Path),
			route.Handlers...,
		)
	}

	core.Log().Error("error message")
	core.Log().Warn("warn message")
	core.Log().Info("info message")
	core.Log().Debug("debug message")

	return a.srv.Listen(a.Port)
}

func (a *API) Stop() error {
	return a.srv.Shutdown()
}

type Route struct {
	Path     string
	Method   string
	Handlers []fiber.Handler
}

func RegisterRoute(route *Route) error {
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

	routesMu.Lock()
	defer routesMu.Unlock()

	if _, ok := routes[route.Path]; ok {
		return fmt.Errorf("route '%s' is already registered", route.Path)
	}

	routes[route.Path] = route

	return nil
}

var (
	routes   = make(map[string]*Route)
	routesMu sync.RWMutex
)

// Interface guards.
var (
	_ core.Module      = (*API)(nil)
	_ core.Provisioner = (*API)(nil)
	_ core.Validator   = (*API)(nil)
	_ core.App         = (*API)(nil)
)
