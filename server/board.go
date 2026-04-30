package main

import "fmt"

type Board struct {
	fields [][]int
	Width  int
	Height int
}

type BoardUpdate struct {
	Id int //new value of the tile
	X  int //x pos of the updated tile
	Y  int //y pos of the updated tile
}

func (b *Board) GetField(x, y int) int {
	return b.fields[x][y]
}

func (b *Board) SetField(x, y, id int) {
	b.fields[y][x] = id
}

func (b *Board) UpdateField(x, y, id int) {
	b.fields[x][y] = id
}

func (b *Board) UpdateBoard(moves []BoardUpdate) {
	for i := range moves {
		m := moves[i]
		fmt.Println("Update Field")
		b.UpdateField(m.X, m.Y, m.Id)

		fmt.Println("Fill Surrounded Areas")
		b.FillSurroundedAreas(m)
	}
}

func (b *Board) FillSurroundedAreas(move BoardUpdate) {
	id := move.Id
	x, y := move.X, move.Y

	val_below, val_above, val_left, val_right := b.GetAdjacentValues(x, y)

	adjacent_player_tiles := 0
	if val_below == move.Id {
		adjacent_player_tiles++
	}

	if val_above == move.Id {
		adjacent_player_tiles++
	}

	if val_left == move.Id {
		adjacent_player_tiles++
	}

	if val_right == move.Id {
		adjacent_player_tiles++
	}

	below := []int{x, y - 1}
	above := []int{x, y + 1}
	left := []int{x - 1, y}
	right := []int{x + 1, y}

	if adjacent_player_tiles >= 2 {
		fmt.Println("Playertiles >= 2")
		fmt.Println("Check Below")
		below_hit_other_players, below_area := b.floodFill(below[0], below[1], id)
		if !below_hit_other_players {
			fmt.Println("Apply Below")
			b.ApplyFill(below_area, id)
		}

		fmt.Println("Check Above")
		above_hit_other_players, above_area := b.floodFill(above[0], above[1], id)
		if !above_hit_other_players {
			fmt.Println("Apply Above")
			b.ApplyFill(above_area, id)
		}

		fmt.Println("Check Left")
		left_hit_other_players, left_area := b.floodFill(left[0], left[1], id)
		if !left_hit_other_players {
			fmt.Println("Apply Left")
			b.ApplyFill(left_area, id)
		}

		fmt.Println("Check Right")
		right_hit_other_players, right_area := b.floodFill(right[0], right[1], id)
		if !right_hit_other_players {
			fmt.Println("Apply Right")
			b.ApplyFill(right_area, id)
		}
	}
}

func (b *Board) ApplyFill(fill [][]bool, value int) {
	for y := range fill {
		for x := range fill[y] {
			if fill[y][x] {
				b.fields[y][x] = value
			}
		}
	}
}

// returns the values of the tiles below, above, left, right from x,y
func (b *Board) GetAdjacentValues(x, y int) (int, int, int, int) {

	below := -1
	above := -1
	left := -1
	right := -1

	if y-1 >= 0 {
		below = b.fields[x][y-1]
	}

	if y+1 < game.Board.Height {
		above = b.fields[x][y+1]
	}

	if x-1 >= 0 {
		left = b.fields[x-1][y]
	}

	if x+1 < game.Board.Width {
		right = b.fields[x+1][y]
	}

	return below, above, left, right
}

func (b *Board) floodFill(x, y, playerID int) (bool, [][]bool) {
	visited := make([][]bool, b.Width)
	for i := range visited {
		visited[i] = make([]bool, b.Height)
	}

	stack := [][]int{{x, y}}
	for len(stack) > 0 {

		field := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		x, y := field[0], field[1]

		if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
			continue
		}

		value := b.GetField(x, y)
		if value != 0 && value != playerID {
			return true, nil
		}

		if visited[x][y] || b.fields[x][y] == playerID {
			continue
		}

		visited[x][y] = true

		below := []int{x, y - 1}
		above := []int{x, y + 1}
		left := []int{x - 1, y}
		right := []int{x + 1, y}

		stack = append(stack, below, above, left, right)
	}

	return false, visited
}
