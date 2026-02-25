package core

const (
	Width  = 64
	Height = 32
)

type Display struct {
	Pixels [Width * Height]bool
	Dirty  bool
}
