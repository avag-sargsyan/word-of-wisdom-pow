package pow_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/avag-sargsyan/word-of-wisdom-pow/pkg/pow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Black-box testing for Proof of Work package

func TestGetPuzzleHash(t *testing.T) {
	date := time.Date(2022, 5, 1, 12, 36, 0, 0, time.UTC)
	secret := "test_1"
	dt := date.Unix()
	phSum := pow.GetPuzzleHash(dt, secret)

	assert.Equal(t, "f6e68548ea45915253c9a2431fb1cd921b122283bcaba014900ad6121a8ad958", fmt.Sprintf("%x", phSum))
}

func TestGetTargetHash(t *testing.T) {
	date := time.Date(2022, 5, 1, 12, 42, 0, 0, time.UTC)
	secret := "test_2"
	dt := date.Unix()
	phSum := pow.GetPuzzleHash(dt, secret)
	thSum := pow.GetTargetHash(phSum)

	assert.Equal(t, "650acdb9b19027697191b8cbc9d358cc02b5e064a3289c03de636abb86b880d3", fmt.Sprintf("%x", thSum))
}

func TestSolvePuzzle(t *testing.T) {
	date := time.Date(2022, 5, 1, 12, 48, 0, 0, time.UTC)
	secret := "test_3"
	strength := 2
	dt := date.Unix()
	phSum := pow.GetPuzzleHash(dt, secret)
	thSum := pow.GetTargetHash(phSum)

	puzzle := pow.ClientPuzzle{
		TargetHash:    thSum,
		PuzzleToSolve: phSum[:len(phSum)-strength],
		Date:          dt,
		Strength:      strength,
		Key:           123,
	}

	solvedPuzzle, err := puzzle.SolvePuzzle()
	require.NoError(t, err)
	assert.Equal(t, solvedPuzzle.PuzzleToSolve, phSum)
}

func TestIsCorrect(t *testing.T) {
	date := time.Date(2022, 5, 1, 12, 59, 0, 0, time.UTC)
	secret := "test_4"
	strength := 2
	dt := date.Unix()
	phSum := pow.GetPuzzleHash(dt, secret)
	thSum := pow.GetTargetHash(phSum)

	puzzle := pow.ClientPuzzle{
		TargetHash:    thSum,
		PuzzleToSolve: phSum,
		Date:          dt,
		Strength:      strength,
		Key:           123,
	}

	isCorrect := puzzle.IsCorrect(secret)
	assert.True(t, isCorrect)
}
