package game

import "github.com/chewxy/math32"

type Vector2 struct {
	X float32
	Y float32
}

func (v Vector2) Add(v2 Vector2) Vector2 {
	return Vector2{X: v.X + v2.X, Y: v.Y + v2.Y}
}

func (v Vector2) Multiply(coefficient float32) Vector2 {
	return Vector2{X: v.X * coefficient, Y: v.Y * coefficient}
}

func (v Vector2) Divide(denom float32) Vector2 {
	return Vector2{X: v.X / denom, Y: v.Y / denom}
}

func (v Vector2) Subtract(v2 Vector2) Vector2 {
	return Vector2{X: v2.X - v.X, Y: v2.Y - v.Y}
}

func (v Vector2) Normalize() Vector2 {
	l := math32.Sqrt(v.X*v.X + v.Y*v.Y)
	if l == 0 {
		return Vector2{X: 1, Y: 0}
	}
	return v.Divide(l)
}

type Player struct {
	Position Vector2
	Pointer  Vector2
	Username string
	Bullet   Ammo
}

func NewPlayer(username string) Player {
	return Player{
		Username: username,
		Bullet: Ammo{
			AxT: Vector2{
				X: 0,
				Y: 0,
			},
			A: Vector2{
				X: 0,
				Y: 0,
			},
		},
	}
}

type Ammo struct {
	AxT Vector2
	A   Vector2
}

// TODO: monoid for configuration to handle attributes?
func (a Ammo) GetVelocity(start, end Vector2) Vector2 {
	return Direction(start, end)
}

func (a Ammo) GetAxT() Vector2 {
	return a.AxT
}

func (a Ammo) GetAcceleration() Vector2 {
	return a.A
}

type Object struct {
	Owner        string
	ID           string
	Size         float32
	Updated      bool
	AxT          Vector2
	Acceleration Vector2
	Velocity     Vector2
	Position     Vector2
}

func NewObject(owner, id string) Object {
	return Object{
		Owner:   owner,
		ID:      id,
		Size:    5,
		Updated: true,
	}
}

func (o *Object) Tick(timePassed float32) {
	o.Acceleration = o.Acceleration.Add(o.AxT.Multiply(timePassed))
	o.Velocity = o.Velocity.Add(o.Acceleration.Multiply(timePassed))
	o.Position = o.Position.Add(o.Velocity.Multiply(timePassed))
}
