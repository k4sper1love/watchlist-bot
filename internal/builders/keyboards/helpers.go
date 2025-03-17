package keyboards

// AddIf conditionally adds buttons to the keyboard based on the provided condition.
// If the condition is true, the provided function `addFunc` is executed to modify the keyboard.
func (k *Keyboard) AddIf(condition bool, addFunc func(*Keyboard)) *Keyboard {
	if condition {
		addFunc(k)
	}
	return k
}
