package constants

// ErrorCodeType represents strongly-typed error codes for operations and commands.
type ErrorCodeType string

const (
	ErrCodeGenericExecution     ErrorCodeType = "E9000"
	ErrCodeSSHKeyEmpty          ErrorCodeType = "E9001"
	ErrCodeSSHKeyMissing        ErrorCodeType = "E9002"
	ErrCodeWorkgroupAddUsage    ErrorCodeType = "E9003"
	ErrCodeWorkgroupBindUsage   ErrorCodeType = "E9004"
	ErrCodeWorkgroupNoBoundKey  ErrorCodeType = "E9005"
	ErrCodeWorkgroupDeleteUsage ErrorCodeType = "E9006"
)
