package events

//go:generate go run github.com/dmarkham/enumer -type=DomainEventType -transform=kebab
type DomainEventType int

const (
	CompanyCreatedEventType DomainEventType = iota
	CompanyUpdatedEventType
	CompanyDeletedEventType
)
