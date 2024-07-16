//go:build js && wasm

package maze

type mazeRunner struct {
	pos                  coord
	emoji, winning_emoji rune
}

func newMazeRunner(c coord) mazeRunner {
	return mazeRunner{pos: c, emoji: '🙂', winning_emoji: '😀'}
}

func (r *mazeRunner) move(d direction) {
	switch d {
	case _LEFT:
		r.pos.y -= 1
	case _RIGHT:
		r.pos.y += 1
	case _UP:
		r.pos.x -= 1
	case _DOWN:
		r.pos.x += 1
	}
}
