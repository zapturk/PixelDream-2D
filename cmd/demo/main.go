package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/zapturk/PixelDream-2D/pkg/display"
	"github.com/zapturk/PixelDream-2D/pkg/engine"
	"github.com/zapturk/PixelDream-2D/pkg/gfx"
	"github.com/zapturk/PixelDream-2D/pkg/input"
	"github.com/zapturk/PixelDream-2D/pkg/physics"
	"github.com/zapturk/PixelDream-2D/pkg/tilemap"
)

// Generate procedural 64x16 texture atlas with 4 tiles of 16x16:
// Tile 1: Grass top
// Tile 2: Dirt
// Tile 3: Stone Brick
// Tile 4: Wood / Platform
func createProceduralTilesetTexture() *gfx.Texture {
	img := image.NewRGBA(image.Rect(0, 0, 64, 16))

	cGrass := color.RGBA{R: 76, G: 175, B: 80, A: 255}
	cGrassDark := color.RGBA{R: 56, G: 142, B: 60, A: 255}
	cDirt := color.RGBA{R: 121, G: 85, B: 72, A: 255}
	cDirtDark := color.RGBA{R: 93, G: 64, B: 55, A: 255}
	cStone := color.RGBA{R: 140, G: 150, B: 160, A: 255}
	cStoneDark := color.RGBA{R: 90, G: 100, B: 110, A: 255}
	cWood := color.RGBA{R: 210, G: 140, B: 70, A: 255}
	cWoodDark := color.RGBA{R: 160, G: 100, B: 40, A: 255}

	// 1. Grass Tile (0..15)
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			if y < 4 {
				if (x+y)%2 == 0 {
					img.Set(x, y, cGrass)
				} else {
					img.Set(x, y, cGrassDark)
				}
			} else {
				if (x*3+y*7)%5 == 0 {
					img.Set(x, y, cDirtDark)
				} else {
					img.Set(x, y, cDirt)
				}
			}
		}
	}

	// 2. Pure Dirt Tile (16..31)
	for y := 0; y < 16; y++ {
		for x := 16; x < 32; x++ {
			if (x*7+y*11)%6 == 0 {
				img.Set(x, y, cDirtDark)
			} else {
				img.Set(x, y, cDirt)
			}
		}
	}

	// 3. Stone Brick Tile (32..47)
	for y := 0; y < 16; y++ {
		for x := 32; x < 48; x++ {
			relX := x - 32
			if y%8 == 0 || (y < 8 && relX == 8) || (y >= 8 && (relX == 0 || relX == 15)) {
				img.Set(x, y, cStoneDark)
			} else {
				img.Set(x, y, cStone)
			}
		}
	}

	// 4. Wood Platform Tile (48..63)
	for y := 0; y < 16; y++ {
		for x := 48; x < 64; x++ {
			if y == 0 || y == 15 || (x-48)%5 == 0 {
				img.Set(x, y, cWoodDark)
			} else {
				img.Set(x, y, cWood)
			}
		}
	}

	return gfx.NewTextureFromImage(ebiten.NewImageFromImage(img))
}

// Generate procedural 64x16 animation sheet:
// Frame 0-1: Idle (subtle bob)
// Frame 2-3: Run (legs step)
func createProceduralCharacterSheet() *gfx.Texture {
	img := image.NewRGBA(image.Rect(0, 0, 64, 16))

	cBody := color.RGBA{R: 54, G: 130, B: 242, A: 255}
	cEye := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	cPupil := color.RGBA{R: 20, G: 20, B: 40, A: 255}
	cCap := color.RGBA{R: 235, G: 87, B: 87, A: 255}
	cFeet := color.RGBA{R: 40, G: 40, B: 60, A: 255}

	drawFrame := func(frameIdx int, bobY int, legSpread int) {
		offsetX := frameIdx * 16

		// Cap
		for y := 2 + bobY; y <= 4 + bobY; y++ {
			for x := 4; x <= 11; x++ {
				img.Set(offsetX+x, y, cCap)
			}
		}
		// Body
		for y := 5 + bobY; y <= 12 + bobY; y++ {
			for x := 3; x <= 12; x++ {
				img.Set(offsetX+x, y, cBody)
			}
		}
		// Eyes
		img.Set(offsetX+5, 7+bobY, cEye)
		img.Set(offsetX+6, 7+bobY, cPupil)
		img.Set(offsetX+9, 7+bobY, cEye)
		img.Set(offsetX+10, 7+bobY, cPupil)

		// Feet
		footY := 13 + bobY
		if footY > 15 {
			footY = 15
		}
		img.Set(offsetX+4-legSpread, footY, cFeet)
		img.Set(offsetX+5-legSpread, footY, cFeet)
		img.Set(offsetX+10+legSpread, footY, cFeet)
		img.Set(offsetX+11+legSpread, footY, cFeet)
	}

	// Frame 0: Idle 1
	drawFrame(0, 0, 0)
	// Frame 1: Idle 2 (1px bob)
	drawFrame(1, 1, 0)
	// Frame 2: Run 1 (spread legs)
	drawFrame(2, 0, 1)
	// Frame 3: Run 2 (tight legs, 1px bob)
	drawFrame(3, 1, 0)

	return gfx.NewTextureFromImage(ebiten.NewImageFromImage(img))
}

type DemoGame struct {
	canvas     *display.Canvas
	camera     *display.Camera
	inputMap   *input.ActionMap
	tilemap    *tilemap.TileMap
	solidLayer *tilemap.Layer

	player      *gfx.AnimatedSprite
	playerBody  *physics.KinematicBody
	showColliders bool
}

func NewDemoGame(canvas *display.Canvas, am *input.ActionMap) *DemoGame {
	// 1. Create TileMap (40x20 tiles = 640x320 pixels)
	tilesetTex := createProceduralTilesetTexture()
	ts := tilemap.NewTileset(tilesetTex, 16, 16)
	tm := tilemap.NewTileMap(40, 20, ts)

	bgLayer := tm.AddLayer("Background")
	bgLayer.Opacity = 0.5
	solidLayer := tm.AddLayer("Solid")

	// Build a fun platformer level
	// Ground floor across the bottom
	for x := 0; x < 40; x++ {
		solidLayer.SetTile(x, 18, 1) // Grass
		solidLayer.SetTile(x, 19, 2) // Dirt
	}

	// Boundary walls
	for y := 0; y < 20; y++ {
		solidLayer.SetTile(0, y, 3)  // Stone left wall
		solidLayer.SetTile(39, y, 3) // Stone right wall
	}

	// Floating platforms
	for x := 6; x <= 11; x++ {
		solidLayer.SetTile(x, 14, 4) // Wood platform
	}
	for x := 15; x <= 20; x++ {
		solidLayer.SetTile(x, 11, 4) // Wood platform higher
	}
	for x := 24; x <= 30; x++ {
		solidLayer.SetTile(x, 14, 1) // Grass floating island
		solidLayer.SetTile(x, 15, 2)
	}

	// Stone steps
	solidLayer.SetTile(32, 17, 3)
	solidLayer.SetTile(33, 17, 3)
	solidLayer.SetTile(33, 16, 3)
	solidLayer.SetTile(34, 17, 3)
	solidLayer.SetTile(34, 16, 3)
	solidLayer.SetTile(34, 15, 3)

	// Background pillars
	for y := 10; y < 18; y++ {
		bgLayer.SetTile(8, y, 3)
		bgLayer.SetTile(17, y, 3)
		bgLayer.SetTile(27, y, 3)
	}

	// 2. Create Animated Player Sprite
	charSheet := createProceduralCharacterSheet()
	playerSprite := gfx.NewAnimatedSprite(charSheet)
	playerSprite.SetOriginPreset(gfx.OriginBottomCenter)

	// Configure animations
	idleAnim := gfx.NewAnimation("idle", gfx.LoopForever)
	idleAnim.AddFrame(image.Rect(0, 0, 16, 16), 0.4)
	idleAnim.AddFrame(image.Rect(16, 0, 32, 16), 0.4)

	walkAnim := gfx.NewAnimation("walk", gfx.LoopForever)
	walkAnim.AddFrame(image.Rect(32, 0, 48, 16), 0.12)
	walkAnim.AddFrame(image.Rect(48, 0, 64, 16), 0.12)

	jumpAnim := gfx.NewAnimation("jump", gfx.LoopOnce)
	jumpAnim.AddFrame(image.Rect(32, 0, 48, 16), 1.0)

	playerSprite.Player.AddAnimation(idleAnim)
	playerSprite.Player.AddAnimation(walkAnim)
	playerSprite.Player.AddAnimation(jumpAnim)
	playerSprite.Player.Play("idle", false)

	// 3. Setup Kinematic Physics Body (12x15 hitbox centered inside 16x16 sprite)
	body := physics.NewKinematicBody(48, 100, 12, 15)
	body.Gravity = 750.0
	body.MaxFallSpeed = 380.0

	// 4. Setup Camera
	cam := display.NewCamera(canvas.VirtualWidth(), canvas.VirtualHeight())
	cam.SetBounds(0, 0, float64(tm.PixelWidth()), float64(tm.PixelHeight()))

	return &DemoGame{
		canvas:     canvas,
		camera:     cam,
		inputMap:   am,
		tilemap:    tm,
		solidLayer: solidLayer,
		player:     playerSprite,
		playerBody: body,
	}
}

func (g *DemoGame) Update(dt float64) error {
	const moveSpeed = 120.0
	const jumpForce = -260.0

	// Horizontal input
	var moveDir float64
	if g.inputMap.IsPressed("left") {
		moveDir -= 1.0
		g.player.FlipX = true
	}
	if g.inputMap.IsPressed("right") {
		moveDir += 1.0
		g.player.FlipX = false
	}

	g.playerBody.VelocityX = moveDir * moveSpeed

	// Jump input (Space or W / Up)
	if g.inputMap.IsJustPressed("jump") && g.playerBody.OnGround {
		g.playerBody.VelocityY = jumpForce
	}

	// Physics Step: Apply gravity and resolve tile collision
	g.playerBody.ApplyGravity(dt)

	isSolid := func(tx, ty int) bool {
		return g.solidLayer.GetTile(tx, ty) > 0
	}
	g.playerBody.MoveAndSlide(dt, isSolid, 16)

	// Sync animation state with movement
	if !g.playerBody.OnGround {
		g.player.Player.Play("jump", false)
	} else if moveDir != 0 {
		g.player.Player.Play("walk", false)
	} else {
		g.player.Player.Play("idle", false)
	}
	g.player.Update(dt)

	// Sync sprite position to physics box (origin bottom center)
	g.player.X = g.playerBody.Box.CenterX()
	g.player.Y = g.playerBody.Box.Bottom()

	// Screen shake on Shake key
	if g.inputMap.IsJustPressed("shake") {
		g.camera.Shake(5.0, 0.3)
	}

	// Toggle collider boundary rendering on C or F1
	if g.inputMap.IsJustPressed("toggle_colliders") {
		g.showColliders = !g.showColliders
	}

	// Camera smoothly tracks player
	g.camera.Follow(g.player.X, g.player.Y-20, 0.12)
	g.camera.Update(dt)

	return nil
}

func (g *DemoGame) Draw(canvas *display.Canvas) {
	// 1. Draw tilemap with camera viewport culling
	g.tilemap.Draw(canvas, g.camera)

	// 2. Draw animated player sprite
	g.player.Draw(canvas, g.camera)

	// 3. Draw debug collider boundaries if enabled
	if g.showColliders {
		// Draw solid tile boundary boxes (neon green outline with translucent fill)
		physics.DrawTileColliders(
			canvas,
			g.camera,
			func(tx, ty int) bool { return g.solidLayer.GetTile(tx, ty) > 0 },
			16,
			0, 0, g.tilemap.Width, g.tilemap.Height,
			color.RGBA{R: 0, G: 255, B: 100, A: 220}, // stroke
			color.RGBA{R: 0, G: 255, B: 100, A: 40},  // fill
		)

		// Draw player kinematic AABB (bright yellow outline with translucent fill)
		physics.DrawAABB(
			canvas,
			g.camera,
			g.playerBody.Box,
			color.RGBA{R: 255, G: 230, B: 0, A: 255}, // stroke
			color.RGBA{R: 255, G: 230, B: 0, A: 70},  // fill
		)
	}

	// 4. Draw HUD
	mx, my, inCanvas := input.CursorCanvasPosition(canvas)
	colStatus := "OFF"
	if g.showColliders {
		colStatus = "ON"
	}
	hudText := fmt.Sprintf("PixelDream-2D [Phase 2: Animated Sprite & Tilemap]\nPos: (%.0f, %.0f) | Grounded: %v | Wall: %v\nColliders: %s (Press C to toggle) | Mouse: (%d, %d) [In: %v]\nMove: A/D or Arrows | Jump: Space/W | Shake: X",
		g.playerBody.Box.X, g.playerBody.Box.Y, g.playerBody.OnGround, g.playerBody.OnWall, colStatus, mx, my, inCanvas)
	ebitenutil.DebugPrint(canvas.Image(), hudText)
}

func main() {
	cfg := engine.DefaultConfig()
	cfg.Title = "PixelDream-2D - Phase 2 Animation & Tilemap Demo"
	cfg.VirtualWidth = 320
	cfg.VirtualHeight = 180
	cfg.WindowWidth = 1280
	cfg.WindowHeight = 720
	cfg.ScaleMode = display.ScaleModeInteger

	runner := engine.NewRunner(nil, cfg)

	// Set up input bindings
	am := runner.InputMap()
	am.BindKey("left", ebiten.KeyArrowLeft, ebiten.KeyA)
	am.BindKey("right", ebiten.KeyArrowRight, ebiten.KeyD)
	am.BindKey("jump", ebiten.KeySpace, ebiten.KeyW, ebiten.KeyArrowUp)
	am.BindKey("shake", ebiten.KeyX)
	am.BindKey("toggle_colliders", ebiten.KeyC, ebiten.KeyF1)

	demo := NewDemoGame(runner.Canvas(), am)

	ebiten.SetWindowTitle(cfg.Title)
	ebiten.SetWindowSize(cfg.WindowWidth, cfg.WindowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	gameWrapper := &gameWrapper{
		runner: runner,
		game:   demo,
	}

	if err := ebiten.RunGame(gameWrapper); err != nil {
		log.Fatalf("Engine failed: %v", err)
	}
}

type gameWrapper struct {
	runner *engine.Runner
	game   *DemoGame
}

func (w *gameWrapper) Update() error {
	return w.game.Update(1.0 / 60.0)
}

func (w *gameWrapper) Draw(screen *ebiten.Image) {
	w.runner.Canvas().Clear()
	w.game.Draw(w.runner.Canvas())
	w.runner.Canvas().RenderToScreen(screen)
}

func (w *gameWrapper) Layout(outsideWidth, outsideHeight int) (int, int) {
	return w.runner.Layout(outsideWidth, outsideHeight)
}
