package uservalidator

import (
	"gameapp/dto"
	"gameapp/pkg/errmsg"
	// "gameapp/pkg/phonenumber"
	"gameapp/pkg/richerror"

	"regexp"
	"fmt"
	
	"github.com/go-ozzo/ozzo-validation/v4"
	// "github.com/go-ozzo/ozzo-validation/v4/is"
)


type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}

// func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) error {
// 	const op = "uservalidator.ValidateRegisterRequest"

// 	if !phonenumber.IsValid(req.PhoneNumber) {
// 		return richerror.New(op).WithMessage(errmsg.ErrorMsgPhoneNumberIsNotValid).
// 		 WithKind(richerror.KindInvalid).WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
// 	}

// 	if isUnique, err := v.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {

// 		if err != nil {
// 			return richerror.New(op).WithErr(err)
// 		}

// 		if !isUnique {
// 			return richerror.New(op).WithMessage(errmsg.ErrorMsgPhoneNumberIsNotUnique).
// 		 	 WithKind(richerror.KindInvalid).WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
// 		}
// 	}

// 	// TODO - add length to config
// 	if len(req.Name) < 3 {
// 		return richerror.New(op).WithMessage(errmsg.ErrorMsgNameLength).
// 		 WithKind(richerror.KindInvalid).WithMeta(map[string]interface{}{"name": req.Name})
// 	}

// 	if len(req.Password) < 8 {
// 		return richerror.New(op).WithMessage(errmsg.ErrorMsgPasswordLength).
// 		 WithKind(richerror.KindInvalid).WithMeta(map[string]interface{}{"password": req.Password})
// 	}

// 	return nil

// }


func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateRegisterRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),
		validation.Field(&req.Password, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z0-9!@#$%^&*]{8,}$"))),
		validation.Field(
			&req.PhoneNumber, validation.Required, 
			validation.Match(regexp.MustCompile("^09[0-9]{9}$")), validation.By(v.checkPhoneNumberUniqueness)),
	); err != nil {
		fieldErrors := make(map[string]string)

		errV, ok := err.(validation.Errors)
		if ok {
			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}

		return fieldErrors, richerror.New(op).WithMessage(errmsg.ErrorMsgInvalidInput).
		 WithKind(richerror.KindInvalid).WithMeta(map[string]interface{}{"req": req}).WithErr(err)
	}

	return nil, nil

}

func (v Validator) checkPhoneNumberUniqueness(value interface{}) error {
	phoneNumber := value.(string)

	if isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber); err != nil || !isUnique {

		if err != nil {
			return err
		}

		if !isUnique {
			return fmt.Errorf(errmsg.ErrorMsgPhoneNumberIsNotUnique)
		}
	}

	return nil
}
