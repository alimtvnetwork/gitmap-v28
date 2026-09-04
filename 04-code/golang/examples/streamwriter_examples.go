package examples

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/streamwriter"
)

// UserAccount represents a domain model used across example demonstrations.
type UserAccount struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
}

// OrderEvent represents a business event payload for streaming.
type OrderEvent struct {
	OrderId   string  `json:"orderId"`
	AccountId string  `json:"accountId"`
	Amount    float64 `json:"amount"`
	ItemCount int     `json:"itemCount"`
}

// RunLoggerExample demonstrates the composite logger, multi-destination fanout, and context tracing.
func RunLoggerExample(dest io.Writer) *appfault.AppError {
	// 1. Create a composite logger supporting any payload
	log := streamwriter.NewLogger[any]()

	// 2. Attach a thread-safe locked destination (e.g. console / file)
	consoleWriter := streamwriter.NewLockedStreamer[any](streamwriter.LockedOptions[any]{
		Name:        "console-streamer",
		Destination: dest,
	})

	// 3. Attach a lockless streamer (e.g. high-throughput in-memory buffer)
	memBuffer := &bytes.Buffer{}
	memStreamer := streamwriter.NewLocklessStreamer[any](streamwriter.LocklessOptions[any]{
		Name:        "mem-streamer",
		Destination: memBuffer,
	})

	// 4. Attach a custom pluggable writer (e.g. HTTP webhook / auditing service)
	auditWriter := streamwriter.NewPluggableWriter[any](streamwriter.WriterOptions[any]{
		Name:        "audit-api-writer",
		Destination: dest,
		WriteMethod: func(ctx context.Context, w *streamwriter.PluggableWriter[any], payload any) *appfault.AppError {
			trace := ""
			if traceVal := ctx.Value("traceId"); traceVal != nil {
				trace = fmt.Sprintf("[%v] ", traceVal)
			}
			outDest := w.Destination()
			if outDest == nil {
				outDest = os.Stdout
			}
			_, err := fmt.Fprintf(outDest, "[%s] %s%s\n", w.Name(), trace, streamwriter.Compile(payload))
			if err != nil {
				return appfault.Wrap(errtype.IO, err, "audit write failed")
			}
			return nil
		},
	})

	// Fluent registration of destinations
	log.AddWriters(consoleWriter, memStreamer).AddWriter(auditWriter)

	// 5. Emit structured logs with context tracing
	ctx := context.WithValue(context.Background(), "traceId", "req-tx-8891")

	appErr := log.Info(ctx, "User account authentication succeeded", map[string]any{
		"userId": "usr-4412",
		"tier":   "enterprise",
	})
	if appErr != nil {
		return appErr
	}

	appErr = log.Warn(ctx, "High memory utilization threshold reached", map[string]any{
		"percent": 84.5,
		"node":    "worker-node-03",
	})
	if appErr != nil {
		return appErr
	}

	appErr = log.Error(ctx, "Payment gateway returned unexpected gateway timeout", map[string]any{
		"gateway": "stripe",
		"attempt": 3,
	})
	if appErr != nil {
		return appErr
	}

	return nil
}

// RunJsonExample demonstrates multi-source JsonResult creation, dynamic validity, unmarshaling, and payload extension.
func RunJsonExample(dest io.Writer) *appfault.AppError {
	account := UserAccount{
		Id:       "acc-901",
		Username: "alim.karim",
		Email:    "alim@riseup.asia",
		Role:     "lead-architect",
		IsActive: true,
	}

	// 1. Multi-source ingestion: From structured Go object
	resFromPayload := streamwriter.JsonSource.FromPayload(account)
	if !resFromPayload.IsValid() {
		return resFromPayload.AppError()
	}

	// 2. Ingestion from raw bytes
	rawJSON := []byte(`{"id":"acc-902","username":"sarah.connor","email":"sarah@riseup.asia","role":"devops","isActive":true}`)
	resFromBytes := streamwriter.JsonSource.FromBytes(rawJSON)
	if !resFromBytes.IsValid() {
		return resFromBytes.AppError()
	}

	// 3. Ingestion from JSON string
	resFromString := streamwriter.JsonSource.FromString(`{"id":"acc-903","username":"john.doe","role":"qa"}`)
	if !resFromString.IsValid() {
		return resFromString.AppError()
	}

	// 4. Ingestion from streaming io.Reader
	reader := strings.NewReader(`{"id":"acc-904","username":"alex.chen","role":"security"}`)
	resFromReader := streamwriter.JsonSource.FromReader(reader)
	if !resFromReader.IsValid() {
		return resFromReader.AppError()
	}

	// 5. Ingestion from on-demand lazy serializer closure
	resFromSerializer := streamwriter.JsonSource.FromSerializer(func() ([]byte, *appfault.AppError) {
		return []byte(`{"id":"acc-905","username":"elena.rostova","role":"sre"}`), nil
	})
	if !resFromSerializer.IsValid() {
		return resFromSerializer.AppError()
	}

	// 6. Formatting: Pretty and Compact JSON outputs
	prettyOut := resFromPayload.Pretty()
	compactOut := resFromPayload.Compact()
	fmt.Fprintf(dest, "--- Pretty JSON Output ---\n%s\n", prettyOut)
	fmt.Fprintf(dest, "--- Compact JSON Output ---\n%s\n", compactOut)

	// 7. Typed Unmarshaling: Method Unmarshal and Generic UnmarshalAs
	var decodedAcc UserAccount
	appErr := resFromBytes.Unmarshal(&decodedAcc)
	if appErr != nil {
		return appErr
	}

	directAcc, appErr := streamwriter.UnmarshalAs[UserAccount](resFromBytes)
	if appErr != nil {
		return appErr
	}
	fmt.Fprintf(dest, "--- Unmarshaled Account: %s (%s) ---\n", directAcc.Username, directAcc.Role)

	// 8. Type-Casting: Convert between matching structures without manual mappings
	type PublicProfile struct {
		Id       string `json:"id"`
		Username string `json:"username"`
	}
	profileRes := streamwriter.Cast[PublicProfile](account)
	if !profileRes.IsValid() {
		return profileRes.AppError()
	}

	var directTarget PublicProfile
	appErr = streamwriter.JsonSource.Cast(account, &directTarget)
	if appErr != nil {
		return appErr
	}
	fmt.Fprintf(dest, "--- Casted Public Profile: %s [%s] ---\n", directTarget.Username, directTarget.Id)

	// 9. Extended JsonPayloadResult: Embedding JsonResult with strongly-typed payload T
	payloadResult := streamwriter.WithPayload(resFromPayload, account)
	fmt.Fprintf(dest, "--- Extended JsonPayloadResult: %s | IsValid: %v | StatusCode: %d ---\n",
		payloadResult.Payload().Username, payloadResult.IsValid(), payloadResult.StatusCode())

	// 10. Scoped factory: JsonSourceOf[T] produces JsonPayloadResult[T]
	scopedResult := streamwriter.JsonSourceOf[UserAccount]().FromPayload(account)
	fmt.Fprintf(dest, "--- Scoped Factory Result: %s | Bytes: %d ---\n",
		scopedResult.Payload().Email, scopedResult.Len())

	return nil
}

// RunStreamerExample demonstrates LockedStreamer, LocklessStreamer, and atomic PluggableWriter batching.
func RunStreamerExample(dest io.Writer) *appfault.AppError {
	ctx := context.Background()

	// 1. Thread-safe LockedStreamer with sync.Locker support
	lockedStreamer := streamwriter.NewLockedStreamer[OrderEvent](streamwriter.LockedOptions[OrderEvent]{
		Name:        "order-event-streamer",
		Destination: dest,
	})

	// Stream individual event
	appErr := lockedStreamer.Stream(ctx, OrderEvent{
		OrderId:   "ord-1001",
		AccountId: "acc-901",
		Amount:    159.99,
		ItemCount: 3,
	})
	if appErr != nil {
		return appErr
	}

	// 2. Concurrent safety demonstration: Multiple goroutines writing through LockedStreamer
	var wg sync.WaitGroup
	goroutines := 8
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_ = lockedStreamer.Stream(ctx, OrderEvent{
				OrderId:   fmt.Sprintf("ord-concurrent-%d", idx),
				AccountId: "acc-901",
				Amount:    float64(idx * 25),
				ItemCount: idx + 1,
			})
		}(i)
	}
	wg.Wait()

	// 3. LocklessStreamer for high-throughput single-producer scenarios
	memBuf := &bytes.Buffer{}
	locklessStreamer := streamwriter.NewLocklessStreamer[OrderEvent](streamwriter.LocklessOptions[OrderEvent]{
		Name:        "fast-memory-streamer",
		Destination: memBuf,
	})
	appErr = locklessStreamer.Stream(ctx, OrderEvent{
		OrderId:   "ord-fast-001",
		AccountId: "acc-902",
		Amount:    999.00,
		ItemCount: 1,
	})
	if appErr != nil {
		return appErr
	}

	// 4. PluggableWriter with compound atomic batch locking
	batchWriter := streamwriter.NewPluggableWriter[string](streamwriter.WriterOptions[string]{
		Name:     "batch-writer",
		Streamer: streamwriter.NewLockedStreamer[string](streamwriter.LockedOptions[string]{Destination: dest}),
	})

	// Atomic batch write under explicit Lock() / Unlock()
	batchWriter.Lock()
	_ = batchWriter.Write(ctx, "=== BEGIN BATCH TRANSACTION ===")
	_ = batchWriter.Write(ctx, "STEP 1: Reserve inventory items")
	_ = batchWriter.Write(ctx, "STEP 2: Capture card authorization")
	_ = batchWriter.Write(ctx, "STEP 3: Dispatch shipping notification")
	_ = batchWriter.Write(ctx, "=== COMMIT BATCH TRANSACTION ===")
	batchWriter.Unlock()

	// 5. Runtime Method Swapping: Hot-swap write behavior dynamically
	batchWriter.SetWriteMethod(func(ctx context.Context, w *streamwriter.PluggableWriter[string], payload string) *appfault.AppError {
		outDest := w.Destination()
		if outDest == nil {
			outDest = os.Stdout
		}
		_, err := fmt.Fprintf(outDest, "[%s][swapped] %s\n", w.Name(), payload)
		if err != nil {
			return appfault.Wrap(errtype.IO, err, "swapped write failed")
		}
		return nil
	})
	_ = batchWriter.Write(ctx, "Message sent via runtime swapped method")

	return nil
}
