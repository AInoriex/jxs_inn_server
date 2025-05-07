package handler

import (
	"encoding/json"
	"eshop_server/dao"
	"eshop_server/middleware"
	"eshop_server/model"
	uerrors "eshop_server/utils/errors"
	"eshop_server/utils/log"
	"eshop_server/utils/uuid"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"strings"
)

// @Summary      用户登陆
// @Description  邮箱+密码登陆
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/home/login [post]
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
		Fail(c, uerrors.Parse(uerrors.ErrorGetAppUserFail.Error()).Code, uerrors.Parse(uerrors.ErrorGetAppUserFail.Error()).Detail)
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
		Fail(c, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrorUserLoginFail.Error()).Detail)
		return
	}

	// 返回token
	dataMap["token_type"] = middleware.TokenType
	dataMap["access_token"] = tokenString
	Success(c, dataMap)
}

// @Summary      用户注册
// @Description  邮箱+密码注册
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/home/register [post]
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
	Success(c, dataMap)
}

// @Summary      用户刷新token
// @Description  传入旧token用于获取新token
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/home/refresh_token [post]
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

	// 校验旧token
	claims, err := middleware.ValidateToken(reqbody.OldToken)
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			// token已过期，生成新token
		} else {
			// token未过期，直接返回旧token
			log.Warn("UserRefreshToken token未过期，返回旧token", zap.String("token", reqbody.OldToken))
			dataMap["token_type"] = middleware.TokenType
			dataMap["access_token"] = reqbody.OldToken
			Success(c, dataMap)
			return
		}
	}

	// Generate new token using the claims from the old token
	newToken, err := middleware.GenerateToken(claims.UserId, claims.Roles)
	if err != nil {
		log.Error("UserRefreshToken 生成新token失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// Return the new token
	dataMap["token_type"] = middleware.TokenType
	dataMap["access_token"] = newToken
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
