package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/KasperJSdeVries/zoned_world_servers/internal/game"
)


func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Test")
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
