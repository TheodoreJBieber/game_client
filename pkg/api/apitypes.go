package api

// TODO: may only want to synchronize updated items? if it slows down then we can explore this
type ApiUpdate struct {
	Player         ApiPlayer   `json:"player"`
	Objects        []ApiObject `json:"objects"`
	RemovedObjects []string    `json:"removed_objects"`
	// Disconnect     bool        `json:"disconnect"`
}

type ApiVector2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type ApiPlayer struct {
	Position ApiVector2 `json:"position"`
	Pointer  ApiVector2 `json:"pointer"`
	Username string     `json:"username"`
}

type ApiObject struct {
	Owner        string     `json:"owner"`
	ID           string     `json:"id"`
	Size         float32    `json:"size"`
	AxT          ApiVector2 `json:"axt"`
	Acceleration ApiVector2 `json:"acceleration"`
	Velocity     ApiVector2 `json:"velocity"`
	Position     ApiVector2 `json:"position"`
}
