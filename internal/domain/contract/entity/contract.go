package entity

import (
	"time"

	"github.com/google/uuid"

	contractErrors "job4j_go_share_trip_contracts/internal/api/errors"
)

type ContractStatus string

const (
	StatusDraft     ContractStatus = "draft"
	StatusActive    ContractStatus = "active"
	StatusSuspended ContractStatus = "suspended"
	StatusTerminated ContractStatus = "terminated"
)

type ServiceType string

const (
	ServiceTripCreation    ServiceType = "trip_creation"
	ServiceTripPublish    ServiceType = "trip_publish"
	ServiceTripStart    ServiceType = "trip_start"
	ServiceTripParticipants ServiceType = "trip_participants"
	ServiceNotifications   ServiceType = "notifications"
	ServicePremiumSupport  ServiceType = "premium_support"
)

type ContractService struct {
	Service ServiceType `json:"service"`
	Enabled bool        `json:"enabled"`
}

type Contract struct {
	ID          uuid.UUID          `json:"id"`
	CompanyID   uuid.UUID          `json:"companyId"`
	Status      ContractStatus     `json:"status"`
	ValidFrom   time.Time          `json:"validFrom"`
	ValidTo     time.Time          `json:"validTo"`
	Services    []ContractService  `json:"services"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

func NewContract(companyID uuid.UUID, validFrom, validTo time.Time) *Contract {
	return &Contract{
		ID:        uuid.New(),
		CompanyID: companyID,
		Status:    StatusDraft,
		ValidFrom: validFrom,
		ValidTo:   validTo,
		Services:  []ContractService{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (c *Contract) ChangeStatus(newStatus ContractStatus) error {
	validStatuses := map[ContractStatus]bool{
		StatusDraft:     true,
		StatusActive:    true,
		StatusSuspended: true,
		StatusTerminated: true,
	}
	if !validStatuses[newStatus] {
		return contractErrors.ErrInvalidStatus
	}

	switch c.Status {
	case StatusDraft:
		if newStatus != StatusActive && newStatus != StatusTerminated {
			return contractErrors.ErrInvalidStatusTransition
		}
	case StatusActive:
		if newStatus != StatusSuspended && newStatus != StatusTerminated {
			return contractErrors.ErrInvalidStatusTransition
		}
	case StatusSuspended:
		if newStatus != StatusActive && newStatus != StatusTerminated {
			return contractErrors.ErrInvalidStatusTransition
		}
	case StatusTerminated:
		return contractErrors.ErrCannotUpdateTerminated
	default:
		return contractErrors.ErrUnknownStatus
	}

	c.Status = newStatus
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Contract) UpdateServices(services []ContractService) error {
	if c.Status == StatusTerminated {
		return contractErrors.ErrCannotUpdateTerminated
	}

	if len(services) == 0 {
		return contractErrors.ErrServicesListEmpty
	}

	seen := make(map[ServiceType]bool)
	for _, s := range services {
		if seen[s.Service] {
			return contractErrors.ErrDuplicateService
		}
		seen[s.Service] = true
	}

	c.Services = services
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Contract) IsActive() bool {
	return c.Status == StatusActive
}

func (c *Contract) IsExpired() bool {
	return time.Now().After(c.ValidTo)
}

func (c *Contract) IsServiceEnabled(serviceType ServiceType) bool {
	for _, s := range c.Services {
		if s.Service == serviceType {
			return s.Enabled
		}
	}
	return false
}