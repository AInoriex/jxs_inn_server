package handler

import (
	"encoding/json"
	"errors"
	"eshop_server/src/common/api"
	"eshop_server/src/common/cache"
	"eshop_server/src/router/dao"
	"eshop_server/src/router/middleware"
	"eshop_server/src/router/model"
	uerrors "eshop_server/src/utils/errors"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/mail"
	"eshop_server/src/utils/uuid"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Title        用户登录
// @Description  邮箱+密码登录
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/auth/login [post]
func UserLogin(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Info("UserLogin 请求参数", zap.String("body", string(req)))

	// TODO 添加IP风控

	// JSON解析
	var reqbody model.UserLoginReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("UserLogin json解析失败, error:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 参数有效性判断
	if !isValidEmail(reqbody.Email) {
		log.Error("UserLogin 邮箱格式无效", zap.String("email", reqbody.Email))
		api.Fail(c, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Detail)
		return
	}
	if !isValidPassword(reqbody.HashedPassword) {
		log.Error("UserLogin 密码格式无效", zap.String("password", reqbody.HashedPassword))
		api.Fail(c, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Detail)
		return
	}

	// 查询用户是否存在
	user, err := dao.GetUserByEmail(reqbody.Email)
	if err != nil {
		log.Error("UserLogin 查询用户失败", zap.Error(err))
		api.Fail(c, uerrors.Parse(uerrors.ErrorUserNotFound.Error()).Code, uerrors.Parse(uerrors.ErrorUserNotFound.Error()).Detail)
		return
	}

	// 验证用户状态
	if user.Status == model.UserStatusBanned {
		log.Error("UserLogin 用户已被禁用", zap.String("user_id", user.Id))
		api.Fail(c, uerrors.Parse(uerrors.ErrorUserBanned.Error()).Code, uerrors.Parse(uerrors.ErrorUserBanned.Error()).Detail)
		return
	}

	// 验证哈希密码是否一致
	if reqbody.HashedPassword != user.Password {
		log.Error("UserLogin 密码不一致", zap.String("请求哈希密码", reqbody.HashedPassword), zap.String("目标哈希密码", user.Password))
		api.Fail(c, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Detail)
		return
	}

	// 生成 JWT Token
	// TODO 新增权限区分生成不同权限token(user/admin)
	tokenString, err := middleware.GenerateToken(user.Id, user.Roles)
	if err != nil {
		log.Error("UserLogin 生成jwt token失败", zap.Error(err))
		api.Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// 缓存保存token
	if err = cache.SaveJxsUserToken(user.Id, tokenString); err != nil {
		log.Errorf("UserLogin 缓存保存JxsUserToken失败, user_id:%v, token:%s, err:%v", user.Id, tokenString, err)
		api.Fail(c, uerrors.Parse(uerrors.ErrRedis.Error()).Code, uerrors.Parse(uerrors.ErrRedis.Error()).Detail)
		return
	}

	// 更新用户最后登录时间
	user.LastLogin = time.Now()
	dao.UpdateUserByField(user, []string{"last_login"})

	// 返回token
	dataMap["token_type"] = middleware.TokenType
	dataMap["access_token"] = tokenString
	api.Success(c, dataMap)
}

// @Title        管理后台登录
// @Description  邮箱+密码登录
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/auth/login [post]
func AdminLogin(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Info("AdminLogin 请求参数", zap.String("body", string(req)))

	// TODO 添加IP风控

	// JSON解析
	var reqbody model.UserLoginReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("AdminLogin json解析失败, error:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 参数有效性判断
	if !isValidEmail(reqbody.Email) {
		log.Error("AdminLogin 邮箱格式无效", zap.String("email", reqbody.Email))
		api.Fail(c, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Detail)
		return
	}
	if !isValidPassword(reqbody.HashedPassword) {
		log.Error("AdminLogin 密码格式无效", zap.String("password", reqbody.HashedPassword))
		api.Fail(c, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Detail)
		return
	}

	// 查询用户是否存在
	user, err := dao.GetUserByEmail(reqbody.Email)
	if err != nil {
		log.Error("AdminLogin 查询用户失败", zap.Error(err))
		api.Fail(c, uerrors.Parse(uerrors.ErrorUserNotFound.Error()).Code, uerrors.Parse(uerrors.ErrorUserNotFound.Error()).Detail)
		return
	}

	// 验证用户状态
	if user.Status == model.UserStatusBanned {
		log.Error("AdminLogin 用户已被禁用", zap.String("user_id", user.Id))
		api.Fail(c, uerrors.Parse(uerrors.ErrorUserBanned.Error()).Code, uerrors.Parse(uerrors.ErrorUserBanned.Error()).Detail)
		return
	}

	// 验证哈希密码是否一致
	if reqbody.HashedPassword != user.Password {
		log.Error("AdminLogin 密码不一致", zap.String("请求哈希密码", reqbody.HashedPassword), zap.String("目标哈希密码", user.Password))
		api.Fail(c, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Detail)
		return
	}

	// 验证管理员权限
	_roleCheckPass := false
	for _, role := range user.Roles {
		if role == model.UserRoleAdmin {
			_roleCheckPass = true
			break
		}
	}
	if !_roleCheckPass {
		log.Errorf("AdminLogin 用户权限不足, user_id:%s, roles:%v", user.Id, user.Roles)
		api.Fail(c, uerrors.Parse(uerrors.ErrorShopUserUnAuthorization.Error()).Code, uerrors.Parse(uerrors.ErrorShopUserUnAuthorization.Error()).Detail)
		return
	}

	// 生成 JWT Token
	// TODO 新增权限区分生成不同权限token(user/admin)
	tokenString, err := middleware.GenerateToken(user.Id, user.Roles)
	if err != nil {
		log.Error("AdminLogin 生成jwt token失败", zap.Error(err))
		api.Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// 缓存保存token
	if err = cache.SaveJxsUserToken(user.Id, tokenString); err != nil {
		log.Errorf("AdminLogin 缓存保存JxsUserToken失败, user_id:%v, token:%s, err:%v", user.Id, tokenString, err)
		api.Fail(c, uerrors.Parse(uerrors.ErrRedis.Error()).Code, uerrors.Parse(uerrors.ErrRedis.Error()).Detail)
		return
	}

	// 更新用户最后登录时间
	user.LastLogin = time.Now()
	dao.UpdateUserByField(user, []string{"last_login"})

	// 返回token
	dataMap["token_type"] = middleware.TokenType
	dataMap["access_token"] = tokenString
	api.Success(c, dataMap)
}

// @Title      	 用户登出
// @Description  登出
// @Produce      json
// @Router       /v1/eshop_api/auth/logout [get]
func UserLogout(c *gin.Context) {
	// 从Header获取Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Warn("UserLogout 未携带Authorization头部信息")
		// 无需返回错误，直接成功
		api.Success(c, nil)
		return
	}

	// 检查Token格式：{TokenType} {TokenString}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == middleware.TokenType) {
		log.Warn("UserLogout 认证头格式错误")
		api.Success(c, nil)
		return
	}

	// 校验Token
	requestToken := parts[1]
	claims, err := middleware.ValidateToken(requestToken)
	if err != nil {
		log.Warn("UserLogout 解析token失败", zap.Error(err))
		api.Success(c, nil)
		return
	}

	// 从Redis中删除token
	if claims.UserId != "" {
		isDel := cache.DelJxsUserToken(claims.UserId)
		if !isDel {
			log.Error("UserLogout 删除缓存JxsUserToken失败", zap.String("user_id", claims.UserId))
		}
	}
	api.Success(c, nil)
}

// @Title        用户注册
// @Description  邮箱+密码注册
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/auth/register [post]
func UserRegister(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Info("UserRegister 请求参数", zap.String("body", string(req)))

	// TODO 添加IP风控

	// JSON解析
	var reqbody model.UserRegisterReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("UserRegister json解析失败, error:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 查询用户是否存在
	user, err := dao.GetUserByEmail(reqbody.Email)
	if err == nil || user.Id != "" {
		log.Error("UserRegister 邮箱已存在，注册失败", zap.String("user_id", user.Id), zap.String("name", user.Name), zap.String("email", user.Email))
		api.Fail(c, uerrors.Parse(uerrors.ErrorRegisterMailExisted.Error()).Code, uerrors.Parse(uerrors.ErrorRegisterMailExisted.Error()).Detail)
		return
	}

	// 创建新用户
	new_user := &model.User{
		Id:        uuid.GetUuid(), // 随机生成用户ID字符串
		Name:      reqbody.Name,
		Email:     reqbody.Email,
		Password:  reqbody.HashedPassword,
		AvatarUrl: "",                           // 默认头像URL为空
		Roles:     []string{model.UserRoleUser}, // 默认普通用户角色
		Status:    model.UserStatusNormal,
	}
	_, err = dao.CreateUser(new_user)
	if err != nil {
		log.Error("UserRegister 创建用户失败", zap.Error(err))
		api.Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// 返回
	log.Infof("UserRegister 用户注册成功，user_id:%v, name:%v, email:%v", new_user.Id, new_user.Name, new_user.Email)
	api.Success(c, dataMap)
}

// @Title        用户注册
// @Description  邮箱+密码+邮箱验证码注册
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/auth/register [post]
func UserRegisterWithVerifyCode(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("UserRegisterWithVerifyCode 请求参数, req:%s", string(req))

	// JSON解析
	var reqbody model.UserRegisterReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("UserRegisterWithVerifyCode json解析失败, error:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 验证码校验
	if reqbody.VerifyCode == "" {
		log.Errorf("UserRegisterWithVerifyCode 验证码为空")
		api.Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":验证码为空")
		return
	}
	flag, verifyCode := cache.GetJxsVerifyMailCode(c.ClientIP(), reqbody.Email)
	if !flag {
		log.Errorf("UserRegisterWithVerifyCode 请求的验证码不存在, clientIp:%v, reqbody.Email:%v", c.ClientIP(), reqbody.Email)
		api.Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":验证码错误")
		return
	}
	if reqbody.VerifyCode != verifyCode {
		log.Errorf("UserRegisterWithVerifyCode 验证码错误, cache.VerifyCode:%v, reqbody.VerifyCode:%v", verifyCode, reqbody.VerifyCode)
		api.Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":验证码错误")
		return
	}
	// 移除缓存验证码
	cache.DelJxsVerifyMailCode(c.ClientIP(), reqbody.Email)

	// 查询用户是否存在
	user, err := dao.GetUserByEmail(reqbody.Email)
	if err == nil || user.Id != "" {
		log.Error("UserRegisterWithVerifyCode 邮箱已存在，注册失败", zap.String("user_id", user.Id), zap.String("name", user.Name), zap.String("email", user.Email))
		api.Fail(c, uerrors.Parse(uerrors.ErrorRegisterMailExisted.Error()).Code, uerrors.Parse(uerrors.ErrorRegisterMailExisted.Error()).Detail)
		return
	}

	// 创建新用户
	new_user := &model.User{
		Id:       uuid.GetUuid(), // 随机生成用户ID字符串
		Name:     reqbody.Name,
		Email:    reqbody.Email,
		Password: reqbody.HashedPassword,
	}
	_, err = dao.CreateUser(new_user)
	if err != nil {
		log.Error("UserRegisterWithVerifyCode 创建用户失败", zap.Error(err))
		api.Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// 返回
	log.Infof("UserRegisterWithVerifyCode 用户注册成功，user_id:%v, name:%v, email:%v", new_user.Id, new_user.Name, new_user.Email)
	api.Success(c, dataMap)
}

// @Title        用户刷新token
// @Description  传入旧token用于获取新token
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/auth/refresh_token [post]
func UserRefreshToken(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Info("UserRefreshToken 请求参数", zap.String("body", string(req)))

	// TODO 添加IP风控

	// JSON解析
	var reqbody model.UserRefreshTokenReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("UserRegister json解析失败, error:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}
	// 校验token是否为空
	if reqbody.OldToken == "" {
		log.Error("UserRefreshToken token为空", zap.String("token", reqbody.OldToken))
		api.Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail)
		return
	}

	// 校验旧token
	claims, err := middleware.GetTokenClaims(reqbody.OldToken)
	if err != nil {
		log.Error("UserRefreshToken 解析token失败", zap.Error(err))
		api.Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// 校验user_id是否有效
	user, err := dao.GetValidUserById(claims.UserId)
	if err != nil {
		log.Errorf("UserRefreshToken 查询用户失败, userId:%s, error:%v", claims.UserId, err)
		api.Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// 校验token是否过期
	if claims.ExpiresAt < time.Now().Unix() { // token已过期，生成新token
		log.Warnf("UserRefreshToken token已过期，生成新token, old_token:%s", reqbody.OldToken)
		newToken, err := middleware.GenerateToken(user.Id, user.Roles)
		if err != nil {
			log.Errorf("UserRefreshToken 生成新token失败, err:%v", err)
			api.Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
			return
		}
		dataMap["access_token"] = newToken
	} else if claims.ExpiresAt-time.Now().Unix() < int64(middleware.TokenRefreshWindow.Seconds()) { // token临界过期，在动态刷新token窗口内，生成新token
		log.Warnf("UserRefreshToken token达刷新临界时间，生成新token, old_token:%s", reqbody.OldToken)
		newToken, err := middleware.GenerateToken(user.Id, user.Roles)
		if err != nil {
			log.Errorf("UserRefreshToken 生成新token失败, err:%v", err)
			api.Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
			return
		}
		dataMap["access_token"] = newToken
	} else { // token未过期，直接返回旧token
		log.Infof("UserRefreshToken token未过期，返回旧token")
		dataMap["access_token"] = reqbody.OldToken
	}
	dataMap["token_type"] = middleware.TokenType
	api.Success(c, dataMap)
}

// @Title		校验用户注册邮箱
// @Description	往邮箱发送验证码
// @Produce      json
// @Router       /v1/eshop_api/user/verify_email [post]
func VerifyEmail(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("VerifyEmail 请求参数, reqbody:%s", string(req))

	// TODO 添加IP风控
	clientIp := c.ClientIP()

	// JSON解析
	var reqbody model.UserVerifyEmailReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("VerifyEmail json解析失败, error:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}
	// 校验邮箱
	if !isValidEmail(reqbody.Email) {
		log.Errorf("VerifyEmail 邮箱格式错误, reqbody.email:%s", reqbody.Email)
		api.Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail)
		return
	}

	// 校验邮箱是否已注册
	user, err := dao.GetUserByEmail(reqbody.Email)
	if err == nil || user.Id != "" {
		log.Errorf("VerifyEmail 注册邮箱已存在, userId:%s, userName:%s, userEmail:%s", user.Id, user.Name, user.Email)
		api.Fail(c, uerrors.Parse(uerrors.ErrorRegisterMailExisted.Error()).Code, uerrors.Parse(uerrors.ErrorRegisterMailExisted.Error()).Detail)
		return
	}

	// 缓存验证码
	code := mail.GenerateRandomEmailCode()
	err = cache.SaveJxsVerifyMailCode(clientIp, reqbody.Email, code)
	if err != nil {
		log.Errorf("VerifyEmail 缓存验证码失败, clientIp:%s, toEmail:%s, error:%v", clientIp, reqbody.Email, err)
		api.Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// 发送验证码
	err = SendEshopVerifyCodeToEmail(reqbody.Email, code)
	if err != nil {
		log.Errorf("VerifyEmail 发送验证码失败, toEmail:%s, error:%v", reqbody.Email, err)
		api.Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// TODO 飞书通知
	log.Infof("VerifyEmail 发送验证码邮件成功, toEmail:%s, code:%s", reqbody.Email, code)
	api.Success(c, dataMap)
}

// 邮箱有效性判断
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	return true
}

// 密码有效性判断
func isValidPassword(password string) bool {
	if password == "" {
		return false
	}
	if len(password) < 6 {
		return false
	}
	return true
}

// 用户有效性判断
func isValidUser(c *gin.Context) (user *model.User, err error) {
	userId := c.GetString("userId")
	if userId == "" {
		return nil, errors.New("gin.Context用户ID为空")
	}
	// user, err = dao.GetUserById(userId)
	user, err = dao.GetValidUserById(userId)
	if err != nil {
		return nil, err
	}
	return
}

// 发送Eshop的邮箱验证码
// 有效时间为cache.KeyJxsVerifyMailCodeMinsLimit分钟，全局搜索
func SendEshopVerifyCodeToEmail(toemail string, code string) (err error) {
	title := "【江心上客栈】请确认您的新账户..."
	text := fmt.Sprintf("欢迎来到江心上客栈。您的验证码为：%s，有效时间%v分钟。祝您入住愉快。", code, cache.KeyJxsVerifyMailCodeMinsLimit)
	err = mail.SendEmail(toemail, title, text)
	if err != nil {
		log.Errorf("SendEshopVerifyCodeToEmail 发送验证码邮件失败, to:%s, title:%s, text:%s, error:%v", toemail, title, text, err)
		return err
	}
	return nil
}
