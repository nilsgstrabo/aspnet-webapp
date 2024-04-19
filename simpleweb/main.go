package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	var timeout time.Duration
	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		fmt.Println(err)
		timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
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

	go func() {
		if err := http.ListenAndServe(":9001", handler); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println("waiting")
	<-ctx.Done()
	fmt.Println("done waiting")
}
