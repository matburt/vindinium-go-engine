package vindinium

import "errors"

type Set struct {
	pos   []Position
	count int
}


func NewSet() Set {
	s := new(Set)
	s.count = 0
	s.pos = make([]Position, 255, 1024)
	return *s
}

func (s *Set) Append(p Position) {
	if s.count>=len(s.pos){
		s.pos=s.pos[:cap(s.pos)]
	}
	if s.count>=cap(s.pos){
		t:=make([]Position, cap(s.pos)*2, cap(s.pos)*2)
		copy(t,s.pos)
		s.pos=t
	}
	s.pos[s.count]=p
	s.Increment()
}

func (s *Set) Increment(){
	s.count++
}

func (s *Set) Decrement(){
	s.count--
}

func (s *Set) Remove(p Position) error {
	i, e := s.Index(p)
	if e != nil {
		return e
	}
	if i < len(s.pos) {
		t:=s.pos[i+1:]
		s.pos = s.pos[0:i]
		for _, aPos:= range t{
			s.pos=append(s.pos, aPos)
		}
		s.Decrement()
		return nil
	}

	s.pos = s.pos[0:i]
	s.Decrement()
	return nil
}

func (s *Set) Pop() (p Position, e error) {
	if len(s.pos) == 0 {
		return *new(Position), errors.New("Pop on empty set failed")
	}
	p = s.pos[0]
	if len(s.pos) > 1 {
		s.pos = s.pos[1:]
		s.Decrement()
	} else {
		s := new(Set)
		s.count = 0
	}
	return p, nil
}

func (s *Set) Index(p Position) (i int, e error) {
	for i, aPos := range s.pos {
		if aPos.X == p.X && aPos.Y == p.Y {
			return i, nil
		}
		if i>=s.count-1{
			break
		}
	}
	return -1, errors.New("Position not found")
}

func (s Set) Find(p Position) bool {
	_, e := s.Index(p)
	if e==nil{
		return true
	}
	return false
}
