package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hyptocrypto/dumbNodes/internal/config"
	"github.com/hyptocrypto/dumbNodes/internal/mproto"
	"github.com/hyptocrypto/dumbNodes/internal/types"
)

type Client struct {
	id   string
	conn net.Conn
}

var dummyData = []string{"Hello world from the client", "Some other data", "This is some really long data will all kinds of words and data. This should fill the buffer up much more.", "This is some really long data will all kinds of words and data. This should fill the buffer up much more.This is some really long data will all kinds of words and data. This should fill the buffer up much more.This is some really long data will all kinds of words and data. This should fill the buffer up much more.This is some really long data will all kinds of words and data. This should fill the buffer up much more."}

func randomData() string {
	randomIndex := rand.Intn(len(dummyData))
	return dummyData[randomIndex]
}

// Create new client with default values
func NewClient(config *config.Config) (*Client, error) {
	newId := fmt.Sprintf("ClientSide--%v", uuid.New())

	fmt.Println("Connecting")
	newConn, err := net.Dial(config.Protocol, config.FullHost())
	if err != nil {
		return nil, fmt.Errorf("could not create connection for client: %v", err)
	}
	fmt.Println("Connection cretated")
	return &Client{id: newId, conn: newConn}, nil
}

func (c *Client) SendRequest() error {
	r := types.Request{Method: "GET", Headers: map[string]string{"request": "test"}, Destination: "www.google.com", Data: []byte(randomData())}
	payload, err := mproto.Serialize(r)
	if err != nil {
		return err
	}
	fmt.Printf("Sending request: %v\n", r)

	c.conn.Write(payload)
	return err
}

func main() {
	config := config.NewConfig()
	c, err := NewClient(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for {
		fmt.Println("Send loop")
		err := c.SendRequest()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		time.Sleep(time.Second * 2)
	}
}
