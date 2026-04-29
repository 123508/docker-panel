package service

import (
	"errors"
	"fmt"
	"log"

	"docker-panel/internal/config"
	"docker-panel/internal/db"
	"docker-panel/internal/models"
	"docker-panel/internal/utils"

	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	InitAdmin() (string, error)
	Login(username, password string) (string, error)
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) InitAdmin() (string, error) {
	adminUsername := "admin"
	configPassword := ""
	if config.AppConfig != nil {
		if config.AppConfig.User.AdminUsername != "" {
			adminUsername = config.AppConfig.User.AdminUsername
		}
		configPassword = config.AppConfig.User.AdminPassword
	}
	if configPassword == "" {
		return "", errors.New("admin password is empty in config.toml")
	}

	var admin models.User
	err := db.DB.Where("username = ?", adminUsername).First(&admin).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("failed to query admin user: %w", err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(configPassword), bcrypt.DefaultCost)
		if hashErr != nil {
			return "", fmt.Errorf("failed to hash password: %w", hashErr)
		}

		admin = models.User{
			Username: adminUsername,
			Password: string(hashedPassword),
		}

		if createErr := db.DB.Create(&admin).Error; createErr != nil {
			return "", fmt.Errorf("failed to create admin user: %w", createErr)
		}

		log.Printf("========================================================\n")
		log.Printf("Initial Admin Account Created\n")
		log.Printf("Username: %s\n", adminUsername)
		log.Printf("Password: %s\n", configPassword)
		log.Printf("Please change the password after logging in.\n")
		log.Printf("========================================================\n")
		return "", nil
	}

	if cmpErr := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(configPassword)); cmpErr == nil {
		return "", nil
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(configPassword), bcrypt.DefaultCost)
	if hashErr != nil {
		return "", fmt.Errorf("failed to hash password: %w", hashErr)
	}

	if updateErr := db.DB.Model(&admin).Update("password", string(hashedPassword)).Error; updateErr != nil {
		return "", fmt.Errorf("failed to update admin password: %w", updateErr)
	}

	log.Printf("admin password in DB is out of sync with config, updated DB password for user: %s", adminUsername)
	return "", nil
}

func (s *userService) Login(username, password string) (string, error) {
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
