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
// user错误码 [33000,34000)
const (
	ErrorCodeStreamServiceUnknownError int32 = 33001
	ErrorCodeStreamFileUploadFailed    int32 = 33002
	ErrorCodeStreamFileStreamingFailed int32 = 33003
)

var (
	ErrorStreamServiceUnknownError = New("", "流媒体服务发生未知错误，请联系管理员", ErrorCodeStreamServiceUnknownError)
	ErrorStreamFileUploadFailed    = New("", "上传流媒体文件失败", ErrorCodeStreamFileUploadFailed)
	ErrorStreamFileStreamingFailed = New("", "文件流媒体处理失败", ErrorCodeStreamFileStreamingFailed)
)
