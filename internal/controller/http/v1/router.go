package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang-diplom-work/internal/service"
	"log"
	"os"
	"path/filepath"
)

func NewRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","method":"${method}","uri":"${uri}","status":"${status}","error":"${error}"` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(200) })
	//r := &authRoutes{
	//	authService: services.Auth,
	//}
	//handler.POST("/sign-up", r.signUp)
	//handler.POST("/sign-in", r.signIn)
	auth := handler.Group("/auth")
	{
		newAuthRoutes(auth, services.Auth)
	}
	// fayllar bilan ishlashga keganimda bundan foydalanaman
	authMiddleware := &AuthMiddleware{services.Auth}
	v1 := handler.Group("/api/v1", authMiddleware.UserIdentity)
	{
		newFileRoutes(v1, services.Files)
	}

}

func setLogsFile() *os.File {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	// Define the path for the log file relative to the working directory
	logFilePath := filepath.Join(cwd, "logs/request.log")

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
