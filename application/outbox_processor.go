package application

import (
	"fmt"
	"log"
	"strings"
	"time"
	"transaction-outbox-practice/models"
	
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type OutboxProcessor struct {
	outboxEventRepo OutboxEventRepository
	kafkaProducer   *kafka.Producer
	kafkaTopic      string
	interval        time.Duration
	batchSize       int
	stopCh          chan struct{}
}

func NewOutboxProcessor(outboxEventRepo OutboxEventRepository, kafkaProducer *kafka.Producer, kafkaTopic string, interval time.Duration, batchSize int) *OutboxProcessor {
	return &OutboxProcessor{
		outboxEventRepo: outboxEventRepo,
		kafkaProducer:   kafkaProducer,
		kafkaTopic:      kafkaTopic,
		interval:        interval,
		batchSize:       batchSize,
		stopCh:          make(chan struct{}),
	}
}

func (p *OutboxProcessor) Start() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	log.Println("Outbox processor started")

	for {
		select {
		case <-ticker.C:
			p.processOutboxEvents()
		case <-p.stopCh:
			log.Println("Outbox processor stopped")
			return
		}
	}
}

func (p *OutboxProcessor) Stop() {
	close(p.stopCh)
	
	// Kafka producer 정리
	if p.kafkaProducer != nil {
		// 남은 메시지 전송 대기
		p.kafkaProducer.Flush(5000)
		p.kafkaProducer.Close()
	}
}

func (p *OutboxProcessor) processOutboxEvents() {
	events, err := p.outboxEventRepo.FindPendingEvents(p.batchSize)
	if err != nil {
		log.Printf("Error fetching outbox events: %v", err)
		return
	}

	for _, event := range events {
		if err := p.processEvent(&event); err != nil {
			log.Printf("Error processing event %d: %v", event.ID, err)
			continue
		}
	}
}

func (p *OutboxProcessor) processEvent(event *models.OutboxEvent) error {
	log.Printf("Processing event: ID=%d, Type=%s, AggregateID=%s",
		event.ID, event.EventType, event.AggregateID)

	// Kafka로 이벤트 전송
	topic := p.kafkaTopic
	// 이벤트 타입에 따라 다른 토픽 사용 가능
	if strings.HasPrefix(event.EventType, "Order") {
		topic = "order-events"
	}

	deliveryChan := make(chan kafka.Event, 1)
	err := p.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(event.AggregateID),
		Value: []byte(event.Payload),
		Headers: []kafka.Header{
			{Key: "event_id", Value: []byte(fmt.Sprintf("%d", event.ID))},
			{Key: "event_type", Value: []byte(event.EventType)},
			{Key: "timestamp", Value: []byte(event.CreatedAt.Format(time.RFC3339))},
		},
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to produce message to Kafka: %w", err)
	}

	// 전송 결과 대기
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
	}

	log.Printf("Event delivered to topic %s [%d] at offset %v",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

	// 성공적으로 전송되면 processed로 표시
	return p.outboxEventRepo.MarkAsProcessed(event)
}
