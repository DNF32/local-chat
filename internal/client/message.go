package client

import (
	"regexp"
	"strings"
)

func ParseContent() (MessageType, string) {
	joinParser := regexp.MustCompile(`^/join\s+(.+)$`)
	leaveParser := regexp.MustCompile(`^/leave\s*$`)

	switch m.Type {
	case Join:
		matches := joinParser.FindStringSubmatch(m.Content)
		if len(matches) == 2 {
			return m.Type, strings.TrimSpace(matches[1])
		}
		return Undefined, ""

	case Leave:
		if leaveParser.MatchString(m.Content) {
			return m.Type, ""
		}
		return Undefined, ""

	case Text, InitUser:
		return m.Type, strings.TrimSpace(m.Content)

	default:
		return Undefined, ""
	}
}
