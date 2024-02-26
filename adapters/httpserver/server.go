package httpserver

import (
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"hexagon/domain/book"
	"hexagon/pkg/config"
	"hexagon/pkg/logger"
	"hexagon/pkg/sentry"
	"net/http"
	"strings"
)

type Options func(s *Server) error

type Server struct {
	router *echo.Echo
	Config *config.Config
	Logger *zap.SugaredLogger

	// storage adapters
	BookStore book.Storage
}

func New(options ...Options) (*Server, error) {
	s := Server{
		router: echo.New(),
		Config: config.Empty,
		Logger: logger.NOOPLogger,
	}

	for _, fn := range options {
		if err := fn(&s); err != nil {
			return nil, err
		}
	}

	s.RegisterGlobalMiddlewares()

	s.RegisterHealthCheck(s.router.Group(""))
	s.RegisterBookRoutes(s.router.Group("/api/books"))

	return &s, nil
}

func (s *Server) RegisterGlobalMiddlewares() {
	s.router.Use(middleware.Recover())
	s.router.Use(middleware.Secure())
	s.router.Use(middleware.RequestID())
	s.router.Use(middleware.Gzip())
	s.router.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	// CORS
	if s.Config.AllowOrigins != "" {
		aos := strings.Split(s.Config.AllowOrigins, ",")
		s.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: aos,
		}))
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) RegisterHealthCheck(router *echo.Group) {
	router.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK!!!")
	})
}

func (s *Server) handleError(c echo.Context, err error, status int) error {
	s.Logger.Errorw(
		err.Error(),
		zap.String("request_id", s.requestID(c)),
	)

	if status >= http.StatusInternalServerError {
		sentry.WithContext(c).Error(err)
	}

	return c.JSON(status, map[string]string{
		"message": http.StatusText(status),
	})
}

func (s *Server) requestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}
