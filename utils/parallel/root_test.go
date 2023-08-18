package parallel

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestZero(t *testing.T) {
	var group Group
	res := make(chan error)
	go func() {
		res <- group.Run()
	}()

	select {
	case err := <-res:
		if err != nil {
			t.Errorf("%v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Log(4)
		t.Error("timeout")
	}
}

func TestOne(t *testing.T) {
	myError := errors.New("foobar")
	var group Group
	group.Add(func() error { return myError }, func(error) {})
	res := make(chan error)
	go func() { res <- group.Run() }()
	select {
	case err := <-res:
		if want, have := myError, err; errors.Is(want, have) {
			t.Errorf("want %v, have %v", want, have)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout")
	}
}

func TestMany(t *testing.T) {
	interrupt := errors.New("interrupt")
	var group Group
	group.Add(func() error { return interrupt }, func(error) {})
	cancel := make(chan struct{})
	group.Add(func() error { <-cancel; return nil }, func(error) { close(cancel) })
	res := make(chan error)
	go func() { res <- group.Run() }()
	select {
	case err := <-res:
		if want, have := interrupt, err; errors.Is(err, want) {
			t.Errorf("want %v, have %v", want, have)
		}
	case <-time.After(100 * time.Millisecond):
		t.Errorf("timeout")
	}
}

func TestAddBasic(t *testing.T) {
	outputs := make([]string, 0)
	var group Group
	{
		cancel := make(chan struct{})
		group.Add(func() error {
			select {
			case <-time.After(time.Second):
				outputs = append(outputs, "The first actor had its time elapsed")
				return nil
			case <-cancel:
				outputs = append(outputs, "The first actor was canceled")
				return nil
			}
		}, func(err error) {
			outputs = append(outputs, fmt.Sprintf("The first actor was interrupted with: %v", err))
			close(cancel)
		})
	}
	{
		group.Add(func() error {
			outputs = append(outputs, "The second actor is returning immediately")
			return errors.New("immediate teardown")
		}, func(err error) {
			// Note that this interrupt function is called, even though the
			// corresponding execute function has already returned.
			outputs = append(outputs, fmt.Sprintf("The second actor was interrupted with: %v", err))
		})
	}

	ret := group.Run()
	outputs = append(outputs, fmt.Sprintf("The group was terminated with: %v", ret))

	assert.Equal(t, outputs, []string{
		"The second actor is returning immediately",
		"The first actor was interrupted with: immediate teardown",
		"The second actor was interrupted with: immediate teardown",
		"The first actor was canceled",
		"The group was terminated with: immediate teardown",
	})
}

func TestAddContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var group Group
	{
		ctx, cancel := context.WithCancel(ctx) // note: shadowed
		group.Add(func() error {
			return runUntilCanceled(ctx)
		}, func(error) {
			cancel()
		})
	}
	go cancel()
	fmt.Printf("The group was terminated with: %v\n", group.Run())
	// Output:
	// The group was terminated with: context canceled
}

func TestAddListener(t *testing.T) {
	var group Group
	{
		ln, _ := net.Listen("tcp", ":0")
		group.Add(func() error {
			defer fmt.Printf("http.Serve returned\n")
			return http.Serve(ln, http.NewServeMux())
		}, func(error) {
			_ = ln.Close()
		})
	}
	{
		group.Add(func() error {
			return errors.New("immediate teardown")
		}, func(error) {
			//
		})
	}
	fmt.Printf("The group was terminated with: %v\n", group.Run())
	// Output:
	// http.Serve returned
	// The group was terminated with: immediate teardown
}

func TestAddSignal(t *testing.T) {
	ctx, cancle := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second)
		cancle()
	}()
	var group Group
	group.Add(SignalActor(ctx, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT))
	t.Log("group exit with err: ", group.Run())
}

func runUntilCanceled(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}
