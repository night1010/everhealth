package transactor

import (
	"context"

	"gorm.io/gorm"
)

type Manager interface {
	Run(ctx context.Context, runner func(c context.Context) error) error
}

type manager struct {
	db *gorm.DB
}

func NewManager(db *gorm.DB) Manager {
	return &manager{db: db}
}

func (t *manager) Run(ctx context.Context, runner func(c context.Context) error) error {
	tx := t.db.Begin()
	ctx = injectTx(ctx, tx)

	err := runner(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func injectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, "tx", tx)
}
