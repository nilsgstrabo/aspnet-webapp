package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type CodeRequest struct {
	Code string `uri:"code" binding:"required"`
}

type FileRequest struct {
	FileName string `uri:"filename" binding:"required"`
}

func main() {
	fmt.Printf("Running server on %s/%s\n\n", runtime.GOOS, runtime.GOARCH)

	// var timeout time.Duration
	// timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	// if err != nil {
	// 	fmt.Println(err)
	// 	timeout = 10 * time.Second
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// defer cancel()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	handler := gin.New()
	handler.RemoveExtraSlash = true
	handler.Use(logRequestInfo)
	handler.GET("/", func(ctx *gin.Context) {
		fmt.Println("")
		host, _, _ := net.SplitHostPort(ctx.Request.RemoteAddr)
		fmt.Printf("Remote address : %s\n", net.ParseIP(host))
		fmt.Println("")
		// for k, v := range ctx.Request.Header {
		// 	fmt.Printf("%q: %v\n", k, v)
		// }

		hostName, err := os.Hostname()
		if err != nil {
			hostName = "N/A"
			fmt.Printf("error getting host name: %v", err)
		}
		// ctx.String(http.StatusOK, fmt.Sprintf("hello from %s", hostName))

		ctx.Header("x-large-header", strings.Repeat("abcdefghijklmnop", 1280))
		ctx.Status(http.StatusOK)

		sleep := time.Duration(rand.Intn(3000) * int(time.Millisecond))
		ctx.Writer.WriteString(fmt.Sprintf("hello from %s\n", hostName))
		ctx.Writer.WriteString(fmt.Sprintf("sleeping %s before sending more data in response\n", sleep.String()))
		ctx.Writer.Flush()
		fmt.Printf("sleeping for %s\n", sleep.String())
		time.Sleep(sleep)
		ctx.Writer.WriteString("this is the last data in the response")

		// Sleep between 0 and 1000 ms

		// ctx.Status(http.StatusOK)
		// for i := range 5 {
		// 	fmt.Printf("sending line %d to client\n", i)
		// 	if _, err := ctx.Writer.Write([]byte(fmt.Sprintf("line %d", i))); err != nil {
		// 		fmt.Printf("error writing response: %v\n", err)
		// 	}
		// 	ctx.Writer.Flush()
		// 	time.Sleep(1 * time.Second)
		// }

	})

	handler.GET("/nils", func(ctx *gin.Context) {
		de, err := os.ReadDir("/mnt/videos/")
		if err != nil {
			fmt.Println(err)
			ctx.AbortWithError(500, err)
			return
		}
		fmt.Printf("found %d files in directory\n", len(de))

		f, err := os.Open("/mnt/videos/nils.txt")
		if err != nil {
			fmt.Println(err)
			ctx.AbortWithError(500, err)
			return
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			fmt.Println(err)
			ctx.AbortWithError(500, err)
			return
		}
		fmt.Println(fs.FormatFileInfo(fi))

		b, err := io.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			ctx.AbortWithError(500, err)
			return
		}
		ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
		ctx.Writer.WriteHeader(200)
		ctx.Writer.Write(b)
		// ctx.String(200, string(b))
	})

	handler.GET("/files/:filename", func(ctx *gin.Context) {
		var f FileRequest
		if err := ctx.ShouldBindUri(&f); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		ctx.File(filepath.Join("/mnt/videos/", filepath.Join("/", f.FileName)))
	})

	handler.GET("/:code", func(ctx *gin.Context) {
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

	handler.GET("/log", initLogHandler())
	handler.POST("/data", readBody)

	go func() {
		if err := http.ListenAndServe(":9001", handler); err != nil {
			fmt.Println(err)
		}
	}()

	// sigCh := make(chan os.Signal, 1)
	// signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	fmt.Println("waiting for signal")
	// 	s := <-sigCh
	// 	fmt.Printf("received signal %v, but we ignore it", s)
	// }()

	// go slowlyWriteToFile(ctx.Done())

	fmt.Println("waiting")
	<-ctx.Done()
	fmt.Println("done waiting")
}

func logRequestInfo(ctx *gin.Context) {
	fmt.Println()
	fmt.Printf("Content length: %d \n", ctx.Request.ContentLength)
	fmt.Printf("Remote address: %s \n", ctx.Request.RemoteAddr)

	fmt.Println("Headers:")
	for k, v := range ctx.Request.Header {
		fmt.Printf("  %q: %v\n", k, v)
	}
}

func readBody(ctx *gin.Context) {
	defer ctx.Request.Body.Close()
	buf := make([]byte, 1024*64)
	code, msg := http.StatusOK, "Ok"
	var total int
	for {
		l, err := ctx.Request.Body.Read(buf)
		total += l
		if l > 0 {
			fmt.Printf("read %d bytes from request, total %d\n", l, total)
		}
		if err != nil && !errors.Is(err, io.EOF) {
			msg = err.Error()
		}
		if l == 0 || err != nil {
			fmt.Println("finished reading request body")
			break
		}
	}
	ctx.String(code, msg)
}

func initLogHandler() func(*gin.Context) {
	cred, err := azidentity.NewWorkloadIdentityCredential(nil)
	if err != nil {
		panic(err)
	}
	client := oauth2.NewClient(
		context.Background(),
		oauth2.ReuseTokenSource(nil, &azTokenSource{cred: cred}),
	)
	return func(ctx *gin.Context) {
		fmt.Println("preparing log request")
		tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req, _ := http.NewRequest("GET", "https://server-radix-log-api-qa.dev.radix.equinor.com/api/v1/applications/oauth-demo/environments/dev/components/simple", nil)
		req = req.WithContext(tctx)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("log request failed: %v", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		d, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("reading response failed: %v", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		fmt.Println("returning log data")
		ctx.String(http.StatusOK, string(d))
	}
}

var _ oauth2.TokenSource = &azTokenSource{}

type azTokenSource struct {
	cred azcore.TokenCredential
}

func (s *azTokenSource) Token() (*oauth2.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	t, err := s.cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"6dae42f8-4368-4678-94ff-3960e28e3630/.default"},
	})
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken: t.Token,
		Expiry:      t.ExpiresOn,
	}, nil
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
