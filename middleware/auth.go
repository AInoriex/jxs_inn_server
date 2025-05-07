package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
	"time"
)

// 自定义Claims结构
type CustomClaims struct {
	UserId string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.StandardClaims
}

// 安全配置（建议通过环境变量配置）
var (
	jwtSecret     = []byte(os.Getenv("JWT_SECRET")) // 32位以上安全密钥
	tokenDuration = 18 * time.Minute                // Access Token有效期
	refreshWindow = 5 * time.Minute                 // 刷新时间窗口
	TokenType     = "Bearer"                        // token类型
)

// 用户接口权限校验
func ParseAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证凭证"})
			return
		}

		// 检查Token格式：{TokenType} {TokenString}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == TokenType) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "认证头格式错误"})
			return
		}

		// 校验Token
		tokenString := parts[1]
		claims, err := ValidateToken(tokenString)
		if err != nil {
			HandleTokenError(c, err)
			return
		}

		// 自动刷新机制
		if claims.ExpiresAt-time.Now().Unix() < int64(refreshWindow.Seconds()) {
			newToken, err := generateToken(claims.UserId, claims.Roles)
			if err == nil {
				c.Header("Set-Access-Token", newToken)
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
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "token已过期",
				"refreshable": true,
			})
			return
		}
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的认证凭证"})
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
