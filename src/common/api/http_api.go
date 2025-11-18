package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"eshop_server/src/utils/log"
	"net/http"
)

const (
	ServerApiVersion = "eshop_api_v1"

	KeyEventTracking          = "eventTracking"
	KeyEventName              = "eventName"
	KeyEventAttrSdkVersion    = "sdk_version"
	KeyEventAttrRequestId     = "request_id"
	KeyEventAttrApiStampBegin = "api_stamp_begin"
	KeyEventAttrApiStampEnd   = "api_stamp_end"
	KeyEventAttrApiStampCost  = "api_stamp_cost"
	KeyEventAttrApiErrorCode  = "api_error_code"
	KeyEventAttrApiErrorMsg   = "api_error_msg"
	KeyEventAttrApiFunc       = "api_func"
	KeyEventAttrApiHost       = "api_host"
	KeyEventAttrStampBegin    = "stamp_begin"
	KeyEventAttrStampEnd      = "stamp_end"
	KeyEventAttrStampCost     = "stamp_cost"
	KeyEventAttrErrorCode     = "error_code"
	KeyEventAttrErrorMsg      = "error_msg"
	KeyEventAttrUid           = "$uid"
	KeyEventAttrAppid         = "$app_id"
)

// 响应结构体
type Response struct {
	ErrorCode int32       `json:"code"` // 自定义错误码
	Data      interface{} `json:"data"` // 数据
	Message   string      `json:"msg"`  // 信息
}

// Success 响应成功 ErrorCode 为 0 表示成功
func Success(c *gin.Context, data interface{}) {
	c.Header("Server-Api-Version", ServerApiVersion)
	c.JSON(http.StatusOK, Response{
		0,
		data,
		"ok",
	})
}

// Fail 响应失败 ErrorCode 不为 0 表示失败
func Fail(c *gin.Context, errorCode int32, msg string) {
	c.Header("Server-Api-Version", ServerApiVersion)
	c.JSON(http.StatusOK, Response{
		errorCode,
		struct{}{},
		msg,
	})
}

// Fail 响应失败 ErrorCode 不为 0 表示失败
func FailWithDataMap(c *gin.Context, errorCode int32, msg string, dataMap interface{}) {
	c.JSON(http.StatusOK, Response{
		errorCode,
		dataMap,
		msg,
	})
}

func FailWithAuthorization(c *gin.Context) {
	c.Header("Server-Api-Version", ServerApiVersion)
	c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
		-1,
		struct{}{},
		"权限验证失败",
	})
}

func FailWithFileNotFound(c *gin.Context) {
	c.Header("Server-Api-Version", ServerApiVersion)
	c.AbortWithStatusJSON(http.StatusNotFound, Response{
		-1,
		struct{}{},
		"文件不存在",
	})	
}

// @Title  获取请求Body参数
// @Description desc
// @Author  wzj  (2022/12/7 17:01)
func GetGinBody(c *gin.Context) (req []byte) {
	var body []byte
	if cb, ok := c.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}
	if body != nil {
		return body
	}

	req, err := c.GetRawData()
	if err != nil {
		log.Error("GetRequestBody fail", zap.Error(err))
		return
	}
	return
}
