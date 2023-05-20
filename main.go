package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorcon/rcon"
)

const (
	LISTPLAYERS_HEADER_SEPARATOR = "==============================================================================="
	LISTPLAYERS_CELL_SEPARATOR   = "|"
)

type Player struct {
	ID    string
	Name  string
	NetID string
	IP    string
	Score string
}

func main() {

	persistantPlayers := make(map[string]*Player)

	tick := time.Tick(15 * time.Second)

	for range tick {
		go func() {
			conn, err := rcon.Dial(os.Getenv("ADDR"), os.Getenv("PASS"))
			if err != nil {
				log.Fatal(err)
			}
			var wg sync.WaitGroup
			for _, player := range getPlayers(conn) {
				if player.ID == "0" || player.ID == "" || player.Name == "" {
					continue
				}
				_, exists := persistantPlayers[player.Name]
				if !exists {
					wg.Add(1)
					go func(player *Player) {
						defer wg.Done()
						time.Sleep(5 * time.Second)
						persistantPlayers[player.Name] = player
						fmt.Printf("new player %s\n", player.Name)
						conn.Execute(fmt.Sprintf("say Hi %s, Welcome! Be advised, this is a high bot-count server. Take it slow", player.Name))
					}(player)
				}
			}
			wg.Wait()
			conn.Close()
		}()
	}

}

func getPlayers(conn *rcon.Conn) []*Player {
	response, _ := conn.Execute("listplayers")
	response = strings.Replace(response, "\t", "", -1)
	response = strings.Replace(response, "\n", "", -1)
	response = strings.Replace(response, " | ", LISTPLAYERS_CELL_SEPARATOR, -1)
	splitResponse := strings.Split(response, LISTPLAYERS_HEADER_SEPARATOR)

	players := make([]*Player, 0)

	playersBlocks := strings.Split(splitResponse[1], LISTPLAYERS_CELL_SEPARATOR)
	for i := 0; i < len(playersBlocks); i++ {
		switch i % 5 {
		case 0:
			if playersBlocks[i] == "0" {
				i += 4
				continue
			}
			players = append(players, &Player{ID: playersBlocks[i]})
		case 1:
			players[len(players)-1].Name = playersBlocks[i]
		case 2:
			players[len(players)-1].NetID = playersBlocks[i]
		case 3:
			players[len(players)-1].IP = playersBlocks[i]
		case 4:
			players[len(players)-1].Score = playersBlocks[i]
		}
	}

	return players
}

// func stdin() {
// 	conn, err := rcon.Dial(os.Getenv("ADDR"), os.Getenv("PASS"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Close()
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		command, _ := reader.ReadString('\n')
// 		response, _ := conn.Execute(command)
// 		fmt.Println(response)
// 	}
// }
