package handler

import (
	"eshop_server/src/router/middleware"
	"eshop_server/src/router/model"
	"eshop_server/src/utils/config"
	"eshop_server/src/utils/log"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	router := gin.Default()
	router.Use(middleware.Logger(), gin.Recovery(), middleware.Cors())

	// 健康路由
	health_api := router.Group("/health")
	{
		health_api.GET("/ping", HealthPing)
	}

	// 设置路由组
	api := router.Group("/v1/eshop_api")
	{
		// 公共路由（需要接口风控限制）
		// 商品路由
		product := api.Group("/product")
		// product.Use(middleware.RateLimitMiddleware())
		{
			product.GET("/list", GetProductList)
			// product.GET("/search", SearchProducts)
		}

		// 登录路由
		auth := api.Group("/auth")
		// auth.Use(middleware.RateLimitMiddleware())
		{
			// 网站
			// auth.POST("/register", UserRegister)
			auth.POST("/register", UserRegisterWithVerifyCode)
			auth.POST("/login", UserLogin)
			auth.GET("/logout", UserLogout)
			auth.POST("/refresh_token", RefreshToken)
			auth.POST("/verify_email", VerifyEmail)
			
			// 管理后台
			auth.POST("/admin_login", AdminLogin)
		}

		// 用户权限路由
		user := api.Group("/user")
		user.Use(middleware.ParseAuthorization(), middleware.RequireRole(model.UserRoleUser))
		{
			// 用户信息
			user.GET("/info", GetUserInfo)
			user.POST("/update_info", UpdateUserInfo)
			user.POST("/reset_password", ResetPassword)
			user.GET("/purchase_history", GetUserPurchaseHistory)

			// 购物车
			user.GET("/cart/list", GetCartList)
			user.POST("/cart/create", CreateCart)
			user.POST("/cart/remove", RemoveCart)
			// user.PUT("/cart/update", UpdateCart)

			// 订单&支付
			user.GET("/order/status", GetUserOrderStatus)
			user.POST("/order/create", CreateUserOrder)
			// user.POST("/order/cancel", CancelUserOrder)
			// user.GET("/order/list", GetUserOrderList)

			// 藏品
			user.GET("/inventory/list", GetInventoryList)
		}

		// 管理员权限路由
		admin := api.Group("/admin")
		admin.Use(middleware.ParseAuthorization(), middleware.RequireRole(model.UserRoleAdmin))
		{
			// 用户操作
			user := admin.Group("/user")
			{
				user.GET("/list", AdminGetUserList)
				user.PUT("/ban/:id", AdminBanUser)
			}

			// 商品操作
			product := admin.Group("/product")
			{
				product.GET("/list", AdminGetProductList)
				product.POST("/create", AdminCreateProduct)
				product.PUT("/remove/:id", AdminRemoveProduct)
				// product.DELETE("/delete/:id", DeleteProduct)
				// product.GET("/search", SearchProducts)
			}

			// 商品资源操作
			product_player := admin.Group("/player")
			{
				// product_player.GET("/list", AdminGetProductPlayerList)
				product_player.POST("/upload_streaming_file", UploadStreamingFile)
				product_player.POST("/update_streaming_file", UpdateStreamingFile)
			}

			// 订单操作
			order := admin.Group("/order")
			{
				order.GET("/list", AdminGetUserOrderList)
			}
		}
	}

	log.Infof("初始化路由成功, URL：%s", config.CommonConfig.HttpServer.Addr)
	router.Run(config.CommonConfig.HttpServer.Addr)
}
