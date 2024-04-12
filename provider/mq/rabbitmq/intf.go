package rabbitmq

//
//import (
//	"context"
//	amqp "github.com/rabbitmq/amqp091-go"
//)
//
//type rabbitMQChannelBroker interface {
//	PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
//	PublishWithDeferredConfirmWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) (*amqp.DeferredConfirmation, error)
//	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
//	QueueBind(name string, key string, exchange string, noWait bool, args amqp.Table) error
//	Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
//	Nack(tag uint64, multiple bool, requeue bool) error
//	Ack(tag uint64, multiple bool) error
//	ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table) error
//	Qos(prefetchCount, prefetchSize int, global bool) error
//	Confirm(noWait bool) error
//	Close() error
//	IsClosed() bool
//}
//
//type rabbitMQConnectionBroker interface {
//	Close() error
//}
