package user

import (
	"telefool/pkg/db"

	"gorm.io/gorm/clause"
)

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{database}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	res := repo.Database.Create(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByUsername(username string) (*User, error) {
	var user User
	res := repo.Database.DB.First(&user, "username = ?", username)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (repo *UserRepository) Update(user *User) error {
	res := repo.Database.DB.Clauses(clause.Returning{}).Updates(user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

//func (repo *UserRepository) UpdateField(user *User, field string, value any) error {
//	res := repo.Database.DB.Model(user).Update(field, value)
//	if res.Error != nil {
//		return res.Error
//	}
//
//	return nil
//}

func (repo *UserRepository) Delete(user *User) error {
	res := repo.Database.DB.Delete(user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
