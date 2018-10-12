package game_of_life_go

import (
	"github.com/buger/goterm"
	"time"
)

var liveCells = make(map[Cell]bool)

var liveCellsToKill = make(map[Cell]bool)
var deadCellsToCheck = make(map[Cell]bool)
var deadCellsToRevive = make(map[Cell]bool)

var xOffset = goterm.Width()/2
var yOffset = goterm.Height()/2

const DEBUG = true

type Cell struct {
	X int
	Y int
}

// Checks if a given cell is in the liveCells list
func  (cell *Cell) IsAlive() bool {
	if _, present := liveCells[*cell]; present {
		return true
	} else {
		return false
	}
}

// Seed function
func RPentomino() map[Cell]bool{
	seed := make(map[Cell]bool)

	// Top row
	seed[Cell{0, -1}] = true
	seed[Cell{1, -1}] = true

	// Middle row
	seed[Cell{-1, 0}] = true
	seed[Cell{0, 0}] = true

	// Bottom row
	seed[Cell{0, 1}] = true

	return seed
}

// Seed function
func Blinker() map[Cell]bool {
	seed := make(map[Cell]bool)

	seed[Cell{1, 0}] = true
	seed[Cell{0, 0}] = true
	seed[Cell{-1, 0}] = true

	return seed
}

// Seed function
func Acorn() map[Cell]bool {
	seed := make(map[Cell]bool)

	// First row
	seed[Cell{1, 0}] = true

	// Second row
	seed[Cell{3, 1}] = true

	// Third row
	seed[Cell{0, 2}] = true
	seed[Cell{1, 2}] = true
	seed[Cell{4, 2}] = true
	seed[Cell{5, 2}] = true
	seed[Cell{6, 2}] = true

	return seed
}

func Seed() map[Cell]bool {

	return Acorn()
}

// Adds dead neighbours to deadCellsToCheck and returns number of live neighbours
func CountLiveNeighbours(cell Cell, addDeadCellsToList bool) (int) {
	var count int

	// Always eight neighbours, hardcode coordinates
	neighbourCells := [8]Cell{
		// Top row
		{cell.X - 1, cell.Y - 1},
		{cell.X, cell.Y - 1},
		{cell.X + 1, cell.Y - 1},

		// Middle row
		{cell.X - 1, cell.Y},
		{cell.X + 1, cell.Y},

		// Bottom row
		{cell.X - 1, cell.Y + 1},
		{cell.X, cell.Y + 1},
		{cell.X + 1, cell.Y + 1}}

	for i := 0; i < len(neighbourCells) ;i++  {
		currentCell := neighbourCells[i]
		if currentCell.IsAlive() {
			count++
		} else if addDeadCellsToList {
			deadCellsToCheck[currentCell] = true
		}
	}

	return count
}

func CheckLiveCells() {
	for currentCell := range liveCells  {
		liveNeighbours := CountLiveNeighbours(currentCell, true)

		if liveNeighbours < 2 || liveNeighbours > 3 {
			liveCellsToKill[currentCell] = true
		}
	}
}

func CheckDeadCells() {
	for currentCell := range deadCellsToCheck{
		liveNeighbours := CountLiveNeighbours(currentCell, false)

		if liveNeighbours == 3 {
			deadCellsToRevive[currentCell] = true
		}
	}
}


func KillLiveCellsInList() {
	for currentCell := range liveCellsToKill  {
		delete(liveCells, currentCell)
	}
}

func ReviveDeadCellsInList() {
	for currentCell := range deadCellsToRevive {
		liveCells[currentCell] = true
	}
}

func CalculateCellChange() {

	// Iterate through all live cells
	// 	 count live neighbours, add cell to list of cells to kill if less than two or more than three live neighbours
	//   add each dead neighbour cell to list of cells to check
	CheckLiveCells()

	// Iterate through all dead cells added in previous step
	//   add cell to list of cells to revive if exactly three live neighbours
	// 	 do NOT add neighbours to list of cells to check, this will result in an infinite loop
	CheckDeadCells()

	// Iterate through list of cells to kill and kill them
	KillLiveCellsInList()

	// Iterate through list of cells to revive and revive them
	ReviveDeadCellsInList()

	// Reset all lists except live cells list
	liveCellsToKill = make(map[Cell]bool)
	deadCellsToCheck = make(map[Cell]bool)
	deadCellsToRevive = make(map[Cell]bool)
}

func DrawCells() {
	goterm.Clear()

	for currentCell := range liveCells {
		if currentCell.X + xOffset > goterm.Width() || currentCell.X + xOffset < 0 || currentCell.Y + yOffset > goterm.Height() || currentCell.Y + yOffset < 0 {
			continue
		}

		goterm.MoveCursor(currentCell.X + xOffset, currentCell.Y + yOffset)
		goterm.Print("X")
	}

	goterm.Flush()
}

func main() {
	liveCells = Seed()

	for {
		DrawCells()
		CalculateCellChange()
		time.Sleep(100 * time.Millisecond)
	}
}
