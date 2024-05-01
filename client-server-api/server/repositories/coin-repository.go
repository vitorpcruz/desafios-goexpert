package repositories

import (
	"github.com/vitorpcruz/desafios-golang/client-server-api/server/models"
	"gorm.io/gorm"
)

type CoinRepositoryInterface interface {
	Save(coin *models.Coin) error
}

type CoinRepository struct {
	db *gorm.DB
}

func Init(gormDb *gorm.DB) *CoinRepository {
	return &CoinRepository{db: gormDb}
}

func (repo *CoinRepository) Save(coin *models.Coin) error {
	if err := repo.db.Save(coin); err != nil {
		return err.Error
	}
	return nil
}
