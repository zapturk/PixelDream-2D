package pixeleditor

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/zapturk/PixelDream-2D/pkg/editor"
)

// PixelArtEditor is the interactive pixel art drawing and animation editor.
type PixelArtEditor struct {
	Doc        *Document
	ActiveTool Tool
	UI         *editor.UIContext

	// Viewport & Zoom
	Zoom          int // e.g. 16, 20, 24, 32
	CanvasOffsetX int
	CanvasOffsetY int
	ShowPixelGrid bool
	OnionSkin     bool

	// Drawing interaction
	isDrawing   bool
	hasLastDraw bool
	lastDrawX   int
	lastDrawY   int
	dragStartX  int
	dragStartY  int

	// Animation Preview
	IsPlaying         bool
	PlaybackFPS       float64
	playbackTimer     float64
	previewFrameIndex int

	statusMsg   string
	statusTimer float64
}

// NewPixelArtEditor creates a new pixel art editor.
func NewPixelArtEditor(width, height int, ui *editor.UIContext) *PixelArtEditor {
	if ui == nil {
		ui = editor.NewUIContext()
	}

	doc := NewDocument(width, height, Pico8Palette)

	return &PixelArtEditor{
		Doc:           doc,
		ActiveTool:    ToolPencil,
		UI:            ui,
		Zoom:          20, // 20x magnification for 16x16 = 320x320 px canvas
		ShowPixelGrid: true,
		OnionSkin:     true,
		PlaybackFPS:   8.0,
	}
}

// SetStatus sets a temporary notification message.
func (pe *PixelArtEditor) SetStatus(msg string, duration float64) {
	pe.statusMsg = msg
	pe.statusTimer = duration
}

// Update handles mouse interactions, tools, and animation playback.
func (pe *PixelArtEditor) Update(dt float64, screenWidth, screenHeight int) {
	if pe.statusTimer > 0 {
		pe.statusTimer -= dt
		if pe.statusTimer <= 0 {
			pe.statusMsg = ""
		}
	}

	// Update Animation Playback Preview
	if pe.IsPlaying && len(pe.Doc.Frames) > 1 {
		pe.playbackTimer += dt
		frameDuration := 1.0 / pe.PlaybackFPS
		if pe.playbackTimer >= frameDuration {
			pe.playbackTimer -= frameDuration
			pe.previewFrameIndex = (pe.previewFrameIndex + 1) % len(pe.Doc.Frames)
		}
	} else {
		pe.previewFrameIndex = pe.Doc.ActiveFrameIndex
	}

	// Center canvas in available work area
	sidebarW := 210
	toolbarH := 36
	timelineH := 90

	availW := screenWidth - sidebarW
	availH := screenHeight - toolbarH - timelineH

	canvasPixelW := pe.Doc.Width * pe.Zoom
	canvasPixelH := pe.Doc.Height * pe.Zoom

	pe.CanvasOffsetX = (availW - canvasPixelW) / 2
	pe.CanvasOffsetY = toolbarH + (availH-canvasPixelH)/2

	// Hotkeys
	ctrlOrCmd := ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta)
	if ctrlOrCmd && inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			if pe.Doc.Redo() {
				pe.SetStatus("Redo", 1.0)
			}
		} else {
			if pe.Doc.Undo() {
				pe.SetStatus("Undo", 1.0)
			}
		}
	}
	if ctrlOrCmd && inpututil.IsKeyJustPressed(ebiten.KeyY) {
		if pe.Doc.Redo() {
			pe.SetStatus("Redo", 1.0)
		}
	}

	// Quick tool hotkeys
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		pe.ActiveTool = ToolPencil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		pe.ActiveTool = ToolEraser
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		pe.ActiveTool = ToolPicker
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		pe.ActiveTool = ToolBucket
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		pe.ActiveTool = ToolLine
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		pe.ActiveTool = ToolRect
	}

	// Mouse canvas coordinate tracking
	mx, my := ebiten.CursorPosition()
	pe.UI.Update(mx, my)

	relX := mx - pe.CanvasOffsetX
	relY := my - pe.CanvasOffsetY

	inCanvas := relX >= 0 && relX < canvasPixelW && relY >= 0 && relY < canvasPixelH
	pixelX := relX / pe.Zoom
	pixelY := relY / pe.Zoom

	// Tool interactions
	if inCanvas {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			pe.Doc.SaveSnapshot()
			pe.isDrawing = true
			pe.dragStartX = pixelX
			pe.dragStartY = pixelY
			pe.lastDrawX = pixelX
			pe.lastDrawY = pixelY
			pe.hasLastDraw = true

			switch pe.ActiveTool {
			case ToolPencil:
				pe.Doc.SetPixel(pixelX, pixelY, pe.Doc.ActiveColor)
			case ToolEraser:
				pe.Doc.SetPixel(pixelX, pixelY, color.RGBA{})
			case ToolPicker:
				c := pe.Doc.GetPixel(pixelX, pixelY)
				if c.A > 0 {
					pe.Doc.ActiveColor = c
					pe.ActiveTool = ToolPencil
				}
			case ToolBucket:
				FloodFill(pe.Doc, pixelX, pixelY, pe.Doc.ActiveColor)
			}
		} else if pe.isDrawing && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			if pe.ActiveTool == ToolPencil {
				if pe.hasLastDraw {
					DrawStrokeLine(pe.Doc, pe.lastDrawX, pe.lastDrawY, pixelX, pixelY, pe.Doc.ActiveColor)
				}
				pe.lastDrawX = pixelX
				pe.lastDrawY = pixelY
			} else if pe.ActiveTool == ToolEraser {
				if pe.hasLastDraw {
					DrawStrokeLine(pe.Doc, pe.lastDrawX, pe.lastDrawY, pixelX, pixelY, color.RGBA{})
				}
				pe.lastDrawX = pixelX
				pe.lastDrawY = pixelY
			}
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if pe.isDrawing {
			if pe.ActiveTool == ToolLine {
				DrawStrokeLine(pe.Doc, pe.dragStartX, pe.dragStartY, pixelX, pixelY, pe.Doc.ActiveColor)
			} else if pe.ActiveTool == ToolRect {
				DrawRectStroke(pe.Doc, pe.dragStartX, pe.dragStartY, pixelX, pixelY, pe.Doc.ActiveColor, false)
			}
		}
		pe.isDrawing = false
		pe.hasLastDraw = false
	}
}

// Draw renders the pixel art canvas, onion skinning, palette, and animation timeline.
func (pe *PixelArtEditor) Draw(screen *ebiten.Image) {
	sw := screen.Bounds().Dx()
	sh := screen.Bounds().Dy()

	// 1. Draw Checkerboard background for transparent pixels
	canvasPixelW := pe.Doc.Width * pe.Zoom
	canvasPixelH := pe.Doc.Height * pe.Zoom

	cCheck1 := color.RGBA{R: 35, G: 38, B: 48, A: 255}
	cCheck2 := color.RGBA{R: 45, G: 49, B: 62, A: 255}

	for py := 0; py < pe.Doc.Height; py++ {
		for px := 0; px < pe.Doc.Width; px++ {
			bgCol := cCheck1
			if (px+py)%2 == 0 {
				bgCol = cCheck2
			}
			bx := float32(pe.CanvasOffsetX + px*pe.Zoom)
			by := float32(pe.CanvasOffsetY + py*pe.Zoom)
			vector.DrawFilledRect(screen, bx, by, float32(pe.Zoom), float32(pe.Zoom), bgCol, false)
		}
	}

	// 2. Draw Onion Skin (previous frame faintly tinted if enabled)
	if pe.OnionSkin && pe.Doc.ActiveFrameIndex > 0 {
		prevFrame := pe.Doc.Frames[pe.Doc.ActiveFrameIndex-1]
		for py := 0; py < pe.Doc.Height; py++ {
			for px := 0; px < pe.Doc.Width; px++ {
				col := prevFrame.Image.RGBAAt(px, py)
				if col.A > 0 {
					// 30% alpha onion skin
					onionCol := color.RGBA{R: col.R, G: col.G, B: col.B, A: 75}
					bx := float32(pe.CanvasOffsetX + px*pe.Zoom)
					by := float32(pe.CanvasOffsetY + py*pe.Zoom)
					vector.DrawFilledRect(screen, bx, by, float32(pe.Zoom), float32(pe.Zoom), onionCol, false)
				}
			}
		}
	}

	// 3. Draw Active Frame Pixels
	activeFrame := pe.Doc.ActiveFrame()
	if activeFrame != nil {
		for py := 0; py < pe.Doc.Height; py++ {
			for px := 0; px < pe.Doc.Width; px++ {
				col := activeFrame.Image.RGBAAt(px, py)
				if col.A > 0 {
					bx := float32(pe.CanvasOffsetX + px*pe.Zoom)
					by := float32(pe.CanvasOffsetY + py*pe.Zoom)
					vector.DrawFilledRect(screen, bx, by, float32(pe.Zoom), float32(pe.Zoom), col, false)
				}
			}
		}
	}

	// 4. Draw Line / Rect preview while dragging
	if pe.isDrawing {
		mx, my := ebiten.CursorPosition()
		currPX := (mx - pe.CanvasOffsetX) / pe.Zoom
		currPY := (my - pe.CanvasOffsetY) / pe.Zoom

		if pe.ActiveTool == ToolLine {
			BresenhamLine(pe.dragStartX, pe.dragStartY, currPX, currPY, func(x, y int) {
				if x >= 0 && x < pe.Doc.Width && y >= 0 && y < pe.Doc.Height {
					bx := float32(pe.CanvasOffsetX + x*pe.Zoom)
					by := float32(pe.CanvasOffsetY + y*pe.Zoom)
					vector.DrawFilledRect(screen, bx, by, float32(pe.Zoom), float32(pe.Zoom), pe.Doc.ActiveColor, false)
				}
			})
		} else if pe.ActiveTool == ToolRect {
			minX := min(pe.dragStartX, currPX)
			maxX := max(pe.dragStartX, currPX)
			minY := min(pe.dragStartY, currPY)
			maxY := max(pe.dragStartY, currPY)
			for y := minY; y <= maxY; y++ {
				for x := minX; x <= maxX; x++ {
					if x == minX || x == maxX || y == minY || y == maxY {
						if x >= 0 && x < pe.Doc.Width && y >= 0 && y < pe.Doc.Height {
							bx := float32(pe.CanvasOffsetX + x*pe.Zoom)
							by := float32(pe.CanvasOffsetY + y*pe.Zoom)
							vector.DrawFilledRect(screen, bx, by, float32(pe.Zoom), float32(pe.Zoom), pe.Doc.ActiveColor, false)
						}
					}
				}
			}
		}
	}

	// 5. Draw Pixel Grid overlay
	if pe.ShowPixelGrid {
		gridCol := color.RGBA{R: 255, G: 255, B: 255, A: 25}
		editor.DrawGridOverlay(screen, float64(pe.CanvasOffsetX), float64(pe.CanvasOffsetY), float64(canvasPixelW), float64(canvasPixelH), float64(pe.Zoom), gridCol)
	}

	// Canvas outer border
	vector.StrokeRect(screen, float32(pe.CanvasOffsetX), float32(pe.CanvasOffsetY), float32(canvasPixelW), float32(canvasPixelH), 2.0, pe.UI.Theme.BorderActive, false)

	// 6. Draw UI: Right Palette Panel & Bottom Timeline
	pe.drawRightPalettePanel(screen, sw, sh)
	pe.drawBottomTimeline(screen, sw, sh)
}

func (pe *PixelArtEditor) drawRightPalettePanel(screen *ebiten.Image, sw, sh int) {
	sidebarW := 210
	sidebarX := sw - sidebarW
	toolbarH := 36
	panelRect := image.Rect(sidebarX, toolbarH, sw, sh)
	pe.UI.DrawPanel(screen, panelRect, "Colors & Tools")

	y := toolbarH + 24

	// Tool Selector Buttons
	ebitenutil.DebugPrintAt(screen, "Drawing Tools:", sidebarX+8, y)
	y += 16

	tools := []struct {
		name string
		tool Tool
	}{
		{"Pencil (B)", ToolPencil},
		{"Eraser (E)", ToolEraser},
		{"Picker (I)", ToolPicker},
		{"Bucket (G)", ToolBucket},
		{"Line (L)", ToolLine},
		{"Rect (U)", ToolRect},
	}

	btnW := 92
	btnH := 24
	for i, t := range tools {
		bx := sidebarX + 8 + (i%2)*(btnW+8)
		by := y + (i/2)*(btnH+6)
		active := pe.ActiveTool == t.tool
		if pe.UI.Button(screen, image.Rect(bx, by, bx+btnW, by+btnH), t.name, active) {
			pe.ActiveTool = t.tool
		}
	}
	y += (len(tools)/2)*(btnH+6) + 16

	// Active Color Display
	ebitenutil.DebugPrintAt(screen, "Active Color:", sidebarX+8, y)
	y += 16
	vector.DrawFilledRect(screen, float32(sidebarX+8), float32(y), 32, 32, pe.Doc.ActiveColor, false)
	vector.StrokeRect(screen, float32(sidebarX+8), float32(y), 32, 32, 1.0, pe.UI.Theme.BorderActive, false)
	colorHex := fmt.Sprintf("#%02X%02X%02X", pe.Doc.ActiveColor.R, pe.Doc.ActiveColor.G, pe.Doc.ActiveColor.B)
	ebitenutil.DebugPrintAt(screen, colorHex, sidebarX+48, y+8)
	y += 42

	// Palette Swatches
	ebitenutil.DebugPrintAt(screen, pe.Doc.Palette.Name+" Palette:", sidebarX+8, y)
	y += 16

	swatchSize := 22
	swatchesPerRow := 7
	for i, col := range pe.Doc.Palette.Colors {
		colIdx := i % swatchesPerRow
		rowIdx := i / swatchesPerRow
		sx := sidebarX + 8 + colIdx*(swatchSize+4)
		sy := y + rowIdx*(swatchSize+4)

		vector.DrawFilledRect(screen, float32(sx), float32(sy), float32(swatchSize), float32(swatchSize), col, false)

		borderCol := pe.UI.Theme.Border
		if col == pe.Doc.ActiveColor {
			borderCol = pe.UI.Theme.BorderActive
			vector.StrokeRect(screen, float32(sx-1), float32(sy-1), float32(swatchSize+2), float32(swatchSize+2), 2.0, borderCol, false)
		} else {
			vector.StrokeRect(screen, float32(sx), float32(sy), float32(swatchSize), float32(swatchSize), 1.0, borderCol, false)
		}

		mx, my := ebiten.CursorPosition()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if mx >= sx && mx < sx+swatchSize && my >= sy && my < sy+swatchSize {
				pe.Doc.ActiveColor = col
				if pe.ActiveTool == ToolEraser {
					pe.ActiveTool = ToolPencil
				}
			}
		}
	}
	y += ((len(pe.Doc.Palette.Colors)+swatchesPerRow-1)/swatchesPerRow)*(swatchSize+4) + 16

	// Options Checkboxes
	pe.OnionSkin = pe.UI.Checkbox(screen, sidebarX+8, y, "Onion Skin", pe.OnionSkin)
	y += 20
	pe.ShowPixelGrid = pe.UI.Checkbox(screen, sidebarX+8, y, "Pixel Grid", pe.ShowPixelGrid)
	y += 24

	// Live Animation Preview Box
	ebitenutil.DebugPrintAt(screen, "Animation Preview:", sidebarX+8, y)
	y += 16

	previewBoxSize := 48
	vector.DrawFilledRect(screen, float32(sidebarX+8), float32(y), float32(previewBoxSize), float32(previewBoxSize), color.RGBA{R: 24, G: 26, B: 32, A: 255}, false)
	vector.StrokeRect(screen, float32(sidebarX+8), float32(y), float32(previewBoxSize), float32(previewBoxSize), 1.0, pe.UI.Theme.Border, false)

	// Render preview frame
	if len(pe.Doc.Frames) > 0 {
		var pFrame *FrameData
		if pe.IsPlaying {
			pFrame = pe.Doc.Frames[pe.previewFrameIndex]
		} else {
			pFrame = pe.Doc.ActiveFrame()
		}

		if pFrame != nil {
			scale := float32(previewBoxSize) / float32(pe.Doc.Width)
			for py := 0; py < pe.Doc.Height; py++ {
				for px := 0; px < pe.Doc.Width; px++ {
					c := pFrame.Image.RGBAAt(px, py)
					if c.A > 0 {
						pxX := float32(sidebarX+8) + float32(px)*scale
						pxY := float32(y) + float32(py)*scale
						vector.DrawFilledRect(screen, pxX, pxY, scale, scale, c, false)
					}
				}
			}
		}
	}

	// Play / Pause preview button
	playLabel := "Play"
	if pe.IsPlaying {
		playLabel = "Pause"
	}
	if pe.UI.Button(screen, image.Rect(sidebarX+64, y+10, sidebarX+130, y+38), playLabel, pe.IsPlaying) {
		pe.IsPlaying = !pe.IsPlaying
	}
}

func (pe *PixelArtEditor) drawBottomTimeline(screen *ebiten.Image, sw, sh int) {
	sidebarW := 210
	availW := sw - sidebarW
	timelineH := 90
	timelineY := sh - timelineH

	timelineRect := image.Rect(0, timelineY, availW, sh)
	pe.UI.DrawPanel(screen, timelineRect, "Animation Timeline (Frames)")

	x := 10
	btnY := timelineY + 28

	// Timeline Controls
	if pe.UI.Button(screen, image.Rect(x, btnY, x+48, btnY+24), "+ Frame", false) {
		pe.Doc.AddFrame()
	}
	x += 54
	if pe.UI.Button(screen, image.Rect(x, btnY, x+40, btnY+24), "Copy", false) {
		pe.Doc.DuplicateFrame()
	}
	x += 46
	if pe.UI.Button(screen, image.Rect(x, btnY, x+40, btnY+24), "Del", false) {
		pe.Doc.DeleteFrame()
	}
	x += 54

	// Frame Thumbnails
	thumbSize := 36
	for i, frame := range pe.Doc.Frames {
		tx := x
		ty := btnY - 6
		tRect := image.Rect(tx, ty, tx+thumbSize, ty+thumbSize)

		isActive := i == pe.Doc.ActiveFrameIndex
		borderCol := pe.UI.Theme.Border
		if isActive {
			borderCol = pe.UI.Theme.BorderActive
		}

		vector.DrawFilledRect(screen, float32(tx), float32(ty), float32(thumbSize), float32(thumbSize), color.RGBA{R: 20, G: 22, B: 28, A: 255}, false)

		// Draw thumbnail pixels
		scale := float32(thumbSize) / float32(pe.Doc.Width)
		for py := 0; py < pe.Doc.Height; py++ {
			for px := 0; px < pe.Doc.Width; px++ {
				c := frame.Image.RGBAAt(px, py)
				if c.A > 0 {
					vector.DrawFilledRect(screen, float32(tx)+float32(px)*scale, float32(ty)+float32(py)*scale, scale, scale, c, false)
				}
			}
		}

		vector.StrokeRect(screen, float32(tx), float32(ty), float32(thumbSize), float32(thumbSize), 1.0, borderCol, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i+1), tx+2, ty+2)

		// Click thumbnail to switch active frame
		mx, my := ebiten.CursorPosition()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if mx >= tRect.Min.X && mx < tRect.Max.X && my >= tRect.Min.Y && my < tRect.Max.Y {
				pe.Doc.ActiveFrameIndex = i
			}
		}

		x += thumbSize + 8
	}

	// Status notice if any
	if pe.statusMsg != "" {
		ebitenutil.DebugPrintAt(screen, pe.statusMsg, 10, sh-20)
	}
}
