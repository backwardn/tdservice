package postgres

import (
	"intel/isecl/threat-detection-service/types"

	"github.com/jinzhu/gorm"
)

type PostgresHostRepository struct {
	db *gorm.DB
}

func (r *PostgresHostRepository) Create(host types.Host) error {
	return r.db.Create(&host).Error
}

func (r *PostgresHostRepository) Retrieve(host types.Host) (*types.Host, error) {
	err := r.db.First(&host).Error
	return &host, err
}

func (r *PostgresHostRepository) Update(host types.Host) error {
	return r.db.Save(&host).Error
}

func (r *PostgresHostRepository) Delete(host types.Host) error {
	return r.db.Delete(&host).Error
}
