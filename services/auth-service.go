package services

import (
	"log"

	"github.com/alonecandies/mysql-gin-gorm-auth/api/dtos"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/entities"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/repositories"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredentials(email string, password string) interface{}
	CreateUser(userCreateDTO *dtos.RegisterDTO) entities.User
	FindByEmail(email string) entities.User
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

func comparePassword(password []byte, hash []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, password) == nil
}

func (service *authService) VerifyCredentials(email string, password string) interface{} {
	res := service.userRepository.VerifyCredentials(email, password)
	v := res.(entities.User)
	if res != nil {
		comparePassword := comparePassword([]byte(password), []byte(v.Password))
		if v.Email == email && comparePassword {
			return res
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (service *authService) CreateUser(userCreateDTO *dtos.RegisterDTO) entities.User {
	user := entities.User{
		Email:    userCreateDTO.Email,
		Password: userCreateDTO.Password,
		Name:     userCreateDTO.Name,
	}
	err := smapping.FillStruct(&user, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed to map: %v", err)
	}
	res := service.userRepository.InsertUser(user)
	return res
}

func (service *authService) FindByEmail(email string) entities.User {
	return service.userRepository.FindByEmail(email)
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)
	if res.Error != nil {
		return true
	} else {
		return false
	}
}
