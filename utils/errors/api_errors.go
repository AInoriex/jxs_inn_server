package errors

type CodeMsg struct {
	Code int32
	Msg  string
}

const (
	ErrCodeDigitalHumanNotGenerate     = 30000
	ErrCodeAiGenerateDigitalHumanFail  = 30001
	ErrCodeArgsFail                    = 30002
	ErrCodeAppConfig                   = 30003
	ErrCodeAccessToken                 = 30004
	ErrCodeAccessTokenExpireTime       = 30005
	ErrCodeCacheNotExist               = 30006
	ErrCodeFrameSequenceNotUpload      = 30007
	ErrCodeFrameSequenceUploadArgsFail = 30008
	ErrCodePushFrameSequenceFail       = 30009
	ErrCodePushFrameSequenceArgsFail   = 30010
	ErrCodeDbQueryFail                 = 30011
	ErrCodeDbQueryNotFound             = 30012
	ErrCodeApiParamSignNotPass         = 30013
	ErrCodeDboperationFail             = 30014
)

var (
	ErrDbQueryFail         = New("api", "查询失败", ErrCodeDbQueryFail)
	ErrDbQueryNotFound     = New("api", "查询为空", ErrCodeDbQueryNotFound)
	ErrApiParamSignNotPass = New("api", "参数签名不通过", ErrCodeApiParamSignNotPass)
	ErrDboperationFail     = New("api", "操作失败", ErrCodeDboperationFail)
)
