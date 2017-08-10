// +build !ci

package structExample

//go:generate golangAnnotations -input-dir .

// @JsonEnum( base = "Color", unknown = "Unknown", stripped = "true", tolerant = "true" )
type Color int

const (
	ColorUnknown Color = iota
	ColorRed
	ColorGreen
	ColorBlue
)

// @JsonEnum( )
type Word int

const (
	Aap Word = iota
	Noot
	Mies
)
