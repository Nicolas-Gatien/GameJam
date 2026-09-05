package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	FLY_SPRITE           *ebiten.Image
	keys                 []ebiten.Key
	player               Gecko
	flies                []Fly
	rocks                []Rock
	framesBtwSpawn       int
	framesUntilNextSpawn int
	frameCount           int
	score                int
	flyspeed             float64

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
	Position
	SPRITE        *ebiten.Image
	ARMS_SPRITE   *ebiten.Image
	LEGS_SPRITE   *ebiten.Image
	TONGUE_SPRITE *ebiten.Image
	speed         float64
	tongueLength  float64
	tongue        Position
	maxLength     float64
	tongueSpeed   float64
	attacking     bool
}

type Fly struct {
	Position
}

type Rock struct {
}

func (g *Game) SpawnNpcs() {
	min := 10
	max := g.layoutWidth - 30

	if rand.IntN(1) == 0 {
		flyPositionX := float64(rand.IntN(max-min+1) + min)

		g.flies = append(g.flies, Fly{Position{positionX: flyPositionX, positionY: 0}})
		fmt.Printf("Spawned a fly: %f\n", flyPositionX)

	} else {
		print("Spanwed a rock\n")
	}

}

func (g *Game) CollisionCheck() error {

	// print("CHECKING FOR COLLISION")

	xStart := g.player.tongue.positionX - 6
	yStart := g.player.tongue.positionY

	xBound := g.player.TONGUE_SPRITE.Bounds().Dx() * 3
	// yBound := g.player.TONGUE_SPRITE.Bounds().Dy()

	xEnd := xStart + float64(xBound)
	yEnd := yStart - g.player.tongueLength

	for index, fly := range g.flies {
		// fmt.Printf("xStart: %f, xEnd: %f, FlyX: %f \n", xStart, xEnd, fly.positionX)
		// fmt.Printf("yStart: %f, yEnd: %f, FlyY: %f \n", yStart, yEnd, fly.positionY)

		if fly.positionX >= xStart && fly.positionX <= xEnd {
			if fly.positionY <= yStart && fly.positionY >= yEnd {
				g.flies = slices.Delete(g.flies, index, index+1)
				g.score += 1
				print(g.flies)
				// print("DEAD FLY!!!")
			}
		}
	}

	return nil
}

func (g *Game) Update() error {
	g.frameCount += 1
	g.framesUntilNextSpawn -= 1
	if g.framesUntilNextSpawn <= 0 {
		g.SpawnNpcs()
		g.framesUntilNextSpawn = g.framesBtwSpawn
		g.framesBtwSpawn -= 1
		g.flyspeed += 0.25
	}

	for i, fly := range g.flies {
		g.flies[i].positionY += 0.5
		g.flies[i].positionX = g.flies[i].positionX + (math.Sin(float64(g.frameCount)/10) * 2)
		if fly.positionY >= float64(g.layoutHeight) {
			g.score = 0
			g.flyspeed = 0.5
			g.framesBtwSpawn = 180
			g.flies[i].positionY = 0
		}
	}

	geckoBounds := g.player.SPRITE.Bounds()
	geckoWidth := geckoBounds.Dx()
	// geckoHeight := geckoBounds.Dy()

	// fmt.Printf("Gecko Width: %d, Gecko Height: %d\n", geckoWidth, geckoHeight)

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

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
		g.CollisionCheck()
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

	geckoBounds := g.player.SPRITE.Bounds()
	geckoWidth := geckoBounds.Dx()

	background := color.RGBA{172, 220, 215, 0xff}
	screen.Fill(background)

	ebitenutil.DebugPrint(screen, strconv.Itoa(g.score))

	tongue := g.player.TONGUE_SPRITE
	tongueOptions := ebiten.DrawImageOptions{}

	tongueOptions.GeoM.Scale(1, -g.player.tongueLength)

	g.player.tongue.positionX = g.player.positionX - float64(g.player.TONGUE_SPRITE.Bounds().Dx()/2)
	g.player.tongue.positionY = g.player.positionY

	tongueOptions.GeoM.Translate(g.player.tongue.positionX, g.player.tongue.positionY)
	screen.DrawImage(tongue, &tongueOptions)

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-float64(g.player.ARMS_SPRITE.Bounds().Dx())/2, -float64(g.player.ARMS_SPRITE.Bounds().Dy()/2))
	if !g.player.attacking {
		opts.GeoM.Rotate(math.Sin(float64(g.frameCount)/3) / 2)
	}
	opts.GeoM.Translate(g.player.positionX, float64(g.layoutHeight)-68)
	screen.DrawImage(g.player.ARMS_SPRITE, &opts)

	opts = ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-float64(g.player.LEGS_SPRITE.Bounds().Dx())/2, -float64(g.player.LEGS_SPRITE.Bounds().Dy()/2))
	if !g.player.attacking {
		opts.GeoM.Rotate(-math.Sin(float64(g.frameCount)/3) / 2)
	}
	opts.GeoM.Translate(g.player.positionX, float64(g.layoutHeight)-48)
	screen.DrawImage(g.player.LEGS_SPRITE, &opts)

	opts = ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.player.positionX-float64(geckoWidth)/2, float64(g.layoutHeight)-float64(g.player.SPRITE.Bounds().Dy()))
	screen.DrawImage(g.player.SPRITE, &opts)

	for _, fly := range g.flies {
		flyOptions := ebiten.DrawImageOptions{}
		flyOptions.GeoM.Translate(fly.positionX-8, fly.positionY-8)
		screen.DrawImage(g.FLY_SPRITE, &flyOptions)
	}

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

	geckoLegsImage, _, err := ebitenutil.NewImageFromFile("assets/gecko_legs.png")
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

	player := Gecko{
		SPRITE:        geckoImage,
		TONGUE_SPRITE: tongueImage,
		ARMS_SPRITE:   geckoArmsImage,
		LEGS_SPRITE:   geckoLegsImage,
		speed:         3,
		maxLength:     float64(windowHeight/2) - 100,
		tongueSpeed:   8,
		Position:      Position{float64(windowWidth) / 4, float64(windowHeight)/2 - 80},
	}

	game := Game{
		FLY_SPRITE:     flyImage,
		player:         player,
		frameCount:     0,
		framesBtwSpawn: 180,
		flyspeed:       0.5,

		windowWidth:  windowWidth,
		windowHeight: windowHeight,
		layoutWidth:  windowWidth / 2,
		layoutHeight: windowHeight / 2,

		leftOutBoundLimit:  0,
		rightOutBoundLimit: (windowWidth / 2),
	}

	ebiten.SetWindowSize(game.windowWidth, game.windowHeight)
	ebiten.SetWindowTitle("Gecko Game")

	err = ebiten.RunGame(&game)

	if err != nil {
		log.Fatal(err)
	}
}
