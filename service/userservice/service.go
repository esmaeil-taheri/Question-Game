package userservice

import (
	"fmt"
	"gameapp/entity"

	"gameapp/pkg/phonenumber"
	"gameapp/pkg/security/password"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
	GetUserByID(userID uint) (entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User) (string, error)
}

type Service struct {
	auth AuthGenerator
	repo Repository
}

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Password	string `json:"password"`
}

type RegisterResponse struct {
	User entity.User
}

func New(authGenerator AuthGenerator, repo Repository) Service {
	return Service{repo: repo, auth: authGenerator}
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO - we should verify phone number by verification code

	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {

		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
		}

		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("phone number is not unique")
		}
	}

	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name should be grater than 3")
	}

	if len(req.Password) < 8 {
		return RegisterResponse{}, fmt.Errorf("password should be grater than 8")
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("can't hash password")
	}

	user := entity.User{
		ID: 0,
		PhoneNumber: req.PhoneNumber,
		Name: req.Name,
		Password: hashedPassword,
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return RegisterResponse{
		User: createdUser,
	}, nil
} 

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	// TODO - it would be better to use two seprate method for existence check and GetUserByPhoneNumber

	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	if !exist {
		return LoginResponse{}, fmt.Errorf("username or password is incorrect")
	}

	if !password.Compare(user.Password, req.Password) {
		return LoginResponse{}, fmt.Errorf("username or password is incorrect")
	}

	accesstoken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshtoken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return LoginResponse{AccessToken: accesstoken, RefreshToken: refreshtoken}, nil
}

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	Name string `json:"name"`
}

func (s Service) Profile(req ProfileRequest) (ProfileResponse, error) {

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		// I Have not expect the repository call return "record not found" error, because I assume the interactor input is sanitized.
		// TODO - we can use Rich Error.
		return ProfileResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return ProfileResponse{Name: user.Name}, nil

}
