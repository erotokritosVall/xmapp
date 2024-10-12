package events

type DomainEvent interface {
	Type() DomainEventType
	UniqueId() string
}
