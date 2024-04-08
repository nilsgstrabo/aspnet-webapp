package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	handler := gin.New()
	handler.RemoveExtraSlash = true
	handler.GET("/", func(ctx *gin.Context) {
		fmt.Println("")
		host, _, _ := net.SplitHostPort(ctx.Request.RemoteAddr)
		fmt.Printf("Remote addr : %s\n", net.ParseIP(host))
		fmt.Println("")
		for k, v := range ctx.Request.Header {
			fmt.Printf("%q: %v\n", k, v)
		}
		ctx.Status(http.StatusOK)
	})

	if err := http.ListenAndServe(":9001", handler); err != nil {
		panic(err)
	}
}
