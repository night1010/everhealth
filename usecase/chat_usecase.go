package usecase

import (
	"context"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/repository"
)

type ChatUsecase interface {
	ValidateTelemedicine(ctx context.Context, telemedicineId uint) (bool, error)
	AddChatMessage(ctx context.Context, chat *entity.Chat) (*entity.Chat, error)
}

type chatUsecase struct {
	chatRepo         repository.ChatRepository
	telemedicineRepo repository.TelemedicineRepository
}

func NewChatUsecase(chatRepo repository.ChatRepository, telemedicineRepo repository.TelemedicineRepository) ChatUsecase {
	return &chatUsecase{
		chatRepo:         chatRepo,
		telemedicineRepo: telemedicineRepo,
	}
}

func (u *chatUsecase) ValidateTelemedicine(ctx context.Context, telemedicineId uint) (bool, error) {
	userId := ctx.Value("user_id").(uint)

	fetchedTelemedicine, err := u.telemedicineRepo.FindById(ctx, telemedicineId)
	if err != nil {
		return false, err
	}
	if fetchedTelemedicine == nil {
		return false, apperror.NewResourceNotFoundError("telemedicine", "id", telemedicineId)
	}

	if fetchedTelemedicine.ProfileId == userId || fetchedTelemedicine.DoctorId == userId {
		return true, nil
	}

	return false, nil
}

func (u *chatUsecase) AddChatMessage(ctx context.Context, chat *entity.Chat) (*entity.Chat, error) {
	return u.chatRepo.Create(ctx, chat)
}
