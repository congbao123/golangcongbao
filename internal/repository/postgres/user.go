package postgres

import (
	"errors"

	"todo-app/domain"
	"todo-app/pkg/clients"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Save(user *domain.UserCreate) error {
	if err := r.db.Create(&user).Error; err != nil {
		return clients.ErrDB(err)
	}

	return nil
}

func (r *userRepo) GetUser(conditions map[string]any) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where(conditions).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, clients.ErrRecordNotFound
		}

		return nil, clients.ErrDB(err)
	}

	return &user, nil
}
func (r *userRepo) GetByIdUser(id uuid.UUID) (domain.User, error) {
    var user domain.User

    if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return domain.User{}, clients.ErrRecordNotFound
        }
        return domain.User{}, clients.ErrDB(err)
    }

    return user, nil
}


func (r *userRepo) Update(id uuid.UUID, user *domain.UserUpdate) error {
	if err := r.db.Where("id = ?", id).Updates(&user).Error; err != nil {
		return clients.ErrDB(err)
	}

	return nil
}
func (r *userRepo) Delete(id uuid.UUID) error {
	if err := r.db.Table(domain.User{}.TableName()).Where("id = ?", id).Delete(nil).Error; err != nil {
		return clients.ErrDB(err)
	}

	return nil
}