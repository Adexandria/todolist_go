package repositories

import (
	"SM/models"
	"context"
	"log/slog"

	"gorm.io/gorm"
)

type userRepository struct {
	Db     *gorm.DB
	Logger *slog.Logger
}

func (u *userRepository) VerifyUserByUsername(username string) (bool, error) {
	var count int64

	err := u.Db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (u *userRepository) GetUserByRefreshToken(refreshToken string) (models.User, error) {
	ctx := context.Background()

	currentUser, err := gorm.G[models.User](u.Db).Where("refresh_token = ?", refreshToken).First(ctx)

	return currentUser, err
}

func (u *userRepository) GetUserById(id int) (models.User, error) {
	ctx := context.Background()

	currentUser, err := gorm.G[models.User](u.Db).Where("id = ?", id).First(ctx)

	return currentUser, err
}

func (u *userRepository) GetUserByUsername(username string) (models.User, error) {
	ctx := context.Background()

	currentUser, err := gorm.G[models.User](u.Db).Where("username = ?", username).First(ctx)

	return currentUser, err
}

func (u *userRepository) GetUserByEmail(email string) (models.User, error) {
	ctx := context.Background()

	currentUser, err := gorm.G[models.User](u.Db).Where("email = ?", email).First(ctx)

	return currentUser, err
}

func (u *userRepository) VerifyUserById(id int) (bool, error) {
	var count int64

	err := u.Db.Model(&models.User{}).Where("id = ?", id).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (u *userRepository) VerifyUserByEmail(email string) (bool, error) {
	var count int64

	err := u.Db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (u *userRepository) CreateUser(user *models.User) error {
	ctx := context.Background()

	err := gorm.G[models.User](u.Db).Create(ctx, user)

	return err
}

func (u *userRepository) UpdateUser(user *models.User) error {
	ctx := context.Background()

	_, err := gorm.G[models.User](u.Db).Where("id = ?", user.ID).Select("*").Updates(ctx, *user)

	return err
}

func (u *userRepository) DeleteUser(user *models.User) error {
	ctx := context.Background()

	_, err := gorm.G[models.User](u.Db).Where("id = ?", user.ID).Select("user_id").Delete(ctx)

	return err
}

func NewUserRepo(db *gorm.DB, handler *slog.JSONHandler) UserRepository {
	return &userRepository{
		Db:     db,
		Logger: slog.New(handler),
	}
}

var _ UserRepository = (*userRepository)(nil)
