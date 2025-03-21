package sql

import (
	"context"
	"errors"
	"sip-monitor/src/entity"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUser creates a new user in the database
func (r *GormRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if user.CreateAt.IsZero() {
		user.CreateAt = time.Now()
	}
	user.UpdateAt = time.Now()
	return r.db.WithContext(ctx).Create(user).Error
}

// GetUserByID retrieves a user by ID
func (r *GormRepository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *GormRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user
func (r *GormRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	user.UpdateAt = time.Now()
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Updates(user).Error
}

// DeleteUser deletes a user by ID
func (r *GormRepository) DeleteUser(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

func (r *GormRepository) CreateDefaultAdminUser(ctx context.Context) error {
	var userNum int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Count(&userNum).Error
	if err != nil {
		return err
	}

	if userNum > 0 {
		return nil
	}

	// Create default admin user

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create the admin user
	adminUser := &entity.User{
		Username: "admin",
		Password: string(hashedPassword),
		Nickname: "Administrator",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	// Save the admin user to the database
	if err := r.CreateUser(ctx, adminUser); err != nil {
		return err
	}

	return nil
}

func (r *GormRepository) GetUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}
