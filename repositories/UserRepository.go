package repositories

import (
	"SM/models"
	"context"
	"log/slog"

	"gorm.io/gorm"
)

type UserRepository struct {
	Db     *gorm.DB
	Logger *slog.Logger
}

func (u *UserRepository) GetUserById(id int) (models.User, error) {
	ctx := context.Background()

	currentUser, err := gorm.G[models.User](u.Db).Where("id = ?", id).First(ctx)

	return currentUser, err
}

func (u *UserRepository) GetUserByUsername(username string) (models.User, error) {
	ctx := context.Background()

	currentUser, err := gorm.G[models.User](u.Db).Where("username = ?", username).First(ctx)

	return currentUser, err
}

func (u *UserRepository) GetUserByEmail(email string) (models.User, error) {
	ctx := context.Background()

	currentUser, err := gorm.G[models.User](u.Db).Where("email = ?", email).First(ctx)

	return currentUser, err
}

func (u *UserRepository) VerifyUserById(id int) (bool, error) {
	var count int64

	err := u.Db.Model(&models.User{}).Where("id = ?", id).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (u *UserRepository) VerifyUserByEmail(email string) (bool, error) {
	var count int64

	err := u.Db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (u *UserRepository) CreateUser(user *models.User) error {
	ctx := context.Background()

	err := gorm.G[models.User](u.Db).Create(ctx, user)

	return err
}

func (u *UserRepository) UpdateUser(user *models.User) error {
	ctx := context.Background()

	_, err := gorm.G[models.User](u.Db).Where("id = ?", user.ID).Updates(ctx, *user)

	return err
}

func (u *UserRepository) DeleteUser(user *models.User) error {
	ctx := context.Background()

	_, err := gorm.G[models.User](u.Db).Where("id = ?", user.ID).Select("user_id").Delete(ctx)

	return err
}

func UserRepoCon(db *gorm.DB, handler *slog.JSONHandler) *UserRepository {
	return &UserRepository{
		Db:     db,
		Logger: slog.New(handler),
	}
}

var _ IUserRepository = (*UserRepository)(nil)
