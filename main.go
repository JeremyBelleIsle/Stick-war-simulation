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
	"Sticks_War/unit"
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
		sm.Move(g.rocks, g.camps, g.stickMans, gamemode.Mode)
		sm.Mine(g.rocks)
		sm.DropMoney(&g.money, g.camps)
	}

	if gamemode.Mode == "attack" {
		stickman.CheckColl(&g.stickMans)
	} else {
		for i := range g.stickMans {
			sm := &g.stickMans[i]

			sm.Attack = false
		}
	}

	for i := range g.offers {
		o := g.offers[i]
		if o.DetectClickOnOffer() {
			stickman.Spawn(&g.stickMans, g.offers, g.camps[0].X, g.camps[0].Y, o.TypeS, "player", &g.money)
		}
	}

	ordiPurchase, err := ordi.ManageMoney(g.offers)
	if err == nil {
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
		money: 10,

		camps: []camp.Camp{
			camp.New(500, screen.Height-350, 1.5),
			camp.New(500*14, screen.Height-350, -1.5),
		},

		rocks: []rock.Rock{
			rock.New(float64(screen.Widht/2), 200),
			rock.New(float64((screen.Widht/2)*4.5), 200),
		},

		offers: store.New(),
	}

	stickman.Spawn(&g.stickMans, g.offers, g.camps[1].X, g.camps[1].Y, unit.Miner, "ordi", &ordi.Money)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
