package main

import (
	"bytes"
	"net"
	"strings"
	"testing"
)

func TestRunClient(t *testing.T) {
	serverConn, clientConn := net.Pipe()

	input := strings.NewReader("hello\n")
	var output bytes.Buffer

	doneReading := make(chan struct{})

	go func() {
		err := RunClient(clientConn, input, &output)
		if err != nil {
			t.Errorf("runClient error: %v", err)
		}
		close(doneReading)
	}()

	readDone := make(chan struct{})
	writeDone := make(chan struct{})
	// Горрутина для чтения из serverConn и проверки полученных данных
	go func() {
		buf := make([]byte, 1024)
		n, err := serverConn.Read(buf)
		if err != nil {
			t.Errorf("serverConn.Read error: %v", err)
			close(readDone)
			return
		}
		got := string(buf[:n])
		want := "hello\r\n"
		if got != want {
			t.Errorf("Expected to receive %q, got %q", want, got)
		}
		close(readDone)
	}()
	// Горрутина для записи в serverConn
	go func() {
		<-readDone
		_, err := serverConn.Write([]byte("world\r\n"))
		if err != nil {
			t.Errorf("serverConn.Write error: %v", err)
		}
		close(writeDone)
	}()

	<-writeDone
	serverConn.Close() // закрываем соединение — клиент получит EOF

	<-doneReading

	if !strings.Contains(output.String(), "world\r\n") {
		t.Errorf("Expected output to contain %q, got %q", "world\r\n", output.String())
	}
}
