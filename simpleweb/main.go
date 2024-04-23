package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
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

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("waiting for signal")
		s := <-sigCh
		fmt.Printf("received signal %v, but we ignore it", s)
	}()

	go slowlyWriteToFile(ctx.Done())

	fmt.Println("waiting")
	<-ctx.Done()
	fmt.Println("done waiting")
}

func slowlyWriteToFile(stop <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fileName := path.Join("/", os.Getenv("TMPDIR"), "file.txt")
	f, err := os.Create(fileName)
	if err != nil {
		panic(fmt.Errorf("failed to create file: %w", err))
	}

	append := strings.Repeat("helloworld", 10000)

	for {
		select {
		case <-ticker.C:
			log.Println("tick")
			_, err = f.WriteString(append)
			if err != nil {
				panic(fmt.Errorf("failed to create file: %w", err))
			}
		case <-stop:
			return
		}

	}
}
