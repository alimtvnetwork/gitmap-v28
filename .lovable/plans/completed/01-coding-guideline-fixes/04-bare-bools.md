# Task: Refactor Bare Bools

According to the Coding Guidelines: Golang functions returning bare bools instead of a wrapped Result object should be refactored.

## Files to fix

- d:\wp-work\riseup-asia\gitmap\.lovable\scratch\find_ignored_errs.go:58: return true
- d:\wp-work\riseup-asia\gitmap\.lovable\scratch\find_monolithic.go:39: return true
- d:\wp-work\riseup-asia\gitmap\.lovable\scratch\find_nested_ifs.go:40: return false
- d:\wp-work\riseup-asia\gitmap\.lovable\scratch\find_nested_ifs.go:44: return true
- d:\wp-work\riseup-asia\gitmap\.lovable\scratch\find_nested_ifs.go:49: return false
- d:\wp-work\riseup-asia\gitmap\.lovable\scratch\find_nested_ifs.go:51: return true
- d:\wp-work\riseup-asia\gitmap\.lovable\scratch\find_nested_ifs.go:64: return false
- d:\wp-work\riseup-asia\gitmap\.lovable\scratch\find_nested_ifs.go:67: return true
- d:\wp-work\riseup-asia\gitmap\.lovable\scratch\find_nested_ifs.go:70: return false
- d:\wp-work\riseup-asia\gitmap\.lovable\scratch\find_nested_ifs.go:72: return true
- d:\wp-work\riseup-asia\gitmap\.lovable\scratch\find_nested_ifs.go:77: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\create.go:156: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\create.go:159: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\create.go:162: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\create.go:182: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\create.go:185: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\create.go:188: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\create.go:200: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\create.go:204: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\source.go:67: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\source.go:78: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\source.go:81: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\source.go:84: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\archive\source.go:87: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\depthflag_format_test.go:98: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\depthflag_format_test.go:102: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\execute_dest.go:57: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\execute_dest.go:62: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\prompt.go:15: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\validate.go:66: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\validate.go:70: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\validate.go:74: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\validate.go:87: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\validate.go:90: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\validate.go:107: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\validate.go:120: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\validate.go:124: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\validate.go:136: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\validate.go:142: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonenext\batch.go:168: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\clonenext\batch.go:172: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonenext\repodetect.go:40: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonenext\repodetect.go:47: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\clonenext\repodetect.go:51: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonenow\execute_idempotent.go:158: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonenow\execute_idempotent.go:161: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\clonenow\execute_idempotent.go:174: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonenow\parsetext.go:119: return true
- d:\wp-work\riseup-asia\gitmap\gitmap\clonenow\parsetext.go:123: return false
- d:\wp-work\riseup-asia\gitmap\gitmap\clonenow\parse_schema.go:121: return true
