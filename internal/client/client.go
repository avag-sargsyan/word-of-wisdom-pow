package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/avag-sargsyan/word-of-wisdom-pow/pkg/pow"
	"github.com/avag-sargsyan/word-of-wisdom-pow/internal/pkg/protocol"
)

// Establish TCP connection by the given host and port
func Run(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	fmt.Println("connected to", address)
	defer conn.Close()

	// Make request endlessly
	for {
		message, err := HandleConnection(conn, conn)
		if err != nil {
			return err
		}
		fmt.Println("word of wisdom:", message)
		time.Sleep(3 * time.Second)
	}
}

// Handle TCP connection
func HandleConnection(readerConn io.Reader, writerConn io.Writer) (string, error) {
	reader := bufio.NewReader(readerConn)

	// Send message with {RequestChallenge} header to receive challenge
	err := sendMsg(protocol.Message{
		Header: protocol.RequestChallenge,
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}

	// Parse received message
	msgStr, err := readMsg(reader)
	if err != nil {
		return "", fmt.Errorf("error reading message: %w", err)
	}
	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("error parsing message: %w", err)
	}

	// Unmarshal the received message body into ClientPuzzle struct
	var puzzle pow.ClientPuzzle
	err = json.Unmarshal([]byte(msg.Payload), &puzzle)
	if err != nil {
		return "", fmt.Errorf("error parsing puzzle message: %w", err)
	}

	// Try to solve the puzzle
	puzzle, err = puzzle.SolvePuzzle()
	if err != nil {
		return "", fmt.Errorf("error solving puzzle: %w", err)
	}

	// Prepare sending to server the solved puzzle
	pm, err := json.Marshal(puzzle)
	if err != nil {
		return "", fmt.Errorf("error marshaling puzzle: %w", err)
	}

	// Send to server the solved puzzle with the {RequestResource} header
	// that is if everything is successful establish the connection
	err = sendMsg(protocol.Message{
		Header:  protocol.RequestResource,
		Payload: string(pm),
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	
	// TCP connection is established, receive data from server
	msgStr, err = readMsg(reader)
	if err != nil {
		return "", fmt.Errorf("error reading message: %w", err)
	}
	msg, err = protocol.ParseMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("error parsing message: %w", err)
	}

	return msg.Payload, nil
}

// Read message from connection
func readMsg(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

// Send protocol message to server
func sendMsg(msg protocol.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	return err
}
