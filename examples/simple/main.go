package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/funatsufumiya/ebiten_layered_scene_manager/layered_manager"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type State int

type CounterLayer struct {
	layered_manager.LayerBase
	Color color.RGBA
}

// func (l *CounterLayer) Enter()  { fmt.Println("enter:", l.Name)}
// func (l *CounterLayer) Exit()   {}
// func (l *CounterLayer) Reset()  {}
// func (l *CounterLayer) Update() {}
// func (l *CounterLayer) DrawFront(screen *ebiten.Image) {}
// func (l *CounterLayer) DrawBack(screen *ebiten.Image) {}

func (l *CounterLayer) Draw(screen *ebiten.Image) {
	c := l.Color
	c.A = uint8(l.Alpha * 255)
	// screen.Fill(c)
	vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, c, false)
	s := fmt.Sprintf("Layer: %s, Count: %d", l.Name)
	var y int
	switch l.Name {
	case "A":
		y = 40
	case "B":
		y = 80
	case "C":
		y = 120
	default:
		y = 20
	}
	ebitenutil.DebugPrintAt(screen, s, 20, y)
}

func main() {
	layerA := &CounterLayer{layered_manager.LayerBase{Name: "A", Alpha: 1}, color.RGBA{255, 0, 0, 255}}
	layerB := &CounterLayer{layered_manager.LayerBase{Name: "B", Alpha: 1}, color.RGBA{0, 255, 0, 255}}
	layerC := &CounterLayer{layered_manager.LayerBase{Name: "C", Alpha: 1}, color.RGBA{0, 0, 255, 255}}

	manager := layered_manager.NewLayeredSceneManager()
	manager.AddLayer(layerA)
	manager.AddLayer(layerB)
	manager.AddLayer(layerC)

	game := &Game{manager: manager}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Layered Scene Manager Example")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	manager *layered_manager.LayeredSceneManager
}

func (g *Game) Update() error {
       now := time.Now()
       sec := float64(now.UnixNano()) / 1e9
       period := 3.0
       phase := sec / period

       for i, name := range []string{"A", "B", "C"} {
	       l := g.manager.GetLayer(name)
	       if cl, ok := l.(*CounterLayer); ok {
		       offset := float64(i) * (2 * math.Pi / 3)
		       alpha := 0.5 + 0.5*math.Sin(phase*math.Pi*2+offset)
		    //    if alpha < 0.1 {
			//        alpha = 0
		    //    }
		       cl.Alpha = float32(alpha)
	       }
       }
       g.manager.Update()
       return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.manager.Draw(screen)
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}
