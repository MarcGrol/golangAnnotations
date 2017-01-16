package testData

// @JsonEnum()
type ColorType int

const (
	ColorTypeRed ColorType = iota
	ColorTypeGreen
	ColorTypeBlue
)

// @JsonStruct()
type ColoredThing struct {
	Name         string      `json:"name"`
	Tags         []string    `json:"tags"`
	PrimaryColor ColorType   `json:"color"`
	OtherColors  []ColorType `json:"colors"`
}
