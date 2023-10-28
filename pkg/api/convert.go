package api

import "main/pkg/game"

func FromVector2(vector2 game.Vector2) ApiVector2 {
	return ApiVector2(vector2)
}

func (a ApiVector2) ToVector2() game.Vector2 {
	return game.Vector2(a)
}

func FromPlayer(player game.Player) ApiPlayer {
	return ApiPlayer{
		Username: player.Username,
		Position: FromVector2(player.Position),
		Pointer:  FromVector2(player.Pointer),
	}
}

func (a ApiPlayer) FillPlayer(p *game.Player) {
	p.Position.X = a.Position.X
	p.Position.Y = a.Position.Y
	p.Pointer.X = a.Pointer.X
	p.Pointer.Y = a.Pointer.Y
}

func (a ApiPlayer) ToPlayer() *game.Player {
	return &game.Player{
		Username: a.Username,
		Position: a.Position.ToVector2(),
		Pointer:  a.Pointer.ToVector2(),
	}
}

func FromObjects(localObjects map[string]*game.Object) []ApiObject {
	objects := make([]ApiObject, 0, len(localObjects))

	for _, o := range localObjects {
		if o.Updated {
			objects = append(objects, FromObject(*o))
		}
	}

	return objects
}

func FromObject(localObject game.Object) ApiObject {
	return ApiObject{
		Owner:        localObject.Owner,
		ID:           localObject.ID,
		Size:         localObject.Size,
		AxT:          FromVector2(localObject.AxT),
		Acceleration: FromVector2(localObject.Acceleration),
		Velocity:     FromVector2(localObject.Velocity),
		Position:     FromVector2(localObject.Position),
	}
}

func (a ApiObject) ToObject() *game.Object {
	return &game.Object{
		Owner:        a.Owner,
		ID:           a.ID,
		Size:         a.Size,
		AxT:          a.AxT.ToVector2(),
		Acceleration: a.Acceleration.ToVector2(),
		Velocity:     a.Velocity.ToVector2(),
		Position:     a.Position.ToVector2(),
	}
}
