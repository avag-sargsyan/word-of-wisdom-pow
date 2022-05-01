package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/avag-sargsyan/word-of-wisdom-pow/pkg/pow"
	"github.com/avag-sargsyan/word-of-wisdom-pow/internal/pkg/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// White-box testing for TCP Client

// MockConnection - mocks tcp connection by two interfaces (reader and writer) and funcs
type MockConnection struct {
	ReadFunc  func([]byte) (int, error)
	WriteFunc func([]byte) (int, error)
}

func (m MockConnection) Read(p []byte) (n int, err error) {
	return m.ReadFunc(p)
}

func (m MockConnection) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

func TestHandleConnection(t *testing.T) {
	t.Parallel()

	t.Run("Write error", func(t *testing.T) {
		mock := MockConnection{
			WriteFunc: func(p []byte) (int, error) {
				return 0, fmt.Errorf("test write error")
			},
		}
		_, err := HandleConnection(mock, mock)
		assert.Error(t, err)
		assert.Equal(t, "error sending request: test write error", err.Error())
	})

	t.Run("Read error", func(t *testing.T) {
		mock := MockConnection{
			WriteFunc: func(p []byte) (int, error) {
				return 0, nil
			},
			ReadFunc: func(p []byte) (int, error) {
				return 0, fmt.Errorf("test read error")
			},
		}
		_, err := HandleConnection(mock, mock)
		assert.Error(t, err)
		assert.Equal(t, "error reading message: test read error", err.Error())
	})

	t.Run("Read response in bad format", func(t *testing.T) {
		mock := MockConnection{
			WriteFunc: func(p []byte) (int, error) {
				return 0, nil
			},
			ReadFunc: func(p []byte) (int, error) {
				return fillTestReadBytes("||\n", p), nil
			},
		}
		_, err := HandleConnection(mock, mock)
		assert.Error(t, err)
		assert.Equal(t, "error parsing message: message doesn't match protocol", err.Error())
	})

	t.Run("Read response in bad format", func(t *testing.T) {
		mock := MockConnection{
			WriteFunc: func(p []byte) (int, error) {
				return 0, nil
			},
			ReadFunc: func(p []byte) (int, error) {
				return fillTestReadBytes(fmt.Sprintf("%d|{wrong_json}\n", protocol.ResponseChallenge), p), nil
			},
		}
		_, err := HandleConnection(mock, mock)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "error parsing puzzle message"))
	})

	t.Run("Success", func(t *testing.T) {
		date := time.Now()

		dt := date.Unix()
		phSum := pow.GetPuzzleHash(dt, "test")

		thSum := pow.GetTargetHash(phSum)

		puzzle := pow.ClientPuzzle{
			TargetHash:    thSum,
			PuzzleToSolve: phSum[:len(phSum)-2],
			Date:          dt,
			Strength:      2,
			Key:           123,
		}

		// counter for reading attempts to change content
		readAttempt := 0

		writeAttempt := 0

		mock := MockConnection{
			WriteFunc: func(p []byte) (int, error) {
				if writeAttempt == 0 {
					writeAttempt++
					assert.Equal(t, "1|\n", string(p))
				} else {
					msg, err := protocol.ParseMessage(string(p))
					require.NoError(t, err)
					var wp pow.ClientPuzzle
					err = json.Unmarshal([]byte(msg.Payload), &wp)
					require.NoError(t, err)
				}
				return 0, nil
			},
			ReadFunc: func(p []byte) (int, error) {
				if readAttempt == 0 {
					marshaled, err := json.Marshal(puzzle)
					require.NoError(t, err)
					readAttempt++
					return fillTestReadBytes(fmt.Sprintf("%d|%s\n", protocol.ResponseChallenge, string(marshaled)), p), nil
				} else {
					// second read, send quote
					return fillTestReadBytes(fmt.Sprintf("%d|test quote\n", protocol.ResponseChallenge), p), nil
				}
			},
		}
		response, err := HandleConnection(mock, mock)
		assert.NoError(t, err)
		assert.Equal(t, "test quote", response)
	})
}

// fillTestReadBytes - helper to easier mock Reader
func fillTestReadBytes(str string, p []byte) int {
	dataBytes := []byte(str)
	counter := 0
	for i := range dataBytes {
		p[i] = dataBytes[i]
		counter++
		if counter >= len(p) {
			break
		}
	}
	return counter
}
