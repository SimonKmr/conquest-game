package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var game Game

func main() {
	fmt.Println("[Create Game]")
	players := GetPlayersFromFile("players.json")
	game = NewGame(1920/8, 1080/8, players)

	fmt.Println("[Configur Websockets]")
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/player", player)
	http.HandleFunc("/spectator", echo)
	http.HandleFunc("/echo", player)
	http.HandleFunc("/", home)
	go http.ListenAndServe(*addr, nil)

	fmt.Println("[Start Game]")
	game.start()
}

func player(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	player := VerifyPlayer(c)

	if player == nil {
		c.WriteMessage(websocket.TextMessage, []byte("Authentication Failed"))
		c.Close()
		return
	}

	player.connection = c
}

func VerifyPlayer(c *websocket.Conn) *Player {
	// if ID and Password match one of the players
	// match connection to player
	// else send player message about non matching credentails
	mt, message, err := c.ReadMessage()

	if mt != websocket.TextMessage {
		fmt.Println("not a text message")
		return nil
	}

	if err != nil {
		fmt.Println("error while receiving message")
		return nil
	}

	message_str := string(message)
	credentials := strings.Split(message_str, " ")
	id, err := strconv.Atoi(credentials[0])

	if err != nil {
		fmt.Println("string conversion error")
		return nil
	}

	player := find(game.players, id)

	if player == nil {
		fmt.Println("no player found")
		return nil
	}

	if player.Password != credentials[1] {
		fmt.Println("password incorrect")
		return nil
	}

	fmt.Printf("Player: %d joined!\n", player.Id)
	return player
}

func find(players []Player, id int) *Player {
	for i := range players {
		if players[i].Id == id {
			return &players[i]
		}
	}
	return nil
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	VerifyPlayer(c)
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
