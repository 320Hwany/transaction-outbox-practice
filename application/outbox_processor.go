package application

import (
	"fmt"
	"log"
	"time"
	"transaction-outbox-practice/models"
)

type OutboxProcessor struct {
	outboxEventRepo OutboxEventRepository
	interval        time.Duration
	batchSize       int
	stopCh          chan struct{}
}

func NewOutboxProcessor(outboxEventRepo OutboxEventRepository, interval time.Duration, batchSize int) *OutboxProcessor {
	return &OutboxProcessor{
		outboxEventRepo: outboxEventRepo,
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

	fmt.Printf("Publishing event to message broker: %s\n", event.Payload)

	return p.outboxEventRepo.MarkAsProcessed(event)
}
