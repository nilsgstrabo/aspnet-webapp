package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

const appInsightsConnectionStringEnv = "APPLICATIONINSIGHTS_CONNECTION_STRING"

var traceParentPattern = regexp.MustCompile(`^[0-9a-f]{2}-([0-9a-f]{32})-([0-9a-f]{16})-[0-9a-f]{2}$`)

// initTelemetry returns a middleware that sends request telemetry to Application Insights,
// and a function that flushes pending telemetry on shutdown.
// Telemetry is disabled when the connection string is missing or invalid.
func initTelemetry(debug bool) (gin.HandlerFunc, func()) {
	disabled := func(*gin.Context) {}

	connectionString := os.Getenv(appInsightsConnectionStringEnv)
	if connectionString == "" {
		fmt.Printf("%s is not set, telemetry is disabled\n", appInsightsConnectionStringEnv)
		return disabled, func() {}
	}

	client, err := newTelemetryClient(connectionString)
	if err != nil {
		fmt.Printf("telemetry is disabled: %v\n", err)
		return disabled, func() {}
	}

	if debug {
		key := client.InstrumentationKey()
		redactKey := strings.NewReplacer(key, "<redacted>", strings.ReplaceAll(key, "-", ""), "<redacted>")
		appinsights.NewDiagnosticsMessageListener(func(msg string) error {
			fmt.Printf("appinsights: %s\n", redactKey.Replace(msg))
			return nil
		})
	}

	fmt.Printf("sending telemetry to %s\n", client.Channel().EndpointAddress())

	flush := func() {
		select {
		case <-client.Channel().Close(5 * time.Second):
		case <-time.After(10 * time.Second):
			fmt.Println("timed out flushing telemetry")
		}
	}

	return trackRequest(client), flush
}

func newTelemetryClient(connectionString string) (appinsights.TelemetryClient, error) {
	instrumentationKey, endpointURL, err := parseAppInsightsConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	config := appinsights.NewTelemetryConfiguration(instrumentationKey)
	if endpointURL != "" {
		config.EndpointUrl = endpointURL
	}
	client := appinsights.NewTelemetryClientFromConfig(config)

	role := os.Getenv("RADIX_COMPONENT")
	if role == "" {
		role = "simpleweb"
	}
	client.Context().Tags.Cloud().SetRole(role)
	client.Context().Tags.Application().SetVer(buildVersion())
	if env := os.Getenv("RADIX_ENVIRONMENT"); env != "" {
		client.Context().CommonProperties["environment"] = env
	}

	return client, nil
}

// parseAppInsightsConnectionString returns the instrumentation key and the URL telemetry is posted to.
// The URL is empty when the connection string has no IngestionEndpoint.
func parseAppInsightsConnectionString(connectionString string) (instrumentationKey string, endpointURL string, err error) {
	var ingestionEndpoint string
	for _, part := range strings.Split(connectionString, ";") {
		key, value, found := strings.Cut(part, "=")
		if !found {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "instrumentationkey":
			instrumentationKey = strings.TrimSpace(value)
		case "ingestionendpoint":
			ingestionEndpoint = strings.TrimSpace(value)
		}
	}

	if instrumentationKey == "" {
		return "", "", errors.New("connection string is missing InstrumentationKey")
	}
	if ingestionEndpoint == "" {
		return instrumentationKey, "", nil
	}

	u, err := url.Parse(ingestionEndpoint)
	if err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return "", "", errors.New("connection string has an invalid IngestionEndpoint")
	}

	return instrumentationKey, strings.TrimRight(ingestionEndpoint, "/") + "/v2/track", nil
}

func trackRequest(client appinsights.TelemetryClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		end := time.Now()

		request := appinsights.NewRequestTelemetry(ctx.Request.Method, requestURL(ctx.Request), end.Sub(start), strconv.Itoa(ctx.Writer.Status()))
		request.MarkTime(start, end)
		// Use the route template instead of the actual path to keep the request name low cardinality.
		request.Name = ctx.Request.Method + " " + ctx.FullPath()
		request.Tags.Operation().SetName(request.Name)
		request.Tags.Location().SetIp(ctx.ClientIP())

		if m := traceParentPattern.FindStringSubmatch(ctx.GetHeader("traceparent")); m != nil {
			request.Tags.Operation().SetId(m[1])
			request.Tags.Operation().SetParentId(m[2])
		} else {
			request.Tags.Operation().SetId(newOperationID())
		}

		client.Track(request)
	}
}

func requestURL(r *http.Request) string {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
		if r.TLS != nil {
			scheme = "https"
		}
	}
	return scheme + "://" + r.Host + r.URL.RequestURI()
}

func newOperationID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
