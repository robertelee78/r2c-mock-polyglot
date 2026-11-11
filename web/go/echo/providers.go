package echoweb

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

var sessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func CsrfRouter() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CSRF())
	e.POST("/entities", EchoInstalledCsrfControl)
	e.GET("/session/options/positive", EchoSessionOptionsPositive)
	e.GET("/session/options/control", EchoSessionOptionsControl)
	return e
}

func requestDatabase(c echo.Context) *gorm.DB {
	return database.WithContext(c.Request().Context())
}

func EchoInstalledCsrfControl(c echo.Context) error {
	if _, err := c.Cookie("auth"); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	database := requestDatabase(c)
	entity := Entity{Name: c.FormValue("name")}
	database.Save(&entity)
	return c.NoContent(http.StatusNoContent)
}

func EchoSessionOptionsPositive(c echo.Context) error {
	session, err := sessionStore.Get(c.Request(), "auth")
	if err != nil {
		return err
	}
	session.Options = &sessions.Options{Secure: false, HttpOnly: true}
	return session.Save(c.Request(), c.Response())
}

func EchoSessionOptionsControl(c echo.Context) error {
	session, err := sessionStore.Get(c.Request(), "auth")
	if err != nil {
		return err
	}
	session.Options = &sessions.Options{Secure: true, HttpOnly: true}
	return session.Save(c.Request(), c.Response())
}

func EchoCorsPositiveRouter() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowCredentials: true, AllowOriginFunc: func(origin string) (bool, error) { return true, nil }}))
	return e
}

func EchoCorsControlRouter() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowCredentials: true, AllowOrigins: []string{"https://trusted.example"}}))
	return e
}

func EchoCorsUncertainRouter() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowCredentials: true, AllowOriginFunc: func(origin string) (bool, error) { return opaqueValue(origin) != "", nil }}))
	return e
}
