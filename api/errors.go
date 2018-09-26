package api

const Status_OK = 0

const (
	ERR_CODE_DISPLAY_ERROR                 = 1    // 自定义提示错误，需要展示给用户
	ERR_CODE_BAD_REQUEST                   = 400  // 参数错误
	ERR_CODE_FORBIDDEN                     = 403  // 请求被拒绝
	ERR_CODE_BAD_NOT_FOUND                 = 404  // 资源不存在
	ERR_CODE_INTERNAL_ERROR                = 500  // 内部错误
)

const (
	ERR_MSG_BAD_REQUEST = "参数错误"
)
