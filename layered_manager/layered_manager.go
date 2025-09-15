package layered_manager

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type LayerBase struct {
	Name      string
	Alpha     float32
}
func (l *LayerBase) IsVisible() bool { return l.Alpha > 0.0 }
func (l *LayerBase) GetName() string { return l.Name }

type BlendMode int

const (
	BlendNone BlendMode = iota
	BlendAlpha
	BlendAdd
)

type Layer interface {
	GetName() string
	Update()
	Draw(screen *ebiten.Image)
	Draw2dFront(screen *ebiten.Image)
	Draw2dBack(screen *ebiten.Image)
	Enter()
	Exit()
	Reset()
}

type LayeredSceneManager struct {
	layers []Layer
}

func NewLayeredSceneManager() *LayeredSceneManager {
	return &LayeredSceneManager{}
}

func (m *LayeredSceneManager) AddLayer(l Layer) {
	m.layers = append(m.layers, l)
	l.Enter()
}

func (m *LayeredSceneManager) RemoveLayer(name string) {
	for i, l := range m.layers {
		if l.GetName() == name {
			l.Exit()
			m.layers = append(m.layers[:i], m.layers[i+1:]...)
			break
		}
	}
}

func (m *LayeredSceneManager) Update() {
       for _, l := range m.layers {
	       // Layer側でVisibleフィールドを持つ場合は型アサーションで判定
	       if lb, ok := l.(interface{ Visible() bool }); ok {
		       if lb.Visible() {
			       l.Update()
		       }
	       } else {
		       l.Update()
	       }
       }
}

func (m *LayeredSceneManager) Draw(screen *ebiten.Image) {
       for _, l := range m.layers {
	       if lb, ok := l.(interface{ IsVisible() bool }); ok {
		       if lb.IsVisible() {
			       l.Draw2dBack(screen)
		       }
	       } else {
		       l.Draw2dBack(screen)
	       }
       }
       for _, l := range m.layers {
	       if lb, ok := l.(interface{ IsVisible() bool }); ok {
		       if lb.IsVisible() {
			       l.Draw(screen)
		       }
	       } else {
		       l.Draw(screen)
	       }
       }
       for _, l := range m.layers {
	       if lb, ok := l.(interface{ IsVisible() bool }); ok {
		       if lb.IsVisible() {
			       l.Draw2dFront(screen)
		       }
	       } else {
		       l.Draw2dFront(screen)
	       }
       }
}

func (m *LayeredSceneManager) SetLayerOrder(names []string) {
	var newLayers []Layer
	for _, name := range names {
		for _, l := range m.layers {
			if l.GetName() == name {
				newLayers = append(newLayers, l)
				break
			}
		}
	}
	m.layers = newLayers
}

func (m *LayeredSceneManager) GetLayer(name string) Layer {
	for _, l := range m.layers {
		if l.GetName() == name {
			return l
		}
	}
	return nil
}
