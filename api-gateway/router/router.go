package router

import (
	"api-gateway/middleware"
	"api-gateway/proto"
	"api-gateway/service"
	"log"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func SetupRouter(userSvc *service.UserService, productSvc *service.ProductService) *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API 路由组
	v1 := r.Group("/api/v1")
	{
		// 用户服务路由 (注册和登录不需要认证)
		userRoutes := v1.Group("/users")
		{
			userRoutes.POST("/register", func(c *gin.Context) {
				var req proto.RegisterRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				resp, err := userSvc.Register(c.Request.Context(), &req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, resp)
			})

			userRoutes.POST("/login", func(c *gin.Context) {
				var req proto.LoginRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				resp, err := userSvc.Login(c.Request.Context(), &req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, resp)
			})
		}

		// 公开的商品服务路由
		productRoutes := v1.Group("/products")
		{
			// 获取商品列表
			productRoutes.GET("", func(c *gin.Context) {
				page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
				pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)

				// 这里不需要 Authorization 头部，但为了传递可能的 Trace ID 等，可以传递 Context
				resp, err := productSvc.ListProducts(c.Request.Context(), &proto.ListProductsRequest{
					Page:     int32(page),
					PageSize: int32(pageSize),
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, resp)
			})

			// 获取商品详情
			productRoutes.GET("/:id", func(c *gin.Context) {
				productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
					return
				}

				// 这里不需要 Authorization 头部，但为了传递可能的 Trace ID 等，可以传递 Context
				resp, err := productSvc.GetProduct(c.Request.Context(), &proto.GetProductRequest{ProductId: productID})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, resp)
			})
		}

		// 需要认证的路由
		authRoutes := v1.Group("")
		authRoutes.Use(middleware.AuthMiddleware())
		{
			// 用户服务需要认证的路由 (获取、更新、删除、列表)
			authUserRoutes := userRoutes.Group("") // 仍然使用 userRoutes 前缀
			{
				// 获取用户信息
				authUserRoutes.GET("/:id", func(c *gin.Context) {
					userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
					if err != nil {
						log.Printf("Invalid user ID in path: %v", err)
						c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
						return
					}
					log.Printf("Calling GetUserInfo for user ID: %d", userID)

					// 从 Gin context 中获取 Authorization 头部
					authHeader := c.GetHeader("Authorization")
					if authHeader == "" {
						// 理论上 AuthMiddleware 已经检查过，这里只是双重确认
						c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing after middleware"})
						return
					}

					// 将 Authorization 头部添加到 gRPC metadata
					md := metadata.Pairs("authorization", authHeader)
					ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

					resp, err := userSvc.GetUserInfo(ctx, &proto.UserInfoRequest{UserId: userID})
					if err != nil {
						log.Printf("GetUserInfo gRPC call failed: %v", err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, resp)
				})

				// 更新用户信息
				authUserRoutes.PUT("/:id", func(c *gin.Context) {
					userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
						return
					}
					var req proto.UpdateUserRequest
					if err := c.ShouldBindJSON(&req); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					req.UserId = userID

					// 从 Gin context 中获取 Authorization 头部
					authHeader := c.GetHeader("Authorization")
					if authHeader == "" {
						c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing after middleware"})
						return
					}

					// 将 Authorization 头部添加到 gRPC metadata
					md := metadata.Pairs("authorization", authHeader)
					ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

					resp, err := userSvc.UpdateUser(ctx, &req)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, resp)
				})

				// 删除用户
				authUserRoutes.DELETE("/:id", func(c *gin.Context) {
					userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
						return
					}

					// 从 Gin context 中获取 Authorization 头部
					authHeader := c.GetHeader("Authorization")
					if authHeader == "" {
						c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing after middleware"})
						return
					}

					// 将 Authorization 头部添加到 gRPC metadata
					md := metadata.Pairs("authorization", authHeader)
					ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

					resp, err := userSvc.DeleteUser(ctx, &proto.DeleteUserRequest{UserId: userID})
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, resp)
				})

				// 获取用户列表
				authUserRoutes.GET("", func(c *gin.Context) {
					page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
					pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)

					// 从 Gin context 中获取 Authorization 头部
					authHeader := c.GetHeader("Authorization")
					if authHeader == "" {
						c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing after middleware"})
						return
					}

					// 将 Authorization 头部添加到 gRPC metadata
					md := metadata.Pairs("authorization", authHeader)
					ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

					resp, err := userSvc.ListUsers(ctx, &proto.ListUsersRequest{
						Page:     int32(page),
						PageSize: int32(pageSize),
					})
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, resp)
				})
			}

			// 商品服务需要认证的路由 (创建、更新、删除)
			authProductRoutes := productRoutes.Group("") // 仍然使用 productRoutes 前缀
			{
				authProductRoutes.POST("", nil)       // TODO: 创建商品
				authProductRoutes.PUT("/:id", nil)    // TODO: 更新商品
				authProductRoutes.DELETE("/:id", nil) // TODO: 删除商品
			}

			// 订单服务路由
			orderRoutes := authRoutes.Group("/orders")
			{
				orderRoutes.POST("", nil)    // TODO: 创建订单
				orderRoutes.GET("/:id", nil) // TODO: 获取订单详情
				orderRoutes.GET("", nil)     // TODO: 获取订单列表
			}
		}
	}

	return r
}
