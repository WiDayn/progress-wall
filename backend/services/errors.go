package services

import "errors"

// 定义业务错误
var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserDisabled       = errors.New("账户已被禁用")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrGenerateToken      = errors.New("生成token失败")
	ErrUpdateLoginTime    = errors.New("更新登录时间失败")
	ErrUserExists         = errors.New("用户名或邮箱已存在")
	ErrInvalidPassword    = errors.New("密码格式不正确")
)
