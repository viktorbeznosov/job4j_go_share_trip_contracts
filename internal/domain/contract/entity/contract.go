package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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
		return fmt.Errorf("invalid status: %s", newStatus)
	}

	// Бизнес-правила для переходов между статусами
	switch c.Status {
	case StatusDraft:
		// Из черновика можно перейти в active или terminated
		if newStatus != StatusActive && newStatus != StatusTerminated {
			return fmt.Errorf("cannot change status from %s to %s. Allowed transitions: active, terminated", c.Status, newStatus)
		}
	case StatusActive:
		// Из активного можно перейти в suspended или terminated
		if newStatus != StatusSuspended && newStatus != StatusTerminated {
			return fmt.Errorf("cannot change status from %s to %s. Allowed transitions: suspended, terminated", c.Status, newStatus)
		}
	case StatusSuspended:
		// Из приостановленного можно перейти в active или terminated
		if newStatus != StatusActive && newStatus != StatusTerminated {
			return fmt.Errorf("cannot change status from %s to %s. Allowed transitions: active, terminated", c.Status, newStatus)
		}
	case StatusTerminated:
		// Из завершенного нельзя перейти никуда
		return fmt.Errorf("cannot change status from terminated to %s", newStatus)
	default:
		return fmt.Errorf("unknown status: %s", c.Status)
	}

	c.Status = newStatus
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Contract) UpdateServices(services []ContractService) error {
	if c.Status == StatusTerminated {
		return fmt.Errorf("cannot update services for terminated contract")
	}

	if len(services) == 0 {
		return fmt.Errorf("services list cannot be empty")
	}

	seen := make(map[ServiceType]bool)
	for _, s := range services {
		if seen[s.Service] {
			return fmt.Errorf("duplicate service: %s", s.Service)
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