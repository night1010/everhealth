package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/appjwt"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/hasher"
	"github.com/night1010/everhealth/imagehelper"
	"github.com/night1010/everhealth/mail"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/valueobject"
)

type AuthUsecase interface {
	Register(context.Context, *entity.User) error
	Verify(context.Context, *entity.User, *entity.Profile, *entity.DoctorProfile) error
	Login(context.Context, *entity.User) (*entity.User, error)
	ForgotPassword(context.Context, *entity.User, *entity.ForgotPasswordToken) error
	ResetPassword(context.Context, *entity.User, *entity.ForgotPasswordToken) error
}

type authUsecase struct {
	manager            transactor.Manager
	userRepo           repository.UserRepository
	profileRepo        repository.ProfileRepository
	doctorProfileRepo  repository.DoctorProfileRepository
	forgotPasswordRepo repository.ForgotPasswordRepository
	cartRepo           repository.CartRepository
	smtpGmail          mail.SmtpGmail
	hash               hasher.Hasher
	jwt                appjwt.Jwt
	imageHelper        imagehelper.ImageHelper
}

func NewAuthUsecase(
	manager transactor.Manager,
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	doctorProfileRepo repository.DoctorProfileRepository,
	forgotPasswordRepo repository.ForgotPasswordRepository,
	cartRepo repository.CartRepository,
	smtpGmail mail.SmtpGmail,
	hash hasher.Hasher,
	jwt appjwt.Jwt,
	imageHelper imagehelper.ImageHelper,
) AuthUsecase {
	return &authUsecase{
		manager:            manager,
		userRepo:           userRepo,
		profileRepo:        profileRepo,
		doctorProfileRepo:  doctorProfileRepo,
		forgotPasswordRepo: forgotPasswordRepo,
		cartRepo:           cartRepo,
		smtpGmail:          smtpGmail,
		hash:               hash,
		jwt:                jwt,
		imageHelper:        imageHelper,
	}
}

func (u *authUsecase) Register(ctx context.Context, user *entity.User) error {
	emailQuery := valueobject.NewQuery().Condition("email", valueobject.Equal, user.Email)
	fetchedUser, err := u.userRepo.FindOne(ctx, emailQuery)
	if err != nil {
		return err
	}
	var token string
	if fetchedUser != nil {
		if fetchedUser.IsVerified {
			return apperror.NewResourceAlreadyExistError("user", "email", user.Email)
		}
		token = fetchedUser.Token
	} else {
		hashedToken, err := u.hash.Hash(user.Email)
		if err != nil {
			return err
		}
		token = string(hashedToken)
		user.Token = token
		_, err = u.userRepo.Create(ctx, user)
		if err != nil {
			return err
		}
	}
	link := fmt.Sprintf("verify/%d?token=%s", user.RoleId, token)
	err = u.smtpGmail.SendEmail(link, user.Email, true)
	if err != nil {
		return err
	}
	return nil
}

func (u *authUsecase) Verify(ctx context.Context, user *entity.User, profile *entity.Profile, doctorProfile *entity.DoctorProfile) error {
	err := u.manager.Run(ctx, func(c context.Context) error {
		tokenQuery := valueobject.NewQuery().
			Condition("token", valueobject.Equal, user.Token).Lock()
		fetchedUser, err := u.userRepo.FindOne(c, tokenQuery)
		if err != nil {
			return err
		}
		if fetchedUser == nil {
			return apperror.NewClientError(apperror.NewInvalidTokenError()).BadRequest()
		}
		if fetchedUser.IsVerified {
			return apperror.NewClientError(apperror.NewResourceStateError("Already Verified"))
		}
		if fetchedUser.RoleId != user.RoleId {
			return apperror.NewClientError(apperror.NewForbiddenActionError("Different Role"))
		}
		hashPass, err := u.hash.Hash(user.Password)
		if err != nil {
			return err
		}
		fetchedUser.Password = string(hashPass)
		fetchedUser.IsVerified = true
		updatedUser, err := u.userRepo.Update(c, fetchedUser)
		if err != nil {
			return err
		}

		profile.UserId = updatedUser.Id
		_, err = u.profileRepo.Create(c, profile)
		if err != nil {
			return err
		}
		if fetchedUser.RoleId == entity.RoleUser {
			var cart entity.Cart
			cart.UserId = updatedUser.Id
			_, err = u.cartRepo.Create(c, &cart)
			if err != nil {
				return err
			}
			return nil
		}
		pdf := c.Value("pdf")
		pdfKey := entity.DoctorCertificatePrefix + generateRandomString(10)
		pdfUrl, err := u.imageHelper.Upload(ctx, pdf.(multipart.File), entity.DoctorCertificateFolder, pdfKey)
		if err != nil {
			return err
		}
		doctorProfile.CertificateKey = pdfKey
		doctorProfile.Certificate = pdfUrl
		doctorProfile.ProfileId = updatedUser.Id
		_, err = u.doctorProfileRepo.Create(c, doctorProfile)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (u *authUsecase) Login(ctx context.Context, user *entity.User) (*entity.User, error) {
	emailQuery := valueobject.NewQuery().Condition("email", valueobject.Equal, user.Email)
	fetchedUser, err := u.userRepo.FindOne(ctx, emailQuery)
	if err != nil {
		return nil, err
	}
	if fetchedUser == nil {
		return nil, apperror.NewResourceNotFoundError("user", "email", user.Email)
	}
	if !(u.hash.Compare(fetchedUser.Password, user.Password)) {
		return nil, apperror.NewInvalidCredentialsError()
	}
	token, err := u.jwt.GenerateToken(fetchedUser)
	if err != nil {
		return nil, err
	}
	fetchedUser.Token = token
	return fetchedUser, nil
}

func (u *authUsecase) ForgotPassword(ctx context.Context, user *entity.User, tokenEntity *entity.ForgotPasswordToken) error {
	emailQuery := valueobject.NewQuery().Condition("email", valueobject.Equal, user.Email)
	fetchedUser, err := u.userRepo.FindOne(ctx, emailQuery)
	if err != nil {
		return err
	}
	if fetchedUser == nil {
		return apperror.NewResourceNotFoundError("user", "email", user.Email)
	}
	if !fetchedUser.IsVerified {
		return apperror.NewResourceNotFoundError("user", "email", user.Email)
	}
	hashedToken, err := u.hash.Hash(user.Email)
	if err != nil {
		return err
	}
	tokenEntity.Token = string(hashedToken)
	tokenEntity.UserId = fetchedUser.Id
	tokenEntity, err = u.forgotPasswordRepo.Create(ctx, tokenEntity)
	if err != nil {
		return err
	}
	tokenLink := fmt.Sprintf("forgot-password/apply?token=%s", tokenEntity.Token)
	err = u.smtpGmail.SendEmail(tokenLink, fetchedUser.Email, false)
	if err != nil {
		return err
	}
	return nil
}

func (u *authUsecase) ResetPassword(ctx context.Context, user *entity.User, tokenEntity *entity.ForgotPasswordToken) error {
	err := u.manager.Run(ctx, func(c context.Context) error {
		tokenQuery := valueobject.NewQuery().
			Condition("token", valueobject.Equal, tokenEntity.Token).Lock()
		fetchedTokenEntity, err := u.forgotPasswordRepo.FindOne(c, tokenQuery)
		if err != nil {
			return err
		}
		if fetchedTokenEntity == nil {
			return apperror.NewClientError(apperror.NewInvalidTokenError()).BadRequest()
		}
		if !fetchedTokenEntity.IsActive {
			return apperror.NewInvalidTokenError()
		}
		if fetchedTokenEntity.ExpiredAt.Before(time.Now()) {
			return apperror.NewInvalidTokenError()
		}
		userQuery := valueobject.NewQuery().
			Condition("id", valueobject.Equal, fetchedTokenEntity.UserId).Lock()
		fetchedUser, err := u.userRepo.FindOne(c, userQuery)
		if err != nil {
			return err
		}
		if u.hash.Compare(fetchedUser.Password, user.Password) {
			return apperror.NewClientError(apperror.NewResourceStateError("can't use the same password"))
		}
		hashPass, err := u.hash.Hash(user.Password)
		if err != nil {
			return err
		}
		fetchedUser.Password = string(hashPass)
		fetchedTokenEntity.IsActive = false
		_, err = u.userRepo.Update(c, fetchedUser)
		if err != nil {
			return err
		}
		_, err = u.forgotPasswordRepo.Update(c, fetchedTokenEntity)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
