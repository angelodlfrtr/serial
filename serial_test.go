package serial

import (
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	port0 := "/dev/tty.usbserial-1410"
	c0 := &Config{Name: port0, Baud: 2000000}

	s1, err := OpenPort(c0)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan int, 1)
	go func() {
		buf := make([]byte, 128)
		var readCount int
		for {
			n, err := s1.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
			readCount++
			t.Logf("Read %v %v bytes: % 02x %s", readCount, n, buf[:n], buf[:n])
			select {
			case <-ch:
				ch <- readCount
				close(ch)
			default:
			}
		}
	}()

	if _, err = s1.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	if _, err = s1.Write([]byte(" ")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	if _, err = s1.Write([]byte("world")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second / 10)

	ch <- 0
	s1.Write([]byte(" ")) // We could be blocked in the read without this
	c := <-ch
	exp := 5
	if c >= exp {
		t.Fatalf("Expected less than %v read, got %v", exp, c)
	}
}
