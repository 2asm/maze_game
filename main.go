package main

import (
	"fmt"
	"math/rand"
)

const (
	LEFT   = 3
	RIGHT  = 1
	TOP    = 0
	BOTTOM = 2
)

const WALL = string('█')
// const WALL = string('▒')
// const WALL = string('░')
// const WALL = string('π')

type Cell struct {
	// top right bottom left
	// 0   1     2      3
	broken_walls [4]int
}

type Maze struct {
	rows int
	cols int
	maze [][]Cell
}

func NewMaze(n int, m int) *Maze {
	return &Maze{
		rows: n,
		cols: m,
		maze: generate_maze(n, m),
	}
}

func (mz *Maze) Print() {
	fmt.Printf("%v%v  %v%v", WALL, WALL, WALL, WALL)
	for i := 1; i < mz.cols; i++ {
		fmt.Printf("%v%v%v%v", WALL, WALL, WALL, WALL)
	}
	fmt.Printf("\n")
	for i := 0; i < mz.rows; i++ {
		fmt.Printf("%v%v", WALL, WALL)
		for j := 0; j < mz.cols; j++ {
			if mz.maze[i][j].broken_walls[RIGHT] == 0 {
				fmt.Printf("  %v%v", WALL, WALL)
			} else {
				fmt.Printf("    ")
			}
		}
		fmt.Printf("\n")
		fmt.Printf("%v%v", WALL, WALL)
		for j := 0; j < mz.cols; j++ {
			if mz.maze[i][j].broken_walls[BOTTOM] == 0 {
				if !(i == mz.rows-1 && j == mz.cols-1) {
					fmt.Printf("%v%v%v%v", WALL, WALL, WALL, WALL)
				} else {
					fmt.Printf("  %v%v", WALL, WALL)
				}
			} else {
				fmt.Printf("  %v%v", WALL, WALL)
			}
		}
		fmt.Printf("\n")
	}
}

func (mz *Maze) Print2() {
	fmt.Printf("%v %v", WALL, WALL)
	for i := 1; i < mz.cols; i++ {
		fmt.Printf("%v%v", WALL, WALL)
	}
	fmt.Printf("\n")
	for i := 0; i < mz.rows; i++ {
		fmt.Printf("%v", WALL)
		for j := 0; j < mz.cols; j++ {
			if mz.maze[i][j].broken_walls[RIGHT] == 0 {
				fmt.Printf(" %v", WALL)
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Printf("\n")
		fmt.Printf("%v", WALL)
		for j := 0; j < mz.cols; j++ {
			if mz.maze[i][j].broken_walls[BOTTOM] == 0 {
				if !(i == mz.rows-1 && j == mz.cols-1) {
					fmt.Printf("%v%v", WALL, WALL)
				} else {
					fmt.Printf(" %v", WALL)
				}
			} else {
				fmt.Printf(" %v", WALL)
			}
		}
		fmt.Printf("\n")
	}
}

func generate_maze(n int, m int) [][]Cell {
	maze := make([][]Cell, n)
	for i := 0; i < n; i++ {
		maze[i] = make([]Cell, m)
	}
	vis := make([][]bool, n)
	for i := 0; i < n; i++ {
		vis[i] = make([]bool, m)
	}
	dxy := [4][2]int{{0, 1}, {1, 0}, {-1, 0}, {0, -1}}
	break_wall := func(i int, j int, x int, y int) {
		if i == x {
			if j < y {
				maze[i][j].broken_walls[RIGHT] = 1
				maze[x][y].broken_walls[LEFT] = 1
			} else {
				maze[i][j].broken_walls[LEFT] = 1
				maze[x][y].broken_walls[RIGHT] = 1
			}
		} else {
			if i < x {
				maze[i][j].broken_walls[BOTTOM] = 1
				maze[x][y].broken_walls[TOP] = 1
			} else {
				maze[i][j].broken_walls[TOP] = 1
				maze[x][y].broken_walls[BOTTOM] = 1
			}
		}
	}
	var dfs func(i int, j int)
	dfs = func(i int, j int) {
		vis[i][j] = true
		nei := make([][2]int, 0)
		for _, xy := range dxy {
			x := i + xy[0]
			y := j + xy[1]
			if x < 0 || x >= n || y < 0 || y >= m {
				continue
			}
			nei = append(nei, [2]int{x, y})
		}
		for ii := len(nei) - 1; ii > 0; ii-- { // Fisher–Yates shuffle
			jj := rand.Intn(ii + 1)
			nei[ii], nei[jj] = nei[jj], nei[ii]
		}
		for _, v := range nei {
			x := v[0]
			y := v[1]
			if !vis[x][y] {
				break_wall(i, j, x, y)
				dfs(x, y)
			}
		}
	}
	dfs(0, 0)
	return maze
}

func main() {
	maze := NewMaze(10, 20)
	maze.Print()
}
