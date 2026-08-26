 # Error Resolution: Plugin Activation Endpoint 404

**Version:** 3.2.0  
**Last Updated:** 2026-04-16

 
 ## Issue ID: ACTIVATE-404-001
 
 ## Summary
 Plugin activation failed with 404 because the Go backend was calling a non-existent endpoint `/plugins/{slug}/enable` on the Riseup Asia Uploader WordPress plugin.
 
 ## Root Cause
 
 The Go backend defined endpoint constants for enable/disable/delete operations:
 ```go
 EndpointEnable = "/plugins/%s/enable"   // Does NOT exist in PHP!
 EndpointDisable = "/plugins/%s/disable" // Does NOT exist in PHP!
 EndpointDelete = "/plugins/%s/delete"   // Does NOT exist in PHP!
 ```
 
 But the WordPress PHP plugin (`riseup-asia-uploader.php`) only registers these endpoints:
 - `/status` - Status check
 - `/upload` - Plugin upload (with optional `activate: true`)
 - `/plugins` - List plugins
 - `/export-self` - Export plugin as ZIP
 - `/posts`, `/categories`, `/logs` - Other features
 
 **No separate enable/disable/delete endpoints exist!**
 
 ## Symptoms
 
 ```
 enable plugin via RiseupAsia Uploader (POST /riseup-asia-uploader/v1/plugins/category-generator/enable): status 404
 {"code":"rest_no_route","message":"No route was found matching the URL and request method."}
 ```
 
 ## Solution
 
 1. **Activation during upload**: The `/upload` endpoint already accepts `activate: true` and activates during upload
 2. **Track activation status**: Return `activated: true/false` from upload function
 3. **Skip if already activated**: Check if activation happened during upload before attempting separate activation
 4. **Fallback to WordPress Core API**: If not activated during upload, use `/wp/v2/plugins/{identifier}` with `PUT {status: "active"}`
 
 ### Updated Flow
 
 ```
 Upload Stage:
   POST /riseup-asia-uploader/v1/upload
   Body: { plugin_zip: base64, slug: "...", activate: true }
   Response: { success: true, activated: true }
   
 Activate Stage:
   IF activated == true:
     Skip (already done during upload)
   ELSE:
     Try: Onboard Plugin /enable endpoint
     Fallback: WordPress Core API PUT /wp/v2/plugins/{identifier}
 ```
 
 ## Key Rules
 
 1. **Backend and PHP constants MUST match** - Never define endpoints in Go that don't exist in PHP
 2. **Verify API contracts** - Before calling an endpoint, confirm it's registered in the PHP plugin
 3. **Use activation during upload** - Riseup Asia Uploader handles activation in the upload flow
 4. **Fallback gracefully** - If companion plugin lacks endpoint, fallback to WordPress Core API
 
 ## Prevention Checklist
 
 - [ ] All Go endpoint constants have corresponding PHP route registrations
 - [ ] Endpoint format (path parameters, HTTP methods) matches between Go and PHP
 - [ ] Check `register_routes()` in PHP plugin before using an endpoint in Go
 - [ ] Use constants from same file, never hardcode endpoint strings
 
 ## Files Changed
 
 - `backend/internal/services/publish/service.go` - Track `alreadyActivated` from upload, skip activation if true
 - `backend/internal/wordpress/uploader.go` - Return `activated` status from upload
 
 ## References
 
 - PHP routes: `wp-plugins/riseup-asia-uploader/riseup-asia-uploader.php:147-232`
 - Go constants: `backend/internal/wordpress/constants.go`
 - PHP constants: `wp-plugins/riseup-asia-uploader/includes/constants.php`
