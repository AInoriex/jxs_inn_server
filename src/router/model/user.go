package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

/*
	按以下sql生成对应struct

-- @Author AInoriex
-- @Desc 目前只支持邮箱一种方式登录
-- @TODO 用户角色：如果未来有管理员、普通用户等不同角色, 可以增加一个role字段, 用于区分用户权限。
-- @TODO 联系方式：除了邮箱, 可以增加手机号字段, 方便用户接收验证码、订单通知等信息。
-- @TODO 账户锁定机制：可以增加login_attempts字段记录登录失败次数, 当连续多次登录失败时, 暂时锁定账户, 防止暴力破解。
-- @TTODO 会员信息：如果计划推出会员制度, 可以增加会员等级、会员积分等字段。
-- @TTODO 登录方式：除了邮箱登录, 可以考虑支持社交媒体账号登录(如微信、QQ、微博等), 增加social_login_id字段存储第三方登录的唯一标识。
CREATE TABLE `users` (

	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识',
	`name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户姓名',
	`email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户邮箱',
	`password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户密码(强加密算法存储, 如bcrypt、scrypt等)',
	`avatar_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户头像URL',
	`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	`roles` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT 'user' COMMENT '用户角色权限（admin:管理员, user:普通用户，逗号分隔）',
	`last_login` datetime DEFAULT NULL COMMENT '最后登录时间',
	`status` tinyint(1) DEFAULT '1' COMMENT '用户状态（1:正常, 0:禁用）',
	`banned_at` datetime DEFAULT NULL COMMENT '账户锁定时间',
	PRIMARY KEY (`id`),
	UNIQUE KEY `email` (`email`),
	KEY `idx_email` (`email`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
*/

// 新增：自定义类型实现 sql.Scanner 和 driver.Valuer 接口
type RoleSlice []string

const (
	UserStatusNormal = 1       // 用户正常
	UserStatusBanned = 0       // 用户禁用
	UserRoleAdmin    = "admin" // 管理员
	UserRoleUser     = "user"  // 普通用户
)

type User struct {
	Id        string    `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL;comment:'自增唯一ID'"`
	Name      string    `json:"name" gorm:"column:name;default:NULL;comment:'用户姓名'"`
	Email     string    `json:"email" gorm:"column:email;NOT NULL;comment:'用户邮箱'"`
	Password  string    `json:"password" gorm:"column:password;NOT NULL;comment:'用户密码(强加密算法存储, 如bcrypt、scrypt等)'"`
	AvatarUrl string    `json:"avatar_url" gorm:"column:avatar_url;default:NULL;comment:'用户头像URL'"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;default:NULL ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'"`
	Roles     RoleSlice  `json:"roles" gorm:"column:roles;type:varchar(255);default:'user';comment:'用户角色（admin:管理员, user:普通用户，逗号分隔）'"`
	LastLogin time.Time `json:"last_login" gorm:"column:last_login;default:NULL;comment:'最后登录时间'"`
	Status    int32     `json:"status" gorm:"column:status;default:1;comment:'用户状态（1:正常, 0:禁用）'"`
	BannedAt  time.Time `json:"banned_at" gorm:"column:banned_at;default:NULL;comment:'账户锁定时间'"`
}

func (User) TableName() string {
	return "users"
}


// Scan 从数据库读取时，将字符串转换为 []string
func (r *RoleSlice) Scan(value interface{}) error {
    if value == nil {
        *r = []string{}
        return nil
    }
    // 数据库返回的 value 可能是 []uint8（字节切片）或 string
    str, ok := value.(string)
    if !ok {
        byteSlice, ok := value.([]byte)
        if !ok {
            return fmt.Errorf("invalid scan type for RoleSlice: %T", value)
        }
        str = string(byteSlice)
    }
    if str == "" {
        *r = []string{}
        return nil
    }
    *r = strings.Split(str, ",")
    return nil
}

// Value 写入数据库时，将 []string 转换为逗号分隔的字符串
func (r RoleSlice) Value() (driver.Value, error) {
    return strings.Join(r, ","), nil
}

// 用户登录请求体
type UserLoginReq struct {
	Email          string `json:"email"`
	HashedPassword string `json:"password"`
}

// 用户注册请求体
type UserRegisterReq struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword string `json:"password"`
	VerifyCode     string `json:"code"`
}

// 用户刷新token请求体
type RefreshTokenReq struct {
	OldToken string `json:"token"`
}

// 校验用户邮箱请求体
type UserVerifyEmailReq struct {
	Email string `json:"email"`
}

// 用户更新个人信息请求体
type UserUpdateInfoReq struct {
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	Email     string `json:"email"`
}

// 用户重置密码请求体
type UserResetPasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
