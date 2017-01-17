// +build !ci

package structExample

//go:generate golangAnnotations -input-dir .

// @JsonEnum()
type Color int

const (
	Red Color = iota
	Green
	Blue
)
