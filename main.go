package main

import (
	rock "Sticks_War/Rock"
	stickman "Sticks_War/Stick_Man"
	"Sticks_War/cam"
	"Sticks_War/camp"
	gamemode "Sticks_War/gameMode"
	"Sticks_War/ordi"
	"Sticks_War/screen"
	"Sticks_War/store"
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Game struct {
	rocks     []rock.Rock
	stickMans []stickman.StickMan
	camps     []camp.Camp
	offers    []store.Offer
	money     int
}

var mplusSource *text.GoTextFaceSource

func (g *Game) Update() error {
	for i := range g.stickMans {
		sm := &g.stickMans[i]
		sm.Move(g.rocks, g.camps, gamemode.Mode)
		sm.Mine(g.rocks)
		sm.DropMoney(&g.money, g.camps)
	}
	for i := range g.offers {
		o := g.offers[i]
		if o.DetectClickOnOffer() {
			stickman.Spawn(&g.stickMans, g.offers, g.camps[0].X, g.camps[0].Y, o.TypeS, "player", &g.money)
		}
	}
	ordiPurchase := ordi.ManageMoney(g.offers)
	if ordiPurchase != "not enough money" {
		stickman.Spawn(&g.stickMans, g.offers, g.camps[1].X, g.camps[1].Y, ordiPurchase, "ordi", &ordi.Money)
	}

	cam.ScrollCamX()
	gamemode.DetectClickToChangeGameMode(&gamemode.Mode)
	return nil
}

func (g *Game) Draw(screenI *ebiten.Image) {
	for _, r := range g.rocks {
		r.Draw(screenI)
	}
	for _, sm := range g.stickMans {
		sm.Draw(screenI, mplusSource)
	}
	for _, c := range g.camps {
		c.Draw(screenI, mplusSource)
	}
	gameutil.DrawText(fmt.Sprintf("Money: %d", g.money), 40, screen.Widht, 0, 0, 0, screenI, color.RGBA{255, 223, 0, 255}, mplusSource)
	gamemode.Draw(screenI)
	for _, o := range g.offers {
		o.Draw(screenI, mplusSource)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screen.Widht, screen.Height
}

func main() {
	ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(screen.Widht, screen.Height)
	ebiten.SetWindowTitle("Sticks War")

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))

	if err != nil {
		log.Fatal(err)
	}

	mplusSource = s

	g := &Game{
		money: 500,
	}
	rock.Init(&g.rocks)
	camp.Init(&g.camps)
	store.Init(&g.offers)
	stickman.Spawn(&g.stickMans, g.offers, g.camps[1].X, g.camps[1].Y, "minner", "ordi", &ordi.Money)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
