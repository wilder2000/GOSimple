package dbmodel

import (
	"gorm.io/gorm"
	"time"
	"github.com/wilder2000/GOSimple/comm"
)

const TableNameSDebug = "s_debug"

// SDebug mapped from table <s_debug>
type SDebug struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Userid     string    `gorm:"column:userid" json:"userid"`
	Envinfo    string    `gorm:"column:envinfo" json:"envinfo"`
	Detail     string    `gorm:"column:detail" json:"detail"`
	Createtime time.Time `gorm:"column:createtime;not null" json:"createtime"`
}

// TableName SLog's table name
func (*SDebug) TableName() string {
	return TableNameSDebug
}

func (ad *SDebug) BeforeCreate(tx *gorm.DB) error {
	ad.Createtime = comm.LocalTime()
	return nil
}
