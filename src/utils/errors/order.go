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
// user错误码 [32000,33000)
const (
	ErrorCodeUserPayUnknown int32 = 32001
	ErrorCodeUserNotPay     int32 = 32002
	ErrorCodeUserPayFailed  int32 = 32003
	ErrorCodeUserPayTimeout int32 = 32004
)

var (
	ErrorUserPayUnknown = New("", "未知错误，请联系管理员", ErrorCodeUserPayUnknown)
	ErrorUserNotPay     = New("", "订单未支付", ErrorCodeUserNotPay)
	ErrorUserPayFailed  = New("", "订单支付失败，如有问题请联系管理员", ErrorCodeUserPayFailed)
	ErrorUserPayTimeout = New("", "订单支付超时，请重新下单", ErrorCodeUserPayTimeout)
)
