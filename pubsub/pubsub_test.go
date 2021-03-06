package pubsub

import (
	"sync"
	"testing"
	"time"

	"github.com/solomonwzs/goxutil/pubsub"
)

const (
	_N = 100
	_M = 100
)

func testPubsub(t testing.TB, timeout time.Duration) {
	var (
		wg           sync.WaitGroup
		readyForSend sync.WaitGroup
	)

	ch := pubsub.NewChannel(32)
	readyForSend.Add(_N)
	for i := 0; i < _N; i++ {
		wg.Add(1)
		go func() {
			sub := ch.NewSubscriber()
			readyForSend.Done()
			for j := 0; j < _M; j++ {
				sub.Recv(timeout)
			}
			wg.Done()
		}()
	}

	readyForSend.Wait()
	pub := ch.NewPublisher()
	for i := 0; i < _M; i++ {
		pub.Send(i*10, 0)
	}

	wg.Wait()
}

func testPubsub2(t testing.TB, timeout time.Duration) {
	var (
		wg           sync.WaitGroup
		readyForSend sync.WaitGroup
	)

	ch := _newChannel()
	readyForSend.Add(_N)
	for i := 0; i < _N; i++ {
		wg.Add(1)
		go func() {
			sub := ch.newSubscriber()
			readyForSend.Done()
			for j := 0; j < _M; j++ {
				sub.recv(timeout)
			}
			wg.Done()
		}()
	}

	readyForSend.Wait()
	pub := ch.newPublisher()
	for i := 0; i < _M; i++ {
		pub.send(i*10, 0)
	}

	wg.Wait()
}

func TestPubsub(t *testing.T) {
	testPubsub2(t, 0)
}

func BenchmarkPubsub(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testPubsub(b, 0)
	}
}

func BenchmarkPubsubWithTimeout(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testPubsub(b, 500*time.Millisecond)
	}
}

func BenchmarkPubsub2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testPubsub2(b, 0)
	}
}

func BenchmarkPubsub2WithTimeout(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testPubsub2(b, 500*time.Millisecond)
	}
}
