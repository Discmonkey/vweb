package image

type Shape struct {
	s []int
}

type Aspect struct {
	Width  int
	Height int
}

func (s *Shape) Aspect() Aspect {
	return Aspect{
		s.s[1], s.s[0],
	}
}
