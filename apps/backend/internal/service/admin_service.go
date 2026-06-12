package service

import (
	"taskmanager/internal/domain"
	"taskmanager/internal/dto"
	"taskmanager/internal/repository"
)

// AdminService exposes administrative read operations.
type AdminService struct {
	users repository.UserRepository
}

// NewAdminService constructs an AdminService.
func NewAdminService(users repository.UserRepository) *AdminService {
	return &AdminService{users: users}
}

// ListUsers returns a paginated list of all users.
func (s *AdminService) ListUsers(q dto.ActivityLogQuery) ([]domain.User, dto.PageMeta, error) {
	q.Normalize()
	users, total, err := s.users.List(q.Offset(), q.Limit)
	if err != nil {
		return nil, dto.PageMeta{}, err
	}
	return users, dto.NewPageMeta(total, q.Page, q.Limit), nil
}
