//go:build js && wasm

package main

import "github.com/2asm/maze_game/maze"

func main() {
	maze.NewMaze(20, 30, 30).Start()
}
