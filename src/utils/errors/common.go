package errors

// 说明：
//
//	错误码主要用于服务返回，方便客户端展示错误信息以及后续的统计
//	每个服务对应一个错误文件 分两部分，第一部分为错误码，
//	第二部分为错误码对应的错误（ps：错误信息需要动态生成的可不放在这里，但是错误码要存放在这里）
//
// 命名格式：
//
//	错误码以ErrorCode开头，紧接服务名，最后加上错误描述
//	错误信息以Error开头，紧接服务名，最后加上错误描述
//
// tt.api 使用 [30000-30100)
// 公共错误定义 [30100-31000)
const (
	CodeOK                          = 0
	CodeSuccess                     = 0
	ErrCodeRedis                    = 30101 //redis错误
	ErrCodeMysql                    = 30102 //mysql错误
	ErrCodeFrequency                = 30103
	ErrCodeParam                    = 30104
	ErrCodeNoMatchObj               = 30105
	ErrCodeBusy                     = 30106
	ErrCodeRedisLockTimeOut         = 30107 //redis错误
	ErrCodeJsonUnmarshal            = 30108
	ErrCodeJsonMarshal              = 30109
	ErrCodeMaintain                 = 39999 //服务器维护
	ErrCodeImageFilterNoResponse    = 30107
	ErrCodeImageFilterResponseError = 30108
	ErrCodeTargetVersionTooLow      = 30109
	ErrCodeWaitAMoment              = 30110
	ErrCodeAINotFace                = 30111
	ErrCodeUploadFileTooLarge       = 30112
	ErrCodeUploadFileEmpty          = 30113
	ErrCodeUploadFileFail           = 30114
	ErrCodeUploadTokenFail          = 30115
)

var (
	Success                = New("", "success", CodeSuccess)
	OK                     = New("", "ok", CodeOK)
	ErrRedis               = New("", "系统繁忙", ErrCodeRedis) //redis错误用系统繁忙替代
	ErrMysql               = New("", "系统繁忙", ErrCodeMysql) //mysql错误用系统繁忙替代
	ErrFrequency           = New("", "请求过于频繁", ErrCodeFrequency)
	ErrParam               = New("", "参数错误", ErrCodeParam)
	ErrBusy                = New("", "系统繁忙", ErrCodeBusy)      //一些极少出现，且不适合返回的错误
	ErrMaintain            = New("", "服务器维护", ErrCodeMaintain) // 服务器维护
	ErrNoMatchObj          = New("", "调用对象不匹配", ErrCodeNoMatchObj)
	ErrJsonUnmarshal       = New("", "JSON反序列化失败", ErrCodeJsonUnmarshal)
	ErrJsonMarshal         = New("", "JSON序列化失败", ErrCodeJsonMarshal)
	ErrRedisLockTimeOut    = New("", "系统繁忙", ErrCodeRedisLockTimeOut) //redis锁超时错误用替代
	ErrTargetVersionTooLow = New("", "对方版本过低，暂不支持此功能", ErrCodeTargetVersionTooLow)
	ErrAINotFace           = New("", "未检测到人脸", ErrCodeAINotFace)
	ErrUploadFileTooLarge  = New("", "上传文件太大", ErrCodeUploadFileTooLarge)
	ErrUploadFileEmpty     = New("", "上传文件不存在", ErrCodeUploadFileEmpty)
	ErrUploadFileFail      = New("", "上传文件失败", ErrCodeUploadFileFail)
	ErrUploadTokenFail     = New("", "获取上传token失败", ErrCodeUploadTokenFail)
)
