package dbmodel

import (
	"time"
)

const TableNameSUser = "s_users"

// SUser mapped from table <s_users>
type SUser struct {
	ID         string    `gorm:"column:id;primaryKey" json:"id"`
	Mobile     string    `gorm:"column:mobile" json:"mobile"`
	Email      string    `gorm:"column:email" json:"email"`
	Password   string    `gorm:"column:password" json:"password"`
	Name       string    `gorm:"column:name" json:"name"`
	Icon       string    `gorm:"column:icon" json:"icon"`                      // 头像的url
	Aliasname  string    `gorm:"column:aliasname" json:"aliasname"`            // 别名
	State      int32     `gorm:"column:state" json:"state"`                    // 状态：1，可用，0禁用,2自动注册的,3vip,4svip,999admin
	Createtime time.Time `gorm:"column:createtime;not null" json:"createtime"` // 创建时间
	Modtime    time.Time `gorm:"column:modtime;not null" json:"modtime"`       // 更新时间
	Sex        int32     `gorm:"column:sex" json:"sex"`
}

// TableName SUser's table name
func (*SUser) TableName() string {
	return TableNameSUser
}

const (
	//状态：1，可用，0禁用,2自动注册的,3vip,4svip,999admin
	UserStateAdmin        = 0 //super admin
	UserStateNormal       = 1
	UserStateAutoRegister = 2 //app auto reg user
	UserStateVIP          = 3
	UserStateSVIP         = 4
	UserStateLocked       = 999
)
