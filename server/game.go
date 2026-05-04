package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Rules of the Game:

// Every turn a player may place a pixel anywhere on the board, simultaniously with all other players
// If two players chose the same pixel, randomize the player who gets it
// If a player surrounds an area, where previously werent any pixel, every pixel in the area gets claimed for the player
// If a player makes an illegal move, ignore the turn

// Bonus Rules:
// A player only can place a pixel on an adjacent pixel of his own

type Game struct {
	Id      int
	Board   *Board
	Round   int
	players []Player
	folder  string
}

func NewGame(width, height int, players []Player) Game {

	if players == nil {
		fmt.Println("no Players configured")
	}

	fields := make([][]int, width)
	for i := range fields {
		fields[i] = make([]int, height)
	}

	board := Board{
		fields: fields,
		Width:  width,
		Height: height,
	}

	folder, _ := CreateUniqueFolder("replays/replay")
	fmt.Println(folder)

	game := Game{
		Board:   &board,
		Round:   0,
		players: players,
		folder:  folder,
	}

	return game
}

func (g *Game) start() {
	g.Init()
	for {
		fmt.Println("Round Started")
		g.PlayRound()
		fmt.Println("Round Ended")
	}
}

func (g *Game) Init() {
	// Wait for Players
	g.waitForPlayers()

	//Create Start Positions
	start_positions := g.getStartPositions()

	// Apply changes to board
	g.Board.UpdateBoard(start_positions)

	// Save board
	path := fmt.Sprintf("%s/%09d.png", g.folder, g.Round)
	SaveAsPng(path, g)

	//Send Players Start Position by
	g.updatePlayers(start_positions)
	g.Round++
}

func (g *Game) PlayRound() {
	// Get player turns
	fmt.Println("Read Moves")
	moves := g.getMoves()
	fmt.Println(moves)
	// Remove invalid moves
	fmt.Println("Validate Moves")
	valid_moves := g.validateMoves(moves)

	// Apply changes to board
	fmt.Println("Apply Changes to Board")
	g.Board.UpdateBoard(valid_moves)

	// Fill Surrounded Areas
	g.Board.UpdateTerritory(valid_moves)

	// Save board
	fmt.Println("Save Changes to Png")
	path := fmt.Sprintf("%s/%010d.png", g.folder, g.Round)
	SaveAsPng(path, g)

	// Send Turns to Players
	fmt.Println("Send Update to Players")
	g.updatePlayers(valid_moves)
	g.Round++
}

func (g *Game) waitForPlayers() {
	fmt.Println("Waiting for Players...")
	is_waiting := true
	for is_waiting {
		is_waiting = false
		for i := range g.players {
			if !g.players[i].IsOnline() {
				is_waiting = true
			}
		}
		time.Sleep(time.Second)
	}
	fmt.Println("Waiting ended!")
}

func (g *Game) getStartPositions() []BoardUpdate {
	starting_positions := make([]BoardUpdate, len(g.players))
	for i := range g.players {
		starting_positions[i] = BoardUpdate{
			Id: g.players[i].Id,
			X:  rand.Intn(g.Board.Width),
			Y:  rand.Intn(g.Board.Height),
		}
	}
	return starting_positions
}

func (g *Game) updatePlayers(moves []BoardUpdate) {
	for _, p := range g.players {
		p.SendBoardUpdates(moves)
	}
}

func (g *Game) validateMoves(moves []BoardUpdate) []BoardUpdate {
	b := g.Board

	// shuffle the order of the moves
	// makes the game less predictable when two players claim the same field
	rand.Shuffle(len(moves), func(i, j int) {
		moves[i], moves[j] = moves[j], moves[i]
	})

	claimed := make(map[[2]int]bool)
	valid_moves := make([]BoardUpdate, len(moves))
	count := 0
	for m := range moves {
		id := moves[m].Id
		x := moves[m].X
		y := moves[m].Y

		// move may not be out of bounds
		if x < 0 || y < 0 || x >= b.Width || y >= b.Height {
			continue
		}

		// if field is already claimed
		if b.GetField(x, y) != 0 {
			continue
		}

		// player must own a tile next to the tile being claimed
		below, above, left, right := b.GetAdjacentValues(x, y)
		if below != id && above != id && left != id && right != id {
			continue
		}

		// if field has been claimed in a previous turn, ignore this claim
		pos := [2]int{x, y}
		if claimed[pos] {
			continue
		}
		claimed[pos] = true

		valid_moves[count] = moves[m]
		count++
	}

	return valid_moves[:count]
}

func (g *Game) getMoves() []BoardUpdate {
	moves := make([]BoardUpdate, len(g.players))
	var wg sync.WaitGroup
	for i := range g.players {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p := g.players[i]
			moves[i] = p.receiveTurn()
			moves[i].Id = p.Id
			fmt.Println(moves[i])
		}(i)
	}

	wg.Wait()
	return moves
}
