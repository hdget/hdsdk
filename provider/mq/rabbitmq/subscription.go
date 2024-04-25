package rabbitmq

import (
	"context"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/provider/mq"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type subscription struct {
	out                chan *mq.Message
	logger             intf.LoggerProvider
	notifyCloseChannel chan *amqp.Error
	channel            *amqp.Channel
	queueName          string
	closing            chan struct{}
	closedChan         chan struct{}
	config             *RabbitMqConfig
}

func (s *subscription) createConsumer(queueName string, amqpChannel *amqp.Channel) (<-chan amqp.Delivery, error) {
	amqpMsgs, err := amqpChannel.Consume(
		queueName,
		"",
		false, // autoAck must be set to false - acks are managed by sdk
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "cannot consume from amqpChannel")
	}

	return amqpMsgs, nil
}

func (s *subscription) ProcessMessages(ctx context.Context) {
	amqpMsgs, err := s.createConsumer(s.queueName, s.channel)
	if err != nil {
		s.logger.Error("start consuming messages", "err", err)
		return
	}

ConsumingLoop:
	for {
		select {
		case amqpMsg := <-amqpMsgs:
			if err = s.processMessage(ctx, amqpMsg, s.out); err != nil {
				s.logger.Error("processing message failed, sending nack", "err", err)
				if err = s.nackMsg(amqpMsg); err != nil {
					s.logger.Error("cannot nack message", "err", err)
					// something went really wrong when we cannot nack, let's reconnect
					break ConsumingLoop
				}
			}
			continue ConsumingLoop

		case <-s.notifyCloseChannel:
			s.logger.Error("channel closed, stopping ProcessMessages")
			break ConsumingLoop

		case <-s.closing:
			s.logger.Info("closing from subscriber received")
			break ConsumingLoop

		case <-s.closedChan:
			s.logger.Info("subscriber closed")
			break ConsumingLoop

		case <-ctx.Done():
			s.logger.Info("closing from ctx received")
			break ConsumingLoop
		}
	}
}

func (s *subscription) processMessage(ctx context.Context, amqpMsg amqp.Delivery, out chan *mq.Message) error {
	msg := mq.NewMessage(amqpMsg.Body)

	ctx, cancelCtx := context.WithCancel(ctx)
	msg.SetContext(ctx)
	defer cancelCtx()

	select {
	case <-s.closing:
		s.logger.Info("message not consumed, pub/sub is closing")
		return s.nackMsg(amqpMsg)
	case <-s.closedChan:
		s.logger.Info("message not consumed, subscriber is closed")
		return s.nackMsg(amqpMsg)
	case out <- msg:
		s.logger.Trace("message sent to consumer")
	}

	select {
	case <-s.closing:
		s.logger.Trace("closing pub/sub, message discarded before ack")
		return s.nackMsg(amqpMsg)
	case <-s.closedChan:
		s.logger.Info("message not consumed, subscriber is closed")
		return s.nackMsg(amqpMsg)
	case <-msg.Acked():
		s.logger.Trace("message acked")
		return amqpMsg.Ack(false)
	case <-msg.Nacked():
		s.logger.Trace("message nacked")
		return s.nackMsg(amqpMsg)
	}
}

func (s *subscription) nackMsg(amqpMsg amqp.Delivery) error {
	return amqpMsg.Nack(false, s.config.RequeueInFailure)
}
