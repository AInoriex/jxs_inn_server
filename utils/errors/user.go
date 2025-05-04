package errors

//说明：
//      错误码主要用于服务返回，方便客户端展示错误信息以及后续的统计
//      每个服务对应一个错误文件 分两部分，第一部分为错误码，
//      第二部分为错误码对应的错误（ps：错误信息需要动态生成的可不放在这里，但是错误码要存放在这里）
//命名格式：
//      错误码以ErrorCode开头，紧接服务名，最后加上错误描述
//      错误信息以Error开头，紧接服务名，最后加上错误描述
//user错误码 [31000,32000)
const (
	ErrorCodeRelogin                  int32 = 31001
	ErrorCodeUserArgsFail             int32 = 31002
	ErrorCodeUserAppCfgFail           int32 = 31003
	ErrorCodeUserTokenFail            int32 = 31004
	ErrorCodeUserTokenExpFail         int32 = 31005
	ErrorCodeUserCacheNotExist        int32 = 31006
	ErrorCodeUserAttrUploadFail       int32 = 31007
	ErrorCodeUserLoginFail            int32 = 31010
	ErrorCodeCreateUserFail           int32 = 31011
	ErrorCodeUserSuitsNotExist        int32 = 31012
	ErrorCodeSaveUserSuitsFail        int32 = 31013
	ErrorCodeSaveCardFail             int32 = 31014
	ErrorCodeGenCartoonFail           int32 = 31015
	ErrorCodeSendSmsCodeFastFail      int32 = 31020
	ErrorCodeSendSmsCodeFail          int32 = 31021
	ErrorCodeCheckSmsCodeFail         int32 = 31022
	ErrorCodeGetAppUserFail           int32 = 31023
	ErrorCodeGetWechatInfoFail        int32 = 31024
	ErrorCodeGetUserBindFail          int32 = 31025
	ErrorCodeUserHasBindFail          int32 = 31026
	ErrorCodeUserBindFail             int32 = 31027
	ErrorCodeSendSmsDeviceLimitFail   int32 = 31028
	ErrorCodeUserRegisterFail         int32 = 31029
	ErrorCodeUserUpdateFail           int32 = 31030
	ErrorCodeSmsCodeRefreshFail       int32 = 31031
	ErrorCodeGetFlashPhoneFail        int32 = 31032
	ErrorCodeGetDigitalHumanFail      int32 = 31033
	ErrorCodeFastLoginDeviceLimitFail int32 = 31034
	ErrorCodeText2PaintingFail        int32 = 31035
	ErrorCodeImgModerationFail        int32 = 31036
	ErrorCodeTextModerationFail       int32 = 31037
	ErrorCodeGetWechatOpenidFail      int32 = 31038
	ErrorCodeShopUserLoginFail        int32 = 31039
	ErrorCodeShopUserDeviceListFail   int32 = 31040
	ErrorCodeShopUserUnAuthorization  int32 = 31041
	ErrorCodeShopUserUnused           int32 = 31042
	ErrorCodePhoneFail                int32 = 31043
	//ErrorCodeUserExist int32 = 30012
)

var (
	ErrorUserLoginFail            = New("", "登录失败", ErrorCodeUserLoginFail)
	ErrorUserArgsFail             = New("user", "参数有误", ErrorCodeUserArgsFail)
	ErrorUserAppCfgFail           = New("user", "配置有误", ErrorCodeUserAppCfgFail)
	ErrorUserTokenFail            = New("user", "access_token有误", ErrorCodeUserTokenFail)
	ErrorCacheNotExist            = New("user", "缓存不存在", ErrorCodeUserCacheNotExist)
	ErrorUserTokenExpFail         = New("user", "access_token过期", ErrorCodeUserTokenExpFail)
	ErrorCreateUserFail           = New("user", "创建用户失败", ErrorCodeCreateUserFail)
	ErrorUserSuitsNotExist        = New("user", "用户套装不存在", ErrorCodeUserSuitsNotExist)
	ErrorSaveUserSuitsFail        = New("user", "更新用户失败", ErrorCodeSaveUserSuitsFail)
	ErrorRelogin                  = New("user", "需要重新登录", ErrorCodeRelogin)
	ErrorUserAttrUploadFail       = New("user", "上传附件失败", ErrorCodeUserAttrUploadFail)
	ErrorSaveCardFail             = New("user", "保存卡片失败", ErrorCodeSaveCardFail)
	ErrorGenCartoonFail           = New("user", "生成卡通失败", ErrorCodeGenCartoonFail)
	ErrorSendSmsCodeFastFail      = New("user", "短信发送太频繁", ErrorCodeSendSmsCodeFastFail)
	ErrorSendSmsDeviceLimitFail   = New("user", "短信发送频繁，设备受限", ErrorCodeSendSmsDeviceLimitFail)
	ErrorSendSmsCodeFail          = New("user", "短信发送失败", ErrorCodeSendSmsCodeFail)
	ErrorCheckSmsCodeFail         = New("user", "短信验证码不正确", ErrorCodeCheckSmsCodeFail)
	ErrorSmsCodeRefreshFail       = New("user", "短信验证码已过期，请重新验证", ErrorCodeSmsCodeRefreshFail)
	ErrorGetAppUserFail           = New("user", "app用户不存在", ErrorCodeGetAppUserFail)
	ErrorGetWechatInfoFail        = New("user", "获取微信用户信息失败", ErrorCodeGetWechatInfoFail)
	ErrorGetUserBindFail          = New("user", "未绑定手机", ErrorCodeGetUserBindFail)
	ErrorUserHasBindFail          = New("user", "用户已绑定", ErrorCodeUserHasBindFail)
	ErrorUserBindFail             = New("user", "用户绑定失败", ErrorCodeUserBindFail)
	ErrorUserRegisterFail         = New("user", "注册失败", ErrorCodeUserRegisterFail)
	ErrorUserUpdateFail           = New("user", "更新用户失败", ErrorCodeUserUpdateFail)
	ErrorGetFlashPhoneFail        = New("user", "手机一键登录失败", ErrorCodeGetFlashPhoneFail)
	ErrorGetDigitalHumanFail      = New("user", "获取数字人信息失败", ErrorCodeGetDigitalHumanFail)
	ErrorFastLoginDeviceLimitFail = New("user", "一键登录频繁，设备受限", ErrorCodeFastLoginDeviceLimitFail)
	ErrorText2PaintingFail        = New("user", "文字绘画失败", ErrorCodeText2PaintingFail)
	ErrorImgModerationFail        = New("user", "这对我来说有点难，你可以试试让我画点别的", ErrorCodeImgModerationFail)  //图片未过审
	ErrorTextModerationFail       = New("user", "这对我来说有点难，你可以试试让我画点别的", ErrorCodeTextModerationFail) //文案未过审
	ErrorGetWechatOpenidFail      = New("user", "微信授权失败", ErrorCodeGetWechatOpenidFail)
	ErrorShopUserLoginFail        = New("user", "登录失败", ErrorCodeShopUserLoginFail)
	ErrorShopUserUnAuthorization  = New("user", "用户未授权", ErrorCodeShopUserUnAuthorization)
	ErrorShopUserDeviceListFail   = New("user", "获取设备列表失败", ErrorCodeShopUserDeviceListFail)
	ErrorShopUserUnused           = New("user", "用户已失效", ErrorCodeShopUserUnused)
	ErrorPhoneFail                = New("user", "手机格式错误", ErrorCodePhoneFail)
)
