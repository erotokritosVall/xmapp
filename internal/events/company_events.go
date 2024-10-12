package events

import (
	"fmt"
	"time"
)

type CompanyCreated struct {
	CompanyId string
}

func NewCompanyCreated(companyId string) DomainEvent {
	return &CompanyCreated{
		CompanyId: companyId,
	}
}

func (e *CompanyCreated) Type() DomainEventType {
	return CompanyCreatedEventType
}
func (e *CompanyCreated) UniqueId() string {
	return fmt.Sprintf("%s:%s", e.Type().String(), e.CompanyId)
}

type CompanyUpdated struct {
	CompanyId string
	Timestamp int64
}

func NewCompanyUpdated(companyId string) DomainEvent {
	return &CompanyUpdated{
		CompanyId: companyId,
		Timestamp: time.Now().UTC().Unix(),
	}
}

func (e *CompanyUpdated) Type() DomainEventType {
	return CompanyUpdatedEventType
}
func (e *CompanyUpdated) UniqueId() string {
	return fmt.Sprintf("%s:%s:%d", e.Type().String(), e.CompanyId, e.Timestamp)
}

type CompanyDeleted struct {
	CompanyId string
}

func NewCompanyDeleted(companyId string) DomainEvent {
	return &CompanyDeleted{
		CompanyId: companyId,
	}
}

func (e *CompanyDeleted) Type() DomainEventType {
	return CompanyDeletedEventType
}

func (e *CompanyDeleted) UniqueId() string {
	return fmt.Sprintf("%s:%s", e.Type().String(), e.CompanyId)
}
