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
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

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
	tick, err := time.ParseDuration(os.Getenv("APPENDTICK"))
	if err != nil {
		fmt.Println(err)
		tick = 5 * time.Second
	}
	log.Printf("appending to file every %v second(s)", tick.Seconds())
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	fileName := path.Join("/", os.Getenv("TMPDIR"), "file.txt")
	f, err := os.Create(fileName)
	if err != nil {
		panic(fmt.Errorf("failed to create file: %w", err))
	}

	append := strings.Repeat("helloworld", 100)

	for {
		select {
		case <-ticker.C:
			_, err = f.WriteString(append)
			if err != nil {
				panic(fmt.Errorf("failed to create file: %w", err))
			}
			stat, err := os.Stat(fileName)
			if err != nil {
				panic(fmt.Errorf("failed to get stats for file: %w", err))
			}
			log.Printf("file size: %s\n", formatFileSize(float64(stat.Size()), 1024.0))
		case <-stop:
			return
		}

	}
}

var sizes = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}

func formatFileSize(s float64, base float64) string {
	unitsLimit := len(sizes)
	i := 0
	for s >= base && i < unitsLimit {
		s = s / base
		i++
	}

	f := "%.0f %s"
	if i > 1 {
		f = "%.2f %s"
	}

	return fmt.Sprintf(f, s, sizes[i])
}
