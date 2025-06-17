package handler

import (
	"encoding/json"
	"eshop_server/src/router/dao"
	"eshop_server/src/router/model"
	uerrors "eshop_server/src/utils/errors"
	"eshop_server/src/utils/log"

	"github.com/gin-gonic/gin"
)

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
		log.Errorf("GetUserInfo 非法用户请求, error:%v", err)
		FailWithAuthorization(c)
		return
	}
	dataMap["name"] = user.Name
	dataMap["email"] = user.Email
	dataMap["avatar_url"] = user.AvatarUrl
	Success(c, dataMap)
}

// @Title        更新用户信息
// @Description  通过token认证身份并更新用户个人信息
// @Produce      json
// @Router       /v1/eshop_api/user/update_info [post]
func UpdateUserInfo(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})

	// JWT用户查询&鉴权
	user, err := isValidUser(c)
	if err != nil {
		log.Errorf("UpdateUserInfo 非法用户请求, error:%v", err)
		FailWithAuthorization(c)
		return
	}

	// JSON解析
	var reqbody model.UserUpdateInfoReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("UpdateUserInfo json解析失败, error:%v", err)
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}
	if reqbody.Name != "" {
		user.Name = reqbody.Name
	}
	// TODO 邮箱校验并更新
	// if reqbody.Email != "" {
	// 	user.Email = reqbody.Email
	// }
	if reqbody.AvatarUrl != "" {
		user.AvatarUrl = reqbody.AvatarUrl
	}

	// 更新用户信息
	if _, err = dao.UpdateUserByField(user, []string{"name", "avatar_url"}); err != nil {
		log.Errorf("UpdateUserInfo 更新用户信息失败, error:%v", err)
		Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	Success(c, dataMap)
}

// @Title        重置用户密码
// @Description  通过token认证身份并重置密码
// @Produce      json
// @Router       /v1/eshop_api/user/reset_password [post]
func ResetPassword(c *gin.Context) {
	var err error
	req := GetGinBody(c)
	dataMap := make(map[string]interface{})

	// JWT用户查询&鉴权
	user, err := isValidUser(c)
	if err != nil {
		log.Errorf("ResetPassword 非法用户请求, error:%v", err)
		FailWithAuthorization(c)
		return
	}

	// JSON解析
	var reqbody model.UserResetPasswordReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("ResetPassword json解析失败, error:%v", err)
		Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 密码校验
	if reqbody.OldPassword != user.Password {
		log.Errorf("ResetPassword 旧密码输入有误, user.Password:%v, req.Password:%v", user.Password, reqbody.OldPassword)
		Fail(c, uerrors.Parse(uerrors.ErrorPasswordNotSame.Error()).Code, uerrors.Parse(uerrors.ErrorPasswordNotSame.Error()).Detail)
		return
	}

	// 更新旧密码
	user.Password = reqbody.NewPassword
	if _, err = dao.UpdateUserByField(user, []string{"password"}); err != nil {
		log.Errorf("ResetPassword 更新用户密码失败, error:%v", err)
		Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	Success(c, dataMap)
}
