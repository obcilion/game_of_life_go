package main

import (
	"github.com/buger/goterm"
	"strconv"
	"time"
)

var liveCells = LiveCells{make(map[Cell]bool)}

var liveCellsToKill = make(map[Cell]bool)
var deadCellsToCheck = make(map[Cell]bool)
var deadCellsToRevive = make(map[Cell]bool)

var xOffset = goterm.Width() / 2
var yOffset = goterm.Height() / 2

type Cell struct {
	X int
	Y int
}

type LiveCells struct {
	Cells map[Cell]bool
}

func (liveCells *LiveCells) killCell(cellToKill Cell) {
	delete(liveCells.Cells, cellToKill)
}

func (liveCells *LiveCells) reviveCell(cellToRevive Cell) {
	liveCells.Cells[cellToRevive] = true
}

func (liveCells *LiveCells) cellIsAlive(cellToCheck Cell) bool {
	if _, present := liveCells.Cells[cellToCheck]; present {
		return true
	}

	return false
}

// Seed function
func rPentomino() map[Cell]bool {
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
func blinker() map[Cell]bool {
	seed := make(map[Cell]bool)

	seed[Cell{1, 0}] = true
	seed[Cell{0, 0}] = true
	seed[Cell{-1, 0}] = true

	return seed
}

// Seed function
func acorn() map[Cell]bool {
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

// Seed function
func glider() map[Cell]bool {
	seed := make(map[Cell]bool)

	// First row
	seed[Cell{0, -1}] = true

	// Second row
	seed[Cell{1, 0}] = true

	// Third row
	seed[Cell{-1, 1}] = true
	seed[Cell{0, 1}] = true
	seed[Cell{1, 1}] = true

	return seed
}

func seed() map[Cell]bool {

	return acorn()
}

// Adds dead neighbours to deadCellsToCheck and returns number of live neighbours
func countLiveNeighbours(cell Cell, addDeadCellsToList bool) int {
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

	for i := 0; i < len(neighbourCells); i++ {
		currentCell := neighbourCells[i]
		if liveCells.cellIsAlive(currentCell) {
			count++
		} else if addDeadCellsToList {
			deadCellsToCheck[currentCell] = true
		}
	}

	return count
}

func checkLiveCells() {
	for currentCell := range liveCells.Cells {
		liveNeighbours := countLiveNeighbours(currentCell, true)

		if liveNeighbours < 2 || liveNeighbours > 3 {
			liveCellsToKill[currentCell] = true
		}
	}
}

func checkDeadCells() {
	for celllToCheck := range deadCellsToCheck {
		liveNeighbours := countLiveNeighbours(celllToCheck, false)

		if liveNeighbours == 3 {
			deadCellsToRevive[celllToCheck] = true
		}
	}
}

func killLiveCellsInList() {
	for cellToKill := range liveCellsToKill {
		liveCells.killCell(cellToKill)
	}
}

func reviveDeadCellsInList() {
	for cellToRevive := range deadCellsToRevive {
		liveCells.reviveCell(cellToRevive)
	}
}

func calculateCellChange() {

	// Iterate through all live cells
	// 	 count live neighbours, add cell to list of cells to kill if less than two or more than three live neighbours
	//   add each dead neighbour cell to list of cells to check
	checkLiveCells()

	// Iterate through all dead cells added in previous step
	//   add cell to list of cells to revive if exactly three live neighbours
	// 	 do NOT add neighbours to list of cells to check, this will result in an infinite loop
	checkDeadCells()

	// Iterate through list of cells to kill and kill them
	killLiveCellsInList()

	// Iterate through list of cells to revive and revive them
	reviveDeadCellsInList()

	// Reset all lists except live cells list
	liveCellsToKill = make(map[Cell]bool)
	deadCellsToCheck = make(map[Cell]bool)
	deadCellsToRevive = make(map[Cell]bool)
}

func drawCells() {

	for currentCell := range liveCells.Cells {
		// Don't draw cells outside the game board
		if currentCell.X+xOffset > goterm.Width() ||
			currentCell.X+xOffset < 0 ||
			currentCell.Y+yOffset > goterm.Height() ||
			currentCell.Y+yOffset < 0 {
			continue
		}

		goterm.MoveCursor(currentCell.X+xOffset, currentCell.Y+yOffset)
		goterm.Print("X")
	}
}

func drawText(text string) {
	goterm.MoveCursor(0, 0)
	goterm.Print(text)
}

func main() {
	liveCells.Cells = seed()
	generation := 0

	for {
		goterm.Clear()
		drawText(strconv.Itoa(generation))
		drawCells()
		goterm.Flush()

		calculateCellChange()
		generation++
		time.Sleep(33 * time.Millisecond)
	}
}
