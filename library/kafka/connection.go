package kafka

import (
	"fmt"
	"users-service/library/kafka/notifications"

	"github.com/IBM/sarama"
)

func Connect(brokerList []string, KafkaUsername, KafkaPassword string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	// Only enable SASL if username and password are provided
	if KafkaUsername != "" && KafkaPassword != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = KafkaUsername
		config.Net.SASL.Password = KafkaPassword
	} else {
		// Disable SASL for local development
		config.Net.SASL.Enable = false
	}

	// Create notification producer
	NotificationProducer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	notifications.Producer = NotificationProducer

	// Start consumers in background
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	return nil
}
