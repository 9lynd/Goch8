package renderer

import (
	"github.com/9lynd/Goch8/core"
)

// Renderer is the interface any display backend must satisfy.
type Renderer interface {
	Run(cpu *core.CPU) error
}
