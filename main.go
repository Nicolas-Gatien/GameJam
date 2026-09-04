package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	GECKO_IMAGE      *ebiten.Image
	GECKO_ARMS_IMAGE *ebiten.Image
	TONGUE_IMAGE     *ebiten.Image
	FLY_IMAGE        *ebiten.Image
	keys             []ebiten.Key
	player           Gecko
	flies            []Fly
	frameCount       int

	windowWidth  int
	windowHeight int
	layoutWidth  int
	layoutHeight int

	leftOutBoundLimit  int
	rightOutBoundLimit int
}

type Position struct {
	positionX float64
	positionY float64
}

type Gecko struct {
	positionX    float64
	speed        float64
	tongueLength float64
	maxLength    float64
	tongueSpeed  float64
	attacking    bool
}

type Fly struct {
	Position
}

func (g *Game) SpawnNpcs() {
	if rand.IntN(2) == 0 {
		print("Spawned a fly")
		g.flies = append(g.flies, Fly{Position{positionX: float64(rand.IntN(g.layoutWidth)), positionY: 40}})
	} else {
		print("Spanwed a rock")
	}

}

func (g *Game) Update() error {

	g.frameCount += 1
	if g.frameCount == 180 {
		g.SpawnNpcs()
		g.frameCount = 0
	}

	geckoBounds := g.GECKO_IMAGE.Bounds()
	geckoWidth := geckoBounds.Dx()
	geckoHeight := geckoBounds.Dy()

	fmt.Printf("Gecko Width: %d, Gecko Height: %d\n", geckoWidth, geckoHeight)

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	fmt.Println(g.player.positionX)

	if g.player.attacking == false {
		if slices.Contains(g.keys, ebiten.KeyA) {
			if g.player.positionX-g.player.speed >= float64(g.leftOutBoundLimit+geckoWidth/2) {
				g.player.positionX -= g.player.speed
			}
		}
		if slices.Contains(g.keys, ebiten.KeyD) {
			if g.player.positionX+g.player.speed <= float64(g.rightOutBoundLimit-geckoWidth/2) {
				g.player.positionX += g.player.speed
			}
		}
	}

	if slices.Contains(g.keys, ebiten.KeySpace) {
		g.player.attacking = true
	} else {
		g.player.attacking = false
	}

	if g.player.attacking {
		if g.player.tongueLength < g.player.maxLength {
			g.player.tongueLength += g.player.tongueSpeed
		}
	} else {
		if g.player.tongueLength > 0 {
			g.player.tongueLength -= g.player.tongueSpeed * 2.5
		} else {
			g.player.tongueLength = 0
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	geckoBounds := g.GECKO_IMAGE.Bounds()
	geckoWidth := geckoBounds.Dx()

	background := color.RGBA{172, 220, 215, 0xff}
	screen.Fill(background)

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Rotate(math.Sin(float64(g.frameCount)))
	opts.GeoM.Translate(float64(g.layoutWidth/2)+g.player.positionX-float64(g.GECKO_ARMS_IMAGE.Bounds().Dx())/2, float64(g.layoutHeight)-float64(g.GECKO_ARMS_IMAGE.Bounds().Dy()+55))
	screen.DrawImage(g.GECKO_ARMS_IMAGE, &opts)

	opts = ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(g.layoutWidth/2)+g.player.positionX-float64(geckoWidth)/2, float64(g.layoutHeight)-float64(g.GECKO_IMAGE.Bounds().Dy()))
	screen.DrawImage(g.GECKO_IMAGE, &opts)

	for _, fly := range g.flies {
		flyOptions := ebiten.DrawImageOptions{}
		flyOptions.GeoM.Translate(fly.positionX, fly.positionY)
		screen.DrawImage(g.FLY_IMAGE, &flyOptions)
	}

	tongue := g.TONGUE_IMAGE
	tongueOptions := ebiten.DrawImageOptions{}
	tongueOptions.GeoM.Scale(1, -g.player.tongueLength)
	tongueOptions.GeoM.Translate(float64(g.layoutWidth/2)+g.player.positionX-float64(g.TONGUE_IMAGE.Bounds().Dx()/2), float64(g.layoutHeight)-40)
	screen.DrawImage(tongue, &tongueOptions)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.layoutWidth, g.layoutHeight
}

func main() {
	geckoImage, _, err := ebitenutil.NewImageFromFile("assets/gecko_body.png")
	if err != nil {
		log.Fatal(err)
	}

	geckoArmsImage, _, err := ebitenutil.NewImageFromFile("assets/gecko_arms.png")
	if err != nil {
		log.Fatal(err)
	}

	tongueImage, _, err := ebitenutil.NewImageFromFile("assets/tongue.png")
	if err != nil {
		log.Fatal(err)
	}

	flyImage, _, err := ebitenutil.NewImageFromFile("assets/fly.png")
	if err != nil {
		log.Fatal(err)
	}

	windowWidth := 480
	windowHeight := 640

	game := Game{
		GECKO_IMAGE:      geckoImage,
		TONGUE_IMAGE:     tongueImage,
		GECKO_ARMS_IMAGE: geckoArmsImage,
		FLY_IMAGE:        flyImage,
		player:           Gecko{speed: 3, maxLength: float64(windowHeight/2) - 80, tongueSpeed: 8},
		frameCount:       120,

		windowWidth:  windowWidth,
		windowHeight: windowHeight,
		layoutWidth:  windowWidth / 2,
		layoutHeight: windowHeight / 2,

		leftOutBoundLimit:  -(windowWidth / 4),
		rightOutBoundLimit: (windowWidth / 4),
	}

	ebiten.SetWindowSize(game.windowWidth, game.windowHeight)
	ebiten.SetWindowTitle("Hello Nico I see you ")

	err = ebiten.RunGame(&game)

	if err != nil {
		log.Fatal(err)
	}
}
