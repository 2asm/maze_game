//go:build js && wasm

package main

import "github.com/2asm/maze_game/maze"

func main() {
	maze.NewMaze(15, 20).Start()
}
