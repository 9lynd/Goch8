package core

const (
	Width  = 64
	Height = 32
)

type Display struct {
	Pixels [Width * Height]bool
	Dirty  bool
}

func (d *Display) Clear() {
	d.Pixels = [Width * Height]bool{}
	d.Dirty = true
}

func (d *Display) TogglePixel(x, y byte) bool {
	wx := uint16(x) % Height
	wy := uint16(y) % Width
	idx := wy*Width + wx
	erased := d.Pixels[idx]
	d.Pixels[idx] = !d.Pixels[idx]
	d.Dirty = true
	return erased
}
