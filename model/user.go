package model

import (
	"time"
)

/*
	按以下sql生成对应struct

CREATE TABLE `users` (

	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识',
	`name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户姓名',
	`email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户邮箱',
	`password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户密码(强加密算法存储, 如bcrypt、scrypt等)',
	`avatar_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户头像URL',
	`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	PRIMARY KEY (`id`),
	UNIQUE KEY `email` (`email`),
	KEY `idx_email` (`email`)

) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
*/
type User struct {
	Id        string    `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL;comment:'自增唯一ID'"`
	Name      string    `json:"name" gorm:"column:name;default:NULL;comment:'用户姓名'"`
	Email     string    `json:"email" gorm:"column:email;NOT NULL;comment:'用户邮箱'"`
	Password  string    `json:"password" gorm:"column:password;NOT NULL;comment:'用户密码(强加密算法存储, 如bcrypt、scrypt等)'"`
	AvatarUrl string    `json:"avatar_url" gorm:"column:avatar_url;default:NULL;comment:'用户头像URL'"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;default:NULL ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'"`
}

func (User) TableName() string {
	return "users"
}

type UserLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
