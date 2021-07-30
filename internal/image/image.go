package image

type Dim int
type Shape = []Dim

type Aspect struct {
	Width  Dim
	Height Dim
}

func (s *Shape) Aspect() Aspect {
	return Aspect{
		(*s)[1], (*s)[0],
	}
}
