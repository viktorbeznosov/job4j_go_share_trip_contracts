package events

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubReader struct {
	messages  []kafka.Message
	commits   []kafka.Message
	fetchErr  error
	commitErr error
	idx       int
}

func (s *stubReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	if s.fetchErr != nil {
		return kafka.Message{}, s.fetchErr
	}
	if s.idx >= len(s.messages) {
		<-ctx.Done()
		return kafka.Message{}, ctx.Err()
	}
	msg := s.messages[s.idx]
	s.idx++
	return msg, nil
}

func (s *stubReader) CommitMessages(_ context.Context, msgs ...kafka.Message) error {
	if s.commitErr != nil {
		return s.commitErr
	}
	s.commits = append(s.commits, msgs...)
	return nil
}

func (s *stubReader) Close() error {
	return nil
}

func TestConsumer_HandlesTripCreated(t *testing.T) {
	t.Parallel()
	event := TripEvent{
		EventID:   "evt-1",
		EventType: EventTypeTripCreated,
		TripID:    "trip-1",
		DriverID:  "driver-1",
	}
	body, err := json.Marshal(event)
	require.NoError(t, err)

	reader := &stubReader{messages: []kafka.Message{{Value: body}}}
	consumer := &Consumer{reader: reader}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	handled := 0
	err = consumer.Run(ctx, func(_ context.Context, got TripEvent) error {
		handled++
		assert.Equal(t, event.EventID, got.EventID)
		assert.Equal(t, EventTypeTripCreated, got.EventType)
		return nil
	})

	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, 1, handled)
	require.Len(t, reader.commits, 1)
}

func TestConsumer_HandlesTripPublished(t *testing.T) {
	t.Parallel()
	event := TripEvent{
		EventID:   "evt-published",
		EventType: EventTypeTripPublished,
		TripID:    "trip-published",
		DriverID:  "driver-published",
	}
	body, err := json.Marshal(event)
	require.NoError(t, err)

	reader := &stubReader{messages: []kafka.Message{{Value: body}}}
	consumer := &Consumer{reader: reader}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	handled := 0
	err = consumer.Run(ctx, func(_ context.Context, got TripEvent) error {
		handled++
		assert.Equal(t, EventTypeTripPublished, got.EventType)
		return nil
	})

	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, 1, handled)
	require.Len(t, reader.commits, 1)
}

func TestConsumer_HandlesTripCompleted(t *testing.T) {
	t.Parallel()
	event := TripEvent{
		EventID:   "evt-completed",
		EventType: EventTypeTripCompleted,
		DriverID:  "driver-2",
	}
	body, err := json.Marshal(event)
	require.NoError(t, err)

	reader := &stubReader{messages: []kafka.Message{{Value: body}}}
	consumer := &Consumer{reader: reader}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	handled := 0
	err = consumer.Run(ctx, func(_ context.Context, got TripEvent) error {
		handled++
		assert.Equal(t, EventTypeTripCompleted, got.EventType)
		return nil
	})

	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, 1, handled)
	require.Len(t, reader.commits, 1)
}

func TestConsumer_HandlesTripCanceled(t *testing.T) {
	t.Parallel()
	event := TripEvent{
		EventID:   "evt-canceled",
		EventType: EventTypeTripCanceled,
		DriverID:  "driver-canceled",
	}
	body, err := json.Marshal(event)
	require.NoError(t, err)

	reader := &stubReader{messages: []kafka.Message{{Value: body}}}
	consumer := &Consumer{reader: reader}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	handled := 0
	err = consumer.Run(ctx, func(_ context.Context, got TripEvent) error {
		handled++
		assert.Equal(t, EventTypeTripCanceled, got.EventType)
		return nil
	})

	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, 1, handled)
	require.Len(t, reader.commits, 1)
}

func TestConsumer_SkipsUnsupportedEventType(t *testing.T) {
	t.Parallel()
	event := TripEvent{
		EventID:   "evt-skip",
		EventType: EventType("trip_preapproved"),
	}
	body, err := json.Marshal(event)
	require.NoError(t, err)

	reader := &stubReader{messages: []kafka.Message{{Value: body}}}
	consumer := &Consumer{reader: reader}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	handled := 0
	err = consumer.Run(ctx, func(_ context.Context, _ TripEvent) error {
		handled++
		return nil
	})

	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, 0, handled)
	require.Len(t, reader.commits, 1)
}

func TestConsumer_DoesNotCommitWhenHandlerFails(t *testing.T) {
	t.Parallel()
	event := TripEvent{
		EventID:   "evt-fail",
		EventType: EventTypeTripCreated,
		DriverID:  "driver-3",
	}
	body, err := json.Marshal(event)
	require.NoError(t, err)

	reader := &stubReader{messages: []kafka.Message{{Value: body}}}
	consumer := &Consumer{reader: reader}

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	err = consumer.Run(ctx, func(_ context.Context, _ TripEvent) error {
		return errors.New("db unavailable")
	})

	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Empty(t, reader.commits)
}

func TestNewConsumerAndClose(t *testing.T) {
	t.Parallel()
	consumer := NewConsumer([]string{"localhost:9092"}, "trip.events", "test-group")
	require.NotNil(t, consumer)
	require.NoError(t, consumer.Close())
}
