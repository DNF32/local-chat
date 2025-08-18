package network

import (
	"io"
	"local-chat/internal/message"
	"log"
	"log/slog"
	"net"
)

// Function to read a byte stream of this protocol
func ReadProtocol(r io.Reader, buf []byte) ([]byte, error) {
	var result []byte
	if len(buf) == 0 {
		buf = make([]byte, 64)
		slog.Debug("ReadProtocol: received zero-length buffer, using default 64-byte buffer")
	}
	slog.Debug("ReadProtocol: starting to read from connection")
	for {
		n, err := r.Read(buf)
		slog.Debug("ReadProtocol: read bytes", "n", n)
		if n > 0 {
			slog.Debug("ReadProtocol: data chunk", "chunk", string(buf[:n]))
			result = append(result, buf[:n]...)
			slog.Debug("ReadProtocol: total data so far", "result", string(result))
			if len(result) >= 2 &&
				result[len(result)-2] == '\n' &&
				result[len(result)-1] == '\n' {
				slog.Debug("ReadProtocol: found delimiter '\\n\\n'")
				return result[:len(result)-2], nil
			}
		}
		if err != nil {
			if err == io.EOF {
				slog.Debug("ReadProtocol: reached EOF")
				break
			}
			slog.Error("ReadProtocol: read error", "error", err)
			return nil, err
		}
	}
	slog.Warn("ReadProtocol: finished reading, no delimiter found")
	return result, nil
}

func HandleInput(conn *net.TCPConn, outgoing chan message.Message) {
	log.Println("DEBUG: HandleInput goroutine started")
	slog.Info("HandleInput: goroutine started", "connection", conn.RemoteAddr())

	for {
		log.Println("DEBUG: HandleInput waiting for message...")
		slog.Debug("HandleInput: waiting for message on outgoing channel")

		msg := <-outgoing
		log.Printf("DEBUG: HandleInput received message: %+v", msg)
		slog.Info("HandleInput: received message",
			"type", msg.Type,
			"username", msg.Username,
			"content", msg.Content,
			"timestamp", msg.Timestamp)

		data, err := msg.Encode()
		if err != nil {
			log.Printf("DEBUG: Failed to encode message: %s", err)
			slog.Error("HandleInput: failed to encode message", "error", err)
			return
		}

		log.Printf("DEBUG: Encoded message, length: %d bytes", len(data))
		slog.Debug("HandleInput: encoded message",
			"length", len(data),
			"data", string(data))

		bytesWritten, err := conn.Write(data)
		if err != nil {
			log.Printf("DEBUG: Failed to write to connection: %s", err.Error())
			slog.Error("HandleInput: failed to write to connection", "error", err)
			return
		}

		log.Println("DEBUG: Successfully wrote message to connection")
		slog.Info("HandleInput: successfully sent message",
			"bytes_written", bytesWritten,
			"remote_addr", conn.RemoteAddr())
	}
}

func HandleOutput(conn *net.TCPConn, incoming chan message.Message) {
	log.Println("DEBUG: HandleOutput goroutine started")
	slog.Info("HandleOutput: goroutine started", "connection", conn.RemoteAddr())

	buf := make([]byte, 1024)
	slog.Debug("HandleOutput: created buffer", "size", len(buf))

	for {
		slog.Debug("HandleOutput: waiting for data from connection")
		log.Println("DEBUG: HandleOutput waiting for data...")

		data, err := ReadProtocol(conn, buf)
		if err != nil {
			log.Printf("DEBUG: Failed to read protocol: %s", err)
			slog.Error("HandleOutput: failed to read protocol", "error", err)
			// Don't return here - continue trying to read
			continue
		}

		if len(data) == 0 {
			slog.Warn("HandleOutput: received empty data, continuing...")
			continue
		}

		log.Printf("DEBUG: Read %d bytes from connection", len(data))
		slog.Debug("HandleOutput: received data",
			"length", len(data),
			"raw_data", string(data))

		var msg message.Message
		err = msg.Decode(data)
		if err != nil {
			log.Printf("DEBUG: Failed to decode message: %s", err)
			slog.Error("HandleOutput: failed to decode message",
				"error", err,
				"raw_data", string(data))
			continue
		}

		log.Printf("DEBUG: Decoded message: %+v", msg)
		slog.Info("HandleOutput: decoded message",
			"type", msg.Type,
			"username", msg.Username,
			"content", msg.Content,
			"timestamp", msg.Timestamp)

		slog.Debug("HandleOutput: sending message to incoming channel")
		select {
		case incoming <- msg:
			log.Println("DEBUG: Successfully sent message to incoming channel")
			slog.Debug("HandleOutput: successfully sent message to incoming channel")
		default:
			log.Println("DEBUG: WARNING - Incoming channel is full!")
			slog.Warn("HandleOutput: incoming channel is full, dropping message")
		}
	}
}
