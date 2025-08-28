package network

import (
	"io"
)

func FrameData(data []byte) []byte {
	return append(data, []byte("\n\n")...)
}

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

func WriteProtocol(r io.Writer, data []byte) error {
	frame := FrameData(data)
	for len(frame) > 0 {
		n, err := r.Write(frame)
		if err != nil {
			return err
		}
		frame = frame[n:]
	}
	return nil
}
