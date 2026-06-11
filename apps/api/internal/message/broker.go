package message

type Broker interface {
	Publish(message Message) error
}
