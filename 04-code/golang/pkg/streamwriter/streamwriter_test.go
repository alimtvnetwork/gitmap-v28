package streamwriter_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/streamwriter"
)

// SafeBuffer wraps bytes.Buffer with mutex for test inspections
type SafeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *SafeBuffer) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *SafeBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func (s *SafeBuffer) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buf.Reset()
}

// Custom compilable struct
type CustomReport struct {
	Title string
	Score int
}

func (c CustomReport) Compile() string {
	return fmt.Sprintf("REPORT[%s]: score=%d", c.Title, c.Score)
}

// Nested struct containing custom compilable
type OuterContainer struct {
	ID     string       `json:"id"`
	Report CustomReport `json:"report"`
}

func TestBytesWrapper(t *testing.T) {
	raw := []byte("hello bytes")
	payload := "original-payload"

	// 1. Successful Bytes satisfying WrappedBytes
	b := streamwriter.NewBytes(raw, payload)
	var wb streamwriter.WrappedBytes[string] = b

	if !wb.IsValid() {
		t.Errorf("expected IsValid() to be true")
	}
	if !wb.IsSuccess() {
		t.Errorf("expected IsSuccess() to be true")
	}
	if !wb.Status() {
		t.Errorf("expected Status() flag to be true")
	}
	if wb.StatusCode() != 200 {
		t.Errorf("expected StatusCode() to be 200, got: %d", wb.StatusCode())
	}
	if wb.HasError() {
		t.Errorf("expected HasError() to be false")
	}
	if wb.String() != "hello bytes" {
		t.Errorf("unexpected string: %s", wb.String())
	}
	if wb.Len() != len(raw) {
		t.Errorf("unexpected len: %d", wb.Len())
	}
	if wb.Payload() != payload {
		t.Errorf("unexpected payload: %v", wb.Payload())
	}
	if wb.Value() != payload {
		t.Errorf("unexpected value: %v", wb.Value())
	}
	if wb.AppError() != nil {
		t.Errorf("expected nil AppError")
	}
	if wb.Error() != nil {
		t.Errorf("expected nil Error()")
	}

	unwrappedData, unwrappedErr := wb.Unwrap()
	if string(unwrappedData) != "hello bytes" || unwrappedErr != nil {
		t.Errorf("unwrap failed")
	}

	// 2. Failed Bytes with AppError
	appErr := appfault.New(errtype.Validation, "validation failed")
	errBytes := streamwriter.NewBytesError[string](appErr)
	if errBytes.IsValid() {
		t.Errorf("expected IsValid() to be false on error")
	}
	if errBytes.IsSuccess() {
		t.Errorf("expected IsSuccess() to be false on error")
	}
	if errBytes.Status() {
		t.Errorf("expected Status() flag to be false on error")
	}
	if !errBytes.HasError() {
		t.Errorf("expected HasError() to be true on error")
	}
	if errBytes.AppError() == nil {
		t.Errorf("expected non-nil AppError")
	}
	if errBytes.Error() == nil {
		t.Errorf("expected non-nil Error()")
	}
}

func TestJSONResult_WrappedBytesFlow(t *testing.T) {
	type UserPayload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	user := UserPayload{Name: "Alice", Age: 30}

	// 1. Successful JSONResult satisfying WrappedBytes and WrappedJSON
	res := streamwriter.NewJSONResult(user)
	var wb streamwriter.WrappedBytes[UserPayload] = res
	var wj streamwriter.WrappedJSON[UserPayload] = res

	if !wb.IsValid() {
		t.Fatalf("expected IsValid() to be true")
	}
	if !wb.IsSuccess() {
		t.Fatalf("expected IsSuccess() to be true")
	}
	if !wb.Status() {
		t.Fatalf("expected Status() flag to be true")
	}
	if wb.StatusCode() != 200 {
		t.Fatalf("expected status code 200, got: %d", wb.StatusCode())
	}
	if wb.Value().Name != "Alice" || wb.Payload().Age != 30 {
		t.Fatalf("expected payload value match")
	}
	if wb.Error() != nil {
		t.Fatalf("expected nil Error()")
	}

	// Test Pretty and Compact
	pretty := wj.Pretty()
	if !strings.Contains(pretty, "\n") || !strings.Contains(pretty, `"name": "Alice"`) {
		t.Fatalf("unexpected pretty JSON: %s", pretty)
	}
	compact := wj.Compact()
	if strings.Contains(compact, " ") || !strings.Contains(compact, `{"name":"Alice","age":30}`) {
		t.Fatalf("unexpected compact JSON: %s", compact)
	}

	// Test Unmarshal
	var decoded UserPayload
	if err := wj.Unmarshal(&decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.Name != "Alice" || decoded.Age != 30 {
		t.Fatalf("unmarshaled data mismatch: %+v", decoded)
	}

	// 2. Pre-marshaled bytes constructor
	validJSON := []byte(`{"name":"Bob","age":25}`)
	byteRes := streamwriter.NewJSONResultWithBytes(validJSON, UserPayload{Name: "Bob", Age: 25})
	if !byteRes.IsValid() || byteRes.StatusCode() != 200 {
		t.Fatalf("expected valid JSONResult from bytes")
	}

	// 3. Failed JSONResult with AppError
	appErr := appfault.New(errtype.Validation, "invalid json schema")
	errRes := streamwriter.NewJSONResultError[UserPayload](appErr)
	if errRes.IsValid() || errRes.IsSuccess() || errRes.Status() {
		t.Fatalf("expected error result to be invalid")
	}
	if errRes.Error() == nil || errRes.AppError() == nil {
		t.Fatalf("expected non-nil Error() and AppError()")
	}
}

func TestJSONSource_MultiSourceCreation(t *testing.T) {
	type Account struct {
		ID      string  `json:"id"`
		Balance float64 `json:"balance"`
	}
	acc := Account{ID: "acc-101", Balance: 500.50}

	// 1. FromPayload (top-level and scoped namespace)
	resPayload := streamwriter.FromPayload(acc)
	if !resPayload.IsSuccess() || resPayload.Value().ID != "acc-101" {
		t.Fatalf("FromPayload failed")
	}
	resScoped := streamwriter.JSONSourceOf[Account]().FromPayload(acc)
	if !resScoped.IsSuccess() || resScoped.Value().ID != "acc-101" {
		t.Fatalf("JSONSourceOf FromPayload failed")
	}

	// 2. FromBytes
	rawBytes := []byte(`{"id":"acc-102","balance":750.0}`)
	resBytes := streamwriter.FromBytes(rawBytes, Account{ID: "acc-102", Balance: 750.0})
	if !resBytes.IsSuccess() || resBytes.StatusCode() != 200 {
		t.Fatalf("FromBytes failed")
	}

	// 3. FromString
	jsonString := `{"id":"acc-103","balance":1200.0}`
	resString := streamwriter.FromString(jsonString, Account{ID: "acc-103", Balance: 1200.0})
	if !resString.IsSuccess() || resString.String() != jsonString {
		t.Fatalf("FromString failed")
	}

	// 4. FromReader
	reader := strings.NewReader(`{"id":"acc-104","balance":90.0}`)
	resReader := streamwriter.FromReader(reader, Account{ID: "acc-104", Balance: 90.0})
	if !resReader.IsSuccess() || resReader.Value().Balance != 90.0 {
		t.Fatalf("FromReader failed")
	}

	// 5. FromSerializer closure
	resSerializer := streamwriter.FromSerializer(func() ([]byte, *appfault.AppError) {
		return []byte(`{"id":"acc-105","balance":300.0}`), nil
	}, Account{ID: "acc-105", Balance: 300.0})
	if !resSerializer.IsSuccess() {
		t.Fatalf("FromSerializer failed")
	}

	// 6. FromBytesEnvelope
	wrappedBytes := streamwriter.NewBytes([]byte(`{"id":"acc-106","balance":450.0}`), Account{ID: "acc-106", Balance: 450.0})
	resEnv := streamwriter.FromBytesEnvelope(wrappedBytes)
	if !resEnv.IsSuccess() || resEnv.Value().ID != "acc-106" {
		t.Fatalf("FromBytesEnvelope failed")
	}

	// 7. FromError
	errRes := streamwriter.FromError[Account](appfault.New(errtype.NotFound, "account not found").WithStatusCode(404))
	if errRes.IsSuccess() || errRes.StatusCode() != 404 {
		t.Fatalf("FromError failed")
	}

	// 8. FromAny polymorphic
	anyMap := streamwriter.FromAny(map[string]string{"env": "prod"})
	if !anyMap.IsSuccess() || !strings.Contains(anyMap.String(), `"env":"prod"`) {
		t.Fatalf("FromAny failed")
	}

	// 9. Cast (round-trip type cast)
	type SimpleAccount struct {
		ID string `json:"id"`
	}
	casted := streamwriter.Cast[SimpleAccount](acc)
	if !casted.IsSuccess() || casted.Value().ID != "acc-101" {
		t.Fatalf("Cast failed: %+v", casted.Value())
	}

	// 10. ToBytes conversion
	bytesConv := resPayload.ToBytes()
	if !bytesConv.IsSuccess() || bytesConv.String() != resPayload.String() {
		t.Fatalf("ToBytes conversion failed")
	}

	// 11. PrettyOrError and CompactOrError
	pStr, pErr := resPayload.PrettyOrError()
	if pErr != nil || !strings.Contains(pStr, "\n") {
		t.Fatalf("PrettyOrError failed")
	}
	cStr, cErr := resPayload.CompactOrError()
	if cErr != nil || strings.Contains(cStr, " ") {
		t.Fatalf("CompactOrError failed")
	}
}

func TestCompiler_Primitives(t *testing.T) {
	// String
	if streamwriter.Compile("hello world") != "hello world" {
		t.Errorf("string compile failed")
	}

	// Numbers
	if streamwriter.Compile(42) != "42" {
		t.Errorf("int compile failed")
	}
	if streamwriter.Compile(3.14) != "3.14" {
		t.Errorf("float compile failed")
	}

	// Boolean
	if streamwriter.Compile(true) != "true" {
		t.Errorf("bool compile failed")
	}
	if streamwriter.Compile(false) != "false" {
		t.Errorf("bool compile failed")
	}

	// Nil
	var nilPtr *string
	if streamwriter.Compile(nilPtr) != "nil" {
		t.Errorf("nil compile failed")
	}
}

func TestCompiler_Maps_OrderWise(t *testing.T) {
	data := map[string]any{
		"zebra":  100,
		"apple":  "pie",
		"mango":  true,
		"banana": 42,
	}

	compiled := streamwriter.Compile(data)
	expected := `{apple: "pie", banana: 42, mango: true, zebra: 100}`
	if compiled != expected {
		t.Fatalf("expected map order: %s, got: %s", expected, compiled)
	}
}

func TestCompiler_Slices_OrderWise(t *testing.T) {
	sliceData := []any{"first", 2, true, "fourth"}
	compiled := streamwriter.Compile(sliceData)
	expected := `["first", 2, true, "fourth"]`
	if compiled != expected {
		t.Fatalf("expected slice order: %s, got: %s", expected, compiled)
	}
}

func TestCompiler_ObjectAndRecursiveCompilable(t *testing.T) {
	container := OuterContainer{
		ID: "box-1",
		Report: CustomReport{
			Title: "Audit",
			Score: 98,
		},
	}

	compiled := streamwriter.Compile(container)
	expected := `{id: "box-1", report: REPORT[Audit]: score=98}`
	if compiled != expected {
		t.Fatalf("expected recursive compilable output:\n%s\ngot:\n%s", expected, compiled)
	}
}

func TestLockedStreamer_Generic_ConcurrentSafe(t *testing.T) {
	buf := &SafeBuffer{}
	streamer := streamwriter.NewLockedStreamer[string](streamwriter.LockedOptions[string]{
		Name:        "concurrent-test",
		Destination: buf,
	})

	if !streamer.IsLocked() {
		t.Fatalf("expected IsLocked() to be true")
	}

	var wg sync.WaitGroup
	ctx := context.Background()
	concurrency := 25

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			appErr := streamer.Stream(ctx, fmt.Sprintf("goroutine-%d", idx))
			if appErr != nil {
				t.Errorf("stream failed: %v", appErr)
			}
		}(i)
	}
	wg.Wait()

	out := buf.String()
	for i := 0; i < concurrency; i++ {
		expected := fmt.Sprintf("goroutine-%d", i)
		if !strings.Contains(out, expected) {
			t.Errorf("missing expected output: %s", expected)
		}
	}
}

func TestLocklessStreamer_Generic_Direct(t *testing.T) {
	buf := &bytes.Buffer{}
	type Event struct {
		Name string `json:"name"`
		Code int    `json:"code"`
	}

	streamer := streamwriter.NewLocklessStreamer[Event](streamwriter.LocklessOptions[Event]{
		Name:        "cli-test",
		Destination: buf,
	})

	if streamer.IsLocked() {
		t.Fatalf("expected IsLocked() to be false")
	}

	ctx := context.Background()
	appErr := streamer.Stream(ctx, Event{Name: "login", Code: 200})
	if appErr != nil {
		t.Fatalf("stream failed: %v", appErr)
	}

	out := buf.String()
	if !strings.Contains(out, `name: "login"`) || !strings.Contains(out, `code: 200`) {
		t.Fatalf("expected compiled event in output, got: %s", out)
	}
}

func TestSelfBinding_GenericContracts(t *testing.T) {
	locked := streamwriter.NewLockedStreamer[any](streamwriter.LockedOptions[any]{Name: "test-locked"})
	lockless := streamwriter.NewLocklessStreamer[any](streamwriter.LocklessOptions[any]{Name: "test-lockless"})
	writer := streamwriter.NewPluggableWriter[any](streamwriter.WriterOptions[any]{Name: "test-writer", Streamer: locked})

	// Verify LockedStreamer self-binding and Locker
	var s1 streamwriter.Streamer[any] = locked.AsStreamer()
	var w1 streamwriter.Writer[any] = locked.AsWriter()
	var l1 sync.Locker = locked
	if s1 == nil || w1 == nil || l1 == nil {
		t.Fatal("locked streamer self-binding failed")
	}

	// Verify LocklessStreamer self-binding and Locker
	var s2 streamwriter.Streamer[any] = lockless.AsStreamer()
	var w2 streamwriter.Writer[any] = lockless.AsWriter()
	var l2 sync.Locker = lockless
	if s2 == nil || w2 == nil || l2 == nil {
		t.Fatal("lockless streamer self-binding failed")
	}

	// Verify PluggableWriter self-binding and Locker
	var w3 streamwriter.Writer[any] = writer.AsWriter()
	var l3 sync.Locker = writer
	if w3 == nil || l3 == nil {
		t.Fatal("pluggable writer self-binding failed")
	}
}

func TestWriter_LockerSynchronization(t *testing.T) {
	buf := &SafeBuffer{}
	lockedStreamer := streamwriter.NewLockedStreamer[string](streamwriter.LockedOptions[string]{
		Name:        "locker-test",
		Destination: buf,
	})
	writer := streamwriter.NewPluggableWriter[string](streamwriter.WriterOptions[string]{
		Name:     "pluggable-locker",
		Streamer: lockedStreamer,
	})

	// Verify Writer satisfies sync.Locker
	var locker sync.Locker = writer
	locker.Lock()
	ctx := context.Background()
	_ = writer.Write(ctx, "atomic-line-1")
	_ = writer.Write(ctx, "atomic-line-2")
	locker.Unlock()

	out := buf.String()
	if !strings.Contains(out, "atomic-line-1") || !strings.Contains(out, "atomic-line-2") {
		t.Fatalf("expected atomic writes in buffer, got: %s", out)
	}

	// Verify Lockless streamer satisfies sync.Locker no-op
	lockless := streamwriter.NewLocklessStreamer[string](streamwriter.LocklessOptions[string]{
		Name:        "lockless-locker",
		Destination: buf,
	})
	var locklessLocker sync.Locker = lockless
	locklessLocker.Lock()
	_ = lockless.Stream(ctx, "lockless-atomic")
	locklessLocker.Unlock()

	if !strings.Contains(buf.String(), "lockless-atomic") {
		t.Fatalf("expected lockless atomic write in buffer")
	}
}

func TestWriter_ConcurrentCompoundBatches(t *testing.T) {
	buf := &SafeBuffer{}
	streamer := streamwriter.NewLockedStreamer[string](streamwriter.LockedOptions[string]{
		Name:        "batch-test",
		Destination: buf,
	})
	writer := streamwriter.NewPluggableWriter[string](streamwriter.WriterOptions[string]{
		Name:     "batch-writer",
		Streamer: streamer,
	})

	ctx := context.Background()
	concurrency := 10
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			writer.Lock()
			_ = writer.Write(ctx, fmt.Sprintf("START-%d", id))
			_ = writer.Write(ctx, fmt.Sprintf("END-%d", id))
			writer.Unlock()
		}(i)
	}
	wg.Wait()

	out := buf.String()
	for i := 0; i < concurrency; i++ {
		startTag := fmt.Sprintf("START-%d", i)
		endTag := fmt.Sprintf("END-%d", i)
		startIdx := strings.Index(out, startTag)
		endIdx := strings.Index(out, endTag)
		if startIdx == -1 || endIdx == -1 {
			t.Fatalf("missing tags for id %d", i)
		}
		sub := out[startIdx : endIdx+len(endTag)]
		if strings.Count(sub, "START-") != 1 {
			t.Fatalf("batch %d was interleaved by another goroutine:\n%s", i, sub)
		}
	}
}

func TestSwappableMethods_GenericRuntime(t *testing.T) {
	buf := &SafeBuffer{}
	streamer := streamwriter.NewLockedStreamer[any](streamwriter.LockedOptions[any]{
		Name:        "swappable-test",
		Destination: buf,
	})

	ctx := context.Background()

	// Initial default stream
	_ = streamer.Stream(ctx, map[string]int{"b": 2, "a": 1})
	if !strings.Contains(buf.String(), "{a: 2, b: 2}") && !strings.Contains(buf.String(), "{a: 1, b: 2}") {
		t.Fatalf("unexpected initial output: %s", buf.String())
	}

	buf.Reset()

	// Hot-swap stream method returning *appfault.AppError
	streamer.SetStreamMethod(func(ctx context.Context, payload any, dest io.Writer) *appfault.AppError {
		_, err := fmt.Fprintf(dest, ">>> SWAPPED: %s <<<\n", streamwriter.Compile(payload))
		if err != nil {
			return appfault.Wrap(errtype.IO, err, "failed swap write")
		}
		return nil
	})

	_ = streamer.Stream(ctx, "dynamic-change")
	if !strings.Contains(buf.String(), ">>> SWAPPED: dynamic-change <<<") {
		t.Fatalf("unexpected swapped output: %s", buf.String())
	}
}

func TestCompositeLogger_FluentChaining(t *testing.T) {
	buf1 := &SafeBuffer{}
	buf2 := &SafeBuffer{}
	buf3 := &SafeBuffer{}

	w1 := streamwriter.NewLockedStreamer[any](streamwriter.LockedOptions[any]{Name: "w1", Destination: buf1})
	w2 := streamwriter.NewLocklessStreamer[any](streamwriter.LocklessOptions[any]{Name: "w2", Destination: buf2})

	customWriter := streamwriter.NewPluggableWriter[any](streamwriter.WriterOptions[any]{
		Name: "custom-api",
		WriteMethod: func(ctx context.Context, payload any) *appfault.AppError {
			_, err := fmt.Fprintf(buf3, "CUSTOM-API: %s\n", streamwriter.Compile(payload))
			if err != nil {
				return appfault.Wrap(errtype.IO, err, "custom api write failed")
			}
			return nil
		},
	})

	// FLUENT REGISTRATION
	log := streamwriter.NewLogger[any]().
		AddWriters(w1, w2).
		AddWriter(customWriter)

	if log.WriterCount() != 3 {
		t.Fatalf("expected 3 writers, got %d", log.WriterCount())
	}

	ctx := context.WithValue(context.Background(), "traceId", "trace-999")
	appErr := log.Info(ctx, "Order placed successfully", map[string]any{"orderId": "ord-77"})
	if appErr != nil {
		t.Fatalf("log.Info failed: %v", appErr)
	}

	// Verify emission across all 3 destinations
	if !strings.Contains(buf1.String(), "Order placed successfully") {
		t.Errorf("w1 did not receive log")
	}
	if !strings.Contains(buf2.String(), "Order placed successfully") {
		t.Errorf("w2 did not receive log")
	}
	if !strings.Contains(buf3.String(), "CUSTOM-API: ") {
		t.Errorf("customWriter did not receive log")
	}

	// Dynamic removal
	log.RemoveWriter("custom-api")
	if log.WriterCount() != 2 {
		t.Fatalf("expected 2 writers after removal, got %d", log.WriterCount())
	}

	// Clear to silent mode
	log.ClearWriters()
	if log.WriterCount() != 0 {
		t.Fatalf("expected 0 writers after clear, got %d", log.WriterCount())
	}

	buf1.Reset()
	_ = log.Info(ctx, "Silent message")
	if buf1.String() != "" {
		t.Fatalf("expected empty buffer in silent mode, got: %s", buf1.String())
	}
}

func TestLogRecord_Compile(t *testing.T) {
	rec := streamwriter.LogRecord{
		Timestamp: time.Unix(0, 0).UTC(),
		Level:     streamwriter.LevelInfo,
		Message:   "hello log",
		TraceID:   "tx-1",
		Fields:    map[string]any{"z": 1, "a": 2},
	}

	compiled := streamwriter.Compile(rec)
	if !strings.Contains(compiled, "[trace=tx-1]") || !strings.Contains(compiled, "fields={a: 2, z: 1}") {
		t.Fatalf("unexpected LogRecord compile: %s", compiled)
	}
}
