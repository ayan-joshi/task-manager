package server

import (
	"log"
	"strings"

	"gorm.io/gorm"

	"taskmanager/internal/config"
	"taskmanager/internal/domain"
	"taskmanager/internal/utils"
)

// SeedAdmin creates an ADMIN user from configuration if one is requested and
// does not already exist. It is idempotent: running it repeatedly is safe.
func SeedAdmin(db *gorm.DB, cfg *config.Config) error {
	if cfg.SeedAdminEmail == "" || cfg.SeedAdminPassword == "" {
		return nil
	}
	email := strings.ToLower(strings.TrimSpace(cfg.SeedAdminEmail))

	var count int64
	if err := db.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := utils.HashPassword(cfg.SeedAdminPassword)
	if err != nil {
		return err
	}
	admin := &domain.User{
		Name:     cfg.SeedAdminName,
		Email:    email,
		Password: hash,
		Role:     domain.RoleAdmin,
	}
	if err := db.Create(admin).Error; err != nil {
		return err
	}
	log.Printf("seeded admin user: %s", email)
	return nil
}
