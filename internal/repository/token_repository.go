package repository

import (
	"time"

	"handler/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenRepository interface {
	Revoke(token *entity.RevokedToken) error
	IsRevoked(jti string) (bool, error)
	DeleteExpired() error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Revoke(token *entity.RevokedToken) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(token).Error
}

func (r *tokenRepository) IsRevoked(jti string) (bool, error) {
	if jti == "" {
		return false, nil
	}

	var count int64
	err := r.db.Model(&entity.RevokedToken{}).
		Where("jti = ? AND expires_at > ?", jti, time.Now()).
		Count(&count).Error

	return count > 0, err
}

func (r *tokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&entity.RevokedToken{}).Error
}
