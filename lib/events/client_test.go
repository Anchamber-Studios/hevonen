package events_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/anchamber-studios/hevonen/lib/events"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
)

type ContainerInfos struct {
	Port string
}

func setupSuite(t *testing.T) (ContainerInfos, func(t *testing.T)) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not construct pool: %v", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		t.Fatalf("Could not connect to Docker: %s", err)
	}
	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "docker.redpanda.com/redpandadata/redpanda",
		Tag:        "v23.3.1",
		Cmd: []string{
			"redpanda",
			"start",
			"--smp 1",
			"--overprovisioned",
			"--memory 1G",
			"--reserve-memory 0M",
			"--node-id 0",
			"--check=false"},
	})
	if err != nil {
		t.Fatalf("Could not start redpanda container: %s", err)
	}
	port := container.GetPort("9092/tcp")
	return ContainerInfos{Port: port}, func(t *testing.T) {
		// err := container.Close()
		// if err != nil {
		// 	t.Errorf("Could not close container: %v", err)
		// }
		// err = pool.Purge(container)
		// if err != nil {
		// 	t.Errorf("Could not purge container: %v", err)
		// }
	}
}

func TestCreateAndDeleteTopics(t *testing.T) {
	container, tearDown := setupSuite(t)
	defer tearDown(t)
	testCases := []struct {
		name    string
		topics  []string
		success bool
	}{
		{
			name:    "With valid topic",
			topics:  []string{"test/" + uuid.New().String()},
			success: true,
		},
		{
			name:    "With empty topic",
			topics:  []string{""},
			success: false,
		},
		{
			name:    "With invalid topic",
			topics:  []string{"test!/" + uuid.New().String()},
			success: false,
		},
		{
			name:    "With multiple valid topics",
			topics:  []string{"test1/" + uuid.New().String(), "test2/" + uuid.New().String()},
			success: true,
		},
		{
			name:    "With multiple invalid topics",
			topics:  []string{"test1!/" + uuid.New().String(), "test2!/" + uuid.New().String()},
			success: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			producer, err := events.NewEventProducerRedpanda(fmt.Sprintf("localhost:%s", container.Port))
			if err != nil {
				t.Errorf("Failed to create producer: %v", err)
			}
			ctx := context.Background()
			for _, topic := range tc.topics {
				err = producer.CreateTopic(ctx, topic, events.TopicConfig{ReplicationFactor: 1, Partitions: 1})
				if err != nil && tc.success {
					t.Errorf("%s: CreateTopic() = %v, want nil", tc.name, err)
				} else if err == nil && !tc.success {
					t.Errorf("%s: CreateTopic() = nil, want error", tc.name)
				}
				if tc.success {
					exists, err := producer.TopicExists(ctx, topic)
					if err != nil && tc.success {
						t.Errorf("%s: CreateTopic() = %v, want nil", tc.name, err)
					} else if err == nil && !tc.success {
						t.Errorf("%s: CreateTopic() = nil, want error", tc.name)
					}
					if !exists {
						t.Errorf("%s: CreateTopic() = %v, want true", tc.name, exists)
					}
				}

			}

		})
	}
	// test
}

func TestDeleteTopics(t *testing.T) {
	_, tearDown := setupSuite(t)
	defer tearDown(t)
	// test
}

func TestSendAndReceive(t *testing.T) {
	_, tearDown := setupSuite(t)
	defer tearDown(t)
}
