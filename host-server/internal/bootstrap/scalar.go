package bootstrap

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

const scalarHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Polyglot Plugin Host Server - API Documentation</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
    <script 
        id="api-reference" 
        data-url="/swagger.json"
        data-configuration='{"theme":"purple","darkMode":false}'
    ></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>
`

func RegisterScalarDocs(e *echo.Echo) {
	e.GET("/swagger.json", func(c echo.Context) error {
		return c.File("docs/swagger.json")
	})
	e.GET("/swagger.yaml", func(c echo.Context) error {
		return c.File("docs/swagger.yaml")
	})

	e.GET("/docs", func(c echo.Context) error {
		tmpl, err := template.New("scalar").Parse(scalarHTML)
		if err != nil {
			return err
		}
		return tmpl.Execute(c.Response().Writer, nil)
	})

	// Redirect root to docs
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/docs")
	})
}
