package vindinium

type Path struct {
	Steps map[Position]Position
}

func NewPath() Path {
	p := new(Path)
	p.Steps = make(map[Position]Position, 20)
	return *p
}

func (p Path) reconstructed(start Position) *Set {
	if _, found := p.Steps[start]; found {
		newSet := p.reconstructed(p.Steps[start])
		newSet.Append(start)
		return newSet
	}
	returnSet := NewSet()
	return &returnSet
}
