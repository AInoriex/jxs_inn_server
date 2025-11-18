package handler

import (
	"encoding/json"
	"eshop_server/src/common/api"
	"eshop_server/src/router/dao"
	"eshop_server/src/router/model"
	uerrors "eshop_server/src/utils/errors"
	"eshop_server/src/utils/log"
	"time"

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
		api.FailWithAuthorization(c)
		return
	}
	dataMap["name"] = user.Name
	dataMap["email"] = user.Email
	dataMap["avatar_url"] = user.AvatarUrl
	api.Success(c, dataMap)
}

// @Title        更新用户信息
// @Description  通过token认证身份并更新用户个人信息
// @Produce      json
// @Router       /v1/eshop_api/user/update_info [post]
func UpdateUserInfo(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})

	// JWT用户查询&鉴权
	user, err := isValidUser(c)
	if err != nil {
		log.Errorf("UpdateUserInfo 非法用户请求, error:%v", err)
		api.FailWithAuthorization(c)
		return
	}

	// JSON解析
	var reqbody model.UserUpdateInfoReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("UpdateUserInfo json解析失败, error:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
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
		api.Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	api.Success(c, dataMap)
}

// @Title        重置用户密码
// @Description  通过token认证身份并重置密码
// @Produce      json
// @Router       /v1/eshop_api/user/reset_password [post]
func ResetPassword(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})

	// JWT用户查询&鉴权
	user, err := isValidUser(c)
	if err != nil {
		log.Errorf("ResetPassword 非法用户请求, error:%v", err)
		api.FailWithAuthorization(c)
		return
	}

	// JSON解析
	var reqbody model.UserResetPasswordReq
	err = json.Unmarshal(req, &reqbody)
	if err != nil {
		log.Errorf("ResetPassword json解析失败, error:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Code, uerrors.Parse(uerrors.ErrJsonUnmarshal.Error()).Detail)
		return
	}

	// 密码校验
	if reqbody.OldPassword != user.Password {
		log.Errorf("ResetPassword 旧密码输入有误, user.Password:%v, req.Password:%v", user.Password, reqbody.OldPassword)
		api.Fail(c, uerrors.Parse(uerrors.ErrorPasswordNotSame.Error()).Code, uerrors.Parse(uerrors.ErrorPasswordNotSame.Error()).Detail)
		return
	}

	// 更新旧密码
	user.Password = reqbody.NewPassword
	if _, err = dao.UpdateUserByField(user, []string{"password"}); err != nil {
		log.Errorf("ResetPassword 更新用户密码失败, error:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	api.Success(c, dataMap)
}

// @Title 获取用户列表
// @Desc 管理员获取用户列表
// @Produce json
// @Router /v1/eshop_api/admin/user/list [get]
func AdminGetUserList(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("AdminGetUserList 请求参数, req:%s", string(req))

	// 获取商品信息
	// TODO 分页查询
	resList, err := dao.GetAllUsers(1, 50, "created_at", "desc")
	if err != nil {
		log.Errorf("AdminGetUserList GetAllUsers fail, err:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail)
		return
	}

	// 返回数据
	dataMap["result"] = resList
	dataMap["len"] = len(resList)
	api.Success(c, dataMap)
}

// @Title 拉黑用户
// @Description 管理后台拉黑用户
// @Produce json
// @Router /v1/eshop_api/admin/user/ban/:id [put]
func AdminBanUser(c *gin.Context) {
	var err error
	req := api.GetGinBody(c)
	dataMap := make(map[string]interface{})
	log.Infof("AdminBanUser 请求参数, req:%s", string(req))

	// 获取用户ID
	userId := c.Param("id")
	if userId == "" {
		log.Errorf("AdminBanUser 用户ID不能为空")
		api.Fail(c, uerrors.Parse(uerrors.ErrParam.Error()).Code, uerrors.Parse(uerrors.ErrParam.Error()).Detail+":用户ID无效")
		return
	}

	// 校验用户是否存在
	user, err := dao.GetUserById(userId)
	if err != nil {
		log.Errorf("AdminBanUser GetUserById fail, err:%v", err)
		api.Fail(c, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Code, uerrors.Parse(uerrors.ErrDbQueryFail.Error()).Detail+":用户不存在")
		return
	}

	// 校验用户是否已被拉黑
	if user.Status == model.UserStatusBanned {
		log.Warnf("AdminBanUser 用户已被拉黑, user_id:%s", userId)
		api.Success(c, dataMap)
		return
	}

	// HARDNEED 管理员权限用户不可拉黑
	for _, role := range user.Roles {
		if role == model.UserRoleAdmin {
			log.Errorf("AdminBanUser 管理员用户不可拉黑, user_id:%s", userId)
			api.Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail+":管理员用户不可拉黑")
			return
		}
	}

	// 拉黑用户
	user.Status = model.UserStatusBanned
	user.BannedAt = time.Now()
	if _, err = dao.UpdateUserByField(user, []string{"status", "banned_at"}); err != nil {
		log.Errorf("AdminBanUser 拉黑用户失败, user_id:%s, error:%v", userId, err)
		api.Fail(c, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Code, uerrors.Parse(uerrors.ErrDboperationFail.Error()).Detail)
		return
	}

	api.Success(c, dataMap)
}
