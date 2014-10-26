package vindinium

import (
	"errors"
	"math"
	"strconv"
)

const (
	WALL = iota - 2
	AIR
	TAVERN

	AIR_TILE    = " "
	WALL_TILE   = "#"
	TAVERN_TILE = "["
	MINE_TILE   = "$"
	HERO_TILE   = "@"
)

var (
	AIM = map[Direction]*Position{
		"North": &Position{-1, 0},
		"East":  &Position{0, 1},
		"South": &Position{1, 0},
		"West":  &Position{0, -1},
	}
)

type Direction string

type Board struct {
	Size    int    `json:"size"`
	Tiles   string `json:"tiles"`
	Tileset [][]interface{}
}

type Position struct {
	X, Y int
}

func NewPosition(x, y int) Position {
	var pos Position
	pos.X = x
	pos.Y = y
	return pos
}

func tileToInt(tiles string, index int) int {
	tile := []rune(tiles)[index]
	str, _ := strconv.Atoi(string(tile))

	return str
}

func (board *Board) parseTile(tile string) interface{} {
	switch string([]rune(tile)[0]) {
	case AIR_TILE:
		return AIR
	case WALL_TILE:
		return WALL
	case TAVERN_TILE:
		return TAVERN
	case MINE_TILE:
		id := string([]rune(tile)[1])
		return &MineTile{id}
	case HERO_TILE:
		char := string([]rune(tile)[1])
		id, _ := strconv.Atoi(char)
		return &HeroTile{id}
	default:
		return -3
	}
	return -3
}

func (board *Board) parseTiles() {
	var vector [][]rune
	var matrix [][][]rune
	ts := make([][]interface{}, board.Size)

	for i := 0; i <= len(board.Tiles)-2; i = i + 2 {
		vector = append(vector, []rune(board.Tiles)[i:i+2])
	}

	for i := 0; i < len(vector); i = i + board.Size {
		matrix = append(matrix, vector[i:i+board.Size])
	}

	for xi, x := range matrix {
		innerList := make([]interface{}, board.Size)
		for xsi, xs := range x {

			innerList[xsi] = board.parseTile(string(xs))
		}
		ts[xi] = innerList
	}

	board.Tileset = ts
}

func (board *Board) Passable(loc Position) bool {
	if loc.X>=board.Size || loc.Y>=board.Size {
		return false
	}
	tile := board.Tileset[loc.X][loc.Y]
	return tile != WALL && tile != TAVERN && ! board.IsMine(board.Tile(loc))
}

func (board *Board) IsTavern(tile interface{}) bool {
	if tile == TAVERN {
		return true
	}
	return false
}

func (board *Board) IsMine(tile interface{}) bool {
	_, ok := tile.(*MineTile)
	return ok
}

func (board *Board) IsHero(tile interface{}) bool {
	_, ok := tile.(*HeroTile)
	return ok	
}

func (board *Board) Tile(pos Position) interface{} {
	if pos.X >= board.Size || pos.Y >= board.Size {
		return nil
	}
	return board.Tileset[pos.X][pos.Y]
}

func (board *Board) GetDirection(startLoc Position, endLoc Position) Direction {
        if startLoc.X < endLoc.X {
                return "South"
        }
        if startLoc.X > endLoc.X {
                return "North"
        }
        if startLoc.Y < endLoc.Y {
                return "East"
        }
        if startLoc.Y > endLoc.Y {
                return "West"
        }
	//if the two locations are the same, then stay
        return "Stay"
}

func (board *Board) ManhattanDistance(a Position, b Position) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}

func (board *Board) Neighbors(pos Position) []Position {
	neighbors := make([]Position, 4, 4)
	count:=0
	addNeighbor := func(n Position){
		if count>len(neighbors){
			neighbors=neighbors[:count]
		}
		neighbors[count]=n
		count++
	}
	if n := NewPosition(pos.X-1, pos.Y); pos.X > 0  {
		addNeighbor(n)
	}
	if n := NewPosition(pos.X, pos.Y-1); pos.Y > 0  {
		addNeighbor(n)
	}
	if n := NewPosition(pos.X+1, pos.Y); pos.X < board.Size  {
		addNeighbor(n)
	}
	if n := NewPosition(pos.X, pos.Y+1); pos.Y < board.Size  {
		addNeighbor(n)
	}
	if count==4{
		return neighbors
	}
	return neighbors[:count]
}

//Find the shortest path between A and B
func (board *Board) ShortestPath(start Position, end Position) (p Set, e error) {
	closedSet := NewSet()
	openSet := NewSet()
	openSet.Append(start)
	cameFrom := NewPath()
	scores := NewScores()
	scores.gScore[start] = 0
	scores.fScore[start] = board.ManhattanDistance(start, end)
	for openSet.count > 0 {
		current := scores.LowestF(openSet)
		if current.X == end.X && current.Y == end.Y  {
			returnPath := cameFrom.reconstructed(end)
			return *returnPath, nil
		}
		openSet.Remove(current)
		closedSet.Append(current)
		for _,n := range board.Neighbors(current) {
			if closedSet.Find(n) {
				continue
			}
			if !(n.X==end.X && n.Y==end.Y) && !board.Passable(n){
				continue
			}
			gScore := scores.gScore[current] + 1
			if !openSet.Find(n) || gScore < scores.gScore[n] {
				cameFrom.Steps[n] = current
				scores.gScore[n] = gScore
				scores.fScore[n] = gScore + board.ManhattanDistance(n, end)
				if !openSet.Find(n) {
					openSet.Append(n)
				}
			}
		}
	}
	return *new(Set), errors.New("Failed to find path")
}

func (board *Board) To(loc Position, direction Direction) *Position {
	row := loc.X
	col := loc.Y
	dLoc := AIM[direction]
	nRow := row + dLoc.X
	if nRow < 0 {
		nRow = 0
	}
	if nRow > board.Size {
		nRow = board.Size
	}
	nCol := col + dLoc.Y
	if nCol < 0 {
		nCol = 0
	}
	if nCol > board.Size {
		nCol = board.Size
	}

	return &Position{nRow, nCol}
}
