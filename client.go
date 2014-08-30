package vindinium

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	MoveTimeout  = 15
	StartTimeout = 10 * 60
)

type Bot interface {
	Move(state *State) Direction
}

type Client struct {
	Server string
	Key string
	Mode string
	Bot Bot
	Turns string
	RandomMap bool
	State *State
	Url string

	Debug bool
	DebugMap bool
	Display_map_chan chan string
	Display_map_started bool
}

func display_map(mapchan chan string, board_size int) {
	for {
		board_repr := <- mapchan
		start := 0
		fmt.Printf("|")
		for i := 0; i < (board_size*2)+2; i++ { fmt.Print("-") }
		fmt.Printf("|\n")
		for start < len(board_repr) {
			fmt.Println("|", board_repr[start:start+board_size*2], "|")
			start += board_size*2
		}
		fmt.Printf("|")
		for i := 0; i < (board_size*2)+2; i++ { fmt.Print("-") }
		fmt.Printf("|\n")
	}
}

func NewClient(server, key, mode, turns string, randomMap bool, debug bool, debug_map bool, bot Bot) (client *Client) {
	client = &Client{
		Server: server,
		Key: key,
		Mode: mode,
		Bot: bot,
		Turns: turns,
		RandomMap: randomMap,
		Debug: debug,
		DebugMap: debug_map,
		Display_map_chan: make(chan string),
		Display_map_started: false,
	}
	client.Url = client.Server + "/api/" + client.Mode
	return
}

func NewTrainingClient(key, turns string, randomMap bool, debug bool, debug_map bool, bot Bot) (client *Client) {
	client = NewClient("http://vindinium.org", key, "training", turns, randomMap, debug, debug_map, bot)
	return
}

func NewArenaClient(key, turns string, randomMap bool, debug bool, debug_map bool, bot Bot) (client *Client) {
	client = NewClient("http://vindinium.org", key, "arena", turns, randomMap, debug, debug_map, bot)
	return
}

func (c *Client) post(uri string, values url.Values, seconds int) error {
	if c.Debug {
		fmt.Printf("Making request to to: %s\n", uri)
	}
	timeout := time.Duration(seconds) * time.Second
	dial := func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, timeout)
	}

	transport := http.Transport{Dial: dial}
	client := http.Client{Transport: &transport}

	response, err := client.PostForm(uri, values)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode >= 500 {
		return errors.New(fmt.Sprintf("Server responded with %s", response.Status))
	} else if response.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("Request error: %s", string(data[:])))
	}

	if err := json.Unmarshal(data, &c.State); err != nil {
		return err
	}

	if c.Display_map_started && c.DebugMap {
		c.Display_map_chan <- c.State.Game.Board.Tiles
	}

	if c.Debug {
		fmt.Printf("Setting data to:\n%s\n", string(data))
	}
	return nil
}

func (c *Client) finished() bool {
	return c.State.Game.Finished
}

func (c *Client) move(dir Direction) error {
	values := make(url.Values)
	values.Set("dir", string(dir))
	return c.post(c.State.PlayUrl, values, MoveTimeout)
}


func (c *Client) Start() error {
	values := make(url.Values)
	values.Set("key", c.Key)
	if c.Mode == "training" {
		values.Set("turns", c.Turns)
		if !c.RandomMap {
			values.Set("map", "m1")
		}
	}

	fmt.Println("Connecting and waiting for other players to join...")
	err := c.post(c.Url, values, StartTimeout)
	if !c.Display_map_started && c.DebugMap {
		go display_map(c.Display_map_chan, c.State.Game.Board.Size)
		c.Display_map_started = true
	}
	return err
}

func (c *Client) Play() error {
	fmt.Printf("Playing at: %s\n", c.State.ViewUrl)
	move := 1
	for c.State.Game.Finished == false {
		fmt.Printf("\rMaking move: %d\n", move)

		if c.Debug {
			fmt.Printf("\nclient: %+v\n", c)
			fmt.Printf("bot: %+v\n", c.Bot)
			fmt.Printf("state: %+v\n", c.State)
		}

		dir := c.Bot.Move(c.State)
		if err := c.move(dir); err != nil {
			return err
		}

		move++
	}

	fmt.Println("\nFinished game.")
	return nil
}
