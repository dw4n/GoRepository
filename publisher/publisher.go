package publisher

type Publisher interface {
	PublishMessage(entity interface{}, action string, id string)
}

type NoopPublisher struct{}

func (n *NoopPublisher) PublishMessage(entity interface{}, action string, id string) {
	// Do nothing
}
