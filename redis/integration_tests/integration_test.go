package integration_tests

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/server"
)

func getSetHelper(t *testing.T, index int) {
	cli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

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

	router, err := server.New("localhost:6379")
	if err != nil {
		t.Errorf("error init server: %v", err)
	}
	go func() {
		err := router.Serve()
		if err != nil {
			t.Errorf("error serve error: %v", err)
		}
	}()

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			getSetHelper(t, i)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
