package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func Test_ParseAppInsightsConnectionString(t *testing.T) {
	tests := map[string]struct {
		connectionString string
		expectedKey      string
		expectedURL      string
		expectError      bool
	}{
		"key and ingestion endpoint": {
			connectionString: "InstrumentationKey=abc;IngestionEndpoint=https://northeurope-2.in.applicationinsights.azure.com/;LiveEndpoint=https://northeurope.livediagnostics.monitor.azure.com/",
			expectedKey:      "abc",
			expectedURL:      "https://northeurope-2.in.applicationinsights.azure.com/v2/track",
		},
		"keys are case insensitive": {
			connectionString: "instrumentationkey=abc; ingestionendpoint=https://example.com",
			expectedKey:      "abc",
			expectedURL:      "https://example.com/v2/track",
		},
		"key only uses default endpoint": {
			connectionString: "InstrumentationKey=abc",
			expectedKey:      "abc",
		},
		"missing key": {
			connectionString: "IngestionEndpoint=https://example.com/",
			expectError:      true,
		},
		"invalid ingestion endpoint": {
			connectionString: "InstrumentationKey=abc;IngestionEndpoint=example.com",
			expectError:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			key, endpointURL, err := parseAppInsightsConnectionString(test.connectionString)
			if (err != nil) != test.expectError {
				t.Fatalf("unexpected error: %v", err)
			}
			if key != test.expectedKey {
				t.Errorf("expected key %q, got %q", test.expectedKey, key)
			}
			if endpointURL != test.expectedURL {
				t.Errorf("expected url %q, got %q", test.expectedURL, endpointURL)
			}
		})
	}
}

func Test_TrackRequest_SendsRequestTelemetry(t *testing.T) {
	type envelope struct {
		IKey string            `json:"iKey"`
		Tags map[string]string `json:"tags"`
		Data struct {
			BaseType string `json:"baseType"`
			BaseData struct {
				Name         string `json:"name"`
				ResponseCode string `json:"responseCode"`
				Success      bool   `json:"success"`
				URL          string `json:"url"`
			} `json:"baseData"`
		} `json:"data"`
	}

	type submission struct {
		path string
		body []byte
	}
	received := make(chan submission, 1)
	ingestion := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			t.Errorf("failed to read gzip body: %v", err)
			return
		}
		body, _ := io.ReadAll(gz)
		received <- submission{path: r.URL.Path, body: body}
		w.Write([]byte(`{"itemsReceived":1,"itemsAccepted":1,"errors":[]}`))
	}))
	defer ingestion.Close()

	t.Setenv("RADIX_COMPONENT", "simple")
	client, err := newTelemetryClient("InstrumentationKey=00000000-0000-0000-0000-000000000001;IngestionEndpoint=" + ingestion.URL + "/")
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/files/:filename", trackRequest(client), func(ctx *gin.Context) {
		ctx.Status(http.StatusNotFound)
	})

	req := httptest.NewRequest(http.MethodGet, "http://simple.test/files/a.txt?x=1", nil)
	req.Header.Set("traceparent", "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
	router.ServeHTTP(httptest.NewRecorder(), req)

	select {
	case <-client.Channel().Close():
	case <-time.After(5 * time.Second):
		t.Fatal("timed out flushing telemetry")
	}

	var sub submission
	select {
	case sub = <-received:
	default:
		t.Fatal("no telemetry received")
	}

	if sub.path != "/v2/track" {
		t.Errorf("unexpected ingestion path %q", sub.path)
	}

	var env envelope
	if err := json.Unmarshal(sub.body, &env); err != nil {
		t.Fatalf("failed to decode envelope %q: %v", sub.body, err)
	}

	if env.IKey != "00000000-0000-0000-0000-000000000001" {
		t.Errorf("unexpected iKey %q", env.IKey)
	}
	if env.Data.BaseType != "RequestData" {
		t.Errorf("unexpected base type %q", env.Data.BaseType)
	}
	if env.Data.BaseData.Name != "GET /files/:filename" {
		t.Errorf("unexpected name %q", env.Data.BaseData.Name)
	}
	if env.Data.BaseData.ResponseCode != "404" || env.Data.BaseData.Success {
		t.Errorf("unexpected response code %q, success %v", env.Data.BaseData.ResponseCode, env.Data.BaseData.Success)
	}
	if env.Data.BaseData.URL != "http://simple.test/files/a.txt?x=1" {
		t.Errorf("unexpected url %q", env.Data.BaseData.URL)
	}
	if env.Tags["ai.operation.id"] != "0af7651916cd43dd8448eb211c80319c" || env.Tags["ai.operation.parentId"] != "b7ad6b7169203331" {
		t.Errorf("unexpected operation tags %v", env.Tags)
	}
	if env.Tags["ai.cloud.role"] != "simple" {
		t.Errorf("unexpected cloud role %q", env.Tags["ai.cloud.role"])
	}
}
