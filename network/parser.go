package network

import (
	"fmt"
	"io"
	"local-chat/models"
	"net"
)

// Function to read a byte stream of this protocol
func ReadProtocol(r io.Reader, buf []byte) ([]byte, error) {
	var result []byte

	for {
		n, err := r.Read(buf)
		if n > 0 {
			result = append(result, buf[:n]...)
			// Check for protocol delimiter (double newline)
			if len(result) >= 2 &&
				result[len(result)-2] == '\n' &&
				result[len(result)-1] == '\n' {
				// Remove the delimiter from the result
				return result[:len(result)-2], nil
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("network read error: %w", err)
		}
	}
	return result, nil
}

func HandleInput(conn *net.TCPConn, outgoing chan models.Message) {
	for {
		msg := <-outgoing
		data, err := msg.Encode()
		if err != nil {
			fmt.Println("Failed to handle the input", err)
			return
		}
		_, err = conn.Write(data)
		if err != nil {
			return
		}
	}
}
func HandleOutput(conn *net.TCPConn, incoming chan models.Message) {
	buf := make([]byte, 1024)
	for {
		data, err := ReadProtocol(conn, buf)
		if err != nil {
			fmt.Println("Failed to handle the output", err)
		}

		var msg models.Message
		_ = msg.Decode(data)
		incoming <- msg
	}
}
