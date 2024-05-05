package integration_tests

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/server"
)

func setup(t testing.TB) func() {
	router, err := server.New("localhost:6379")
	if err != nil {
		t.Errorf("error init server: %v", err)
	}
	go func() {
		// ignore errors
		router.Serve()
	}()

	return func() {
		router.Stop()
	}
}

func getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func getSetHelper(t *testing.T, index int) {
	cli := getClient()
	defer cli.Close()

	ctx := context.Background()

	key := "foo" + strconv.Itoa(index)
	err := cli.Set(ctx, key, "bar", 0).Err()
	if err != nil {
		t.Errorf("expected no err, but got %+v", err)
	}

	val, err := cli.Get(ctx, key).Result()
	if err != nil {
		t.Errorf("expected no err, but got %+v", err)
	}
	if val != "bar" {
		t.Errorf("expected %q, but got %q", "bar", val)
	}
}

func TestGetSetIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration")
	}

	teardown := setup(t)
	defer teardown()

	var wg sync.WaitGroup

	const clients = 50
	for i := 0; i < clients; i++ {
		wg.Add(1)
		go func(i int) {
			getSetHelper(t, i)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestPingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration")
	}

	teardown := setup(t)
	defer teardown()

	cli := getClient()
	defer cli.Close()

	status := cli.Ping(context.Background())
	v, err := status.Result()

	if err != nil {
		t.Errorf("expected no err, but got %v", err)
	}
	if v != "PONG" {
		t.Errorf("expected %q, but got %q", "PONG", v)
	}
}

func TestExistsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration")
	}

	teardown := setup(t)
	defer teardown()

	cli := getClient()
	defer cli.Close()
	ctx := context.Background()

	c := cli.Exists(ctx, "k")
	v, err := c.Result()
	if err != nil {
		t.Errorf("expected no err, but got %v", err)
	}
	if v != 0 {
		t.Errorf("expected %v, but got %v", 0, v)
	}

	cli.Set(ctx, "k", "v", 0)
	c = cli.Exists(ctx, "k")
	v, err = c.Result()
	if err != nil {
		t.Errorf("expected no err, but got %v", err)
	}
	if v != 1 {
		t.Errorf("expected %v, but got %v", 1, v)
	}
	c = cli.Exists(ctx, "k", "k")
	v, err = c.Result()
	if err != nil {
		t.Errorf("expected no err, but got %v", err)
	}
	if v != 2 {
		t.Errorf("expected %v, but got %v", 2, v)
	}
}

func TestDelIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration")
	}

	teardown := setup(t)
	defer teardown()

	cli := getClient()
	defer cli.Close()
	// ctx := context.Background()

	// TODO:
}
