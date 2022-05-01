package pow

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
)

// Structure of the Client Puzzle
type ClientPuzzle struct {
	TargetHash    []byte
	PuzzleToSolve []byte
	Date          int64
	Strength      int
	Key           int
}

// Stringify the Client Puzzle by specific format
func (h ClientPuzzle) Stringify() string {
	return fmt.Sprintf("%s:%s:%d:%d:%d", h.TargetHash, h.PuzzleToSolve, h.Date, h.Strength, h.Key)
}

// Verify the puzzle hash which provided client by matching it with initial
// hash by generating it from the params in ClientPuzzle returned from client 
// and the previously stored secret. Returns true if puzzle solved correctly
func (h ClientPuzzle) IsCorrect(secret string) bool {
	dt := h.Date

	// Firsly regenerate the puzzle hash
	phSum := GetPuzzleHash(dt, secret)

	// Then generate the target hash
	targetRec := []byte(h.TargetHash)
	targetCalc := GetTargetHash(phSum)

	// And finally compare what we have with what client provided
	res := bytes.Compare(targetRec[:], targetCalc[:])

	return res == 0
}

// Find missing part of the puzzle by brute forcing missing hash chunk
func (h ClientPuzzle) SolvePuzzle() (ClientPuzzle, error) {
	// Get target hash and partial puzzle hash
	target := h.TargetHash
	pts := h.PuzzleToSolve
	bs := make([]byte, 9)
	// Brute force hash chunk generated from sequential numbers (nonce alternative)
	for i := 0; i < math.MaxInt64; i++ {
		binary.PutVarint(bs, int64(i))
		h3 := sha256.Sum256(bs)
		// Combine puzzle by each generated hash chunk
		hash1Bytes := append(pts, h3[:h.Strength]...)
		s := sha256.Sum256(hash1Bytes)
		// Compare generated target with the actual target to see if puzzle solved correctly
		res := bytes.Compare(target[:], s[:])

		// If target hashes matching the solution is correct, break the loop
		if res == 0 {
			fmt.Printf("Solved puzzle: %x\n", hash1Bytes)
			h.PuzzleToSolve = hash1Bytes

			return h, nil
		}
	}

	return h, fmt.Errorf("failed solve puzzle")
}

// Generate puzzle hash by the given datatime and secret based on SHA-2 algorithm
func GetPuzzleHash(datetime int64, secret string) []byte {
	dtb := make([]byte, 9)
	binary.PutVarint(dtb, datetime)

	ph := sha256.New()
	ph.Write([]byte(secret))
	ph.Write([]byte(dtb))

	return ph.Sum(nil)
}

// Generate target hash by the given puzzle hash based on SHA-2 algorithm
func GetTargetHash(puzzleSum []byte) []byte {
	th := sha256.New()
	th.Write(puzzleSum)

	return th.Sum(nil)
}
