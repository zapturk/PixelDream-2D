package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zapturk/PixelDream-2D/pkg/display"
	"github.com/zapturk/PixelDream-2D/pkg/editor"
	"github.com/zapturk/PixelDream-2D/pkg/engine"
	"github.com/zapturk/PixelDream-2D/pkg/gfx"
	"github.com/zapturk/PixelDream-2D/pkg/input"
	"github.com/zapturk/PixelDream-2D/pkg/physics"
	"github.com/zapturk/PixelDream-2D/pkg/pixeleditor"
	"github.com/zapturk/PixelDream-2D/pkg/tilemap"
)

// Generate procedural tileset with 6 tiles (16x16 each):
// 1: Grass Top
// 2: Dirt
// 3: Stone Brick
// 4: Wood Platform
// 5: Water / Hazard
// 6: Metal Spikes / Trap
func createEditorTilesetTexture() *gfx.Texture {
	img := image.NewRGBA(image.Rect(0, 0, 96, 16))

	cGrass := color.RGBA{R: 76, G: 175, B: 80, A: 255}
	cGrassDark := color.RGBA{R: 56, G: 142, B: 60, A: 255}
	cDirt := color.RGBA{R: 121, G: 85, B: 72, A: 255}
	cDirtDark := color.RGBA{R: 93, G: 64, B: 55, A: 255}
	cStone := color.RGBA{R: 140, G: 150, B: 160, A: 255}
	cStoneDark := color.RGBA{R: 90, G: 100, B: 110, A: 255}
	cWood := color.RGBA{R: 210, G: 140, B: 70, A: 255}
	cWoodDark := color.RGBA{R: 160, G: 100, B: 40, A: 255}
	cWater := color.RGBA{R: 41, G: 128, B: 185, A: 220}
	cWaterLight := color.RGBA{R: 52, G: 152, B: 219, A: 240}
	cSpike := color.RGBA{R: 189, G: 195, B: 199, A: 255}
	cSpikeDark := color.RGBA{R: 127, G: 140, B: 141, A: 255}

	// 1. Grass
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

	// 2. Dirt
	for y := 0; y < 16; y++ {
		for x := 16; x < 32; x++ {
			if (x*7+y*11)%6 == 0 {
				img.Set(x, y, cDirtDark)
			} else {
				img.Set(x, y, cDirt)
			}
		}
	}

	// 3. Stone
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

	// 4. Wood
	for y := 0; y < 16; y++ {
		for x := 48; x < 64; x++ {
			if y == 0 || y == 15 || (x-48)%5 == 0 {
				img.Set(x, y, cWoodDark)
			} else {
				img.Set(x, y, cWood)
			}
		}
	}

	// 5. Water
	for y := 0; y < 16; y++ {
		for x := 64; x < 80; x++ {
			if (x+y)%3 == 0 {
				img.Set(x, y, cWaterLight)
			} else {
				img.Set(x, y, cWater)
			}
		}
	}

	// 6. Spikes
	for y := 0; y < 16; y++ {
		for x := 80; x < 96; x++ {
			relX := x - 80
			triHeight := 15 - int(float64(abs(relX-8))*1.8)
			if y >= triHeight {
				if relX < 8 {
					img.Set(x, y, cSpike)
				} else {
					img.Set(x, y, cSpikeDark)
				}
			}
		}
	}

	return gfx.NewTextureFromImage(ebiten.NewImageFromImage(img))
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func createPlayerCharacter() (*gfx.AnimatedSprite, *physics.KinematicBody) {
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

	drawFrame(0, 0, 0)
	drawFrame(1, 1, 0)
	drawFrame(2, 0, 1)
	drawFrame(3, 1, 0)

	tex := gfx.NewTextureFromImage(ebiten.NewImageFromImage(img))
	sprite := gfx.NewAnimatedSprite(tex)
	sprite.SetOriginPreset(gfx.OriginBottomCenter)

	idle := gfx.NewAnimation("idle", gfx.LoopForever)
	idle.AddFrame(image.Rect(0, 0, 16, 16), 0.4)
	idle.AddFrame(image.Rect(16, 0, 32, 16), 0.4)

	walk := gfx.NewAnimation("walk", gfx.LoopForever)
	walk.AddFrame(image.Rect(32, 0, 48, 16), 0.12)
	walk.AddFrame(image.Rect(48, 0, 64, 16), 0.12)

	jump := gfx.NewAnimation("jump", gfx.LoopOnce)
	jump.AddFrame(image.Rect(32, 0, 48, 16), 1.0)

	sprite.Player.AddAnimation(idle)
	sprite.Player.AddAnimation(walk)
	sprite.Player.AddAnimation(jump)
	sprite.Player.Play("idle", false)

	body := physics.NewKinematicBody(48, 100, 12, 15)
	body.Gravity = 750.0
	body.MaxFallSpeed = 380.0

	return sprite, body
}

type EditorApp struct {
	canvas     *display.Canvas
	editor     *editor.Editor
	inputMap   *input.ActionMap
	player     *gfx.AnimatedSprite
	playerBody *physics.KinematicBody
}

func NewEditorApp(canvas *display.Canvas, am *input.ActionMap) *EditorApp {
	// Initialize Tileset & Map
	tilesetTex := createEditorTilesetTexture()
	ts := tilemap.NewTileset(tilesetTex, 16, 16)
	tm := tilemap.NewTileMap(40, 20, ts)

	bg := tm.AddLayer("Background")
	bg.Opacity = 0.5
	solid := tm.AddLayer("Solid")

	// Pre-seed some default level terrain
	for x := 0; x < 40; x++ {
		solid.SetTile(x, 18, 1) // Grass
		solid.SetTile(x, 19, 2) // Dirt
	}
	for y := 0; y < 20; y++ {
		solid.SetTile(0, y, 3)  // Stone
		solid.SetTile(39, y, 3) // Stone
	}
	for x := 6; x <= 11; x++ {
		solid.SetTile(x, 14, 4) // Wood
	}

	cam := display.NewCamera(canvas.VirtualWidth(), canvas.VirtualHeight())
	cam.SetBounds(0, 0, float64(tm.PixelWidth()), float64(tm.PixelHeight()))

	ed := editor.NewEditor(tm, cam)
	ed.FilePath = "level_data.json"

	player, playerBody := createPlayerCharacter()

	// Initialize Pixel Art Editor (16x16 sprite animation canvas)
	pixelEd := pixeleditor.NewPixelArtEditor(16, 16, ed.UI)

	// Pre-seed an animated 16x16 golden coin in the pixel editor
	cGold := color.RGBA{R: 255, G: 236, B: 39, A: 255}
	cDarkGold := color.RGBA{R: 255, G: 163, B: 0, A: 255}
	cShine := color.RGBA{R: 255, G: 241, B: 232, A: 255}

	// Frame 1: Full round coin
	for y := 4; y <= 11; y++ {
		for x := 4; x <= 11; x++ {
			if (x == 4 && (y == 4 || y == 11)) || (x == 11 && (y == 4 || y == 11)) {
				continue
			}
			pixelEd.Doc.SetPixel(x, y, cGold)
		}
	}
	for y := 5; y <= 10; y++ {
		pixelEd.Doc.SetPixel(6, y, cDarkGold)
		pixelEd.Doc.SetPixel(9, y, cDarkGold)
	}
	pixelEd.Doc.SetPixel(5, 5, cShine)

	// Frame 2: Slightly turned coin
	pixelEd.Doc.AddFrame()
	for y := 4; y <= 11; y++ {
		for x := 5; x <= 10; x++ {
			pixelEd.Doc.SetPixel(x, y, cGold)
		}
	}
	pixelEd.Doc.SetPixel(6, 6, cShine)

	// Frame 3: Side edge coin
	pixelEd.Doc.AddFrame()
	for y := 4; y <= 11; y++ {
		for x := 7; x <= 8; x++ {
			pixelEd.Doc.SetPixel(x, y, cDarkGold)
		}
	}

	// Reset active frame to 0
	pixelEd.Doc.ActiveFrameIndex = 0
	pixelEd.IsPlaying = true

	// Hook Pixel Art Editor callbacks into Editor Shell
	ed.OnPixelEditorUpdate = func(dt float64, sw, sh int) {
		pixelEd.Update(dt, sw, sh)
	}
	ed.OnPixelEditorDraw = func(screen *ebiten.Image) {
		pixelEd.Draw(screen)
	}

	app := &EditorApp{
		canvas:     canvas,
		editor:     ed,
		inputMap:   am,
		player:     player,
		playerBody: playerBody,
	}

	// Play Mode update callback: runs live game physics and controls!
	ed.OnPlayUpdate = func(dt float64) {
		const moveSpeed = 120.0
		const jumpForce = -260.0

		var moveDir float64
		if am.IsPressed("left") {
			moveDir -= 1.0
			player.FlipX = true
		}
		if am.IsPressed("right") {
			moveDir += 1.0
			player.FlipX = false
		}
		playerBody.VelocityX = moveDir * moveSpeed

		if am.IsJustPressed("jump") && playerBody.OnGround {
			playerBody.VelocityY = jumpForce
		}

		playerBody.ApplyGravity(dt)

		// Check collision dynamically against the live edited Solid layer!
		isSolid := func(tx, ty int) bool {
			sLayer := tm.GetLayer("Solid")
			return sLayer != nil && sLayer.GetTile(tx, ty) > 0
		}
		playerBody.MoveAndSlide(dt, isSolid, 16)

		if !playerBody.OnGround {
			player.Player.Play("jump", false)
		} else if moveDir != 0 {
			player.Player.Play("walk", false)
		} else {
			player.Player.Play("idle", false)
		}
		player.Update(dt)

		player.X = playerBody.Box.CenterX()
		player.Y = playerBody.Box.Bottom()

		cam.Follow(player.X, player.Y-20, 0.12)
		cam.Update(dt)
	}

	// Play Mode draw callback
	ed.OnPlayDraw = func(cv *display.Canvas, c *display.Camera) {
		tm.Draw(cv, c)
		player.Draw(cv, c)
	}

	return app
}

func (a *EditorApp) Update() error {
	a.editor.Update(1.0/60.0, a.canvas)
	return nil
}

func (a *EditorApp) Draw(screen *ebiten.Image) {
	a.canvas.Clear()
	a.editor.Draw(a.canvas, screen)
}

func (a *EditorApp) Layout(outsideWidth, outsideHeight int) (int, int) {
	a.canvas.UpdateLayout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func main() {
	cfg := engine.DefaultConfig()
	cfg.Title = "PixelDream-2D - In-Engine Map Editor (Phase 3)"
	cfg.VirtualWidth = 320
	cfg.VirtualHeight = 180
	cfg.WindowWidth = 1280
	cfg.WindowHeight = 720
	cfg.ScaleMode = display.ScaleModeInteger

	canvas := display.NewCanvas(cfg.VirtualWidth, cfg.VirtualHeight)
	canvas.SetScaleMode(cfg.ScaleMode)

	am := input.NewActionMap()
	am.BindKey("left", ebiten.KeyArrowLeft, ebiten.KeyA)
	am.BindKey("right", ebiten.KeyArrowRight, ebiten.KeyD)
	am.BindKey("jump", ebiten.KeySpace, ebiten.KeyW, ebiten.KeyArrowUp)

	app := NewEditorApp(canvas, am)

	ebiten.SetWindowTitle(cfg.Title)
	ebiten.SetWindowSize(cfg.WindowWidth, cfg.WindowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(app); err != nil {
		log.Fatalf("Editor failed: %v", err)
	}
}
