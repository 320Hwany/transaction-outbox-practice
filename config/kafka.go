package config

import (
	"fmt"
	"strings"
	
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func InitKafkaProducer(cfg *KafkaConfig) (*kafka.Producer, error) {
	// Kafka producer 설정
	configMap := kafka.ConfigMap{
		"bootstrap.servers": strings.Join(cfg.Brokers, ","),
	}
	
	// 추가 설정 적용
	for key, value := range cfg.ProducerConfig {
		configMap[key] = value
	}
	
	producer, err := kafka.NewProducer(&configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	
	// 비동기 에러 처리를 위한 고루틴
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	
	return producer, nil
}