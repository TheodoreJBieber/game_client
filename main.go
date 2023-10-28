package main

import (
	"flag"
	"fmt"
	"log"
	"main/pkg/api"
	"main/pkg/game"
	"main/pkg/netsync"
	"main/pkg/serial"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const windowWidth = 640
const windowHeight = 480
const antialias = false
const sendEvery = 1
const moveSpeed = 2
const sprintSpeed = moveSpeed * 1.5
const sprintingEnabled = false
const interpolate = true

// Idea: Skill trees / level up system and perks

// TODO: only send updates on changed objects
// use interpolation to reduce the need for synchronization
// TODO: fix network synchronization

func main() {
	username := flag.String("username", "ted", "username of the player")
	flag.Parse()
	// addr := "54.147.150.253"
	addr := "localhost"
	port := "8087"

	client := netsync.NewClient(addr, port, *username)

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Network Synchronized Game")

	game := NewGame(client)
	client.HandleUpdate = func(ew *serial.EventWrapper) error {
		return game.handleUpdate(ew)
	}

	client.Start()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) handleUpdate(event *serial.EventWrapper) error {
	log.Println(*event)

	if event.Type == serial.Update {
		e := event.Event
		g.Lock()
		p, ok := g.players[e.Player.Username]
		if ok {
			e.Player.FillPlayer(p)
		} else {
			g.players[e.Player.Username] = e.Player.ToPlayer()
		}
		g.Unlock()

		// TODO: may be able to convert logic to be similar to above. synchronize map instead of list

		found := make(map[string]bool, len(e.Objects))
		for _, v := range e.Objects {
			g.Lock()
			g.networkObjects[v.ID] = v.ToObject()
			g.Unlock()
			found[v.ID] = true
		}

		// If we didn't receive an object in an update, delete it from our map
		// for k := range g.networkObjects {
		// 	if !found[k] {
		// 		g.Lock()
		// 		delete(g.networkObjects, k)
		// 		g.Unlock()
		// 	}
		// }

		for _, v := range e.RemovedObjects {
			g.RLock()
			_, ok := g.networkObjects[v]
			g.RUnlock()
			if ok {
				g.Lock()
				delete(g.networkObjects, v)
				g.Unlock()
			}
		}

	}

	return nil
}

type Game struct {
	sync.RWMutex
	client         *netsync.Client
	player         *game.Player
	players        map[string]*game.Player
	count          int
	keys           []ebiten.Key
	networkObjects map[string]*game.Object
	localObjects   map[string]*game.Object
	nextObjectId   uint
	deleted        []string
}

func NewGame(client *netsync.Client) *Game {
	g := &Game{
		client:  client,
		players: make(map[string]*game.Player, 1),
	}

	// Intentionally setting separate references so we keep network and local objects separate
	p := game.NewPlayer(client.Username)
	g.players[client.Username] = &p
	p2 := game.NewPlayer(client.Username)
	g.player = &p2

	g.localObjects = make(map[string]*game.Object)
	g.networkObjects = make(map[string]*game.Object)
	g.nextObjectId = 0
	g.deleted = make([]string, 0)

	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	skipLocal := false
	g.RLock()
	for username, p := range g.players {
		if skipLocal && username == g.player.Username {
			continue
		}
		p.Draw(screen, true, antialias)
	}
	g.RUnlock()
	g.player.Draw(screen, false, antialias)

	g.RLock()
	for _, v := range g.networkObjects {
		if skipLocal {
			// TODO: logic wouldn't work if 2 players have similar names: John, Joh, Jo, J, or even ""
			if strings.HasPrefix(v.ID, g.player.Username) {
				continue
			}
		}
		v.Draw(screen, true, antialias)
	}
	g.RUnlock()

	for _, v := range g.localObjects {
		v.Draw(screen, false, antialias)
	}
}

func (g *Game) GetPlayer() *game.Player {
	return g.players[g.client.Username]
}

func (g *Game) Update() error {
	// 60 ticks per second
	var timePassed float32 = 1

	g.handleKeyboardInput()

	mouseV2 := getMouseVector2()
	x := mouseV2.X
	y := mouseV2.Y
	g.player.Pointer.X = float32(x)
	g.player.Pointer.Y = float32(y)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		g.shoot()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		for _, v := range g.localObjects {
			g.deleteLocalObject(*v)
		}
	}

	for _, v := range g.localObjects {
		v.Tick(timePassed)
	}

	if interpolate {
		for _, v := range g.networkObjects {
			v.Tick(timePassed)
		}
	}

	g.count++
	if g.count%sendEvery == 0 {
		g.sendUpdateToServer()
		g.resetLocalObjectUpdateStatus()
	}
	return nil
}

func (g *Game) shoot() {
	mouseV2 := getMouseVector2()

	id := g.getNextId()
	object := game.NewObject(g.player.Username, id)
	object.Position = g.player.Position

	b := g.player.Bullet
	object.Velocity = b.GetVelocity(g.player.Position, mouseV2)
	object.Acceleration = b.GetAcceleration()
	object.AxT = b.GetAxT()

	g.localObjects[id] = &object
}

func (g *Game) handleKeyboardInput() {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	pressed := make(map[string]bool)
	for _, k := range g.keys {
		if name := ebiten.KeyName(k); name != "" {
			pressed[name] = true
		}
		pressed[k.String()] = true
	}

	var movement game.Vector2 = game.Vector2{}

	speed := float32(moveSpeed)
	if sprintingEnabled && pressed["Shift"] {
		speed = sprintSpeed
	}

	if pressed["W"] {
		movement.Y -= speed
	}
	if pressed["S"] {
		movement.Y += speed
	}
	if pressed["A"] {
		movement.X -= speed
	}
	if pressed["D"] {
		movement.X += speed
	}
	p2 := g.player.Position.Add(movement)
	g.player.Position.X = p2.X
	g.player.Position.Y = p2.Y
}

func (g *Game) sendUpdateToServer() {
	g.RLock()
	player := api.FromPlayer(*g.player)
	g.RUnlock()

	objects := api.FromObjects(g.localObjects)

	update := api.ApiUpdate{
		Player:         player,
		Objects:        objects,
		RemovedObjects: g.deleted,
	}

	event := *serial.WrapEvent(serial.Update, update)
	g.client.WriteMessage(event)
}

func (g *Game) resetLocalObjectUpdateStatus() {
	for _, v := range g.localObjects {
		v.Updated = false
	}

	g.deleted = make([]string, 0)
}

func (g *Game) getNextId() string {
	id := g.nextObjectId
	g.nextObjectId++

	name := g.player.Username

	return fmt.Sprintf("%s%d", name, id)
}

func (g *Game) deleteLocalObject(o game.Object) {
	g.deleted = append(g.deleted, o.ID)
	delete(g.localObjects, o.ID)
}

func getMouseVector2() game.Vector2 {
	x, y := ebiten.CursorPosition()
	return game.Vector2{
		X: float32(x),
		Y: float32(y),
	}
}
