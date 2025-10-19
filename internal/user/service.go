package user

import "gorm.io/gorm"

type UserService struct {
	UserRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{UserRepository: userRepository}
}

func (s *UserService) AddUserToWhiteList(username string) error {
	user, err := s.UserRepository.FindByUsername(username)
	if err != nil {
		return err
	}

	err = s.UserRepository.Update(&User{
		Model:       gorm.Model{ID: user.ID},
		InWhitelist: true,
	})
	if err != nil {
		return err
	}

	return nil
}
