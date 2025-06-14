package handler

import (
	"encoding/json"
	"errors"
	"eshop_server/src/common/cache"
	"eshop_server/src/router/dao"
	"eshop_server/src/router/middleware"
	"eshop_server/src/router/model"
	uerrors "eshop_server/src/utils/errors"
	"eshop_server/src/utils/log"
	"eshop_server/src/utils/uuid"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// @Title        用户登录
// @Description  邮箱+密码登录
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/auth/login [post]
func UserLogin(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	// attrMap := make(map[string]interface{})
	log.Info("UserLogin 请求参数", zap.String("body", string(req)))

	// TODO 添加IP风控

	// JSON解析
	var reqbody model.UserLoginReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Error("UserLogin json解析失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 参数有效性判断
	if !isValidEmail(reqbody.Email) {
		log.Error("UserLogin 邮箱格式无效", zap.String("email", reqbody.Email))
		Fail(c, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Detail)
		return
	}
	if !isValidPassword(reqbody.HashedPassword) {
		log.Error("UserLogin 密码格式无效", zap.String("password", reqbody.HashedPassword))
		Fail(c, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Detail)
		return
	}

	// 查询用户是否存在
	user, err := dao.GetUserByEmail(reqbody.Email)
	if err != nil {
		log.Error("UserLogin 查询用户失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrorUserNotFound.Error()).Code, uerrors.Parse(uerrors.ErrorUserNotFound.Error()).Detail)
		return
	}

	// 验证哈希密码是否一致
	if reqbody.HashedPassword != user.Password {
		log.Error("UserLogin 密码不一致", zap.String("请求哈希密码", reqbody.HashedPassword), zap.String("目标哈希密码", user.Password))
		Fail(c, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Detail)
		return
	}

	// 生成 JWT Token
	// TODO 新增权限区分生成不同权限token(user/admin)
	tokenString, err := middleware.GenerateToken(user.Id, []string{"user"})
	if err != nil {
		log.Error("UserLogin 生成jwt token失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// 缓存保存token
	if err = cache.SaveJxsUserToken(user.Id, tokenString); err != nil {
		log.Errorf("UserLogin 缓存保存JxsUserToken失败, user_id:%v, token:%s, err:%v", user.Id, tokenString, err)
		Fail(c, uerrors.Parse(uerrors.ErrRedis.Error()).Code, uerrors.Parse(uerrors.ErrRedis.Error()).Detail)
		return
	}

	// 返回token
	dataMap["token_type"] = middleware.TokenType
	dataMap["access_token"] = tokenString
	Success(c, dataMap)
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
		Success(c, nil)
		return
	}

	// 检查Token格式：{TokenType} {TokenString}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == middleware.TokenType) {
		log.Warn("UserLogout 认证头格式错误")
		Success(c, nil)
		return
	}

	// 校验Token
	requestToken := parts[1]
	claims, err := middleware.ValidateToken(requestToken)
	if err != nil {
		log.Warn("UserLogout 解析token失败", zap.Error(err))
		Success(c, nil)
		return
	}

	// 从Redis中删除token
	if claims.UserId != "" {
		isDel := cache.DelJxsUserToken(claims.UserId)
		if !isDel {
			log.Error("UserLogout 删除缓存JxsUserToken失败", zap.String("user_id", claims.UserId))
		}
	}
	Success(c, nil)
}

// @Title        用户注册
// @Description  邮箱+密码注册
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/auth/register [post]
func UserRegister(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	// attrMap := make(map[string]interface{})
	log.Info("UserRegister 请求参数", zap.String("body", string(req)))

	// TODO 添加IP风控

	// JSON解析
	var reqbody model.UserRegisterReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Error("UserRegister json解析失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 查询用户是否存在
	user, err := dao.GetUserByEmail(reqbody.Email)
	if err == nil || user.Id != "" {
		log.Error("UserRegister 邮箱已存在，注册失败", zap.String("user_id", user.Id), zap.String("name", user.Name), zap.String("email", user.Email))
		Fail(c, uerrors.Parse(uerrors.ErrorRegisterMailExisted.Error()).Code, uerrors.Parse(uerrors.ErrorRegisterMailExisted.Error()).Detail)
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
		log.Error("UserRegister 创建用户失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// 返回
	log.Infof("UserRegister 用户注册成功，user_id:%v, name:%v, email:%v", new_user.Id, new_user.Name, new_user.Email)
	Success(c, dataMap)
}

// @Title        用户刷新token
// @Description  传入旧token用于获取新token
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/auth/refresh_token [post]
func UserRefreshToken(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Info("UserRefreshToken 请求参数", zap.String("body", string(req)))

	// TODO 添加IP风控

	// JSON解析
	var reqbody struct {
		OldToken string `json:"token"`
	}
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Error("UserRegister json解析失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}
	// 校验token是否为空
	if reqbody.OldToken == "" {
		log.Error("UserRefreshToken token为空", zap.String("token", reqbody.OldToken))
		Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail)
		return
	}

	// 校验旧token
	claims, err := middleware.GetTokenClaims(reqbody.OldToken)
	if err != nil {
		log.Error("UserRefreshToken 解析token失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}
	if claims.ExpiresAt < time.Now().Unix() { // token已过期，生成新token
		log.Warnf("UserRefreshToken token已过期，生成新token, old_token:%s", reqbody.OldToken)
		newToken, err := middleware.GenerateToken(claims.UserId, claims.Roles)
		if err != nil {
			log.Errorf("UserRefreshToken 生成新token失败, err:%v", err)
			Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
			return
		}
		dataMap["access_token"] = newToken
	} else if claims.ExpiresAt-time.Now().Unix() < int64(middleware.TokenRefreshWindow.Seconds()) { // token临界过期，在动态刷新token窗口内，生成新token
		log.Warnf("UserRefreshToken token达刷新临界时间，生成新token, old_token:%s", reqbody.OldToken)
		newToken, err := middleware.GenerateToken(claims.UserId, claims.Roles)
		if err != nil {
			log.Errorf("UserRefreshToken 生成新token失败, err:%v", err)
			Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
			return
		}
		dataMap["access_token"] = newToken
	} else { // token未过期，直接返回旧token
		log.Infof("UserRefreshToken token未过期，返回旧token")
		dataMap["access_token"] = reqbody.OldToken
	}
	dataMap["token_type"] = middleware.TokenType
	Success(c, dataMap)
}

// @Title        获取用户信息
// @Description  通过token认证身份并获取本人用户信息
// @Produce      json
// @Router       /v1/eshop_api/user/info [get]
func GetUserInfo(c *gin.Context) {
	var err error
	dataMap := make(map[string]interface{})

	// JWT用户查询&鉴权
	user, err := isValidUser(c)
	if err != nil {
		log.Error("GetUserInfo 非法用户请求", zap.Error(err))
		FailWithAuthorization(c)
		return
	}

	dataMap["name"] = user.Name
	dataMap["email"] = user.Email
	dataMap["avatar_url"] = user.AvatarUrl
	Success(c, dataMap)
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
	user, err = dao.GetUserById(userId)
	if err != nil {
		return
	}
	return
}
