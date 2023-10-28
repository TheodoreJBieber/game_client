package game

import (
	"image/color"

	"github.com/chewxy/math32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const playerSize = 10
const gunLength = playerSize + gunWidth + 1
const gunWidth = 3

func (p *Player) Draw(screen *ebiten.Image, network bool, antialias bool) {
	start, end := newRay(p.Position, p.Pointer, gunLength)
	if network {
		vector.DrawFilledCircle(screen, p.Position.X, p.Position.Y, playerSize, color.RGBA{190, 190, 190, 100}, antialias)
		vector.StrokeLine(screen, start.X, start.Y, end.X, end.Y, gunWidth, color.RGBA{255, 100, 100, 0}, antialias)
	} else {
		vector.DrawFilledCircle(screen, p.Position.X, p.Position.Y, playerSize, color.RGBA{0, 175, 0, 255}, antialias)
		vector.StrokeLine(screen, start.X, start.Y, end.X, end.Y, gunWidth, color.RGBA{255, 0, 0, 255}, antialias)
	}
}

func newRay(from, to Vector2, length float32) (start, end Vector2) {
	angle := math32.Atan2(to.Y-from.Y, to.X-from.X)
	start = from
	offset := Vector2{
		X: length * math32.Cos(angle),
		Y: length * math32.Sin(angle),
	}
	end = start.Add(offset)
	return start, end
}

func (o *Object) Draw(screen *ebiten.Image, network bool, antialias bool) {
	if network {
		vector.DrawFilledCircle(screen, o.Position.X, o.Position.Y, o.Size, color.RGBA{100, 100, 100, 255}, antialias)
	} else {
		vector.DrawFilledCircle(screen, o.Position.X, o.Position.Y, o.Size, color.RGBA{0, 0, 255, 255}, antialias)
	}
}
