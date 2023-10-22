package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func foo() (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	/*
			WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).

		Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this [Context] complete:

		func slowOperationWithTimeout(ctx context.Context) (Result, error) {
		    ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		    defer cancel()  // releases resources if slowOperation completes before timeout elapses
		    return slowOperation(ctx)
		}
	*/
	ctx.Done()
	defer cancel()
	/*
		defer cancel() гарантирует, что после выхода из функции или
		горутины контекст будёт отменён, и таким образом вы избежите
		утекания горутины — явления, когда горутина продолжает выполняться
		и существовать в памяти, но результат её работы больше никого не интересует.
	*/

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request with ctx: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform http request: %w", err)
	}

	return res, nil

}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doWork(ctx)
}

func doWork(ctx context.Context) {
	newCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	log.Println("starting working...")

	for {
		select {
		case <-newCtx.Done():
			log.Printf("ctx done: %v", ctx.Err())
			return
		default:
			log.Println("working...")
			time.Sleep(1 * time.Second)
		}
	}
}
