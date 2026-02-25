package core

type CPU struct {
	v      [16]byte
	I      uint8
	PC     uint16
	SP     uint8
	Stack  [16]uint16
	Memory [4096]byte

	Display  *Display
	Keyboard *Keyboard
}
