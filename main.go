package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Lobby struct {
	Players    []string
	Enemies    []string
	Events     []string
	LobbyId    string
	Difficulty string
	Map        string
	Host       string
	Seed       int
	Tick       int
}

var lobbies = []Lobby{}
var TPS = 1

func NetHandler(w http.ResponseWriter, r *http.Request) {
	eventParam := r.URL.Query().Get("e")

	if eventParam == "" {
		http.Error(w, "Missing event parameter", http.StatusBadRequest)
		return
	}

	rText := HandleEvent(eventParam)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, rText)
}

func Tick(lobbyId string) {
	for {

		for i, v := range lobbies {

			if v.LobbyId == lobbyId {
				for _, e := range v.Events {
					HandleEvent(e)
				}

				lobbies[i].Events = []string{}
				lobbies[i].Tick++
			}
		}

		time.Sleep(time.Second / time.Duration(TPS))
	}
}

func GetParam(base string, param string) string {
	for _, v := range strings.Split(base, " ") {
		key := strings.Split(v, "=")[0]
		value := strings.Split(v, "=")[1]

		if key == param {
			return value
		}
	}

	return "nil"
}

func HandleEvent(event string) string {
	s := strings.Split(event, " ")

	if s[0] == "host" { // host Username
		lobbyId := uuid.NewString()
		lobbies = append(lobbies, Lobby{
			LobbyId:    lobbyId,
			Difficulty: "Easy",
			Map:        "Warehouse",
			Host:       s[1],
			Seed:       rand.Int(),
			Events:     []string{"join " + lobbyId + " " + s[1]},
		})

		go Tick(lobbyId)

		return lobbyId

	} else if s[0] == "join" { // join LobbyId Username
		for i := range lobbies {
			v := &lobbies[i]
			if v.LobbyId == s[1] {
				v.Players = append(v.Players, "id="+s[2]+"hp=0 x=0 y=0 z=0")
			}
		}

	} else if s[0] == "leave" { // leave LobbyId Username
		for i := range lobbies {
			v := &lobbies[i]
			if v.LobbyId == s[1] {
				newPlayers := []string{}
				isHost := false

				for _, p := range v.Players {
					id := GetParam(p, "id")
					if id == s[2] {
						if id == v.Host {
							isHost = true
						}
						continue
					}
					newPlayers = append(newPlayers, p)
				}

				if isHost {
					if len(newPlayers) > 0 {
						v.Host = GetParam(newPlayers[0], "id")
					} else {

						newLobbies := []Lobby{}

						for _, lobby := range lobbies {
							if lobby.LobbyId != v.LobbyId {
								newLobbies = append(newLobbies, lobby)
							}
						}

						lobbies = newLobbies
					}
				}

				v.Players = newPlayers
			}
		}
	} else if s[0] == "gets" { // gets LobbyId
		for i, v := range lobbies {
			if v.LobbyId == s[1] {
				data, _ := json.MarshalIndent(lobbies[i], "", "  ")
				return string(data)
			}
		}
	} else if s[0] == "psync" { // psync LobbyId Username Health X Y Z
		for i := range lobbies {
			for i2 := range lobbies[i].Players {
				lobbies[i].Players[i2] = "id=" + s[2] + " hp=" + s[3] + " x=" + s[4] + " y=" + s[5] + " z=" + s[6]
			}
		}
	}

	return "nil"
}

func PrintLobbies() {
	data, _ := json.MarshalIndent(lobbies, "", "  ")
	fmt.Println(string(data))
}

func main() {
	http.HandleFunc("/gnh", NetHandler)
	fmt.Println("Listening on localhost:8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		panic(err)
	}
}
