package serviceUseCase

import (
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	s "github.com/Rzhevskydd/techno-db-forum/project/app/units/service/repository"
)

type IServiceUseCase interface {
	Clear() error
	GetStatus() (*models.Status, error)
}

type ServiceUseCase struct {
	ServiceRep s.ServiceRepository
}

func (s *ServiceUseCase) Clear() error {
	return s.ServiceRep.Clear()
}

func (s *ServiceUseCase) GetStatus() (*models.Status, error) {
	return s.ServiceRep.GetStatus()
}


