# RCA: `captureStderr` Pipe Drain Race Condition in `cmdchromeprofile/testhelpers_test.go`

## 1. Why it happened
When `cmdchromeprofile` was modularized into its own standalone package, its test suite needed a `captureStderr` helper. An un-synchronized implementation was placed into `testhelpers_test.go` that called `_ = r.Close()` before reading from `outC`, and lacked an `io` mutex. This caused premature closure of the pipe reader while the draining goroutine was executing, resulting in empty captured stderr and test assertion failure in `TestHandleChromeFileOpenErrorSkipsLockFile`.

## 2. How it happened
1. `TestHandleChromeFileOpenErrorSkipsLockFile` invoked `captureStderr(t, func() { ... })`.
2. `captureStderr` spawned a goroutine copying from pipe reader `r` to `bytes.Buffer`.
3. The test body executed `handleChromeFileOpenError(...)`, writing a warning message to `os.Stderr`.
4. `captureStderr` immediately closed writer `w`, restored `os.Stderr`, and executed `_ = r.Close()`.
5. On Windows, closing the reader before `outC` had finished reading from `r` caused `io.Copy` to return with a closed-handle error before buffering the warning message.
6. The test asserted `strings.Contains(stderr, "skipped volatile Chrome lock file")`, which failed because `stderr` was `""`.

## 3. Root Cause
- File: `gitmap/cmdchromeprofile/testhelpers_test.go`
- Lines: 11-33
- Premature `r.Close()` invocation before waiting for `<-outC`, lack of channel buffering `make(chan string, 1)`, and absence of mutex synchronization (`stdIOMutex`).

## 4. Code Fix
Updated `captureStderr` in `gitmap/cmdchromeprofile/testhelpers_test.go` to match the canonical synchronized pattern in `gitmap/cmd/capturestderr_testhelper_test.go`:
- Added `stdIOMutex sync.Mutex` protection across stdout/stderr redirection.
- Used buffered channel `make(chan string, 1)`.
- Waited for `res := <-outC` before closing `r`.

### Before:
```go
	fn()
	_ = w.Close()
	os.Stderr = origStderr
	_ = r.Close()

	return <-outC
```

### After:
```go
	stdIOMutex.Lock()
	defer stdIOMutex.Unlock()
    ...
	fn()
	_ = w.Close()
	os.Stderr = origStderr
	res := <-outC
	_ = r.Close()

	return res
```
