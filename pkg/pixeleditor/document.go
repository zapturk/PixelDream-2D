package pixeleditor

import (
	"image"
	"image/color"
)

// Palette represents an indexed color palette.
type Palette struct {
	Name   string
	Colors []color.RGBA
}

// Pico8Palette is the classic 16-color PICO-8 retro palette.
var Pico8Palette = Palette{
	Name: "PICO-8",
	Colors: []color.RGBA{
		{0, 0, 0, 255},       // Black
		{29, 43, 83, 255},    // Dark Blue
		{126, 37, 83, 255},   // Dark Purple
		{0, 135, 81, 255},    // Dark Green
		{171, 82, 54, 255},   // Brown
		{95, 87, 79, 255},    // Dark Gray
		{194, 195, 199, 255}, // Light Gray
		{255, 241, 232, 255}, // White
		{255, 0, 77, 255},    // Red
		{255, 163, 0, 255},   // Orange
		{255, 236, 39, 255},  // Yellow
		{0, 228, 54, 255},    // Green
		{41, 173, 255, 255},  // Blue
		{131, 118, 156, 255}, // Indigo
		{255, 119, 168, 255}, // Pink
		{255, 204, 170, 255}, // Peach
	},
}

// DawnBringer16Palette is a popular 16-color pixel art palette.
var DawnBringer16Palette = Palette{
	Name: "DawnBringer 16",
	Colors: []color.RGBA{
		{20, 12, 28, 255},
		{68, 36, 52, 255},
		{48, 52, 109, 255},
		{78, 74, 78, 255},
		{133, 76, 48, 255},
		{52, 101, 36, 255},
		{208, 70, 72, 255},
		{117, 113, 97, 255},
		{89, 125, 206, 255},
		{210, 125, 44, 255},
		{133, 149, 161, 255},
		{109, 170, 44, 255},
		{218, 212, 94, 255},
		{222, 158, 65, 255},
		{255, 255, 255, 255},
	},
}

// FrameData represents a single pixel image in a multi-frame animation.
type FrameData struct {
	Image *image.RGBA
}

// Clone creates a deep copy of the frame data.
func (f *FrameData) Clone() *FrameData {
	bounds := f.Image.Bounds()
	newImg := image.NewRGBA(bounds)
	copy(newImg.Pix, f.Image.Pix)
	return &FrameData{Image: newImg}
}

// Document represents an editable multi-frame pixel art project.
type Document struct {
	Width            int
	Height           int
	Frames           []*FrameData
	ActiveFrameIndex int
	Palette          Palette
	ActiveColor      color.RGBA

	// Undo / Redo history stacks (stores deep copies of the active frame image)
	undoStack [][]*FrameData
	redoStack [][]*FrameData
	maxUndo   int
}

// NewDocument creates a new pixel art document with initial frame and palette.
func NewDocument(width, height int, palette Palette) *Document {
	if width <= 0 {
		width = 16
	}
	if height <= 0 {
		height = 16
	}
	if len(palette.Colors) == 0 {
		palette = Pico8Palette
	}

	initialFrame := &FrameData{
		Image: image.NewRGBA(image.Rect(0, 0, width, height)),
	}

	return &Document{
		Width:            width,
		Height:           height,
		Frames:           []*FrameData{initialFrame},
		ActiveFrameIndex: 0,
		Palette:          palette,
		ActiveColor:      palette.Colors[8], // default to Red or another vibrant color
		maxUndo:          30,
	}
}

// ActiveFrame returns the current active frame.
func (doc *Document) ActiveFrame() *FrameData {
	if len(doc.Frames) == 0 {
		return nil
	}
	if doc.ActiveFrameIndex < 0 || doc.ActiveFrameIndex >= len(doc.Frames) {
		doc.ActiveFrameIndex = 0
	}
	return doc.Frames[doc.ActiveFrameIndex]
}

// SaveSnapshot saves the current frames state to the undo stack before a mutation.
func (doc *Document) SaveSnapshot() {
	snapshot := make([]*FrameData, len(doc.Frames))
	for i, f := range doc.Frames {
		snapshot[i] = f.Clone()
	}

	doc.undoStack = append(doc.undoStack, snapshot)
	if len(doc.undoStack) > doc.maxUndo {
		doc.undoStack = doc.undoStack[1:]
	}
	// Clear redo on new action
	doc.redoStack = nil
}

// Undo reverts the frames to the previous snapshot.
func (doc *Document) Undo() bool {
	if len(doc.undoStack) == 0 {
		return false
	}

	// Save current state to redo
	currentSnapshot := make([]*FrameData, len(doc.Frames))
	for i, f := range doc.Frames {
		currentSnapshot[i] = f.Clone()
	}
	doc.redoStack = append(doc.redoStack, currentSnapshot)

	// Pop undo
	prevSnapshot := doc.undoStack[len(doc.undoStack)-1]
	doc.undoStack = doc.undoStack[:len(doc.undoStack)-1]

	doc.Frames = make([]*FrameData, len(prevSnapshot))
	for i, f := range prevSnapshot {
		doc.Frames[i] = f.Clone()
	}
	if doc.ActiveFrameIndex >= len(doc.Frames) {
		doc.ActiveFrameIndex = len(doc.Frames) - 1
	}

	return true
}

// Redo reapplies a previously undone action.
func (doc *Document) Redo() bool {
	if len(doc.redoStack) == 0 {
		return false
	}

	// Save current state to undo
	currentSnapshot := make([]*FrameData, len(doc.Frames))
	for i, f := range doc.Frames {
		currentSnapshot[i] = f.Clone()
	}
	doc.undoStack = append(doc.undoStack, currentSnapshot)

	// Pop redo
	nextSnapshot := doc.redoStack[len(doc.redoStack)-1]
	doc.redoStack = doc.redoStack[:len(doc.redoStack)-1]

	doc.Frames = make([]*FrameData, len(nextSnapshot))
	for i, f := range nextSnapshot {
		doc.Frames[i] = f.Clone()
	}
	if doc.ActiveFrameIndex >= len(doc.Frames) {
		doc.ActiveFrameIndex = len(doc.Frames) - 1
	}

	return true
}

// AddFrame adds a new blank frame after the current frame.
func (doc *Document) AddFrame() {
	doc.SaveSnapshot()
	newFrame := &FrameData{
		Image: image.NewRGBA(image.Rect(0, 0, doc.Width, doc.Height)),
	}
	doc.Frames = append(doc.Frames, newFrame)
	doc.ActiveFrameIndex = len(doc.Frames) - 1
}

// DuplicateFrame clones the active frame and appends it.
func (doc *Document) DuplicateFrame() {
	doc.SaveSnapshot()
	active := doc.ActiveFrame()
	if active != nil {
		doc.Frames = append(doc.Frames, active.Clone())
		doc.ActiveFrameIndex = len(doc.Frames) - 1
	}
}

// DeleteFrame removes the current frame (unless it's the only one).
func (doc *Document) DeleteFrame() bool {
	if len(doc.Frames) <= 1 {
		return false
	}
	doc.SaveSnapshot()
	idx := doc.ActiveFrameIndex
	doc.Frames = append(doc.Frames[:idx], doc.Frames[idx+1:]...)
	if doc.ActiveFrameIndex >= len(doc.Frames) {
		doc.ActiveFrameIndex = len(doc.Frames) - 1
	}
	return true
}

// GetPixel returns the color at (x, y) on the active frame.
func (doc *Document) GetPixel(x, y int) color.RGBA {
	frame := doc.ActiveFrame()
	if frame == nil || x < 0 || x >= doc.Width || y < 0 || y >= doc.Height {
		return color.RGBA{}
	}
	return frame.Image.RGBAAt(x, y)
}

// SetPixel sets the color at (x, y) on the active frame.
func (doc *Document) SetPixel(x, y int, col color.RGBA) {
	frame := doc.ActiveFrame()
	if frame == nil || x < 0 || x >= doc.Width || y < 0 || y >= doc.Height {
		return
	}
	frame.Image.SetRGBA(x, y, col)
}
