package dbmodel

const TableNameSItems = "s_items"

// SItem mapped from table <s_items>

type SItem struct {
	ID    int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name  string `gorm:"column:name" json:"name"`
	IType int32  `gorm:"column:type" json:"type"`
}

// TableName SItem's table name
func (*SItem) TableName() string {
	return TableNameSItems
}
