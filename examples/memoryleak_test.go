package examples

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

// Example 1: No sender for channel, gorutine waits indefinitely causing leak
func TestLeakWithChannelWithNoSender(t *testing.T) {
	ch1 := make(chan string)

	go func() {
		val := <-ch1
		fmt.Println("We received a value:", val)
	}()
}


// Example 2: context done before search complete
func TestGoroutineLeak(t *testing.T) {

	channel := make(chan string)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	// timer := time.NewTimer(time.Duration(50) * time.Millisecond) // option to context with timeout

	before := runtime.NumGoroutine()
	defer func() {
		cancel()
		// timer.Stop()
		after := runtime.NumGoroutine()
		if after == before {
			t.Fatalf("Failed to identify goroutine leak due to early return from function")
		}
		t.Logf("Goroutine leak. Earlier routines: %d, after routines: %d", before, after)
	}()

	go func() {
		record, _ := search("something")
		channel <- record // waits indefinitely when context is done before search is complete
	}()

	select {
	// case <-timer.C:
	// 	t.Fatalf("error: timer expired")
	// 	return
	case <-ctx.Done():
		t.Log("error: search canceled")
		return
	case result := <-channel:
		t.Logf("Received: %s", result)
	}
}

func search(term string) (string, error) {
	time.Sleep(200 * time.Millisecond)
	return "some value", nil
}
