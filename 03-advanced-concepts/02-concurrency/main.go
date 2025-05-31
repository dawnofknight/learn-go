package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Go Concurrency Examples ===")

	// Basic goroutine
	fmt.Println("\n1. Basic Goroutine:")
	basicGoroutineExample()

	// Multiple goroutines
	fmt.Println("\n2. Multiple Goroutines:")
	multipleGoroutinesExample()

	// Channels
	fmt.Println("\n3. Basic Channel:")
	basicChannelExample()

	// Buffered channels
	fmt.Println("\n4. Buffered Channel:")
	bufferedChannelExample()

	// Channel directions
	fmt.Println("\n5. Channel Directions:")
	channelDirectionsExample()

	// Select statement
	fmt.Println("\n6. Select Statement:")
	selectExample()

	// Mutex for synchronization
	fmt.Println("\n7. Mutex for Synchronization:")
	mutexExample()

	// WaitGroup
	fmt.Println("\n8. WaitGroup:")
	waitGroupExample()

	// Worker pool pattern
	fmt.Println("\n9. Worker Pool:")
	workerPoolExample()

	// Context for cancellation
	fmt.Println("\n10. Context for Cancellation:")
	contextExample()

	// Fan-out, fan-in pattern
	fmt.Println("\n11. Fan-out, Fan-in:")
	fanOutFanInExample()

	// Rate limiting
	fmt.Println("\n12. Rate Limiting:")
	rateLimitingExample()

	// Atomic operations
	fmt.Println("\n13. Atomic Operations:")
	atomicExample()

	// Closing channels
	fmt.Println("\n14. Closing Channels:")
	closingChannelsExample()

	// Timeout pattern
	fmt.Println("\n15. Timeout Pattern:")
	timeoutExample()

	fmt.Println("\nAll concurrency examples completed!")
}

// 1. Basic goroutine example
func basicGoroutineExample() {
	// Start a goroutine
	go func() {
		fmt.Println("Hello from goroutine!")
	}()

	// Sleep to give the goroutine time to execute
	// In real code, you would use proper synchronization
	time.Sleep(100 * time.Millisecond)
}

// 2. Multiple goroutines example
func multipleGoroutinesExample() {
	for i := 0; i < 5; i++ {
		// Capture the loop variable to avoid closure problems
		local := i
		go func() {
			fmt.Printf("Goroutine %d executing\n", local)
		}()
	}

	// Sleep to give goroutines time to execute
	time.Sleep(100 * time.Millisecond)
}

// 3. Basic channel example
func basicChannelExample() {
	// Create a channel
	ch := make(chan string)

	// Send a value to the channel in a goroutine
	go func() {
		ch <- "Hello from channel!"
	}()

	// Receive the value from the channel
	message := <-ch
	fmt.Println(message)
}

// 4. Buffered channel example
func bufferedChannelExample() {
	// Create a buffered channel with capacity 3
	ch := make(chan string, 3)

	// Send values without blocking
	ch <- "First"
	ch <- "Second"
	ch <- "Third"

	// Receive values
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

// 5. Channel directions example
func channelDirectionsExample() {
	ch := make(chan string)

	// Start sender and receiver goroutines
	go sender(ch)
	go receiver(ch)

	// Wait for goroutines to finish
	time.Sleep(100 * time.Millisecond)
}

// Sender only sends to the channel
func sender(ch chan<- string) {
	ch <- "Message from sender"
}

// Receiver only receives from the channel
func receiver(ch <-chan string) {
	message := <-ch
	fmt.Println(message)
}

// 6. Select statement example
func selectExample() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// Send values on both channels
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- "Channel 1"
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- "Channel 2"
	}()

	// Use select to wait for either channel
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received from", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received from", msg2)
		}
	}
}

// 7. Mutex for synchronization example
func mutexExample() {
	var counter int
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Launch 10 goroutines that increment the counter
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Lock the mutex before accessing the counter
			mutex.Lock()
			counter++
			mutex.Unlock()
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("Counter:", counter)
}

// 8. WaitGroup example
func waitGroupExample() {
	var wg sync.WaitGroup

	// Launch 5 goroutines
	for i := 1; i <= 5; i++ {
		wg.Add(1) // Increment the WaitGroup counter

		// Capture the loop variable
		id := i

		// Launch a goroutine
		go func() {
			defer wg.Done() // Decrement the counter when the goroutine completes

			// Simulate work
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			fmt.Printf("Worker %d done\n", id)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All workers done")
}

// 9. Worker pool pattern
func workerPoolExample() {
	const numJobs = 10
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results
	for a := 1; a <= numJobs; a++ {
		<-results
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, j)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		fmt.Printf("Worker %d finished job %d\n", id, j)
		results <- j * 2
	}
}

// 10. Context for cancellation
func contextExample() {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel() // Ensure all paths cancel the context to avoid context leak

	// Start a goroutine that uses the context
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine cancelled due to:", ctx.Err())
				return
			default:
				fmt.Println("Goroutine working...")
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	// Wait for the context to timeout
	time.Sleep(300 * time.Millisecond)
}

// 11. Fan-out, fan-in pattern
func fanOutFanInExample() {
	// Generate some input
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Fan-out to 3 workers
	c1 := square(input[:3])
	c2 := square(input[3:6])
	c3 := square(input[6:])

	// Fan-in the results
	for result := range merge(c1, c2, c3) {
		fmt.Println("Result:", result)
	}
}

// square returns a channel that emits the squares of the numbers
func square(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n * n
			time.Sleep(50 * time.Millisecond) // Simulate work
		}
	}()
	return out
}

// merge combines multiple channels into a single channel
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel
	wg.Add(len(cs))
	for _, c := range cs {
		go func(ch <-chan int) {
			defer wg.Done()
			for n := range ch {
				out <- n
			}
		}(c)
	}

	// Start a goroutine to close the output channel once all input channels are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// 12. Rate limiting example
func rateLimitingExample() {
	// Create a rate limiter that allows 2 operations per second
	rate := time.Second / 2
	limiter := time.Tick(rate)

	// Simulate 5 requests
	for i := 1; i <= 5; i++ {
		<-limiter // Wait for the rate limiter
		fmt.Printf("Request %d processed at %s\n", i, time.Now().Format("15:04:05.000"))
	}
}

// 13. Atomic operations example
func atomicExample() {
	var counter int64
	var wg sync.WaitGroup

	// Launch 1000 goroutines that increment the counter
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Use atomic operation to increment the counter
			// atomic.AddInt64(&counter, 1)
			// For simplicity, we're using a non-atomic operation here
			// which may result in race conditions
			counter++
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("Counter (may have race conditions):", counter)
}

// 14. Closing channels example
func closingChannelsExample() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	// Start a worker that processes jobs
	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("Received job", j)
			} else {
				fmt.Println("All jobs received")
				done <- true
				return
			}
		}
	}()

	// Send 3 jobs, then close the channel
	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("Sent job", j)
	}
	close(jobs)
	fmt.Println("Sent all jobs")

	// Wait for the worker to finish
	<-done
}

// 15. Timeout pattern example
func timeoutExample() {
	// Create a channel for the result
	resultCh := make(chan string, 1)

	// Start a goroutine that simulates a slow operation
	go func() {
		// Simulate work that takes time
		time.Sleep(200 * time.Millisecond)
		resultCh <- "Operation completed"
	}()

	// Wait for the result with a timeout
	select {
	case result := <-resultCh:
		fmt.Println(result)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Operation timed out")
	}
}

// Additional examples (not called in main for brevity)

// HTTP server with graceful shutdown
func httpServerExample() {
	// Create a server
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, World!")
		}),
	}

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Wait for a signal to shutdown
	time.Sleep(5 * time.Second) // Simulate waiting for a signal

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("HTTP server shutdown error: %v\n", err)
	}

	fmt.Println("HTTP server stopped gracefully")
}

// Pipeline pattern
func pipelineExample() {
	// Generate numbers
	numbers := make(chan int)
	go func() {
		defer close(numbers)
		for i := 1; i <= 10; i++ {
			numbers <- i
		}
	}()

	// Square the numbers
	squares := make(chan int)
	go func() {
		defer close(squares)
		for n := range numbers {
			squares <- n * n
		}
	}()

	// Filter even squares
	even := make(chan int)
	go func() {
		defer close(even)
		for s := range squares {
			if s%2 == 0 {
				even <- s
			}
		}
	}()

	// Print the results
	for e := range even {
		fmt.Println(e)
	}
}

// Semaphore pattern for limiting concurrency
func semaphoreExample() {
	const maxConcurrent = 3
	const totalTasks = 10

	// Create a semaphore using a buffered channel
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for i := 1; i <= totalTasks; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }() // Release semaphore

			fmt.Printf("Task %d starting\n", id)
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			fmt.Printf("Task %d completed\n", id)
		}(i)
	}

	wg.Wait()
}
