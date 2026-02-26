package core

type Keyboard struct {
	Keys [16]bool
}

func (k *Keyboard) SetKey(key byte, pressed bool) {
	if key <= 0xF {
		k.Keys[key] = pressed
	}
}

func (k *Keyboard) WaitForKey() (bool, byte) {
	for i := byte(0); i <= 0xF; i++ {
		if k.Keys[i] {
			return true, i
		}
	}
	return false, 0
}

func (k *Keyboard) IsPressed(key byte) bool {
	if key > 0xF {
		return false
	}
	return k.Keys[key]
}
