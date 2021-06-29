package kgb

import (
	"kgb/girov"
	"net/http"
)

func Run() {
	r := girov.New()

	r.GET("/", func(c *girov.Context) {
		c.String(http.StatusOK, "Hello KGB!\n")
	})

	r.Run(":9999")
}
