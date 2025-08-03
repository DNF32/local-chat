package network

import (
	"encoding/json"
	"fmt"
	"io"
	"local-chat/models"
	"net"
)

// This seems super familiar with a io.Reader and writer interface
// This things should have typed channels like Using the User struct and Message struct this ways this acts as a form
// of transport layer
func Encode(s any) ([]byte, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	data = append(data, []byte("\n\n")...)

	return data, nil
}

// This function parses the json data coming from ReadProtocol to the object `s`, which will be a Message or a User
func Decode(data []byte, s any) error {
	return json.Unmarshal(data, s)
}

// Function to read a byte stream of this protocol
func ReadProtocol(r io.Reader) ([]byte, error) {
	var result []byte
	buf := make([]byte, 1024)

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

func HandleInput(conn *net.TCPConn, outgoing chan string) {
	for {
		msg := <-outgoing
		msg = msg + "\n\n"
		_, err := conn.Write([]byte(msg))

		if err != nil {
			return
		}
	}
}
func HandleOutput(conn *net.TCPConn, incoming chan string) {
	buf := make([]byte, 1024)
	var result []byte
	for {
		for {
			n, err := conn.Read(buf)

			if n > 0 {
				result = append(result, buf[:n]...)
				if len(result) >= 2 &&
					result[len(result)-2] == '\n' &&
					result[len(result)-1] == '\n' {
					break
				}
			}

			if err != nil {
				return
			}
		}
		incoming <- string(result)
		result = nil
	}
}
