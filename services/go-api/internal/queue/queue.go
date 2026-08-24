package queue

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	connectRetries = 10
	connectDelay   = 2 * time.Second
)

const clickEventsQueue = "click_events"

type ClickEvent struct {
	Code        string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	ClickedAt   time.Time `json:"clicked_at"`
}

// Publisher skickar klick-events till RabbitMQ. En nil *Publisher är giltig
// och gör alla operationer till no-ops — redirect-endpointen ska fungera
// även om kön är nere, analytics är best-effort.
type Publisher struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger *slog.Logger
}

// Connect ansluter till RabbitMQ och deklarerar kön för klick-events, med
// några återförsök — RabbitMQ:s Docker-healthcheck kan rapportera "healthy"
// strax innan AMQP-porten faktiskt tar emot anslutningar. Om anslutningen
// ändå misslyckas loggas ett fel och nil returneras istället för att stoppa
// hela API:et.
func Connect(url string, logger *slog.Logger) *Publisher {
	var lastErr error
	for attempt := 1; attempt <= connectRetries; attempt++ {
		conn, err := amqp.Dial(url)
		if err != nil {
			lastErr = err
			logger.Warn("rabbitmq not ready yet, retrying", "attempt", attempt, "error", err)
			time.Sleep(connectDelay)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			conn.Close()
			lastErr = err
			logger.Warn("failed to open rabbitmq channel, retrying", "attempt", attempt, "error", err)
			time.Sleep(connectDelay)
			continue
		}

		if _, err := ch.QueueDeclare(clickEventsQueue, true, false, false, false, nil); err != nil {
			ch.Close()
			conn.Close()
			lastErr = err
			logger.Warn("failed to declare click_events queue, retrying", "attempt", attempt, "error", err)
			time.Sleep(connectDelay)
			continue
		}

		return &Publisher{conn: conn, ch: ch, logger: logger}
	}

	logger.Error("failed to connect to rabbitmq after retries, click events will not be published", "error", lastErr)
	return nil
}

func (p *Publisher) Close() {
	if p == nil {
		return
	}
	p.ch.Close()
	p.conn.Close()
}

// PublishClickEvent är best-effort: fel loggas men returneras aldrig, så
// anroparen (redirect-handlern) aldrig behöver hantera köfel.
func (p *Publisher) PublishClickEvent(ctx context.Context, evt ClickEvent) {
	if p == nil {
		return
	}

	body, err := json.Marshal(evt)
	if err != nil {
		p.logger.Error("failed to marshal click event", "error", err)
		return
	}

	err = p.ch.PublishWithContext(ctx, "", clickEventsQueue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		p.logger.Error("failed to publish click event", "error", err)
	}
}
