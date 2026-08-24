package queue

import (
	"context"
	"testing"
	"time"
)

// En nil *Publisher ska vara helt säker att använda — redirect-handlern
// ska inte behöva bry sig om kön är nere.
func TestNilPublisher_IsSafe(t *testing.T) {
	var p *Publisher

	p.PublishClickEvent(context.Background(), ClickEvent{
		Code:        "abc1234",
		OriginalURL: "https://example.com",
		ClickedAt:   time.Now(),
	})
	p.Close()
}
