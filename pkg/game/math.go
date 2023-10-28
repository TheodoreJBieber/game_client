package game

// return a normalized direction vector that points from 'from' to 'to'
func Direction(from, to Vector2) Vector2 {
	// angle := math32.Atan2(to.Y-from.Y, to.X-from.X)
	// offset := Vector2{
	// 	X: math32.Cos(angle),
	// 	Y: math32.Sin(angle),
	// }
	offset := from.Subtract(to)
	return offset.Normalize()
}
