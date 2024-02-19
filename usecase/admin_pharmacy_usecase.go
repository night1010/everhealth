package usecase

import (
	"context"
	"errors"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/hasher"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/valueobject"
)

type AdminPharmacyUsecase interface {
	FindOneAdminPharmacy(ctx context.Context, adminPharmacy *entity.User) (*entity.User, error)
	FindAllAdminPharmacy(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	CreateAdminPharmacy(ctx context.Context, adminPharmacy *entity.User, contact *entity.AdminContact) (*entity.User, error)
	UpdateAdminPharmacy(ctx context.Context, adminPharmacy *entity.User, contact *entity.AdminContact) (*entity.User, error)
	DeleteAdminPharmacy(ctx context.Context, adminPharmacy *entity.User) error
}

type adminPharmacyUsecase struct {
	adminPharmacyRepository repository.UserRepository
	pharmacyRepository      repository.PharmacyRepository
	adminContactRepository  repository.AdminContactRepository
	hash                    hasher.Hasher
	manager                 transactor.Manager
}

func NewAdminPharmacyUsecase(r repository.UserRepository, h hasher.Hasher, p repository.PharmacyRepository, c repository.AdminContactRepository, m transactor.Manager) AdminPharmacyUsecase {
	return &adminPharmacyUsecase{adminPharmacyRepository: r, hash: h, pharmacyRepository: p, adminContactRepository: c, manager: m}
}

func (u *adminPharmacyUsecase) FindAllAdminPharmacy(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.adminPharmacyRepository.FindAllAdminPharmacy(ctx, query)
}

func (u *adminPharmacyUsecase) FindOneAdminPharmacy(ctx context.Context, adminPharmacy *entity.User) (*entity.User, error) {
	selectAdmin, err := u.adminPharmacyRepository.FindOne(ctx, valueobject.NewQuery().Condition("\"users\".id", valueobject.Equal, adminPharmacy.Id).Condition("\"users\".role_id", valueobject.Equal, entity.RoleAdmin).WithPreload("AdminContact"))
	if err != nil {
		return nil, err
	}
	if selectAdmin == nil {
		return nil, apperror.NewResourceNotFoundError("admin pharmacy", "id", adminPharmacy.Id)
	}
	return selectAdmin, err
}

func (u *adminPharmacyUsecase) CreateAdminPharmacy(ctx context.Context, adminPharmacy *entity.User, contact *entity.AdminContact) (*entity.User, error) {
	err := u.manager.Run(ctx, func(c context.Context) error {
		hashPass, err := u.hash.Hash(adminPharmacy.Password)
		if err != nil {
			return err
		}
		adminPharmacy.RoleId = entity.RoleAdmin
		adminPharmacy.Password = string(hashPass)
		adminPharmacy.IsVerified = true
		_, err = u.adminPharmacyRepository.Create(c, adminPharmacy)
		if err != nil {
			return err
		}
		contact.UserId = adminPharmacy.Id
		_, err = u.adminContactRepository.Create(c, contact)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return adminPharmacy, nil
}

func (u *adminPharmacyUsecase) UpdateAdminPharmacy(ctx context.Context, adminPharmacy *entity.User, contact *entity.AdminContact) (*entity.User, error) {
	err := u.manager.Run(ctx, func(c context.Context) error {
		selectAdmin, err := u.adminPharmacyRepository.FindById(ctx, adminPharmacy.Id)
		if err != nil {
			return err
		}
		if selectAdmin == nil {
			return apperror.NewResourceNotFoundError("admin pharmacy", "id", adminPharmacy.Id)
		}
		adminPharmacy.RoleId = selectAdmin.RoleId
		adminPharmacy.Password = selectAdmin.Password
		_, err = u.adminPharmacyRepository.Update(c, adminPharmacy)
		if err != nil {
			return err
		}
		contact.UserId = adminPharmacy.Id
		_, err = u.adminContactRepository.Update(c, contact)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return adminPharmacy, nil
}

func (u *adminPharmacyUsecase) DeleteAdminPharmacy(ctx context.Context, adminPharmacy *entity.User) error {
	checkAdmin, err := u.pharmacyRepository.FindOne(ctx, valueobject.NewQuery().WithJoin("Admin").Condition("\"Admin\".id", valueobject.Equal, adminPharmacy.Id))
	if err != nil {
		return err
	}
	if checkAdmin != nil {
		return apperror.NewClientError(errors.New("cannot delete admin pharmacy because this admin already manage pharmacy"))
	}
	err = u.manager.Run(ctx, func(c context.Context) error {
		err = u.adminContactRepository.HardDelete(c, &entity.AdminContact{UserId: adminPharmacy.Id})
		if err != nil {
			return err
		}
		err = u.adminPharmacyRepository.HardDelete(c, adminPharmacy)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
