package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/imagehelper"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/util"
	"github.com/night1010/everhealth/valueobject"
)

type TelemedicineUsecase interface {
	CreateTelemedicine(ctx context.Context, telemedicine *entity.Telemedicine) (*entity.Telemedicine, error)
	PaymentProofTelemedicine(ctx context.Context, telemedicine *entity.Telemedicine) (*entity.Telemedicine, error)
	SickLeaveTelemedicine(ctx context.Context, telemedicine *entity.Telemedicine, sickLeave *dto.SickLeave) (*entity.Telemedicine, error)
	FindAllTelemedicine(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	FindTelemedicine(ctx context.Context, telemedicineId uint) (*entity.Telemedicine, error)
	PrescriptionTelemedicine(ctx context.Context, telemedicine *entity.Telemedicine, prescription *dto.Prescription) (*entity.Telemedicine, error)
	EndTelemedicine(ctx context.Context, telemedicineId uint) error
}

type telemedicineUsecase struct {
	imageHelper            imagehelper.ImageHelper
	manager                transactor.Manager
	doctorRepository       repository.DoctorProfileRepository
	telemedicineRepository repository.TelemedicineRepository
}

func NewTelemedicineUsecase(rp repository.TelemedicineRepository, dr repository.DoctorProfileRepository, m transactor.Manager, img imagehelper.ImageHelper) TelemedicineUsecase {
	return &telemedicineUsecase{telemedicineRepository: rp, doctorRepository: dr, imageHelper: img, manager: m}
}

func (u *telemedicineUsecase) FindAllTelemedicine(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	page, err := u.telemedicineRepository.FindAllTelemedicine(ctx, query)
	if err != nil {
		return nil, err
	}
	for _, telemedecine := range page.Data.([]*entity.Telemedicine) {
		if telemedecine.ExpiredAt.Before(time.Now()) && telemedecine.Status == entity.Waiting {
			telemedecine.Status = entity.Cancel
			_, err := u.telemedicineRepository.Update(ctx, telemedecine)
			if err != nil {
				return nil, err
			}
		}
	}
	return page, nil
}

func (u *telemedicineUsecase) FindTelemedicine(ctx context.Context, telemedicineId uint) (*entity.Telemedicine, error) {
	roleId := ctx.Value("role_id").(entity.RoleId)
	userId := ctx.Value("user_id").(uint)
	telemedicineQuery := valueobject.NewQuery().Condition("id", valueobject.Equal, telemedicineId)
	if roleId == entity.RoleDoctor {
		telemedicineQuery = telemedicineQuery.Condition("doctor_id", valueobject.Equal, userId)
	} else if roleId == entity.RoleUser {
		telemedicineQuery = telemedicineQuery.Condition("profile_id", valueobject.Equal, userId)
	}
	telemedicineQuery = telemedicineQuery.
		WithPreload("Doctor.Profile").
		WithPreload("Profile")
	fetchedTelemedicine, err := u.telemedicineRepository.FindByTelemedicineId(ctx, telemedicineQuery)
	if err != nil {
		return nil, err
	}
	if fetchedTelemedicine == nil {
		return nil, apperror.NewClientError(apperror.NewResourceNotFoundError("telemedicine", "id", telemedicineId))
	}
	return fetchedTelemedicine, nil
}

func (u *telemedicineUsecase) CreateTelemedicine(ctx context.Context, telemedicine *entity.Telemedicine) (*entity.Telemedicine, error) {
	doctor, err := u.doctorRepository.FindOne(ctx, valueobject.NewQuery().Condition("profile_id", valueobject.Equal, telemedicine.DoctorId))
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, apperror.NewClientError(fmt.Errorf("doctor id %v not found", telemedicine.DoctorId))
	}
	if doctor.Status != entity.Online {
		return nil, apperror.NewClientError(fmt.Errorf("doctor id %v %v", telemedicine.DoctorId, doctor.Status))
	}
	profileId := ctx.Value("user_id").(uint)
	telemedicine.ProfileId = profileId
	checkTelemedicine, err := u.telemedicineRepository.FindTelemedicineWhereStatusNotCancelAndEnd(ctx, telemedicine)
	if err != nil {
		return nil, err
	}
	if checkTelemedicine != nil {
		return nil, apperror.NewClientError(errors.New("cannot create another telemedicine, because you have ongoing telemedicine"))
	}
	telemedicine.ProfileId = profileId
	telemedicine.OrderedAt = time.Now()
	telemedicine.ExpiredAt = telemedicine.OrderedAt.Add(10 * time.Minute)
	telemedicine.TotalPayment = doctor.Fee
	telemedicine.Status = entity.Waiting
	return u.telemedicineRepository.Create(ctx, telemedicine)
}

func (u *telemedicineUsecase) PaymentProofTelemedicine(ctx context.Context, telemedicine *entity.Telemedicine) (*entity.Telemedicine, error) {
	var imageKey string
	var newTelemedicine *entity.Telemedicine
	var err error
	err = u.manager.Run(ctx, func(c context.Context) error {
		newTelemedicine, err = u.telemedicineRepository.FindOne(c, valueobject.NewQuery().Condition("id", valueobject.Equal, telemedicine.Id).Lock())
		if err != nil {
			return err
		}
		if newTelemedicine == nil {
			return apperror.NewResourceNotFoundError("telemedicine", "id", telemedicine.Id)
		}
		if newTelemedicine.ProfileId != ctx.Value("user_id").(uint) {
			return apperror.NewForbiddenActionError("cannot pay this telemedicine becase is belong to another user")
		}
		if newTelemedicine.Status != entity.Waiting {
			return apperror.NewClientError(fmt.Errorf("cannot pay because telemedicine already %v", newTelemedicine.Status))
		}
		image := ctx.Value("image")
		imageKey = entity.TelemedicineProofPrefix + generateRandomString(10)
		imgUrl, err := u.imageHelper.Upload(ctx, image.(multipart.File), entity.TelemedicineProofFolder, imageKey)
		if err != nil {
			return err
		}
		newTelemedicine.Proof = imgUrl
		newTelemedicine.ProofKey = imageKey
		newTelemedicine.Status = entity.Ongoing
		if newTelemedicine.ExpiredAt.Before(time.Now()) {
			newTelemedicine.Status = entity.Cancel
		}
		_, err = u.telemedicineRepository.Update(c, newTelemedicine)
		if err != nil {
			err2 := u.imageHelper.Destroy(ctx, entity.TelemedicineProofFolder, imageKey)
			if err2 != nil {
				return err2
			}
			return err
		}
		fetchedDoctor, err := u.doctorRepository.FindOne(c, valueobject.NewQuery().Condition("profile_id", valueobject.Equal, newTelemedicine.DoctorId))
		if err != nil {
			err2 := u.imageHelper.Destroy(ctx, entity.TelemedicineProofFolder, imageKey)
			if err2 != nil {
				return err2
			}
			return err
		}
		fetchedDoctor.Status = entity.Busy
		_, err = u.doctorRepository.Update(c, fetchedDoctor)
		if err != nil {
			err2 := u.imageHelper.Destroy(ctx, entity.TelemedicineProofFolder, imageKey)
			if err2 != nil {
				return err2
			}
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if newTelemedicine.Status == entity.Cancel {
		return nil, apperror.NewClientError(errors.New("cannot pay because telemedecine already expired"))
	}
	return telemedicine, nil
}

func (u *telemedicineUsecase) SickLeaveTelemedicine(ctx context.Context, telemedicine *entity.Telemedicine, sickLeave *dto.SickLeave) (*entity.Telemedicine, error) {
	var pdfKey string
	var newTelemedicine *entity.Telemedicine
	var err error
	err = u.manager.Run(ctx, func(c context.Context) error {
		newTelemedicine, err = u.telemedicineRepository.FindOne(c, valueobject.NewQuery().Condition("id", valueobject.Equal, telemedicine.Id).WithJoin("Profile").WithJoin("Doctor.Profile"))
		if err != nil {
			return err
		}
		if newTelemedicine == nil {
			return apperror.NewResourceNotFoundError("telemedicine", "id", telemedicine.Id)
		}
		if newTelemedicine.DoctorId != ctx.Value("user_id").(uint) {
			return apperror.NewForbiddenActionError("cannot create sick leave because this telemedicine belong to another doctor")
		}
		if newTelemedicine.Status != entity.Ongoing {
			return apperror.NewClientError(fmt.Errorf("telemedicine status already %v", newTelemedicine.Status))
		}
		if newTelemedicine.SickLeavePdf != "" {
			return apperror.NewClientError(errors.New("sick leave document already created"))
		}
		sickLeavePdf, err := util.CreateSickLeave(sickLeave.StartDate, sickLeave.EndDate, newTelemedicine.Profile, newTelemedicine.Doctor, sickLeave.Diagnosa)
		if err != nil {
			return err
		}
		pdfKey = entity.SickLeavePrefix + generateRandomString(10)
		pdfUrl, err := u.imageHelper.Upload(ctx, sickLeavePdf, entity.SickLeaveFolder, pdfKey)
		if err != nil {
			return err
		}
		newTelemedicine.SickLeavePdf = pdfUrl
		newTelemedicine.SickLeavePdfKey = pdfKey
		_, err = u.telemedicineRepository.Update(c, newTelemedicine)
		if err != nil {
			err2 := u.imageHelper.Destroy(ctx, entity.TelemedicineProofFolder, pdfKey)
			if err2 != nil {
				return err2
			}
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return telemedicine, nil
}

func (u *telemedicineUsecase) PrescriptionTelemedicine(ctx context.Context, telemedicine *entity.Telemedicine, prescription *dto.Prescription) (*entity.Telemedicine, error) {
	var pdfKey string
	var newTelemedicine *entity.Telemedicine
	var err error
	err = u.manager.Run(ctx, func(c context.Context) error {
		newTelemedicine, err = u.telemedicineRepository.FindOne(c, valueobject.NewQuery().Condition("id", valueobject.Equal, telemedicine.Id).WithJoin("Profile").WithJoin("Doctor.Profile"))
		if err != nil {
			return err
		}
		if newTelemedicine == nil {
			return apperror.NewResourceNotFoundError("telemedicine", "id", telemedicine.Id)
		}
		if newTelemedicine.DoctorId != ctx.Value("user_id").(uint) {
			return apperror.NewForbiddenActionError("cannot create prescription because this telemedicine belong to another doctor")
		}
		if newTelemedicine.Status != entity.Ongoing {
			return apperror.NewClientError(fmt.Errorf("telemedicine status already %v", newTelemedicine.Status))
		}
		if newTelemedicine.PrescriptionPdf != "" {
			return apperror.NewClientError(errors.New("prescription already created"))
		}
		prescriptionPdf, err := util.CreatePrescriptionLeave(prescription.Prescription, newTelemedicine.Profile, newTelemedicine.Doctor)
		if err != nil {
			return err
		}
		pdfKey = entity.PrescriptionPrefix + generateRandomString(10)
		pdfUrl, err := u.imageHelper.Upload(ctx, prescriptionPdf, entity.PrescriptionFolder, pdfKey)
		if err != nil {
			return err
		}
		newTelemedicine.PrescriptionPdf = pdfUrl
		newTelemedicine.PrescriptionPdfKey = pdfKey
		_, err = u.telemedicineRepository.Update(c, newTelemedicine)
		if err != nil {
			err2 := u.imageHelper.Destroy(ctx, entity.TelemedicineProofFolder, pdfKey)
			if err2 != nil {
				return err2
			}
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return telemedicine, nil
}

func (u *telemedicineUsecase) EndTelemedicine(ctx context.Context, telemedicineId uint) error {
	query := valueobject.NewQuery().Condition("id", valueobject.Equal, telemedicineId)
	fetchedTelemedicine, err := u.telemedicineRepository.FindOne(ctx, query)
	if err != nil {
		return err
	}
	if fetchedTelemedicine == nil {
		return apperror.NewResourceNotFoundError("telemedicine", "id", telemedicineId)
	}

	if fetchedTelemedicine.Status != entity.Ongoing {
		return apperror.NewResourceStateError("can only end ongoing telemedicine")
	}

	userId := ctx.Value("user_id").(uint)

	if fetchedTelemedicine.DoctorId != userId && fetchedTelemedicine.ProfileId != userId {
		return apperror.NewForbiddenActionError("insufficient access")
	}

	err = u.manager.Run(ctx, func(c context.Context) error {
		fetchedTelemedicine.Status = entity.End
		_, err = u.telemedicineRepository.Update(c, fetchedTelemedicine)
		if err != nil {
			return err
		}

		doctorQuery := valueobject.NewQuery().
			Condition("profile_id", valueobject.Equal, fetchedTelemedicine.DoctorId).
			Lock()
		fetchedDoctor, err := u.doctorRepository.FindOne(c, doctorQuery)
		if err != nil {
			return err
		}

		fetchedDoctor.Status = entity.Online
		_, err = u.doctorRepository.Update(c, fetchedDoctor)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
