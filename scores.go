package vindinium

type Scores struct {
	gScore map[Position]int
	fScore map[Position]int
}

func NewScores() Scores{
	s:=new(Scores)
	s.gScore=make(map[Position]int)
	s.fScore=make(map[Position]int)
	return *s
}

func (s Scores) LowestF(openSet Set) Position {
	var lowestPos Position
	lowestScore := 10000000
	for pos, score := range s.fScore {
		if score < lowestScore && openSet.Find(pos){
			lowestPos = pos
			lowestScore = score
		}
	}
	return lowestPos
}
