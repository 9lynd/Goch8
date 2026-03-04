// Key mapping: CHIP-8's 16-key hex pad is mapped to a standard QWERTY layout:

//   CHIP-8:
// 1 2 3 C
// 4 5 6 D
// 7 8 9 E
// A 0 B F

// Keyboard:
// 1 2 3 4
// Q W E R
// A S D F
// Z X C V

package ebiten

import (
	"image/color"

	"github.com/9lynd/Goch8/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	Scale          = 10
	CyclesPerFrame = 10
)

// keyMap maps ebiten keys to CHIP-8 key indices.
var keyMap = map[ebiten.Key]byte{
	ebiten.Key1: 0x1, ebiten.Key2: 0x2, ebiten.Key3: 0x3, ebiten.Key4: 0xC,
	ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0xD,
	ebiten.KeyA: 0x7, ebiten.KeyS: 0x8, ebiten.KeyD: 0x9, ebiten.KeyF: 0xE,
	ebiten.KeyZ: 0xA, ebiten.KeyX: 0x0, ebiten.KeyC: 0xB, ebiten.KeyV: 0xF,
}

type Game struct {
	cpu        *core.CPU
	backbuffer *ebiten.Image
}

func (g *Game) Update() error {
	g.handleInput()

	for i := 0; i < CyclesPerFrame; i++ {
		if err := g.cpu.Cycle(); err != nil {
			return err
		}
	}

	g.cpu.TickTimers()
	return nil
}

// The backbuffer is only repainted when Display.Dirty is true.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.cpu.Display.Dirty {
		g.backbuffer.Fill(color.Black)
		for y := 0; y < core.Height; y++ {
			for x := 0; x < core.Width; x++ {
				if g.cpu.Display.Pixels[y*core.Width+x] {
					drawPixel(g.backbuffer, x, y)
				}
			}
		}
		g.cpu.Display.Dirty = false
	}

	screen.DrawImage(g.backbuffer, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return core.Width * Scale, core.Height * Scale
}

func drawPixel(screen *ebiten.Image, x, y int) {
	px := x * Scale
	py := y * Scale
	for dy := 0; dy < Scale; dy++ {
		for dx := 0; dx < Scale; dx++ {
			screen.Set(px+dx, py+dy, color.White)
		}
	}
}

func (g *Game) handleInput() {
	for ebitenKey, chip8Key := range keyMap {
		if inpututil.IsKeyJustPressed(ebitenKey) {
			g.cpu.Keyboard.SetKey(chip8Key, true)
		}
		if inpututil.IsKeyJustReleased(ebitenKey) {
			g.cpu.Keyboard.SetKey(chip8Key, false)
		}
	}
}

type Renderer struct{}

func (r *Renderer) Run(cpu *core.CPU) error {
	ebiten.SetWindowSize(core.Width*Scale, core.Height*Scale)
	ebiten.SetWindowTitle("Goch8")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := &Game{
		cpu:        cpu,
		backbuffer: ebiten.NewImage(core.Width*Scale, core.Height*Scale),
	}
	return ebiten.RunGame(game)
}
