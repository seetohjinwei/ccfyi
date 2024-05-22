package integration_tests

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"

	. "github.com/seetohjinwei/ccfyi/redis/internal/pkg/assert"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/server"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
)

func setup(t testing.TB) func() {
	store.ResetSingleton()

	// TODO: instead of a global singleton, let the store be associated with a server / router?
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

func TestIncrDecrIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration")
	}

	teardown := setup(t)
	defer teardown()

	cli := getClient()
	defer cli.Close()
	ctx := context.Background()

	c := cli.Incr(ctx, "k")
	NoError(t, c.Err())
	Equal(t, V(c.Result()), V(int64(1), nil))
	c = cli.Decr(ctx, "k")
	NoError(t, c.Err())
	Equal(t, V(c.Result()), V(int64(0), nil))

	cli.Set(ctx, "k", "notaninteger", 0)
	c = cli.Incr(ctx, "k")
	t.Logf("%v", c.Err())
	HasError(t, c.Err())
}

func TestListIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration")
	}

	teardown := setup(t)
	defer teardown()

	cli := getClient()
	defer cli.Close()
	ctx := context.Background()

	c := cli.LPush(ctx, "k", "c", "b", "a")
	NoError(t, c.Err())
	Equal(t, V(c.Result()), V(int64(3), nil))
	c = cli.RPush(ctx, "k", "d", "e", "f")
	NoError(t, c.Err())
	Equal(t, V(c.Result()), V(int64(6), nil))
	c = cli.LLen(ctx, "k")
	NoError(t, c.Err())
	Equal(t, V(c.Result()), V(int64(6), nil))
	r := cli.LRange(ctx, "k", 0, 5)
	NoError(t, r.Err())
	Equal(t, V(r.Result()), V([]string{"a", "b", "c", "d", "e", "f"}, nil))

	c = cli.LLen(ctx, "dontexist")
	NoError(t, c.Err())
	Equal(t, V(c.Result()), V(int64(0), nil))
	c = cli.Exists(ctx, "dontexist")
	NoError(t, c.Err())
	Equal(t, V(c.Result()), V(int64(0), nil))

	r = cli.LRange(ctx, "dontexist", 0, 0)
	NoError(t, r.Err())
	Equal(t, V(r.Result()), V([]string{}, nil))
	c = cli.Exists(ctx, "dontexist")
	NoError(t, c.Err())
	Equal(t, V(c.Result()), V(int64(0), nil))
}

func TestDelIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration")
	}

	teardown := setup(t)
	defer teardown()

	cli := getClient()
	defer cli.Close()
	ctx := context.Background()

	// attempt to data race
	go func() {
		for range 100 {
			cli.Set(ctx, "race", "me", 0)
		}
	}()

	cli.Set(ctx, "k", "1", 0)

	r := cli.Del(ctx, "k")
	NoError(t, r.Err())
	Equal(t, V(r.Result()), V(int64(1), nil))

	Equal(t, V(cli.Get(ctx, "k").Result()), V("", AnyError{}))
}
