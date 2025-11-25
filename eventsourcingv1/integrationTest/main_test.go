package integrationtest

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ooqls/go-db/sqlx"
	"github.com/ooqls/go-db/testutils"
	"github.com/ooqls/go-eventsource/eventsourcingv1"
	"github.com/ooqls/go-eventsource/eventsourcingv1/tablesv1"
)

type TestEntity struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (e *TestEntity) GetId() uuid.UUID {
	return e.Id
}

func (e *TestEntity) Apply(event eventsourcingv1.Event) error {
	// Apply the event to the entity
	switch event.Key {
	case "name":
		e.Name = event.Value["name"].(string)
	case "id":
		e.Id = event.Value["id"].(uuid.UUID)
	}

	return nil
}
func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cont := testutils.StartPostgres(ctx)

	sqlx.SeedSQLX((tablesv1.GetCreateTableStmts(eventsourcingv1.EventSource("test"))), []string{})

	timeout := time.Second * 30
	defer cont.Stop(ctx, &timeout)

	redisCont := testutils.StartRedis(ctx)
	defer redisCont.Stop(context.Background(), &timeout)

	m.Run()
}
