package repository

import (
	"github.com/night1010/everhealth/entity"
	"gorm.io/gorm"
)

type ChatRepository interface {
	BaseRepository[entity.Chat]
}

type chatRepository struct {
	*baseRepository[entity.Chat]
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Chat]{db: db},
	}
}
