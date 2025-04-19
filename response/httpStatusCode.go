package response

const (
	// Success
	ErrCodeSuccess       = 20001 // Thành công
	ErrCodeRemoveSuccess = 20002 // Xóa thành công

	// Validation Errors
	ErrCodeParamInvalid    = 40001 // Tham số không hợp lệ
	ErrTokenHeadersInvalid = 40002 // Header token không hợp lệ
	ErrEmailAlreadyExists  = 40003 // Email đã tồn tại
	ErrorDataNotExists     = 40004 // Dữ liệu không tồn tại
	ErrorListFailed        = 40005 // Lấy danh sách thất bại
	ErrorInsert            = 40006 // Thêm dữ liệu thất bại
	ErrorUpdate            = 40007 // Cập nhật dữ liệu thất bại
	ErrorDelete            = 40008 // Xóa dữ liệu thất bại

	// Authentication & Authorization
	ErrCodeUserHasExists   = 40101 // Người dùng đã tồn tại
	ErrCodeUserNotRegister = 40102 // Người dùng chưa đăng ký
	ErrCodeAuthenError     = 40103 // Lỗi xác thực
	ErrInvalidOtp          = 40104 // OTP không hợp lệ
	ErrSendOtp             = 40105 // Gửi OTP thất bại
	ErrCodeOtpNotExists    = 40106 // OTP không tồn tại
	ErrCodeOtpNotVerify    = 40107 // OTP chưa được xác minh

	// Redis Errors
	ErrCodeRedisSetFailed = 50001 // Lưu Redis thất bại
	ErrCodeRedisGetFailed = 50002 // Lấy dữ liệu Redis thất bại

	// Two-Factor Authentication (2FA)
	ErrTwoFactorAuthSetUpFailed  = 60001 // Thiết lập xác thực hai yếu tố thất bại
	ErrTwoFactorAuthVerifyFailed = 60002 // Xác minh xác thực hai yếu tố thất bại

	// Kafka Errors
	ErrKafkaSendMessageFailed = 70001 // Gửi tin nhắn Kafka thất bại

	//UUID
	ErrParseUUID = 80002
)

var msg = map[int]string{
	// Success
	ErrCodeSuccess:       "Success",
	ErrCodeRemoveSuccess: "Deletion successful",

	// Validation Errors
	ErrCodeParamInvalid:    "Invalid parameters",
	ErrTokenHeadersInvalid: "Invalid token headers",
	ErrEmailAlreadyExists:  "Email already exists",
	ErrorDataNotExists:     "Data not found",
	ErrorListFailed:        "Failed to retrieve data list",
	ErrorInsert:            "Failed to insert data",
	ErrorUpdate:            "Failed to update data",
	ErrorDelete:            "Failed to delete data",

	// Authentication & Authorization
	ErrCodeUserHasExists:   "User already exists",
	ErrCodeUserNotRegister: "User not registered",
	ErrCodeAuthenError:     "Authentication error",
	ErrInvalidOtp:          "Invalid OTP code",
	ErrSendOtp:             "Failed to send OTP",
	ErrCodeOtpNotExists:    "OTP code does not exist",
	ErrCodeOtpNotVerify:    "OTP code not verified",

	// Redis Errors
	ErrCodeRedisSetFailed: "Failed to store data in Redis",
	ErrCodeRedisGetFailed: "Failed to retrieve data from Redis",

	// Two-Factor Authentication (2FA)
	ErrTwoFactorAuthSetUpFailed:  "Two-factor authentication setup failed",
	ErrTwoFactorAuthVerifyFailed: "Two-factor authentication verification failed",

	// Kafka Errors
	ErrKafkaSendMessageFailed: "Failed to send Kafka message",

	//uuid
	ErrParseUUID: "Failed to parse UUID",
}
