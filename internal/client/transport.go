package client

import (
	"local-chat/internal/network"
	"local-chat/internal/protocol"
	"local-chat/internal/transport"
	"log/slog"
	"net"
)

func HandleInput(conn *net.TCPConn, outgoing chan protocol.ClientMessage, serde *transport.Serde, logger *slog.Logger) {
	logger.Info("HandleInput started", "addr", conn.RemoteAddr())

	for {
		msg := <-outgoing

		data, err := serde.EncodeEncrypted(&msg)
		if err != nil {
			logger.Error("Encode failed", "error", err)
			return
		}

		logger.Info("Sending message", "data", string(data))

		err = network.WriteProtocol(conn, data)
		if err != nil {
			logger.Error("Write failed", "error", err)
			return
		}
	}
}

func HandleOutput(conn *net.TCPConn, incoming chan protocol.ClientMessage, serde *transport.Serde, logger *slog.Logger) {
	logger.Info("HandleOutput started", "addr", conn.RemoteAddr())
	buf := make([]byte, 1024)

	for {
		data, err := network.ReadProtocol(conn, buf)
		if err != nil {
			logger.Error("Read failed", "error", err)
			continue
		}
		if len(data) == 0 {
			continue
		}

		var msg protocol.ClientMessage
		err = serde.DecodeEncrypted(data, &msg)
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
