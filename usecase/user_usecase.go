package usecase

import (
	"context"
	"mime/multipart"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/hasher"
	"github.com/night1010/everhealth/imagehelper"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/valueobject"
)

type UserUsecase interface {
	UserProfile(context.Context) (*entity.User, *entity.Profile, *entity.DoctorProfile, error)
	ResetPassword(context.Context, string, string) error
	UpdateProfile(context.Context, *entity.Profile, *entity.DoctorProfile) error
	ListAllDoctors(context.Context, *valueobject.Query) (*valueobject.PagedResult, error)
	DoctorDetail(context.Context, uint) (*entity.Profile, *entity.DoctorProfile, error)
	UpdateStatus(context.Context, entity.StatusDoctor) error
	GetAllUser(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
}

type userUsecase struct {
	manager           transactor.Manager
	userRepo          repository.UserRepository
	profilRepo        repository.ProfileRepository
	doctorProfileRepo repository.DoctorProfileRepository
	hash              hasher.Hasher
	imageHelper       imagehelper.ImageHelper
}

func NewUserUsecase(
	manager transactor.Manager,
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	doctorProfileRepo repository.DoctorProfileRepository,
	hash hasher.Hasher,
	imageHelper imagehelper.ImageHelper,

) UserUsecase {
	return &userUsecase{
		manager:           manager,
		userRepo:          userRepo,
		profilRepo:        profileRepo,
		doctorProfileRepo: doctorProfileRepo,
		hash:              hash,
		imageHelper:       imageHelper,
	}
}

func (u *userUsecase) UserProfile(ctx context.Context) (*entity.User, *entity.Profile, *entity.DoctorProfile, error) {
	userId := ctx.Value("user_id").(uint)
	fetchedUser, err := u.userRepo.FindById(ctx, userId)
	if err != nil {
		return nil, nil, nil, err
	}
	if fetchedUser == nil {
		return nil, nil, nil, apperror.NewClientError(apperror.NewInvalidCredentialsError())
	}
	uidQuery := valueobject.NewQuery().Condition("user_id", valueobject.Equal, userId)
	fetchProfile, err := u.profilRepo.FindOne(ctx, uidQuery)
	if err != nil {
		return nil, nil, nil, err
	}
	uidDoctorQuery := valueobject.NewQuery().Condition("profile_id", valueobject.Equal, userId).WithJoin("Specialist")
	var fetchDoctorProfile *entity.DoctorProfile
	if fetchedUser.RoleId == entity.RoleDoctor {
		fetchDoctorProfile, err = u.doctorProfileRepo.FindOne(ctx, uidDoctorQuery)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	return fetchedUser, fetchProfile, fetchDoctorProfile, nil
}

func (u *userUsecase) ResetPassword(ctx context.Context, oldPassword, newPassword string) error {
	userId := ctx.Value("user_id").(uint)
	fetchedUser, err := u.userRepo.FindById(ctx, userId)
	if err != nil {
		return err
	}
	if fetchedUser == nil {
		return apperror.NewClientError(apperror.NewInvalidCredentialsError())
	}
	if oldPassword == newPassword {
		return apperror.NewClientError(apperror.NewResourceStateError("can't change to the same password"))
	}
	if !(u.hash.Compare(fetchedUser.Password, oldPassword)) {
		return apperror.NewClientError(apperror.NewResourceStateError("incorrect old password"))
	}
	hashedPassword, err := u.hash.Hash(newPassword)
	if err != nil {
		return err
	}
	fetchedUser.Password = string(hashedPassword)
	_, err = u.userRepo.Update(ctx, fetchedUser)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) UpdateProfile(ctx context.Context, profile *entity.Profile, doctorProfile *entity.DoctorProfile) error {
	userId := ctx.Value("user_id").(uint)
	updatedProfileQuery := valueobject.NewQuery().
		Condition("user_id", valueobject.Equal, userId).Lock()
	updatedProfile, err := u.profilRepo.FindOne(ctx, updatedProfileQuery)
	if err != nil {
		return err
	}
	var imageKey string
	err = u.manager.Run(ctx, func(c context.Context) error {
		fetchedUser, err := u.userRepo.FindById(c, userId)
		if err != nil {
			return err
		}
		if fetchedUser == nil {
			return apperror.NewClientError(apperror.NewInvalidCredentialsError())
		}
		profileQuery := valueobject.NewQuery().
			Condition("user_id", valueobject.Equal, fetchedUser.Id).Lock()
		fetchedProfile, err := u.profilRepo.FindOne(c, profileQuery)
		if err != nil {
			return err
		}
		var imgUrl string
		image := c.Value("image")
		if image != nil {
			imageKey = fetchedProfile.ImageKey
			if fetchedProfile.ImageKey == "" {
				imageKey = entity.ProfilePhotoKeyPrefix + generateRandomString(10)
				fetchedProfile.ImageKey = imageKey
			}
			imgUrl, err = u.imageHelper.Upload(c, image.(multipart.File), entity.ProfilePhotoFolder, entity.ProfilePhotoKeyPrefix+generateRandomString(10))
			if err != nil {
				return err
			}
			fetchedProfile.Image = imgUrl
		}
		fetchedProfile.Name = profile.Name
		updatedProfile, err = u.profilRepo.Update(c, fetchedProfile)
		if err != nil {
			return err
		}
		if fetchedUser.RoleId == entity.RoleUser {
			return nil
		}
		doctorProfileQuery := valueobject.NewQuery().
			Condition("profile_id", valueobject.Equal, fetchedUser.Id).Lock()
		fetchedDoctorProfile, err := u.doctorProfileRepo.FindOne(c, doctorProfileQuery)
		if err != nil {
			return err
		}
		pdf := c.Value("pdf")
		var pdfUrl string
		if pdf != nil {
			pdfUrl, err = u.imageHelper.Upload(c, pdf.(multipart.File), entity.DoctorCertificateFolder, entity.DoctorCertificatePrefix+generateRandomString(10))
			if err != nil {
				return err
			}
			fetchedDoctorProfile.Certificate = pdfUrl
		}
		fetchedDoctorProfile.YearOfExperience = doctorProfile.YearOfExperience
		fetchedDoctorProfile.Fee = doctorProfile.Fee
		_, err = u.doctorProfileRepo.Update(c, fetchedDoctorProfile)
		if err != nil {
			return err
		}
		return nil
	})
	if updatedProfile.Image == "" {
		u.imageHelper.Destroy(ctx, entity.ProfilePhotoFolder, imageKey)
	}
	return err
}

func (u *userUsecase) ListAllDoctors(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.doctorProfileRepo.FindAllDoctors(ctx, query)
}

func (u *userUsecase) DoctorDetail(ctx context.Context, doctorId uint) (*entity.Profile, *entity.DoctorProfile, error) {
	fetchedUser, err := u.userRepo.FindById(ctx, doctorId)
	if err != nil {
		return nil, nil, err
	}
	if fetchedUser == nil {
		return nil, nil, apperror.NewClientError(apperror.NewResourceStateError("Doctor not found"))
	}
	if fetchedUser.RoleId != entity.RoleDoctor {
		return nil, nil, apperror.NewClientError(apperror.NewResourceStateError("Doctor not found"))
	}
	uidQuery := valueobject.NewQuery().Condition("user_id", valueobject.Equal, fetchedUser.Id)
	fetchProfile, err := u.profilRepo.FindOne(ctx, uidQuery)
	if err != nil {
		return nil, nil, err
	}
	uidDoctorQuery := valueobject.NewQuery().Condition("profile_id", valueobject.Equal, fetchedUser.Id).WithJoin("Specialist")
	var fetchDoctorProfile *entity.DoctorProfile
	if fetchedUser.RoleId == entity.RoleDoctor {
		fetchDoctorProfile, err = u.doctorProfileRepo.FindOne(ctx, uidDoctorQuery)
		if err != nil {
			return nil, nil, err
		}
	}
	return fetchProfile, fetchDoctorProfile, nil
}
func (u *userUsecase) UpdateStatus(ctx context.Context, status entity.StatusDoctor) error {
	userId := ctx.Value("user_id").(uint)
	doctorProfileQuery := valueobject.NewQuery().
		Condition("profile_id", valueobject.Equal, userId)
	fetchedProfile, err := u.doctorProfileRepo.FindOne(ctx, doctorProfileQuery)
	if err != nil {
		return err
	}
	fetchedProfile.Status = status
	_, err = u.doctorProfileRepo.Update(ctx, fetchedProfile)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) GetAllUser(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.userRepo.FindAllUser(ctx, query)
}
