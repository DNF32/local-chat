package shared

import (
	"fmt"
	"local-chat/internal/protocol"
	"regexp"
	"strings"
)

var (
	joinRegex   = regexp.MustCompile(`^/join\s+(.+)$`)
	leaveRegex  = regexp.MustCompile(`^/leave`)
	logoutRegex = regexp.MustCompile(`^/logout`)
)

func ParseMsgType(s string) (protocol.EventType, string, error) {
	s = strings.TrimSpace(s)

	switch {
	case joinRegex.MatchString(s):
		matches := joinRegex.FindStringSubmatch(s)
		if len(matches) > 1 {
			room := strings.TrimSpace(matches[1])
			return protocol.EventTypeJoin, room, nil
		}
		return protocol.EventTypeJoin, "", fmt.Errorf("no room specified after /join")

	case leaveRegex.MatchString(s):
		return protocol.EventTypeLeave, "", nil

	case logoutRegex.MatchString(s):
		return protocol.EventTypeLogout, "", nil

	case strings.HasPrefix(s, "/"):
		return protocol.EventTypeError, "", fmt.Errorf("invalid command: %s", s)

	default:
		// Normal chat message
		return protocol.EventTypeChat, s, nil
	}
}


func BuildClientMessage(msgType protocol.EventType, username, content, roomName string) protocol.ClientMessage {
	switch msgType {
	case protocol.EventTypeJoin:
		return protocol.ClientMessage{
			Type:     protocol.EventTypeJoin,
			Username: username,
			Content:  roomName, // room name
		}

	case protocol.EventTypeLeave:
		return protocol.ClientMessage{
			Type:     protocol.EventTypeLeave,
			Username: username,
		}

	case protocol.EventTypeLogout:
		return protocol.ClientMessage{
			Type:     protocol.EventTypeLogout,
			Username: username,
		}

	case protocol.EventTypeChat:
		return protocol.ClientMessage{
			Type:     protocol.EventTypeChat,
			Username: username,
			Content:  content, // chat message
		}

	default:
		return protocol.ClientMessage{
			Type:     protocol.EventTypeError,
			Username: username,
			Content:  "unsupported client message type",
		}
	}
}
