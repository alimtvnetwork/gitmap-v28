package setup

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// HTTPDoer abstracts HTTP requests for unit test injection.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// TestConnection executes an endpoint verification ping using http.DefaultClient.
func TestConnection(req ConnectionTestRequest) ConnectionTestResult {
	return TestConnectionWithClient(http.DefaultClient, req)
}

// RunConnectionTest executes an endpoint verification ping from setup options.
func RunConnectionTest(opts IDESetupOptions) (*ConnectionTestResult, error) {
	ApplyPresetDefaults(&opts)
	req := buildConnectionTestRequest(opts)
	result := TestConnection(req)
	if !result.IsSuccess {
		return &result, apperror.NewExecutionError(result.ErrorMessage)
	}

	return &result, nil
}

func buildConnectionTestRequest(opts IDESetupOptions) ConnectionTestRequest {
	return ConnectionTestRequest{
		Provider: opts.Provider,
		BaseURL:  opts.BaseURL,
		APIKey:   opts.APIKey,
		Model:    opts.Model,
	}
}

// TestConnectionWithClient runs a lightweight HTTP GET against the provider models endpoint.
func TestConnectionWithClient(client HTTPDoer, req ConnectionTestRequest) ConnectionTestResult {
	if strings.TrimSpace(req.BaseURL) == "" {
		return ConnectionTestResult{IsSuccess: false, ErrorMessage: "empty base url"}
	}

	start := time.Now()
	httpReq, cancel, err := createProbeRequest(req)
	if err != nil {
		return ConnectionTestResult{IsSuccess: false, ErrorMessage: err.Error()}
	}

	defer cancel()

	return executeProbe(client, httpReq, req.Model, start)
}

func createProbeRequest(req ConnectionTestRequest) (*http.Request, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	url := buildTestURL(req.BaseURL, req.Provider)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		cancel()

		return nil, nil, err
	}

	applyTestHeaders(httpReq, req)

	return httpReq, cancel, nil
}

func executeProbe(client HTTPDoer, httpReq *http.Request, model string, start time.Time) ConnectionTestResult {
	resp, err := client.Do(httpReq)
	if err != nil {
		return ConnectionTestResult{IsSuccess: false, LatencyMs: time.Since(start).Milliseconds(), Model: model, ErrorMessage: err.Error()}
	}

	defer resp.Body.Close()

	return buildProbeResponse(resp.StatusCode, model, start)
}

func buildProbeResponse(statusCode int, model string, start time.Time) ConnectionTestResult {
	latency := time.Since(start).Milliseconds()
	isSuccess := statusCode >= 200 && statusCode < 300

	return ConnectionTestResult{
		IsSuccess:      isSuccess,
		LatencyMs:      latency,
		Model:          model,
		HTTPStatusCode: statusCode,
		ErrorMessage:   formatHTTPStatusError(statusCode),
	}
}

func buildTestURL(baseURL string, provider IDEProviderType) string {
	clean := strings.TrimRight(baseURL, "/")
	if provider == ProviderAnthropic && !strings.HasSuffix(clean, "/v1") {
		return clean + "/v1/models"
	}

	return clean + "/models"
}

func applyTestHeaders(httpReq *http.Request, req ConnectionTestRequest) {
	if req.Provider == ProviderAnthropic {
		httpReq.Header.Set("x-api-key", req.APIKey)
		httpReq.Header.Set("anthropic-version", "2023-06-01")

		return
	}

	if req.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+req.APIKey)
	}
}

func formatHTTPStatusError(code int) string {
	if code >= 200 && code < 300 {
		return ""
	}

	if code == 401 || code == 403 {
		return fmt.Sprintf("Authentication failed (HTTP %d): invalid or missing API key", code)
	}

	return fmt.Sprintf("Provider responded with HTTP %d", code)
}
