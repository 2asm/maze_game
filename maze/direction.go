//go:build js && wasm

package maze

type direction int

const (
	_LEFT direction = iota + 1
	_UP
	_RIGHT
	_DOWN
)

func (d direction) opposite() direction {
	switch d {
	case _LEFT:
		return _RIGHT
	case _RIGHT:
		return _LEFT
	case _UP:
		return _DOWN
	case _DOWN:
		return _UP
	default:
		return 0
	}
}

func (d direction) String() string {
	switch d {
	case _LEFT:
		return "←"
	case _RIGHT:
		return "→"
	case _UP:
		return "↑"
	case _DOWN:
		return "↓"
	}
	return "•"
}
