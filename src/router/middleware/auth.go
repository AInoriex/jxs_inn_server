package middleware

import (
	"eshop_server/src/common/cache"
	"eshop_server/src/utils/config"
	"eshop_server/src/utils/log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// 自定义Claims结构
type CustomClaims struct {
	UserId string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.StandardClaims
}

var (
	jwtSecret          = []byte(config.CommonConfig.JwtSecret)      // 32位以上安全密钥
	tokenDuration      = cache.KeyJxsUserTokenTimeout * time.Second // Token有效期
	TokenRefreshWindow = 5 * time.Minute                            // 刷新时间窗口
	TokenType          = "Bearer"                                   // token类型
)

// 用户接口权限校验
func ParseAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "权限不足"})
			return
		}

		// 检查Token格式：{TokenType} {TokenString}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == TokenType) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "认证头格式错误"})
			return
		}

		// 校验Token
		requestToken := parts[1]
		claims, err := ValidateToken(requestToken)
		if err != nil {
			HandleTokenError(c, err)
			return
		}

		// 校验Redis中的token是否有效（与当前请求token一致）
		flag, cachedToken := cache.GetJxsUserToken(claims.UserId)
		if !flag {
			log.Errorf("ParseAuthorization 读取缓存失败, userId:%s", claims.UserId)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "认证凭证无效"})
			return
		}
		if cachedToken != requestToken {
			log.Errorf("ParseAuthorization 缓存token与请求token不一致, userId:%s, requestToken:%s, cachedToken:%s", claims.UserId, requestToken, cachedToken)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "认证凭证已失效"})
			return
		}

		// 自动刷新机制
		if claims.ExpiresAt-time.Now().Unix() < int64(TokenRefreshWindow.Seconds()) {
			newToken, err := generateToken(claims.UserId, claims.Roles)
			if err == nil {
				// 刷新token时更新Redis缓存
				if err := cache.SaveJxsUserToken(claims.UserId, newToken); err != nil {
					log.Errorf("ParseAuthorization 更新Redis新token失败, userId:%s, token:%s, err:%s", claims.UserId, newToken, err.Error())
				} else {
					c.Header("Set-Access-Token", newToken) // 设置请求头要求前端刷新token
				}
			}
		}

		// 存储用户上下文
		c.Set("userId", claims.UserId)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

// 权限验证中间件
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists || !containsRole(roles.([]string), requiredRole) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			return
		}
		c.Next()
	}
}

// 校验token
func ValidateToken(tokenString string) (*CustomClaims, error) {
	// 使用`jwt.ParseWithClaims`函数解析 token 字符串并提取其声明。
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 检查token的签名方法是否为`HMAC`（一种对称密钥加密）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid sign method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 检查token是否有效，以及其声明是否为`CustomClaims`类型
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.NewValidationError("invalid token claims", jwt.ValidationErrorClaimsInvalid)
}

func GetTokenClaims(tokenString string) (*CustomClaims, error) {
	// 使用`jwt.ParseWithClaims`函数解析 token 字符串并提取其声明。
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 检查token的签名方法是否为`HMAC`（一种对称密钥加密）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid sign method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtSecret, nil
	})
	// 忽略过期错误
	if err != nil && err.(*jwt.ValidationError).Errors != jwt.ValidationErrorExpired {
		return nil, err
	}
	// 检查token是否有效，以及其声明是否为`CustomClaims`类型
	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims, nil
	}
	return nil, jwt.NewValidationError("invalid token claims", jwt.ValidationErrorClaimsInvalid)
}

// Token生成函数（供登录成功后调用）
func GenerateToken(userId string, roles []string) (string, error) {
	return generateToken(userId, roles)
}

// 私有生成函数
func generateToken(userId string, roles []string) (string, error) {
	claims := CustomClaims{
		UserId: userId,
		Roles:  roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
			Issuer:    "eshop_server",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Token错误处理
func HandleTokenError(c *gin.Context, err error) {
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			c.Header("Refresh-Token", "1") // 设置请求头要求前端刷新token
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token已过期",
			})
			return
		}
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "认证凭证无效"})
}

// 辅助函数：检查角色权限
func containsRole(roles []string, target string) bool {
	for _, role := range roles {
		if role == target {
			return true
		}
	}
	return false
}
