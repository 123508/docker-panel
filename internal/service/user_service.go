package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"docker-panel/internal/db"
	"docker-panel/internal/models"
	"docker-panel/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	InitAdmin() error
	Login(username, password string) (string, error)
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) InitAdmin() error {
	var count int64
	db.DB.Model(&models.User{}).Where("username = ?", "admin").Count(&count)

	if count > 0 {
		return nil
	}

	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Errorf("failed to generate random password: %w", err)
	}
	randomPassword := hex.EncodeToString(bytes)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin := models.User{
		Username: "admin",
		Password: string(hashedPassword),
	}

	if err := db.DB.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Printf("========================================================\n")
	log.Printf("Initial Admin Account Created\n")
	log.Printf("Username: admin\n")
	log.Printf("Password: %s\n", randomPassword)
	log.Printf("Please change the password after logging in.\n")
	log.Printf("========================================================\n")

	return nil
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
