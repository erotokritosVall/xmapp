package events

//go:generate go run github.com/dmarkham/enumer -type=DomainEventType -transform=kebab -trimprefix=EventType
type DomainEventType int

const (
	EventTypeCompanyCreated DomainEventType = iota
	EventTypeCompanyUpdated
	EventTypeCompanyDeleted
)
