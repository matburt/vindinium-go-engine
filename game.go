package vindinium
import "strconv"
const (
	PLAYER1 = iota
	PLAYER2
	PLAYER3
	PLAYER4
)

type Game struct {
	State       *State
	Board       *Board  `json:"board"`
	Heroes      []*Hero `json:"heroes"`
	Id          string  `json:"id"`
	Finished    bool    `json:"finished"`
	Turn        int     `json:"turn"`
	MaxTurns    int     `json:"maxTurns"`
	MinesLocs   map[Position]string
	HeroesLocs  map[Position]int
	TavernsLocs map[Position]Tavern
}

type Tavern struct {
}

func (game *Game) FetchHero(id int) *Hero {
	for _, hero := range game.Heroes{
		if hero.Id == id {
			return hero
		}
	}
	return nil
}

func NewGame(state *State) (game *Game) {
	game = state.Game
	game.State = state
	game.Board.parseTiles()
	game.MinesLocs=make(map[Position]string,16)
	game.HeroesLocs=make(map[Position]int,3)
	game.TavernsLocs=make(map[Position]Tavern,16)
	for x, xval := range game.Board.Tileset {
		for y, tile := range xval {
			pos := Position{x,y}
			isTavern := game.Board.IsTavern(tile)
			mine, isMine := tile.(*MineTile)
			hero, isHero := tile.(*HeroTile)
			switch {
			case isTavern:
				game.TavernsLocs[pos]=Tavern{}
			case isMine:				
				id, nobody := strconv.ParseInt(mine.HeroId, 10, 0)
				if nobody!=nil || int(id) != game.State.Hero.Id {					
					game.MinesLocs[pos]=mine.HeroId
				}
			case isHero:
				if hero.Id != game.State.Hero.Id {
					game.HeroesLocs[pos]=hero.Id
				}
			}
		}
	}
	return
}
