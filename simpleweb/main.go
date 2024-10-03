package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type CodeRequest struct {
	Code string `uri:"code" binding:"required"`
}

func main() {
	fmt.Printf("Running server on %s/%s\n\n", runtime.GOOS, runtime.GOARCH)

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
		if b, err := io.ReadAll(ctx.Request.Body); err == nil {
			fmt.Printf("body: \n%s", string(b))
		} else {
			fmt.Printf("error reading body: %v", err)
		}

		ctx.Status(400)
		ctx.Stream(func(w io.Writer) bool {
			for i := range 5 {
				fmt.Printf("sending line %d to client\n", i)
				w.Write([]byte(fmt.Sprintf("line %d", i)))
				time.Sleep(1 * time.Second)
			}
			return false
		})
	})

	handler.Any("/:code", func(ctx *gin.Context) {
		var c CodeRequest
		if err := ctx.ShouldBindUri(&c); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		code, err := strconv.Atoi(c.Code)
		if err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		ctx.JSON(code, c)
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

	// go slowlyWriteToFile(ctx.Done())

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
	fmt.Printf("appending to file every %v second(s)", tick.Seconds())
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
			fmt.Printf("file size: %s\n", formatFileSize(float64(stat.Size()), 1024.0))
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
