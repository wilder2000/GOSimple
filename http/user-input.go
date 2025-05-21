package http

import (
	"time"

	"github.com/wilder2000/GOSimple/dbmodel"
	"github.com/wilder2000/GOSimple/config"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Icon      string
	State     int
	CreatedAt time.Time
	UpdatedAt time.Time
}
type RegisterUserInput struct {
	Email    string `json:"email" alias:"邮件地址" binding:"required,email,max=20"`
	Password string `json:"password" alias:"密码" binding:"required,min=4,max=12"`
}
type RegistUserInput struct {
	AccessKey string `json:"accesskey" alias:"AccessKey" binding:"required"`
	SecretKey string `json:"secretkey" alias:"SecretKey" binding:"required"`
	Uid       string `json:"uuid" alias:"Identifier" binding:"required"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}
type LoginExistInput struct {
	AccessKey string `json:"accesskey"`
	SecretKey string `json:"secretkey"`
	Uid       string `json:"uuid" alias:"Identifier" binding:"required"`
}
type RequestUserInput struct {
	Uid string `json:"uuid" alias:"Identifier" binding:"required"`
}
type DeleteAccountInput struct {
	Uid string `json:"uuid" alias:"Identifier" binding:"required"`
}
type UpdateAliasInput struct {
	Uid       string `json:"uuid" alias:"Identifier" binding:"required"`
	AliasName string `json:"aliasname" alias:"AliasName" binding:"required,max=50"`
}
type ChangePWD struct {
	Email      string `alias:"邮件地址" binding:"required,email,max=20"`
	Password   string `alias:"密码" binding:"required,eqfield=RePassword,min=4,max=12"`
	RePassword string `alias:"确认密码" binding:"required,min=4,max=12"`
}
type CheckPWD struct {
	Email    string `alias:"邮件地址" binding:"required,email,max=20"`
	Password string `alias:"密码" binding:"required,min=4,max=12"`
}

// LoginInput Email不再是email,必须是userid
type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CheckEmailInput struct {
	Email string `json:"email" binding:"required"`
}
type ErrorsInput struct {
	Uid     string `json:"uuid" binding:"required"`
	Envinfo string `json:"envinfo" binding:"required"`
	Detail  string `json:"detail" binding:"required"`
}
type UserFormatter struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Department  []dbmodel.SDepartment `json:"department"`
	Email       string           `json:"email"`
	Icon        string           `json:"icon"`
	Sex         int32            `json:"sex"`
	Aliasname   string           `json:"aliasname"`
	Mobile      string           `json:"mobile"`
	Password    string           `json:"password"`
	State       int32            `json:"state"`
	ReportError bool             `json:"reporterror"`
	//Token      string           `json:"token"`
}

type UserAvatar struct {
	Name  string `form:"title"`
	Email string `form:"email"`
}

func FormatUser(user dbmodel.SUser) UserFormatter {

	formatter := UserFormatter{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Icon:        user.Icon,
		Sex:         user.Sex,
		Aliasname:   user.Aliasname,
		State:       user.State,
		Mobile:      user.Mobile,
		Password:    user.Password,
		ReportError: config.AConfig.ReportError,
	}
	return formatter
}

const (
	// StateAdmin db.SUser.State 管理员账号
	StateAdmin = 0
	// StateNormal db.SUser.State 一般账号
	StateNormal = 2
)
