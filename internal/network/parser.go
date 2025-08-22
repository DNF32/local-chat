package network

import (
	"io"
	"local-chat/internal/message"
	"log/slog"
	"net"
)

// Function to read a byte stream of this protocol
func ReadProtocol(r io.Reader, buf []byte) ([]byte, error) {
	var result []byte
	if len(buf) == 0 {
		buf = make([]byte, 64)
	}
	for {
		n, err := r.Read(buf)
		if n > 0 {
			result = append(result, buf[:n]...)
			if len(result) >= 2 &&
				result[len(result)-2] == '\n' &&
				result[len(result)-1] == '\n' {
				return result[:len(result)-2], nil
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return result, nil
}

func HandleInput(conn *net.TCPConn, outgoing chan message.Message, logger *slog.Logger) {
	logger.Info("HandleInput started", "addr", conn.RemoteAddr())

	for {
		msg := <-outgoing

		data, err := msg.Encode()
		if err != nil {
			logger.Error("Encode failed", "error", err)
			return
		}

		logger.Info("Sending message", "data", string(data))

		_, err = conn.Write(data)
		if err != nil {
			logger.Error("Write failed", "error", err)
			return
		}
	}
}

func HandleOutput(conn *net.TCPConn, incoming chan message.Message, logger *slog.Logger) {
	logger.Info("HandleOutput started", "addr", conn.RemoteAddr())
	buf := make([]byte, 1024)

	for {
		data, err := ReadProtocol(conn, buf)
		if err != nil {
			logger.Error("Read failed", "error", err)
			continue
		}
		if len(data) == 0 {
			continue
		}

		var msg message.Message
		err = msg.Decode(data)
		if err != nil {
			logger.Error("Decode failed", "error", err, "data", string(data))
			continue
		}

		logger.Info("Received message", "type", msg.Type, "username", msg.Username, "content", msg.Content)

		select {
		case incoming <- msg:
			// Message sent successfully
		default:
			logger.Warn("Channel full, dropping message")
		}
	}
}
