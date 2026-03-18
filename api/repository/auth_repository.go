package repository

import (
	"smartwatch-server/api/models"

	"gorm.io/gorm"
)

// AuthRepository 认证仓储
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository 创建认证仓储
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser 创建用户
func (r *AuthRepository) CreateUser(u *models.AdminUser) error {
	return r.db.Create(u).Error
}

// GetByUsername 按用户名查询
func (r *AuthRepository) GetByUsername(username string) (*models.AdminUser, error) {
	var u models.AdminUser
	err := r.db.Where("username = ? AND status = 1", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByPhone 按手机号查询
func (r *AuthRepository) GetByPhone(phone string) (*models.AdminUser, error) {
	var u models.AdminUser
	err := r.db.Where("phone = ? AND status = 1", phone).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByEmail 按邮箱查询
func (r *AuthRepository) GetByEmail(email string) (*models.AdminUser, error) {
	var u models.AdminUser
	err := r.db.Where("email = ? AND status = 1", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ExistsUsername 用户名是否存在
func (r *AuthRepository) ExistsUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.AdminUser{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsPhone 手机号是否存在
func (r *AuthRepository) ExistsPhone(phone string) (bool, error) {
	if phone == "" {
		return false, nil
	}
	var count int64
	err := r.db.Model(&models.AdminUser{}).Where("phone = ?", phone).Count(&count).Error
	return count > 0, err
}

// ExistsEmail 邮箱是否存在
func (r *AuthRepository) ExistsEmail(email string) (bool, error) {
	if email == "" {
		return false, nil
	}
	var count int64
	err := r.db.Model(&models.AdminUser{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
