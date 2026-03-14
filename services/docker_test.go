package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitProgress_OvershootRequeue(t *testing.T) {
	dock := &Docker{
		progChan: make(chan int, 5),
	}

	ctx := context.Background()

	// Simulate receiving a 100 progress directly (e.g. from a snapshot)
	dock.progChan <- 100

	done := make(chan struct{})
	go func() {
		// We expect WaitProgress for 50 to return immediately, and put 100 back in the channel
		err := dock.WaitProgress(ctx, 50)
		assert.NoError(t, err)

		// WaitProgress for 80 should also return immediately and put 100 back in the channel
		err = dock.WaitProgress(ctx, 80)
		assert.NoError(t, err)

		// WaitProgress for 100 should return immediately
		err = dock.WaitProgress(ctx, 100)
		assert.NoError(t, err)

		close(done)
	}()

	select {
	case <-done:
		// Test finished successfully
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Test timed out: WaitProgress blocked indefinitely instead of returning instantly. You might be missing the re-queuing fix.")
	}
}

func TestWaitProgress_NormalProgression(t *testing.T) {
	dock := &Docker{
		progChan: make(chan int, 5),
	}

	ctx := context.Background()

	go func() {
		time.Sleep(10 * time.Millisecond)
		dock.progChan <- 50
		time.Sleep(10 * time.Millisecond)
		dock.progChan <- 80
		time.Sleep(10 * time.Millisecond)
		dock.progChan <- 100
	}()

	err := dock.WaitProgress(ctx, 50)
	assert.NoError(t, err)

	err = dock.WaitProgress(ctx, 80)
	assert.NoError(t, err)

	err = dock.WaitProgress(ctx, 100)
	assert.NoError(t, err)
}

func TestWaitProgress_ChannelClosed(t *testing.T) {
	dock := &Docker{
		progChan: make(chan int, 5),
	}

	ctx := context.Background()

	dock.progChan <- 10
	close(dock.progChan)

	err := dock.WaitProgress(ctx, 50)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Channel is closed without reaching")
}
