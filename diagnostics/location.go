package diagnostics

type Location interface {
	LineRange() LineRange
	String() string
}

type LineRange interface {
	StartLine() LinePosition
	EndLine() LinePosition
}

type LinePosition interface {
	Line() int
	Offset() int
}

type locationImpl struct {
	lineRange LineRange
}

func NewLocation(lineRange LineRange) Location {
	return &locationImpl{lineRange: lineRange}
}

func (l *locationImpl) LineRange() LineRange {
	return l.lineRange
}

func (l *locationImpl) String() string {
	return ""
}
