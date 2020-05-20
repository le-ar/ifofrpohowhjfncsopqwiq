import (
	"fmt"
	"testing"
	"time"
)

func slowSquare(a int) int {
	time.Sleep(50 * time.Millisecond)
	return a * a
}

func fastSquare(a int) int {
	return a * a
}

func TestNonBlocking(t *testing.T) {
	capacity := 100

	ch1 := make(chan int, capacity)
	a1 := make([]int, capacity)

	ch2 := make(chan int, capacity)
	a2 := make([]int, capacity)

	chOut := make(chan int, capacity)
	a3 := make([]int, capacity)

	i := 0
	for i < capacity {
		a1[i] = i + 9
		ch1 <- a1[i]

		a2[i] = i*3 + 289
		ch2 <- a2[i]

		a3[i] = fastSquare(a1[i]) + fastSquare(a2[i])
		i++
	}

	done := make(chan struct{})

	portion := 30

	go func() {
		MergeTChannels(slowSquare, ch1, ch2, chOut, portion)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(time.Millisecond * 100):
		t.Fail()
		panic("Function should be non-blocking")
	}
}

func TestSlowSquare(t *testing.T) {
	capacity := 100

	ch1 := make(chan int, capacity)
	a1 := make([]int, capacity)

	ch2 := make(chan int, capacity)
	a2 := make([]int, capacity)

	chOut := make(chan int, capacity)
	a3 := make([]int, capacity)

	i := 0
	for i < capacity {
		a1[i] = i + 9
		ch1 <- a1[i]

		a2[i] = i*3 + 289
		ch2 <- a2[i]

		a3[i] = fastSquare(a1[i]) + fastSquare(a2[i])
		i++
	}

	done := make(chan struct{})

	portion := 30

	go func() {
		MergeTChannels(slowSquare, ch1, ch2, chOut, portion)
		close(done)
	}()

	<-done

	i = 0
	for i < portion {
		ans, ok := <-chOut
		if !ok {
			t.Fail()
			panic("Output channel closed prematurely")
		}
		if ans != a3[i] {
			t.Fail()
			panic(fmt.Errorf("Got %d from output channel, should be %d", ans, a3[i]))
		}
		i++
	}

	if len(ch1) != capacity-portion {
		t.Fail()
		panic(fmt.Errorf("First channel has %d numbers in it, should have %d", len(ch1), capacity-portion))
	}
	if len(ch2) != capacity-portion {
		t.Fail()
		panic(fmt.Errorf("Second channel has %d numbers in it, should have %d", len(ch2), capacity-portion))
	}
}
