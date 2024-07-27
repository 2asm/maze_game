//go:build js && wasm

package maze

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"
)

type maze struct {
	height, width int
	scale         int // pixel size
	cell          [][][]direction
	runner        mazeRunner
	isOver        bool
}

func (m maze) outOfMaze(c coord) bool {
	return c.x < 0 || c.x >= m.height || c.y < 0 || c.y >= m.width
}

func NewMaze(height, width, scale int) *maze {
	_maze := make([][][]direction, height)
	for i := range height {
		_maze[i] = make([][]direction, width)
	}
	ret := &maze{
		height: height,
		width:  width,
		scale:  scale,
		cell:   _maze,
		runner: newMazeRunner(coord{0, 0}),
	}
	ret.fillMaze()
	ret.init()
	ret.setAllBorders()
	ret.fillTextCell(0, 0, ret.runner.emoji)
	ret.clearCell(ret.height-1, ret.width-1)
	ret.fillTextCell(ret.height-1, ret.width-1, "💻")
	return ret
}

func (m *maze) Start() {
	go m.listenForArrowKeys()
	for {
		select {
		case <-restartChan:
			m.clearAll()
			m = NewMaze(m.height, m.width, m.scale)
		case d := <-moveChan:
			if m.isOver {
				continue
			}
			c := m.runner.pos
			for _, di := range m.cell[c.x][c.y] {
				if di == d {
					m.clearCell(c.x, c.y)
					m.runner.move(d)
					break
				}
			}
			c = m.runner.pos
			if m.won() {
				m.isOver = true
				m.clearCell(c.x, c.y)
				m.fillTextCell(c.x, c.y, m.runner.winning_emoji)
			} else {
				m.fillTextCell(c.x, c.y, m.runner.emoji)
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func (m *maze) won() bool {
	return m.runner.pos.x == m.height-1 && m.runner.pos.y == m.width-1
}

// maze creation logic
func (m *maze) fillMaze() {
	vis := make([][]bool, m.height)
	for i := range m.height {
		vis[i] = make([]bool, m.width)
	}
	cnt := 0
	var dfs func(i, j int, prev direction)
	dfs = func(i, j int, prev direction) {
		vis[i][j] = true
		cnt += 1
		if od := prev.opposite(); od != 0 {
			m.cell[i][j] = append(m.cell[i][j], prev.opposite())
		}
		dir := []direction{_RIGHT, _LEFT, _UP, _DOWN}
		rand.Shuffle(len(dir), func(i, j int) { dir[i], dir[j] = dir[j], dir[i] })
		for _, d := range dir {
			x, y := i, j
			switch d {
			case _LEFT:
				y -= 1
			case _RIGHT:
				y += 1
			case _UP:
				x -= 1
			case _DOWN:
				x += 1
			}
			if m.outOfMaze(coord{x, y}) || vis[x][y] {
				continue
			}
			m.cell[i][j] = append(m.cell[i][j], d)
			dfs(x, y, d)
		}
	}
	dfs(0, 0, 0)
}

func (m *maze) listenForArrowKeys() {
	for {
		switch {
		case keyDownState[_LEFT]:
			moveChan <- _LEFT
		case keyDownState[_UP]:
			moveChan <- _UP
		case keyDownState[_RIGHT]:
			moveChan <- _RIGHT
		case keyDownState[_DOWN]:
			moveChan <- _DOWN
		}
		time.Sleep(time.Millisecond * 40)
	}
}

var (
	_INIT        bool
	mazeCanvas   js.Value
	moveChan     = make(chan direction)
	restartChan  = make(chan struct{})
	wallSize     int
	keyDownState = []bool{
		_LEFT:  false,
		_UP:    false,
		_RIGHT: false,
		_DOWN:  false,
	}
)

func (m *maze) init() {
	if _INIT {
		return
	}
	_INIT = true
	// todo: change scale accoding to resolution

	scl := m.scale / 10
	scl += scl % 2
	wallSize = max(4, m.scale/10)
	wallSize = 4 // testing

	c := js.Global().Get("document").Call("getElementById", "mazeCanvas")
	c.Set("height", m.scale*m.height)
	c.Set("width", m.scale*m.width)
	mazeCanvas = c.Call("getContext", "2d")

	js.Global().Get("document").Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) any {
		keyCode := args[0].Get("keyCode").Int()
		switch keyCode {
		case 37, 38, 39, 40:
			keyDownState[keyCode-37+1] = true
		case 82: // r
			restartChan <- struct{}{}
		}
		return nil
	}))

	js.Global().Get("document").Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) any {
		keyCode := args[0].Get("keyCode").Int()
		switch keyCode {
		case 37, 38, 39, 40:
			keyDownState[keyCode-37+1] = false
		}
		return nil
	}))
}

func (m *maze) clearAll() {
	for x := range m.height {
		for y := range m.width {
			mazeCanvas.Call("clearRect", y*m.scale, x*m.scale, m.scale, m.scale)
		}
	}
}

func (m *maze) setAllBorders() {
	mazeCanvas.Set("fillStyle", "black")
	for x := range m.height + 1 {
		mazeCanvas.Call("fillRect", 0, x*m.scale-wallSize/2, m.width*m.scale, wallSize)
	}
	for y := range m.width + 1 {
		mazeCanvas.Call("fillRect", y*m.scale-wallSize/2, 0, wallSize, m.height*m.scale)
	}
	mazeCanvas.Call("fill")
	m.clearPaths()
}

func (m *maze) clearPaths() {
	for x := range m.height {
		for y := range m.width {
			for _, d := range m.cell[x][y] {
				m.clearBorder(x, y, d)
			}
		}
	}
}

func (m *maze) clearBorder(x, y int, d direction) {
	switch d {
	case _LEFT:
		mazeCanvas.Call("clearRect", y*m.scale-wallSize/2, x*m.scale+wallSize/2, wallSize, m.scale-wallSize)
	case _UP:
		mazeCanvas.Call("clearRect", y*m.scale+wallSize/2, x*m.scale-wallSize/2, m.scale-wallSize, wallSize)
	case _RIGHT:
		mazeCanvas.Call("clearRect", y*m.scale+m.scale-wallSize/2, x*m.scale+wallSize/2, wallSize, m.scale-wallSize)
	case _DOWN:
		mazeCanvas.Call("clearRect", y*m.scale+wallSize/2, x*m.scale+m.scale-wallSize/2, m.scale-wallSize, wallSize)
	}
}

func (m *maze) fillTextCell(x, y int, s string) {
	mazeCanvas.Set("font", fmt.Sprintf("%vpx arial", m.scale*3/4))
	mazeCanvas.Set("textAlign", "left")
	mazeCanvas.Set("textBaseline", "top")
	mazeCanvas.Call("fillText", s, y*m.scale+wallSize/2, x*m.scale+m.scale/5)
	mazeCanvas.Call("fill")
}

func (m *maze) clearCell(x, y int) {
	mazeCanvas.Call("clearRect", y*m.scale+wallSize/2, x*m.scale+wallSize/2, m.scale-wallSize, m.scale-wallSize)
}
