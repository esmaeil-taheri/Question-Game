package userservice

import (
	"fmt"
	"gameapp/dto"
	"gameapp/entity"

	"gameapp/pkg/richerror"
	"gameapp/pkg/security/password"
)

type Repository interface {
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





func New(authGenerator AuthGenerator, repo Repository) Service {
	return Service{repo: repo, auth: authGenerator}
}

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// TODO - we should verify phone number by verification code

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("can't hash password")
	}

	user := entity.User{
		ID: 0,
		PhoneNumber: req.PhoneNumber,
		Name: req.Name,
		Password: hashedPassword,
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return dto.RegisterResponse{User: dto.UserInfo{
		ID: createdUser.ID,
		PhoneNumber: createdUser.PhoneNumber,
		Name: createdUser.Name,
	}}, nil
} 

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password string `json:"password"`
}

type Tokens struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User dto.UserInfo `json:"user"`
	Tokens Tokens `json:"tokens"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	const op = "userservice.Login"

	// TODO - it would be better to use two seprate method for existence check and GetUserByPhoneNumber
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, richerror.New(op).WithErr(err).WithMeta(
			map[string]interface{}{"phone_number": req.PhoneNumber},
		)
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

	return LoginResponse{
		User: dto.UserInfo{
			ID: user.ID,
			PhoneNumber: user.PhoneNumber,
			Name: user.Name,
		},
		Tokens: Tokens{
			AccessToken: accesstoken,
			RefreshToken: refreshtoken,
		},
	}, nil
}

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	Name string `json:"name"`
}

func (s Service) Profile(req ProfileRequest) (ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	return ProfileResponse{Name: user.Name}, nil

}
