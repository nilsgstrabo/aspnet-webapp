package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
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
		// for k, v := range ctx.Request.Header {
		// 	fmt.Printf("%q: %v\n", k, v)
		// }

		hostName, err := os.Hostname()
		if err != nil {
			hostName = "N/A"
			fmt.Printf("error getting host name: %v", err)
		}
		ctx.String(http.StatusOK, fmt.Sprintf("hello from %s", hostName))

		// Sleet between 0 and 1000 ms
		time.Sleep(time.Duration(rand.Intn(1000) * int(time.Millisecond)))

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

	cred, err := azidentity.NewManagedIdentityCredential(nil)
	if err != nil {
		panic(err)
	}
	azts := &azTokenSource{cred: cred}
	client := oauth2.NewClient(context.Background(), oauth2.ReuseTokenSource(nil, oauth2.ReuseTokenSource(nil, azts)))
	handler.GET("/log", getLog(*client))

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

func getLog(client http.Client) func(*gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("preparing log request")
		resp, err := client.Get("https://server-radix-log-api-qa.dev.radix.equinor.com/api/v1/applications/oauth-demo/environments/dev/components/simple")
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
	t, err := s.cred.GetToken(context.Background(), policy.TokenRequestOptions{
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
