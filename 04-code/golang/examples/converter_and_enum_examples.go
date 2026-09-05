package examples

import (
	"context"
	"fmt"
	"path/filepath"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/fileutil"
	"coding-guidelines/common/pkg/streamwriter"
)

// CustomerProfile represents a customer model for converter and enum demonstrations.
type CustomerProfile struct {
	ID     string                   `json:"id"`
	Name   string                   `json:"name"`
	Email  string                   `json:"email"`
	Status errtype.ProcessStateType `json:"status"`
	Tier   string                   `json:"tier"`
}

// RunBytesLifecycleExample demonstrates the full lifecycle of Bytes[T] envelope:
// creation, status tracking, null checking, deep cloning, and safe concatenation.
func RunBytesLifecycleExample() (streamwriter.Bytes[CustomerProfile], *appfault.AppError) {
	profile := CustomerProfile{
		ID:     "cust-100",
		Name:   "Alice Smith",
		Email:  "alice@example.com",
		Status: errtype.ProcessStateRunning,
		Tier:   "Gold",
	}

	rawBytes := []byte(`{"id":"cust-100","name":"Alice Smith","email":"alice@example.com","status":"running","tier":"Gold"}`)

	// 1. Construct Bytes envelope
	env := streamwriter.NewBytes(rawBytes, profile)
	if !env.IsSuccess() {
		return env, appfault.New(errtype.Generic, "expected successful Bytes envelope")
	}

	// 2. Clone envelope immutably
	cloned := env.Clone()
	if cloned.Len() != env.Len() {
		return env, appfault.New(errtype.Generic, "cloned length mismatch")
	}

	// 3. Null and zero checks
	emptyEnv := streamwriter.Bytes[CustomerProfile]{}
	if !emptyEnv.IsEmpty() || !emptyEnv.IsNull() {
		return env, appfault.New(errtype.Generic, "empty envelope must be empty and null")
	}

	// 4. Safe concat
	concatenated := env.Concat(cloned)
	if concatenated.Len() != env.Len()*2 {
		return env, appfault.New(errtype.Generic, "concatenation length mismatch")
	}

	return env, nil
}

// RunJsonResultExample demonstrates JsonResult unmarshaling, pretty printing, and typecasting.
func RunJsonResultExample() (string, *appfault.AppError) {
	rawJSON := `{"id":"cust-200","name":"Bob Jones","email":"bob@example.com","status":"pending","tier":"Silver"}`

	// 1. Build JsonResult via JsonSource factory
	jr := streamwriter.JsonSource.FromString(rawJSON)
	if !jr.IsValid() {
		return "", jr.Fault()
	}

	// 2. Pretty print JSON
	prettyOutput := jr.Pretty()
	if len(prettyOutput) == 0 {
		return "", appfault.New(errtype.Serialization, "failed to format pretty JSON")
	}

	// 3. Unmarshal directly into domain struct
	var profile CustomerProfile
	unmarshalErr := jr.Unmarshal(&profile)
	if unmarshalErr != nil {
		return "", unmarshalErr
	}

	return prettyOutput, nil
}

// RunReflectConverterExample demonstrates dynamic unmarshaling, pointer reduction, and reflection inspection.
func RunReflectConverterExample() (*CustomerProfile, *appfault.AppError) {
	rawJSON := []byte(`{"id":"cust-300","name":"Carol White","email":"carol@example.com","status":"completed","tier":"Platinum"}`)

	// 1. Dynamic unmarshaling into pointer using Reflect singleton
	var profile CustomerProfile
	appErr := streamwriter.Reflect.UnmarshalTo(rawJSON, &profile)
	if appErr != nil {
		return nil, appErr
	}

	// 2. Generic typed unmarshaler
	genericProfile, genErr := streamwriter.UnmarshalToType[CustomerProfile](rawJSON)
	if genErr != nil {
		return nil, genErr
	}

	if genericProfile.ID != profile.ID {
		return nil, appfault.New(errtype.Generic, "generic and reflect unmarshaling mismatch")
	}

	// 3. Pointer unwrapping (***CustomerProfile -> CustomerProfile)
	p1 := &profile
	p2 := &p1
	p3 := &p2
	reduced := streamwriter.Reflect.ReducePointer(p3)
	unwrapped, ok := reduced.(CustomerProfile)
	if !ok || unwrapped.ID != "cust-300" {
		return nil, appfault.New(errtype.Generic, "pointer reduction failed")
	}

	return &profile, nil
}

// RunEnumOperationsExample demonstrates BaseEnum, NumberEnum, and generic ToEnum helpers.
func RunEnumOperationsExample() *appfault.AppError {
	// 1. String-backed BaseEnum (ProcessStateType)
	state := errtype.ProcessStateRunning
	if state.Name() != "Running" || state.ValueString() != "Running" {
		return appfault.New(errtype.Generic, fmt.Sprintf("invalid string enum: %s", state))
	}

	// 2. Lookup via generic ToEnum helper
	allStates := errtype.AllProcessStates()
	foundState, ok := errtype.ToEnum("pending", allStates)
	if !ok || foundState != errtype.ProcessStatePending {
		return appfault.New(errtype.Generic, "ToEnum lookup failed for pending")
	}

	// 3. Number-backed NumberEnum (LogLevelType)
	level := errtype.LogLevelWarn
	if level.Code() != 3 || level.Name() != "Warn" {
		return appfault.New(errtype.Generic, fmt.Sprintf("invalid number enum: %s", level))
	}

	allLevels := errtype.AllLogLevels()
	foundLevel, ok := errtype.ToEnum("Error", allLevels)
	if !ok || foundLevel != errtype.LogLevelError {
		return appfault.New(errtype.Generic, "ToEnum lookup failed for Error")
	}

	return nil
}

// RunFileWriterAndAppenderExample demonstrates behavior shifting in FileWriter and continuous appending in FileAppender.
func RunFileWriterAndAppenderExample(baseDir string) *appfault.AppError {
	ctx := context.Background()
	writerPath := filepath.Join(baseDir, "output-strategy.txt")
	appenderPath := filepath.Join(baseDir, "journal", "audit.log")

	// 1. Behavior-shifting FileWriter (Direct -> Atomic -> Truncate)
	writer := fileutil.NewFileWriterEngine(writerPath)
	if appErr := writer.WriteString(ctx, "Version 1: Direct Write\n"); appErr != nil {
		return appErr
	}

	// Shift behavior to Atomic swap
	writer.SetMode(fileutil.FileWriteModeAtomic)
	if appErr := writer.WriteString(ctx, "Version 2: Atomic Swap\n"); appErr != nil {
		return appErr
	}

	// Shift behavior to Truncate with fsync
	writer.SetMode(fileutil.FileWriteModeTruncate).SetSyncOnWrite(true)
	if appErr := writer.WriteString(ctx, "Version 3: Truncated Clean State\n"); appErr != nil {
		return appErr
	}

	// 2. Continuous FileAppender (auto directory creation, auto-sync, bytes counter)
	appender := fileutil.NewFileAppender(appenderPath, fileutil.FilePermStandard)
	appender.SetAutoSync(true)

	if appErr := appender.AppendString(ctx, "LOG ENTRY 1: Service started\n"); appErr != nil {
		return appErr
	}

	if appErr := appender.AppendString(ctx, "LOG ENTRY 2: Config loaded\n"); appErr != nil {
		return appErr
	}

	if appender.BytesAppended() <= 0 {
		return appfault.New(errtype.IO, "expected positive bytes appended")
	}

	if appErr := appender.Close(); appErr != nil {
		return appErr
	}

	return nil
}

// RunBoundFileWriterExample demonstrates the file-specific BoundFileWriter with automatic locking,
// appending, auto-closing, and transactional lock blocks.
func RunBoundFileWriterExample(baseDir string) *appfault.AppError {
	ctx := context.Background()
	filePath := filepath.Join(baseDir, "bound-state.txt")

	// 1. Create file-specific writer
	writer := fileutil.NewBoundFileWriter(filePath)

	// 2. Automatic lock write and append
	if err := writer.WriteString(ctx, "State: Initialized\n"); err != nil {
		return err
	}

	if err := writer.AppendString(ctx, "Event: Connection established\n"); err != nil {
		return err
	}

	// 3. AutoClose mode: closes handle immediately after write
	writer.SetAutoClose(true)
	if err := writer.AppendString(ctx, "Audit: Checkpoint recorded\n"); err != nil {
		return err
	}

	// 4. Transactional lock batch: multiple writes under single lock
	err := writer.WithLock(ctx, func(w *fileutil.BoundFileWriter) *appfault.AppError {
		_ = w.AppendLocked(ctx, []byte("--- Batch Header ---\n"))
		_ = w.AppendLocked(ctx, []byte("Action: Sync A\n"))
		_ = w.AppendLocked(ctx, []byte("Action: Sync B\n"))
		_ = w.AppendLocked(ctx, []byte("--- Batch Footer ---\n"))

		return nil
	})
	if err != nil {
		return err
	}

	// 5. Explicit write and close
	if err := writer.AppendAndClose(ctx, []byte("Final: Terminated\n")); err != nil {
		return err
	}

	return nil
}
