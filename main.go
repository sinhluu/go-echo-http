package main

import (
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func main() {
	e := echo.New()

	e.GET("/ifconfig", func(c echo.Context) error {
		return c.String(200, c.RealIP())
	})

	e.GET("/request", func(c echo.Context) error {
		u := c.QueryParam("url")
		if u == "" {
			return c.String(200, "empty url")
		}
		h := http.Client{}
		res, err := h.Get(u)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return c.String(200, "Result: \n"+string(b))
	})

	e.Logger.Fatal(e.Start(":4000"))
}
