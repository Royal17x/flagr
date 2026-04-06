package service

import (
	"github.com/Royal17x/backend/internal/domain"
	"github.com/Royal17x/flagr/backend/internal/domain"
)

type FlagService struct {
	flags    domain.FlagRepository
	projects domain.ProjectRepository
}

func NewFlagService(flags domain.FlagRepository, projects domain.ProjectRepository)
