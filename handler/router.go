package handler

import (
    "github.com/gin-gonic/gin"
    "eshop_server/middleware"
    "eshop_server/utils/config"
    "eshop_server/utils/log"
)

func InitRouter() {
	router := gin.Default()

	// 设置路由组
	api := router.Group("/v1/eshop_api")
	{
		// 商品路由
		product := api.Group("/product")
		product.Use(middleware.Logger(), gin.Recovery(), middleware.Cors())
		{
			product.GET("/list", GetProductList)
			product.POST("/create", CreateProduct) // TODO Warning
			product.PUT("/remove/:id", RemoveProduct) // TODO Warning
		}

		// 登陆路由
        user := api.Group("/user")
        user.Use(middleware.Logger(), gin.Recovery(), middleware.Cors())
        {
            user.POST("/register", UserRegister)
            user.POST("/login", UserLogin)
			user.POST("/refresh_token", UserRefreshToken)
        }
		
		// 用户权限路由
		auth := api.Group("/auth")
		auth.Use(middleware.Logger(), gin.Recovery(), middleware.AuthUser(), middleware.RequireRole("user"))
		{
			// 用户信息
			// auth.GET("/user/info", GetUserInfo)
			// auth.PUT("/user/info", UpdateUserInfo)

			// 购物车
			// auth.GET("/cart/list", GetCartList)
			// auth.POST("/cart/create", CreateCart)
			// auth.PUT("/cart/update", UpdateCart)
			// auth.DELETE("/cart/delete", DeleteCart)

			// 订单&支付
			// auth.GET("/order/list", GetOrderList)
			// auth.POST("/order/create", CreateOrder)
			// auth.PUT("/order/update", UpdateOrder)
			// auth.DELETE("/order/delete", DeleteOrder)
			// auth.POST("/order/pay", PayOrder)

			// 库存
			// auth.GET("/stock/list", GetStockList)
			// auth.POST("/stock/create", CreateStock)
			// auth.PUT("/stock/update", UpdateStock)
		}

		// 管理员路由
		admin := api.Group("/admin")
        admin.Use(middleware.Logger(), gin.Recovery(), middleware.Cors(), middleware.AuthAdmin(), middleware.RequireRole("admin"))
        {
			product := api.Group("/product")
            {
				product.GET("/list", GetProductList)
				product.POST("/create", CreateProduct)
				product.PUT("/remove/:id", RemoveProduct)
				// product.DELETE("/delete/:id", DeleteProduct)
			}
        }
	}

	log.Infof("初始化路由成功, URL：%s", config.CommonConfig.HttpServer.Addr)
	router.Run(config.CommonConfig.HttpServer.Addr)
}