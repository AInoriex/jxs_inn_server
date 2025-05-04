package handler

import (
	"encoding/json"
	"eshop_server/dao"
	"eshop_server/middleware"
	"eshop_server/model"
	uerrors "eshop_server/utils/errors"
	"eshop_server/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// @Summary      用户登陆
// @Description  邮箱+密码登陆
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/user/login [post]
func UserLogin(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	// attrMap := make(map[string]interface{})
	log.Info("UserLogin params", zap.String("body", string(req)))

	// TODO 添加IP风控

	// JSON解析
	var reqbody model.UserLoginReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Error("UserLogin JSON解析失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 查询用户是否存在
	user, err := dao.GetUserByEmail(reqbody.Email)
	if err != nil {
		log.Error("UserLogin 查询用户失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrUserNotFound.Error()).Code, uerrors.Parse(uerrors.ErrUserNotFound.Error()).Detail)
		return
	}

	// 验证密码
	if reqbody.Password != user.Password {
		log.Error("UserLogin 密码不一致", zap.String("请求密码", reqbody.Password), zap.String("目标密码", user.Password))
		Fail(c, uerrors.Parse(uerrors.ErrUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrUserLoginFail.Error()).Detail)
		return
	}

	// 生成JWT Token
	tokenString, err := middleware.GenerateToken(user.Id, []string{"user"})
	if err != nil {
		log.Error("UserLogin 生成jwt token失败", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrUserLoginFail.Error()).Code, uerrors.Parse(uerrors.ErrUserLoginFail.Error()).Detail)
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
// @Router       /v1/eshop_api/user/register [post]
func UserRegister(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	// attrMap := make(map[string]interface{})
	log.Info("UserRegister params", zap.String("body", string(req)))

	// TODO 添加IP风控

	// JSON解析
	var reqbody model.UserRegisterReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Error("UserRegister Unmarshal fail", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 查询用户是否存在
	user, err := dao.GetUserByEmail(reqbody.Email)
	if err == nil || user.Id != "" {
		log.Error("UserRegister 用户已存在，注册失败", zap.String("id", user.Id), zap.String("name", user.Name))
		Fail(c, uerrors.Parse(uerrors.ErrUserExisted.Error()).Code, uerrors.Parse(uerrors.ErrUserExisted.Error()).Detail)
		return
	}

	// 创建新用户
	new_user := &model.User{
		Name:     reqbody.Name,
		Email:    reqbody.Email,
		Password: reqbody.Password,
	}
	_, err = dao.CreateUser(new_user)
	if err != nil {
		log.Error("UserRegister CreateUser fail", zap.Error(err))
		Fail(c, uerrors.Parse(uerrors.ErrBusy.Error()).Code, uerrors.Parse(uerrors.ErrBusy.Error()).Detail)
		return
	}

	// 返回
	Success(c, dataMap)
}

// @Summary      用户刷新token
// @Description  old_token
// @Param        json
// @Produce      json
// @Router       /v1/eshop_api/user/refresh_token [post]
func UserRefreshToken(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Info("UserRefreshToken params", zap.String("body", string(req)))

	// JSON解析
	var reqbody struct {
		OldToken string `json:"token"`
	}
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Error("UserRegister Unmarshal fail", zap.Error(err))
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
			log.Warn("UserRefreshToken token未过期", zap.Error(err))
			dataMap["token_type"] = middleware.TokenType
			dataMap["access_token"] = reqbody.OldToken
			Success(c, dataMap)
			return
		}
	}

	// Generate new token using the claims from the old token
	newToken, err := middleware.GenerateToken(claims.UserID, claims.Roles)
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
