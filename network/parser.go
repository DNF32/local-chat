package network

import "net"

// This seems super familiar with a io.Reader and writer interface
// This things should have typed channels like Using the User struct and Message struct this ways this acts as a form
// of transport layer

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
