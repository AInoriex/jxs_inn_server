package handler

import (
	"eshop_server/src/router/middleware"
	"eshop_server/src/utils/config"
	"eshop_server/src/utils/log"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	router := gin.Default()
	router.Use(middleware.Logger(), gin.Recovery(), middleware.Cors())

	// 设置路由组
	api := router.Group("/v1/eshop_api")
	{
		// 商品路由
		product := api.Group("/product")
		{
			product.GET("/list", GetProductList)
			// product.POST("/create", CreateProduct) // TODO Warning
			// product.PUT("/remove/:id", RemoveProduct) // TODO Warning
		}

		// 登陆路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", UserRegister)
			auth.POST("/login", UserLogin)
			auth.POST("/refresh_token", UserRefreshToken)
		}

		// 用户权限路由
		user := api.Group("/user")
		user.Use(middleware.ParseAuthorization(), middleware.RequireRole("user"))
		{
			// 用户信息
			user.GET("/info", GetUserInfo)
			// user.PUT("/update", UpdateUserInfo)

			// 购物车
			user.GET("/cart/list", GetCartList)
			user.POST("/cart/create", CreateCart)
			user.POST("/cart/remove", RemoveCart)
			// user.PUT("/cart/update", UpdateCart)

			// 订单&支付
			user.GET("/order/status", GetUserOrderStatus)
			// user.GET("/order/list", GetUserOrderList)
			user.POST("/order/create", CreateUserOrder)
			// user.PUT("/order/update", UpdateUserOrder)
			// user.DELETE("/order/delete", DeleteUserOrder)

			// 库存
			// user.GET("/stock/list", GetStockList)
			// user.POST("/stock/create", CreateStock)
			// user.PUT("/stock/update", UpdateStock)
		}

		// 管理员路由
		admin := api.Group("/admin")
		admin.Use(middleware.ParseAuthorization(), middleware.RequireRole("admin"))
		{
			// 商品操作
			product := api.Group("/product")
			{
				// product.GET("/list", GetProductList)
				product.POST("/create", CreateProduct)
				product.PUT("/remove/:id", RemoveProduct)
				// product.DELETE("/delete/:id", DeleteProduct)
			}
		}
	}

	log.Infof("初始化路由成功, URL：%s", config.CommonConfig.HttpServer.Addr)
	router.Run(config.CommonConfig.HttpServer.Addr)
}
