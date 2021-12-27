package examples

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {

	// Channel and buffer channel
	messages := make(chan string) //
	// messages := make(chan string, 4) // buffered channel that can accomodate max 4 messages without blocking on receiver to consume them

	go func() {
		messages <- "hello"
		messages <- "hello2"
		messages <- "hello3"
		messages <- "hello4"
	}()

	fmt.Println("Received messages: ", <-messages)

	// ##################### Channel Synchronization #######################################
	done := make(chan bool, 1)
	go worker(done)
	<-done
	fmt.Println("resuming...")

	// ############################# Channel Directions ####################################
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)

	// ############################# Channel Select ####################################
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	// select lets you wait on multiple channel operations.
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}

func worker(done chan bool) {
	fmt.Println("working...")
	time.Sleep(2 * time.Second)
	fmt.Println("done")
	done <- true
}

// When using channels as function parameters, you can specify if a channel is meant to only send or receive values.
// This specificity increases the type-safety of the program.
func ping(pings chan<- string, msg string) {
	pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func TestTimeouts(t *testing.T) {
	c1 := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(5 * time.Second):
		fmt.Println("timeout 1")
	}

	// c2 := make(chan string, 1)
	// go func() {
	// 	time.Sleep(2 * time.Second)
	// 	c2 <- "result 2"
	// }()
	// select {
	// case res := <-c2:
	// 	fmt.Println(res)
	// case <-time.After(3 * time.Second):
	// 	fmt.Println("timeout 2")
	// }
}

func TestNonBlockingChannels(t *testing.T) {
	messages := make(chan string)
	signals := make(chan bool)
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}
	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

func TestDeadlock(t *testing.T) {
	ch := make(chan string, 2)
	ch <- "within limit"
	ch <- "reached capacity" // post this buffer is full to its capacity

	// Will block for receiver to read the buffered elements resulting a deadlock since there is no receiver.
	// This Test will not report panic but if run with go run, the situation results into PANIC.
	// fatal error: all goroutines are asleep - deadlock!
	// ch <- "exceeded capacity" 	(Test: blocks indefinitely, Main: PANIC)
	
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

var waitGroup sync.WaitGroup

func TestWorkgroup(t *testing.T) {
	rand.Seed(time.Now().Unix())

	// Create a buffered channel to manage the employee vs project load.
	projects := make(chan string, 10)

	// Launch 5 goroutines to handle the projects.
	waitGroup.Add(5)
	for i := 1; i <= 5; i++ {
		go goRoutine(projects, i)
	}

	for j := 1; j <= 10; j++ {
		projects <- fmt.Sprintf("Project :%d", j)
	}

	// Close the channel so the goroutines will quit
	close(projects)
	waitGroup.Wait()
}

func goRoutine(projects chan string, employee int) {
	defer waitGroup.Done()
	for {
		// Wait for project to be assigned.
		project, result := <-projects

		if result == false {
			// This means the channel is empty and closed.
			fmt.Printf("Employee : %d : Exit\n", employee)
			return
		}

		fmt.Printf("Employee : %d : Started   %s\n", employee, project)

		// Randomly wait to simulate work time.
		sleep := rand.Int63n(50)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		// Display time to wait
		fmt.Println("Time to sleep", sleep, "ms")

		// Display project completed by employee.
		fmt.Printf("Employee : %d : Completed %s\n", employee, project)
	}

}
