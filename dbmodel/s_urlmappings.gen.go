package dbmodel

type SUrlMapping struct {
	ID         int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Operatorid int32  `gorm:"column:operatorid" json:"operatorid"`
	Url        string `gorm:"column:url" json:"url"`
}

func (*SUrlMapping) TableName() string {
	return "s_urlmappings"
}
