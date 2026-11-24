package editor

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/zapturk/PixelDream-2D/pkg/display"
	"github.com/zapturk/PixelDream-2D/pkg/physics"
	"github.com/zapturk/PixelDream-2D/pkg/tilemap"
)

// Mode represents whether the editor is in editing mode or playtesting mode.
type Mode int

const (
	ModeEdit Mode = iota
	ModePixelArt
	ModePlay
)

// Editor manages the in-engine level editor shell.
type Editor struct {
	Mode             Mode
	TileMap          *tilemap.TileMap
	ActiveLayerIndex int
	SelectedTileID   int
	ActiveTool       ToolType

	ShowGrid      bool
	ShowColliders bool

	Camera   *display.Camera
	UI       *UIContext
	FilePath string

	// Rect tool drag state
	isRectDragging bool
	rectStartX     int
	rectStartY     int

	// Camera panning state
	isPanning  bool
	panStartX  int
	panStartY  int
	camStartX  float64
	camStartY  float64

	// Playtest callbacks
	OnPlayUpdate func(dt float64)
	OnPlayDraw   func(canvas *display.Canvas, cam *display.Camera)

	// Pixel Art Editor callbacks
	OnPixelEditorUpdate func(dt float64, screenWidth, screenHeight int)
	OnPixelEditorDraw   func(screen *ebiten.Image)

	// Layout dimensions
	windowWidth  int
	windowHeight int
	viewportRect image.Rectangle
	toolbarRect  image.Rectangle
	sidebarRect  image.Rectangle

	statusMessage string
	statusTimer   float64
}

// NewEditor creates a new in-engine editor for the given tilemap.
func NewEditor(tm *tilemap.TileMap, camera *display.Camera) *Editor {
	if camera == nil {
		camera = display.NewCamera(320, 180)
	}

	return &Editor{
		Mode:             ModeEdit,
		TileMap:          tm,
		ActiveLayerIndex: 0,
		SelectedTileID:   1,
		ActiveTool:       ToolPencil,
		ShowGrid:         true,
		ShowColliders:    false,
		Camera:           camera,
		UI:               NewUIContext(),
		FilePath:         "level.json",
		windowWidth:      1280,
		windowHeight:     720,
	}
}

// SetStatus sets a temporary notification message in the status bar.
func (e *Editor) SetStatus(msg string, durationSeconds float64) {
	e.statusMessage = msg
	e.statusTimer = durationSeconds
}

// Update handles editor input, tool actions, and playtest updates.
func (e *Editor) Update(dt float64, canvas *display.Canvas) {
	if e.statusTimer > 0 {
		e.statusTimer -= dt
		if e.statusTimer <= 0 {
			e.statusMessage = ""
		}
	}

	// Global hotkeys: F5 toggles Play / Edit mode
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		e.TogglePlayMode()
	}

	// Tab toggles between Map Editor and Pixel Art Editor (when not in Play mode)
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) && e.Mode != ModePlay {
		if e.Mode == ModeEdit {
			e.Mode = ModePixelArt
			e.SetStatus("PIXEL ART EDITOR (Tab for Map)", 2.0)
		} else {
			e.Mode = ModeEdit
			e.SetStatus("MAP EDITOR (Tab for Pixel Art)", 2.0)
		}
	}

	// Delegate based on mode
	if e.Mode == ModePlay {
		if e.OnPlayUpdate != nil {
			e.OnPlayUpdate(dt)
		}
		return
	}

	if e.Mode == ModePixelArt {
		sw, sh := ebiten.WindowSize()
		if e.OnPixelEditorUpdate != nil {
			e.OnPixelEditorUpdate(dt, sw, sh)
		}
		return
	}

	// Edit Mode Update
	mx, my := ebiten.CursorPosition()
	e.UI.Update(mx, my)

	// Quick hotkeys in Edit mode
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		e.ShowGrid = !e.ShowGrid
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		e.ShowColliders = !e.ShowColliders
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		e.ActiveTool = ToolPencil
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		e.ActiveTool = ToolEraser
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		e.ActiveTool = ToolRect
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		e.ActiveTool = ToolBucket
	}
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		e.ActiveTool = ToolPicker
	}

	// Save hotkey (Ctrl+S or Cmd+S)
	ctrlOrCmd := ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta)
	if ctrlOrCmd && inpututil.IsKeyJustPressed(ebiten.KeyS) {
		e.SaveMap()
	}

	// Check if mouse is inside the canvas viewport (not on UI panels)
	canvasX, canvasY, inCanvas := canvas.ScreenToCanvas(mx, my)

	// Camera Panning with Middle Click or Space+LeftClick
	isPanButton := ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) ||
		(ebiten.IsKeyPressed(ebiten.KeySpace) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft))

	if isPanButton {
		if !e.isPanning {
			e.isPanning = true
			e.panStartX = mx
			e.panStartY = my
			e.camStartX = e.Camera.X
			e.camStartY = e.Camera.Y
		} else {
			dx := float64(mx - e.panStartX)
			dy := float64(my - e.panStartY)
			scaleX, scaleY, _, _ := canvas.LayoutMetrics()
			if scaleX > 0 && scaleY > 0 {
				e.Camera.X = e.camStartX - (dx / (scaleX * e.Camera.Zoom))
				e.Camera.Y = e.camStartY - (dy / (scaleY * e.Camera.Zoom))
			}
		}
		return
	} else {
		e.isPanning = false
	}

	// Handle Viewport Tile Editing
	if inCanvas && !e.isPanning && e.TileMap != nil && len(e.TileMap.Layers) > 0 {
		worldX, worldY := e.Camera.ScreenToWorld(float64(canvasX), float64(canvasY))
		tileX, tileY := e.TileMap.WorldToTile(worldX, worldY)
		activeLayer := e.ActiveLayer()

		// Tool execution
		switch e.ActiveTool {
		case ToolPencil:
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				ApplyPencil(activeLayer, tileX, tileY, e.SelectedTileID)
			}
		case ToolEraser:
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				ApplyEraser(activeLayer, tileX, tileY)
			}
		case ToolPicker:
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				picked := ApplyPicker(activeLayer, tileX, tileY)
				if picked > 0 {
					e.SelectedTileID = picked
					e.ActiveTool = ToolPencil
					e.SetStatus(fmt.Sprintf("Picked Tile #%d", picked), 1.5)
				}
			}
		case ToolBucket:
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				ApplyBucket(activeLayer, tileX, tileY, e.SelectedTileID)
			}
		case ToolRect:
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				e.isRectDragging = true
				e.rectStartX = tileX
				e.rectStartY = tileY
			}
			if e.isRectDragging && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				e.isRectDragging = false
				ApplyRect(activeLayer, e.rectStartX, e.rectStartY, tileX, tileY, e.SelectedTileID)
			}
		}
	}
}

// TogglePlayMode switches between editing and live playtesting.
func (e *Editor) TogglePlayMode() {
	if e.Mode == ModeEdit {
		e.Mode = ModePlay
		e.SetStatus("PLAY MODE (F5 to exit)", 2.0)
	} else {
		e.Mode = ModeEdit
		e.SetStatus("EDIT MODE (F5 to play)", 2.0)
	}
}

// ActiveLayer returns the currently selected layer.
func (e *Editor) ActiveLayer() *tilemap.Layer {
	if e.TileMap == nil || len(e.TileMap.Layers) == 0 {
		return nil
	}
	if e.ActiveLayerIndex < 0 || e.ActiveLayerIndex >= len(e.TileMap.Layers) {
		e.ActiveLayerIndex = 0
	}
	return e.TileMap.Layers[e.ActiveLayerIndex]
}

// SaveMap writes the current tilemap to e.FilePath.
func (e *Editor) SaveMap() {
	if e.TileMap == nil {
		return
	}
	if err := e.TileMap.SaveToFile(e.FilePath); err != nil {
		e.SetStatus(fmt.Sprintf("Save Failed: %v", err), 3.0)
	} else {
		e.SetStatus(fmt.Sprintf("Saved to %s", e.FilePath), 2.5)
	}
}

// LoadMap reads the tilemap from e.FilePath.
func (e *Editor) LoadMap(tileset *tilemap.Tileset) {
	tm, err := tilemap.LoadFromFile(e.FilePath, tileset)
	if err != nil {
		e.SetStatus(fmt.Sprintf("Load Failed: %v", err), 3.0)
	} else {
		e.TileMap = tm
		e.SetStatus(fmt.Sprintf("Loaded %s", e.FilePath), 2.5)
	}
}

// Draw renders the editor viewports, tilemaps, grid, and in-engine UI panels.
func (e *Editor) Draw(canvas *display.Canvas, screen *ebiten.Image) {
	// 1. In Play Mode, render gameplay directly
	if e.Mode == ModePlay {
		if e.OnPlayDraw != nil {
			e.OnPlayDraw(canvas, e.Camera)
		}
		// Play HUD banner on screen
		bannerRect := image.Rect(0, 0, screen.Bounds().Dx(), 24)
		vector.DrawFilledRect(screen, 0, 0, float32(bannerRect.Dx()), 24, color.RGBA{R: 46, G: 204, B: 113, A: 220}, false)
		ebitenutil.DebugPrintAt(screen, "PLAYTESTING - Press F5 to return to Editor", 12, 5)
		return
	}

	// 2. In Pixel Art Mode, delegate rendering to pixel editor
	if e.Mode == ModePixelArt {
		if e.OnPixelEditorDraw != nil {
			e.OnPixelEditorDraw(screen)
		}
		e.drawTopBar(screen)
		return
	}

	// 2. In Edit Mode: Render world on the virtual canvas
	if e.TileMap != nil {
		e.TileMap.Draw(canvas, e.Camera)

		// Grid Overlay
		if e.ShowGrid {
			pw := float64(e.TileMap.PixelWidth())
			ph := float64(e.TileMap.PixelHeight())
			tw := float64(e.TileMap.TileWidth)

			// Transform tilemap grid lines through camera
			minX, minY := e.Camera.WorldToScreen(0, 0)
			maxX, maxY := e.Camera.WorldToScreen(pw, ph)
			w := maxX - minX
			h := maxY - minY

			DrawGridOverlay(canvas.Image(), minX, minY, w, h, tw*e.Camera.Zoom, color.RGBA{R: 255, G: 255, B: 255, A: 40})
		}

		// Collider Boundaries Overlay
		if e.ShowColliders {
			physics.DrawTileColliders(
				canvas,
				e.Camera,
				func(tx, ty int) bool {
					layer := e.ActiveLayer()
					return layer != nil && layer.GetTile(tx, ty) > 0
				},
				e.TileMap.TileWidth,
				0, 0, e.TileMap.Width, e.TileMap.Height,
				color.RGBA{R: 0, G: 255, B: 120, A: 200},
				color.RGBA{R: 0, G: 255, B: 120, A: 45},
			)
		}

		// Cursor Tile Hover / Rect Preview
		mx, my := ebiten.CursorPosition()
		cx, cy, inCanvas := canvas.ScreenToCanvas(mx, my)
		if inCanvas {
			wx, wy := e.Camera.ScreenToWorld(float64(cx), float64(cy))
			hoverTX, hoverTY := e.TileMap.WorldToTile(wx, wy)

			if e.isRectDragging && e.ActiveTool == ToolRect {
				// Draw dragging rectangle highlight
				minX := min(e.rectStartX, hoverTX)
				maxX := max(e.rectStartX, hoverTX)
				minY := min(e.rectStartY, hoverTY)
				maxY := max(e.rectStartY, hoverTY)

				for y := minY; y <= maxY; y++ {
					for x := minX; x <= maxX; x++ {
						rBox := physics.NewAABB(float64(x*e.TileMap.TileWidth), float64(y*e.TileMap.TileHeight), float64(e.TileMap.TileWidth), float64(e.TileMap.TileHeight))
						physics.DrawAABB(canvas, e.Camera, rBox, color.RGBA{R: 74, G: 144, B: 226, A: 255}, color.RGBA{R: 74, G: 144, B: 226, A: 80})
					}
				}
			} else if hoverTX >= 0 && hoverTX < e.TileMap.Width && hoverTY >= 0 && hoverTY < e.TileMap.Height {
				// Single tile hover box
				tBox := physics.NewAABB(float64(hoverTX*e.TileMap.TileWidth), float64(hoverTY*e.TileMap.TileHeight), float64(e.TileMap.TileWidth), float64(e.TileMap.TileHeight))
				physics.DrawAABB(canvas, e.Camera, tBox, color.RGBA{R: 255, G: 255, B: 255, A: 200}, color.RGBA{R: 255, G: 255, B: 255, A: 35})
			}
		}
	}

	// 3. Render Canvas onto Screen Window
	canvas.RenderToScreen(screen)

	// 4. Render Editor In-Engine GUI Overlay on Screen
	e.drawEditorUI(screen, canvas)
}

func (e *Editor) drawTopBar(screen *ebiten.Image) {
	sw := screen.Bounds().Dx()
	toolbarH := 36
	vector.DrawFilledRect(screen, 0, 0, float32(sw), float32(toolbarH), e.UI.Theme.PanelHeader, false)
	vector.StrokeLine(screen, 0, float32(toolbarH), float32(sw), float32(toolbarH), 1.0, e.UI.Theme.Border, false)

	btnX := 8
	// 1. Play Button
	if e.UI.Button(screen, image.Rect(btnX, 6, btnX+72, 30), "Play (F5)", false) {
		e.TogglePlayMode()
	}
	btnX += 80

	// 2. Mode Tabs: Map Editor vs Pixel Art
	if e.UI.Button(screen, image.Rect(btnX, 6, btnX+80, 30), "Map Editor", e.Mode == ModeEdit) {
		e.Mode = ModeEdit
	}
	btnX += 86

	if e.UI.Button(screen, image.Rect(btnX, 6, btnX+100, 30), "Pixel Art (Tab)", e.Mode == ModePixelArt) {
		e.Mode = ModePixelArt
	}
	btnX += 108

	// Mode-specific toolbar buttons
	if e.Mode == ModeEdit {
		// Save Button
		if e.UI.Button(screen, image.Rect(btnX, 6, btnX+50, 30), "Save", false) {
			e.SaveMap()
		}
		btnX += 56

		// Tools Buttons
		tools := []struct {
			name string
			tool ToolType
		}{
			{"Pencil", ToolPencil},
			{"Eraser", ToolEraser},
			{"Rect", ToolRect},
			{"Bucket", ToolBucket},
			{"Picker", ToolPicker},
		}
		for _, t := range tools {
			active := e.ActiveTool == t.tool
			if e.UI.Button(screen, image.Rect(btnX, 6, btnX+54, 30), t.name, active) {
				e.ActiveTool = t.tool
			}
			btnX += 58
		}

		btnX += 8
		if e.UI.Button(screen, image.Rect(btnX, 6, btnX+48, 30), "Grid", e.ShowGrid) {
			e.ShowGrid = !e.ShowGrid
		}
		btnX += 52
		if e.UI.Button(screen, image.Rect(btnX, 6, btnX+70, 30), "Colliders", e.ShowColliders) {
			e.ShowColliders = !e.ShowColliders
		}
		btnX += 76
	} else if e.Mode == ModePixelArt {
		ebitenutil.DebugPrintAt(screen, "Pixel Art Mode (B: Pencil, E: Eraser, G: Bucket, L: Line, U: Rect, Ctrl+Z: Undo)", btnX+12, 11)
	}

	// Status message in toolbar
	if e.statusMessage != "" {
		ebitenutil.DebugPrintAt(screen, e.statusMessage, btnX+12, 11)
	}
}

func (e *Editor) drawEditorUI(screen *ebiten.Image, canvas *display.Canvas) {
	sw := screen.Bounds().Dx()
	sh := screen.Bounds().Dy()
	toolbarH := 36

	// Top toolbar
	e.drawTopBar(screen)

	// --- Right Sidebar: Tileset & Layers ---
	sidebarW := 210
	sidebarX := sw - sidebarW
	sidebarRect := image.Rect(sidebarX, toolbarH, sw, sh)
	e.UI.DrawPanel(screen, sidebarRect, "Tileset & Layers")

	contentY := toolbarH + 24

	// Tileset Palette
	ebitenutil.DebugPrintAt(screen, "Tileset Palette:", sidebarX+8, contentY)
	contentY += 16

	if e.TileMap != nil && e.TileMap.Tileset != nil && e.TileMap.Tileset.Texture != nil {
		ts := e.TileMap.Tileset
		tsImg := ts.Texture.Image()
		totalTiles := ts.TotalTiles()

		tilesPerRow := 6
		tileSize := 28

		for id := 1; id <= totalTiles; id++ {
			row := (id - 1) / tilesPerRow
			col := (id - 1) % tilesPerRow

			tx := sidebarX + 10 + col*(tileSize+4)
			ty := contentY + row*(tileSize+4)

			// Draw tile preview
			srcRect := ts.TileRect(id)
			if !srcRect.Empty() {
				tileSub := tsImg.SubImage(srcRect).(*ebiten.Image)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(float64(tileSize)/float64(ts.TileWidth), float64(tileSize)/float64(ts.TileHeight))
				op.GeoM.Translate(float64(tx), float64(ty))
				screen.DrawImage(tileSub, op)
			}

			// Border & selection
			borderCol := e.UI.Theme.Border
			if id == e.SelectedTileID {
				borderCol = e.UI.Theme.BorderActive
				vector.StrokeRect(screen, float32(tx-1), float32(ty-1), float32(tileSize+2), float32(tileSize+2), 2.0, borderCol, false)
			} else {
				vector.StrokeRect(screen, float32(tx), float32(ty), float32(tileSize), float32(tileSize), 1.0, borderCol, false)
			}

			// Click to select tile
			mx, my := ebiten.CursorPosition()
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if mx >= tx && mx < tx+tileSize && my >= ty && my < ty+tileSize {
					e.SelectedTileID = id
					e.ActiveTool = ToolPencil
				}
			}
		}

		numRows := int(math.Ceil(float64(totalTiles) / float64(tilesPerRow)))
		contentY += numRows*(tileSize+4) + 16
	}

	// Layers Section
	ebitenutil.DebugPrintAt(screen, "Layers:", sidebarX+8, contentY)
	contentY += 16

	if e.TileMap != nil {
		for i, layer := range e.TileMap.Layers {
			isActive := i == e.ActiveLayerIndex
			layerRect := image.Rect(sidebarX+10, contentY, sidebarX+sidebarW-10, contentY+24)

			label := layer.Name
			if !layer.Visible {
				label += " (Hidden)"
			}

			if e.UI.Button(screen, layerRect, label, isActive) {
				e.ActiveLayerIndex = i
			}

			// Visibility toggle mini button
			eyeRect := image.Rect(sidebarX+sidebarW-32, contentY+3, sidebarX+sidebarW-14, contentY+21)
			eyeLabel := "V"
			if !layer.Visible {
				eyeLabel = "-"
			}
			if e.UI.Button(screen, eyeRect, eyeLabel, layer.Visible) {
				layer.Visible = !layer.Visible
			}

			contentY += 28
		}
	}

	// --- Bottom Status Bar ---
	statusH := 24
	statusY := sh - statusH
	vector.DrawFilledRect(screen, 0, float32(statusY), float32(sw), float32(statusH), e.UI.Theme.PanelHeader, false)
	vector.StrokeLine(screen, 0, float32(statusY), float32(sw), float32(statusY), 1.0, e.UI.Theme.Border, false)

	mx, my := ebiten.CursorPosition()
	cx, cy, inCanvas := canvas.ScreenToCanvas(mx, my)
	infoText := fmt.Sprintf("Tool: %s | Tile: #%d | Layer: %s | Middle-Drag: Pan",
		e.ActiveTool, e.SelectedTileID, e.ActiveLayer().Name)

	if inCanvas && e.TileMap != nil {
		wx, wy := e.Camera.ScreenToWorld(float64(cx), float64(cy))
		tx, ty := e.TileMap.WorldToTile(wx, wy)
		infoText += fmt.Sprintf(" | Tile: (%d, %d) World: (%.0f, %.0f)", tx, ty, wx, wy)
	}
	ebitenutil.DebugPrintAt(screen, infoText, 8, statusY+5)
}
