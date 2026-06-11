package message

import "context"

type Broker interface {
	Publish(ctx context.Context, message Message) error
}
