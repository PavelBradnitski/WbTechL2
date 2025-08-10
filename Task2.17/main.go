package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// RunClient запускает клиентское приложение, которое читает данные из входного Reader и отправляет их в сокет,
// а также читает данные из сокета и выводит их в выходной Writer.
func RunClient(conn net.Conn, input io.Reader, output io.Writer) error {
	done := make(chan struct{})
	var once sync.Once
	closeDone := func() {
		once.Do(func() { close(done) })
	}

	// Обработчик сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		conn.Close()
		closeDone()
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	stdinReader := bufio.NewReader(input)
	stdoutWriter := bufio.NewWriter(output)

	errCh := make(chan error, 2)

	// Горрутина чтения из сокета и вывода в stdout
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
					errCh <- nil
					return
				}
				log.Printf("Error reading from socket: %v", err)
				errCh <- err
				return
			}
			_, err = stdoutWriter.WriteString(line)
			if err != nil {
				errCh <- err
				return
			}
			stdoutWriter.Flush()
		}
	}()

	// Горрутина чтения из stdin и записи в сокет
	go func() {
		for {
			line, err := stdinReader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("EOF reached, closing connection")
					if c, ok := conn.(interface{ CloseWrite() error }); ok {
						c.CloseWrite()
					}
					errCh <- nil
					return
				}
				errCh <- err
				return
			}

			if strings.HasSuffix(line, "\n") && !strings.HasSuffix(line, "\r\n") {
				line = strings.TrimSuffix(line, "\n") + "\r\n"
			}

			_, err = writer.WriteString(line)
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					errCh <- nil
					return
				}
				errCh <- err
				return
			}
			writer.Flush()
		}
	}()

	var firstErr error
	for i := 0; i < 2; i++ {
		err := <-errCh
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	closeDone()

	return firstErr
}

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout duration] host port\n", os.Args[0])
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		if os.IsTimeout(err) {
			fmt.Fprintf(os.Stderr, "Error: connection to %s timed out after %s\n", address, *timeout)
		} else {
			fmt.Fprintf(os.Stderr, "Error: failed to connect to %s: %v\n", address, err)
		}
		os.Exit(1)
	}
	defer conn.Close()

	if err := RunClient(conn, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Client error: %v\n", err)
		os.Exit(1)
	}
}
