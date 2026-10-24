package cmddb

import (
	"encoding/json"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"net/http"
	"reflect"
	"time"
)

// RestStatus represents the status metadata of a REST response.
type RestStatus struct {
	IsSuccess bool   `json:"IsSuccess"`
	IsFailed  bool   `json:"IsFailed"`
	Code      int    `json:"Code"`
	Message   string `json:"Message"`
	Timestamp string `json:"Timestamp"`
}

// RestEnvelope represents the universal PascalCase REST response envelope.
type RestEnvelope struct {
	Status     RestStatus     `json:"Status"`
	Attributes map[string]any `json:"Attributes"`
	Results    any            `json:"Results"`
	Errors     any            `json:"Errors"`
}

// writeRestEnvelope writes a universal response envelope to the HTTP response writer.
func writeRestEnvelope(w http.ResponseWriter, code int, results any, attributes map[string]any, err error) {
	status := createRestStatus(code, err)
	resArray, total := sanitizeResultsArray(results)
	attrs := mergeRestAttributes(attributes, total, err)
	errPayload := formatRestErrors(err)

	envelope := RestEnvelope{
		Status:     status,
		Attributes: attrs,
		Results:    resArray,
		Errors:     errPayload,
	}

	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	w.WriteHeader(status.Code)
	_ = json.NewEncoder(w).Encode(envelope)
}

// WriteRestEnvelope exposes writeRestEnvelope for public callers and test suites.
func WriteRestEnvelope(w http.ResponseWriter, code int, results any, attributes map[string]any, err error) {
	writeRestEnvelope(w, code, results, attributes, err)
}

func createRestStatus(code int, err error) RestStatus {
	statusCode := resolveStatusCode(code, err)
	isSuccess := err == nil && statusCode < 400
	isFailed := err != nil || statusCode >= 400
	msg := resolveStatusMessage(statusCode, err)

	return RestStatus{
		IsSuccess: isSuccess,
		IsFailed:  isFailed,
		Code:      statusCode,
		Message:   msg,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
	}
}

func resolveStatusCode(code int, err error) int {
	if code > 0 {
		return code
	}
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func resolveStatusMessage(code int, err error) string {
	if err != nil {
		return err.Error()
	}
	text := http.StatusText(code)
	if text == "" {
		return "OK"
	}

	return text
}

func sanitizeResultsArray(results any) (any, int) {
	if results == nil {
		return []any{}, 0
	}

	val := reflect.ValueOf(results)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return []any{results}, 1
	}

	if val.Len() == 0 {
		return []any{}, 0
	}

	return results, val.Len()
}

func mergeRestAttributes(custom map[string]any, count int, err error) map[string]any {
	out := defaultRestAttributes(count, err)
	for k, v := range custom {
		out[k] = v
	}

	return out
}

func defaultRestAttributes(count int, err error) map[string]any {
	return map[string]any{
		"RequestedAt":  time.Now().UTC().Format(time.RFC3339Nano),
		"HasAnyErrors": err != nil,
		"IsSingle":     count == 1,
		"TotalRecords": count,
		"PerPage":      20,
		"CurrentPage":  1,
	}
}

func formatRestErrors(err error) any {
	if err == nil {
		return nil
	}

	return []string{err.Error()}
}
