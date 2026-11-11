package echoweb

import (
	"html"
	"html/template"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Entity struct {
	ID   uint   `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
	Role string `json:"role" form:"role"`
}

var database *gorm.DB

var policy *casbin.Enforcer

func RegisterCases(e *echo.Echo) {
	e.GET("/html/positive", EchoHtmlPositive)
	e.GET("/html/control", EchoHtmlControl)
	e.GET("/html/uncertain", EchoHtmlUncertain)
	e.POST("/csrf/positive", EchoCsrfPositive)
	e.POST("/csrf/control", EchoCsrfControl)
	e.POST("/csrf/uncertain", EchoCsrfUncertain)
	e.GET("/redirect/positive", EchoRedirectPositive)
	e.GET("/redirect/control", EchoRedirectControl)
	e.GET("/redirect/uncertain", EchoRedirectUncertain)
	e.POST("/policy/positive", EchoPolicyPositive)
	e.POST("/policy/control", EchoPolicyControl)
	e.POST("/policy/uncertain", EchoPolicyUncertain)
	e.POST("/session/positive", EchoSessionPositive)
	e.POST("/session/control", EchoSessionControl)
	e.POST("/session/uncertain", EchoSessionUncertain)
	e.GET("/cors/positive", EchoCorsPositive)
	e.GET("/cors/control", EchoCorsControl)
	e.GET("/cors/uncertain", EchoCorsUncertain)
	e.POST("/binding/positive", EchoBindingPositive)
	e.POST("/binding/control", EchoBindingControl)
	e.POST("/binding/uncertain", EchoBindingUncertain)
}

func EchoHtmlPositive(c echo.Context) error {
	t := template.Must(template.New("page").Parse(`<main>{{.}}</main>`))
	return t.Execute(c.Response(), template.HTML(c.QueryParam("value")))
}

func EchoHtmlControl(c echo.Context) error {
	t := template.Must(template.New("page").Parse(`<main>{{.}}</main>`))
	return t.Execute(c.Response(), template.HTML(html.EscapeString(c.QueryParam("value"))))
}

func EchoHtmlUncertain(c echo.Context) error {
	t := template.Must(template.New("page").Parse(`<main>{{.}}</main>`))
	return t.Execute(c.Response(), template.HTML(opaqueValue(c.QueryParam("value"))))
}

func EchoCsrfPositive(c echo.Context) error {
	if _, err := c.Cookie("auth"); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	entity := Entity{Name: c.FormValue("name")}
	database.Save(&entity)
	return c.NoContent(http.StatusNoContent)
}

func EchoCsrfControl(c echo.Context) error {
	if c.Request().Header.Get("Authorization") == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	entity := Entity{Name: c.FormValue("name")}
	database.Save(&entity)
	return c.NoContent(http.StatusNoContent)
}

func EchoCsrfUncertain(c echo.Context) error {
	if _, err := c.Cookie("auth"); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	if !verifyRequestToken(c) {
		return echo.NewHTTPError(http.StatusForbidden)
	}
	entity := Entity{Name: c.FormValue("name")}
	database.Save(&entity)
	return c.NoContent(http.StatusNoContent)
}

func EchoRedirectPositive(c echo.Context) error {
	_ = c.Redirect(http.StatusFound, c.QueryParam("value"))
	return nil
}

func EchoRedirectControl(c echo.Context) error {
	_ = c.Redirect(http.StatusFound, "/local")
	return nil
}

func EchoRedirectUncertain(c echo.Context) error {
	_ = c.Redirect(http.StatusFound, opaqueValue(c.QueryParam("value")))
	return nil
}

func EchoPolicyPositive(c echo.Context) error {
	allowed, err := policy.Enforce("subject", "object", "write")
	if err != nil || !allowed {
		c.Logger().Warn("policy rejected request")
	}
	return c.NoContent(http.StatusNoContent)
}

func EchoPolicyControl(c echo.Context) error {
	allowed, err := policy.Enforce("subject", "object", "write")
	if err != nil {
		return err
	}
	if !allowed {
		return echo.NewHTTPError(http.StatusForbidden)
	}
	return c.NoContent(http.StatusNoContent)
}

func EchoPolicyUncertain(c echo.Context) error {
	if !c.Echo().Debug {
		allowed, err := policy.Enforce("subject", "object", "write")
		if err != nil || !allowed {
			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func EchoSessionPositive(c echo.Context) error {
	if _, err := c.Cookie("auth"); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	c.SetCookie(&http.Cookie{Name: "auth", Value: "credential", Secure: false, HttpOnly: true})
	return c.NoContent(http.StatusNoContent)
}

func EchoSessionControl(c echo.Context) error {
	if _, err := c.Cookie("auth"); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	c.SetCookie(&http.Cookie{Name: "auth", Value: "credential", Secure: true, HttpOnly: true})
	return c.NoContent(http.StatusNoContent)
}

func EchoSessionUncertain(c echo.Context) error {
	if _, err := c.Cookie("auth"); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	if c.Echo().Debug {
		c.SetCookie(&http.Cookie{Name: "auth", Value: "credential", Secure: false, HttpOnly: true})
	}
	return c.NoContent(http.StatusNoContent)
}

func EchoCorsPositive(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", c.Request().Header.Get("Origin"))
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func EchoCorsControl(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "https://trusted.example")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func EchoCorsUncertain(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", opaqueValue(c.Request().Header.Get("Origin")))
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func EchoBindingPositive(c echo.Context) error {
	var entity Entity
	entity.Role = "member"
	if err := c.Bind(&entity); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	database.Save(&entity)
	return c.JSON(http.StatusCreated, entity)
}

func EchoBindingControl(c echo.Context) error {
	var entity Entity
	entity.Role = "member"
	if err := c.Bind(&entity); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	entity.Role = "member"
	database.Save(&entity)
	return c.JSON(http.StatusCreated, entity)
}

func EchoBindingUncertain(c echo.Context) error {
	var entity Entity
	entity.Role = "member"
	if err := c.Bind(&entity); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if !isAdministrator(c) {
		entity.Role = "member"
	}
	database.Save(&entity)
	return c.JSON(http.StatusCreated, entity)
}
