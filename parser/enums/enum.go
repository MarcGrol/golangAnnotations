package enums

// @Enum()
type ColorType int

const (
	Red ColorType = iota
	Green
	Blue
)

// @Enum()
type Profession string

const (
	Teacher Profession = "_teacher"
	Cleaner Profession = "_cleaner"
)
