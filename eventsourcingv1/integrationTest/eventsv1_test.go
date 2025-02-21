package integrationtest

import (
	"context"
	"testing"

	"github.com/google/uuid"
	db "github.com/ooqls/go-db/postgres"
	"github.com/ooqls/go-eventsource/eventsourcingv1"
	"github.com/ooqls/go-eventsource/eventsourcingv1/eventsv1"
	"github.com/stretchr/testify/assert"
)

func TestEventWriter(t *testing.T) {
	sqldb := db.Get()
	store := eventsv1.NewSQLWriter(sqldb, eventsourcingv1.EventSource("test"))
	obj := &TestEntity{Id: uuid.New(), Name: "test"}

	err := store.Append(context.Background(), eventsourcingv1.Event{
		EntityId: obj.GetId(),
		Key:      "name",
		Value:    map[string]interface{}{"name": "test1"},
	}, eventsourcingv1.Event{
		EntityId: obj.GetId(),
		Key:      "name",
		Value:    map[string]interface{}{"name": "test2"},
	},
	)
	assert.Nilf(t, err, "Append should not return an error: %v", err)
}

func TestEventReader(t *testing.T) {
	sqldb := db.Get()
	ent := TestEntity{}
	source := eventsourcingv1.EventSource("test")
	reader := eventsv1.NewSQLReader(sqldb, source)
	writer := eventsv1.NewSQLWriter(sqldb, source)

	obj := &TestEntity{Id: uuid.New(), Name: "test"}

	err := writer.Append(context.Background(), eventsourcingv1.Event{
		EntityId: obj.GetId(),
		Key:      "name",
		Value:    map[string]interface{}{"name": "test1"},
	}, eventsourcingv1.Event{
		EntityId: obj.GetId(),
		Key:      "name",
		Value:    map[string]interface{}{"name": "test2"},
	})
	assert.Nilf(t, err, "Append should not return an error: %v", err)
	next, err := reader.Get(context.Background(), obj.GetId())
	assert.Nilf(t, err, "Get should not return an error: %v", err)
	assert.NotNilf(t, next, "Get should return an iterator")

	var ev *eventsourcingv1.Event
	ev = next()
	for ev != nil {
		ent.Apply(*ev)
		ev = next()
	}

	assert.Equalf(t, "test2", ent.Name, "events should be in order")
}
