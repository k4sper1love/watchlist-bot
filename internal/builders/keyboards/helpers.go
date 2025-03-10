package keyboards

func (k *Keyboard) AddIf(condition bool, addFunc func(*Keyboard)) *Keyboard {
	if condition {
		addFunc(k)
	}
	return k
}
