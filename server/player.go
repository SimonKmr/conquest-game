package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
)

type Player struct {
	Id         int
	Password   string
	connection *websocket.Conn
}

func (p *Player) receiveTurn() BoardUpdate {
	mt, message, err := p.connection.ReadMessage()
	if mt != websocket.TextMessage {
		// error
	}

	if err != nil {
		// error
	}

	var move BoardUpdate
	err = json.Unmarshal(message, &move)
	fmt.Println(err)
	return move
}

func (p *Player) SendBoardUpdates(updates []BoardUpdate) {
	json, _ := json.Marshal(updates)
	p.connection.WriteMessage(websocket.TextMessage, json)
}

func (p *Player) IsOnline() bool {
	return p.connection != nil
}

func (p *Player) CloseConnection() {
	p.connection.Close()
}

func GetPlayersFromFile(path string) []Player {

	bytes, _ := os.ReadFile(path)
	var players []Player
	json.Unmarshal(bytes, &players)

	return players
}
