 # Error Resolution: ZIP File Not Finalized Before Return

**Version:** 3.2.0  
**Last Updated:** 2026-04-16

 
 ## Issue ID: ZIP-FINALIZE-001
 
 ## Summary
 ZIP files created with Go's `archive/zip` package were returned before being properly finalized, causing upload failures with "Could not find plugin file after extraction" errors.
 
 ## Root Cause
 
 The `createFullZip` and `createSelectiveZip` functions used `defer` statements to close the ZIP writer and file:
 
 ```go
 // WRONG: defer runs AFTER return, but caller may read file immediately
 zipFile, err := os.Create(absZipPath)
 defer zipFile.Close()
 
 zipWriter := zip.NewWriter(zipFile)
 defer zipWriter.Close()
 
 // ... write files ...
 
 return absZipPath, nil  // ZIP not finalized yet!
 ```
 
 The critical issue: `zip.Writer.Close()` writes the **central directory** to the ZIP file, which is required for the archive to be valid. With `defer`, this happens *after* the function returns, creating a race condition where the caller reads an incomplete file.
 
 ## Symptoms
 
 1. Upload API returns 500 with "Could not find plugin file after extraction"
 2. ZIP file exists but is corrupted or incomplete
 3. Manual extraction of the ZIP fails or shows missing files
 4. Intermittent failures (depends on timing)
 
## Solution

Explicitly close the ZIP writer and file BEFORE returning the path, and use maximum compression:

```go
import "wp-plugin-publish/pkg/ziputil"

// CORRECT: Close explicitly before return, with max compression
zipFile, err := os.Create(absZipPath)

if err != nil {
    return "", err
}
// NO defer for zipFile.Close()

zipWriter := zip.NewWriter(zipFile)
ziputil.RegisterBestCompression(zipWriter) // flate.BestCompression (level 9)
// NO defer for zipWriter.Close()

// ... write files ...

// CRITICAL: Close writer FIRST (writes central directory), then file
if err := zipWriter.Close(); err != nil {
    zipFile.Close()
    pathutil.Remove(absZipPath)

    return "", err
}

if err := zipFile.Close(); err != nil {
    pathutil.Remove(absZipPath)

    return "", err
}

// Verify file exists and has content
isValid := pathutil.IsFileValid(absZipPath)

if !isValid {
    pathutil.Remove(absZipPath)

    return "", apperror.New(apperror.ErrZipInvalid, "zip file invalid after creation")
}

return absZipPath, nil
```
 
 ## Key Rules
 
 1. **Never use `defer` for ZIP finalization** when the returned path will be immediately used
 2. **Close `zip.Writer` before `os.File`** - the writer writes to the file
 3. **Verify file existence and size** before returning the path
 4. **Clean up on any error** - remove partial ZIP files
 
 ## Related Issue: Temp File Cleanup
 
 A secondary issue was that temp ZIP files were deleted even on publish failure:
 
 ```go
 // WRONG: Deletes ZIP on failure, preventing debugging
 defer func() {
     if zipPath != "" && !options.KeepZipFiles {
         pathutil.Remove(zipPath)
     }

 }()
 ```
 
 **Solution**: Track success/failure and ALWAYS keep ZIP on failure:
 
 ```go
 publishFailed := false
 
 defer func() {
     if publishFailed {
         // ALWAYS keep for debugging
         log.Info("Keeping ZIP for debugging (publish failed)")

         return
     }

      if !options.KeepZipFiles {
          pathutil.Remove(zipPath)
      }

 }()
 
 // ... later, on failure:
 publishFailed = true

 return result, nil
 ```
 
 ## Affected Files
 
 - `backend/internal/services/publish/service.go` - `createFullZip`, `createSelectiveZip`, cleanup logic
 
 ## Prevention Checklist
 
 - [ ] Any function creating a ZIP that returns a path must explicitly close before return
 - [ ] ZIP writer closed before file handle
 - [ ] File existence/size verified before returning path
 - [ ] Temp files preserved on failure for debugging
 - [ ] Compare behavior with working PowerShell script when debugging upload issues
 
 ## References
 
 - PowerShell reference: `wp-plugins/scripts/upload-plugin.ps1` (working implementation)
 - Go zip package: https://pkg.go.dev/archive/zip
