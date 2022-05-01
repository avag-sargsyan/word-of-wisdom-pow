package server

import (
	"bufio"

	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"

	"time"

	"github.com/avag-sargsyan/word-of-wisdom-pow/internal/pkg/protocol"
	"github.com/avag-sargsyan/word-of-wisdom-pow/pkg/cache"
	"github.com/avag-sargsyan/word-of-wisdom-pow/pkg/pow"
	"github.com/avag-sargsyan/word-of-wisdom-pow/pkg/helper"
)

// Slice of Wise Words which should be served one by one in random order by the established TCP connection
var WiseWords = []string{
	"Many of life’s failures are people who did not realize how close they were to success when they gave up",
	"The whole secret of a successful life is to find out what is one’s destiny to do, and then do it.",
	"Life is like riding a bicycle. To keep your balance, you must keep moving.",
	"You have brains in your head. You have feet in your shoes. You can steer yourself any direction you choose.",
	"Life’s tragedy is that we get old too soon and wise too late.",
	"Every moment is a fresh beginning.",
	"When you cease to dream you cease to live.",
	"The best way to predict your future is to create it.",
	"There are no mistakes, only opportunities.",
	"Sometimes you can’t see yourself clearly until you see yourself through the eyes of others.",
	"All life is an experiment. The more experiments you make, the better.",
	"Do not dwell in the past, do not dream of the future, concentrate the mind on the present moment.",
	"Life is a dream for the wise, a game for the fool, a comedy for the rich, a tragedy for the poor.",
}

// Strength of puzzle, best 1-3
const Strength = 2

// Initiate the in memory key value storage
var store = cache.NewStore()

// Establish TCP connection by the given host and port
func Run(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Println("listening", listener.Addr())
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accept connection: %w", err)
		}
		// Handle connections in a new goroutine.
		go handleConnection(conn)
	}
}

// Handle TCP connection
func handleConnection(conn net.Conn) {
	fmt.Println("new client:", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Serve connection
	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("err read connection:", err)
			return
		}
		msg, err := ProcessRequest(req, conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("err process request:", err)
			return
		}
		if msg != nil {
			err := sendMsg(*msg, conn)
			if err != nil {
				fmt.Println("err send message:", err)
			}
		}
	}
}

// Process and response client request establishing connection once the client puzzle is solved
// https://en.wikipedia.org/wiki/Client_Puzzle_Protocol
func ProcessRequest(msgStr string, clientInfo string) (*protocol.Message, error) {
	msg, err := protocol.ParseMessage(msgStr)

	if err != nil {
		return nil, err
	}

	// Take action accordingly based on received request header
	switch msg.Header {
	case protocol.Quit:
		return nil, errors.New("closing the connection")
	case protocol.RequestChallenge:
		// Prepare Challenge for the client
		date := time.Now()

		secret := helper.RandomString(20)
		key := rand.Intn(1000000)
		// Store the generated secret in memory storage
		store.Put(key, secret)

		dt := date.Unix()

		// Puzzle Hash this is what the client should find out by brute forcing the missing part
		phSum := pow.GetPuzzleHash(dt, secret)

		// Target hash is the second tier hash which generated from puzzle hash
		thSum := pow.GetTargetHash(phSum)

		// Instansiate the client puzzle
		puzzle := pow.ClientPuzzle{
			TargetHash:    thSum,
			PuzzleToSolve: phSum[:len(phSum)-Strength], // Chop from puzzle hash by the given strength
			Date:          dt,
			Strength:      Strength,
			Key:           key,
		}

		// Prepare puzzle data to send to client
		pm, err := json.Marshal(puzzle)
		if err != nil {
			return nil, fmt.Errorf("err marshaling puzzle: %v", err)
		}

		// Send over puzzle to client
		msg := protocol.Message{
			Header:  protocol.ResponseChallenge,
			Payload: string(pm),
		}
		return &msg, nil
	case protocol.RequestResource:
		// Before serving requested resource verify the puzzle solution

		// Instansiate CLientPuzzle by the received data
		var puzzle pow.ClientPuzzle
		err := json.Unmarshal([]byte(msg.Payload), &puzzle)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling puzzle: %w", err)
		}

		// Fetch from in memory storage previously saved secret to verify the puzzle solution
		secret, err := store.Get(puzzle.Key)
		if err != nil {
			return nil, fmt.Errorf("couldn't fetch key from cache")
		}

		// Verify the puzzle solution and proceed only if it's correct
		if !puzzle.IsCorrect(secret) {
			return nil, fmt.Errorf("puzzle isn't solved")
		}

		// Free memory from the stored secret, we don't need it anymore
		store.Delete(puzzle.Key)

		// Serve requested data to client
		msg := protocol.Message{
			Header:  protocol.ResponseResource,
			Payload: WiseWords[rand.Intn(len(WiseWords))],
		}
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown header")
	}
}

// Send protocol message to client
func sendMsg(msg protocol.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	return err
}
