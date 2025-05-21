package dbmodel

type SGroupWithUser struct {
	SGroup
	Selected bool `gorm:"column:selected"`
}
type SUserWithGroup struct {
	SUser
	Selected bool `gorm:"column:selected"`
}
type SRoleWithGroup struct {
	SRole
	Selected bool `gorm:"column:selected"`
}
type SOperatorWithRole struct {
	SOperator
	Selected bool `gorm:"column:selected"`
}
type SUserWithDep struct {
	SUser
	Selected bool `gorm:"column:selected"`
}
