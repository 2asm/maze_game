//go:build js && wasm

package maze

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

type maze struct {
	height int
	width  int
	cell   [][][]direction
	runner mazeRunner
	isOver bool
}

func (m maze) outOfMaze(c coord) bool {
	return c.x < 0 || c.x >= m.height || c.y < 0 || c.y >= m.width
}

func NewMaze(n int, m int) *maze {
	_maze := make([][][]direction, n)
	for i := range n {
		_maze[i] = make([][]direction, m)
	}
	ret := &maze{
		height: n,
		width:  m,
		cell:   _maze,
		runner: newMazeRunner(coord{0, 0}),
	}
	ret.fillMaze()
	if !INIT {
		INIT = true
		ret.init()
	}
	ret.style()
	grid[0][0].Set("innerText", string(ret.runner.emoji))
	grid[0][0].Get("style").Set("color", "grey")
	grid[ret.height-1][ret.width-1].Set("innerText", direction(0).String())
	grid[ret.height-1][ret.width-1].Get("style").Set("color", "green")
	return ret
}

func (m *maze) Start() {
	for {
		select {
		case <-restartChan:
			m = NewMaze(m.height, m.width)
		case d := <-moveChan:
			if m.isOver {
				continue
			}
			c := m.runner.pos
			for _, di := range m.cell[c.x][c.y] {
				if di == d {
					m.runner.move(d)
					break
				}
			}
			grid[c.x][c.y].Set("innerText", direction(0).String())
			// grid[c.x][c.y].Get("style").Set("color", "transparent")

			c = m.runner.pos

			grid[c.x][c.y].Set("innerText", string(m.runner.emoji))
			grid[c.x][c.y].Get("style").Set("color", "grey")
			if c.x == m.height-1 && c.y == m.width-1 {
				grid[c.x][c.y].Set("innerText", string(m.runner.winning_emoji))
				m.isOver = true
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
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

var (
	grid         = make([][]js.Value, 0)
	INIT         bool
	restart      js.Value
	moveChan     = make(chan direction)
	restartChan  = make(chan struct{})
	idToCoordMap = map[string]coord{}
	keyState     = map[string]bool{}
)

func (m *maze) coordToId(c coord) string {
	return fmt.Sprintf("%v-%v", c.x, c.y)
}

func (m *maze) idToCoord(id string) coord {
	if c, ok := idToCoordMap[id]; ok {
		return c
	}
	parts := strings.Split(id, "-")
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		panic("atoi")
	}
	y, err := strconv.Atoi(parts[0])
	if err != nil {
		panic("atoi")
	}
	out := coord{x, y}
	idToCoordMap[id] = out
	return out
}

func (m *maze) init() {
	for i := range m.height {
		gi := []js.Value{}
		for j := range m.width {
			e := js.Global().Get("document").Call("getElementById", m.coordToId(coord{i, j}))
			gi = append(gi, e)
		}
		grid = append(grid, gi)
	}

	restart = js.Global().Get("document").Call("getElementById", "restart")

	restart.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		restartChan <- struct{}{}
		return nil
	}))

	js.Global().Get("document").Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) any {
		key := args[0].Get("key").String()
		keyState[key] = true
		if key == "r" || key == "R" {
			restartChan <- struct{}{}
			return nil
		}
		go func() {
		top:
			for keyState[key] {
				switch key {
				case "ArrowUp":
					moveChan <- _UP
				case "ArrowDown":
					moveChan <- _DOWN
				case "ArrowLeft":
					moveChan <- _LEFT
				case "ArrowRight":
					moveChan <- _RIGHT
				default:
					break top
				}
				time.Sleep(time.Millisecond * 60)
			}
		}()
		return nil
	}))

	js.Global().Get("document").Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) any {
		key := args[0].Get("key").String()
		keyState[key] = false
		return nil
	}))
}

func (m *maze) style() {
	for i := range m.height {
		for j := range m.width {
			str := "border: 0.1em solid black;"
			for _, d := range m.cell[i][j] {
				switch d {
				case _LEFT:
					str += "border-left: 0.1em solid #ddd;"
				case _RIGHT:
					str += "border-right: 0.1em solid #ddd;"
				case _UP:
					str += "border-top: 0.1em solid #ddd;"
				case _DOWN:
					str += "border-bottom: 0.1em solid #ddd;"
				}
			}
			grid[i][j].Call("setAttribute", "style", str)
		}
	}
}
