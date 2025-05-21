package http

import (
	"github.com/wilder2000/GOSimple/dbmodel"
	"gorm.io/gorm"
)

type Repository interface {
	Save(user dbmodel.SUser) (dbmodel.SUser, error)
	FindByEmail(email string) (dbmodel.SUser, error)
	FindByID(ID string) (dbmodel.SUser, error)
	Update(user dbmodel.SUser) (dbmodel.SUser, error)
	UpdatePwd(email string, pwd string) bool
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {

	return &repository{db}
}

func (r *repository) Save(user dbmodel.SUser) (dbmodel.SUser, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// FindByEmail 只用于管理员
func (r *repository) FindByEmail(email string) (dbmodel.SUser, error) {
	var user dbmodel.SUser
	err := r.db.Where("id = ? and state = ?", email, StateAdmin).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) FindByID(ID string) (dbmodel.SUser, error) {
	var user dbmodel.SUser
	err := r.db.Where("ID = ?", ID).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) Update(user dbmodel.SUser) (dbmodel.SUser, error) {
	err := r.db.Save(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
func (r *repository) UpdatePwd(email string, pwd string) bool {
	tb := r.db.Model(&dbmodel.SUser{}).Where("email=?", email).Update("password", pwd)
	return tb.RowsAffected == 1
}
