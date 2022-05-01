package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/avag-sargsyan/word-of-wisdom-pow/pkg/cache"
	"github.com/avag-sargsyan/word-of-wisdom-pow/pkg/pow"
	"github.com/avag-sargsyan/word-of-wisdom-pow/internal/pkg/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// White-box testing for TCP Server

// MockClock - mock for Clock interface (to work with predefined Now)
type MockClock struct {
	NowFunc func() time.Time
}

func (m *MockClock) Now() time.Time {
	if m.NowFunc != nil {
		return m.NowFunc()
	}
	return time.Now()
}

func TestProcessRequest(t *testing.T) {
	store := cache.NewStore()

	const randKey = 123123123

	t.Run("Quit request", func(t *testing.T) {
		input := fmt.Sprintf("%d|", protocol.Quit)
		msg, err := ProcessRequest(input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, errors.New("closing the connection"), err)
	})

	t.Run("Invalid request", func(t *testing.T) {
		input := "||"
		msg, err := ProcessRequest(input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "message doesn't match protocol", err.Error())
	})

	t.Run("Unknown header", func(t *testing.T) {
		input := "111|"
		msg, err := ProcessRequest(input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "unknown header", err.Error())
	})

	t.Run("Request resource without solution", func(t *testing.T) {
		input := fmt.Sprintf("%d|", protocol.RequestResource)
		msg, err := ProcessRequest(input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.True(t, strings.Contains(err.Error(), "error unmarshaling puzzle"))
	})

	t.Run("Request resource with wrong resource", func(t *testing.T) {
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
		marshaled, err := json.Marshal(puzzle)
		require.NoError(t, err)
		input := fmt.Sprintf("%d|%s", protocol.RequestResource, string(marshaled))
		msg, err := ProcessRequest(input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "couldn't fetch key from cache", err.Error())
	})

	t.Run("Request resource with invalid solution", func(t *testing.T) {
		store.Put(randKey, "test")

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
		marshaled, err := json.Marshal(puzzle)
		require.NoError(t, err)
		input := fmt.Sprintf("%d|%s", protocol.RequestResource, string(marshaled))
		msg, err := ProcessRequest(input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "couldn't fetch key from cache", err.Error())
	})
}
