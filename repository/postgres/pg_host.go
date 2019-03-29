package postgres

import (
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/types"

	"github.com/jinzhu/gorm"
)

type PostgresHostRepository struct {
	db *gorm.DB
}

func (r *PostgresHostRepository) Create(host types.Host) (*types.Host, error) {

	uuid, err := repository.UUID()
	if err == nil {
		host.ID = uuid
	} else {
		return &host, err
	}
	err = r.db.Create(&host).Error
	return &host, err
}

func (r *PostgresHostRepository) Retrieve(host types.Host) (*types.Host, error) {
	err := r.db.First(&host).Error
	return &host, err
}

func (r *PostgresHostRepository) RetrieveAll(host types.Host) ([]types.Host, error) {
	var hosts []types.Host
	err := r.db.Where(&host).Find(&hosts).Error
	return hosts, err
}

func (r *PostgresHostRepository) Update(host types.Host) error {
	return r.db.Save(&host).Error
}

func (r *PostgresHostRepository) Delete(host types.Host) error {
	return r.db.Delete(&host).Error
}
