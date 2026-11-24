package pixeleditor

import (
	"image/color"
	"testing"
)

func TestDocumentUndoRedoAndFrames(t *testing.T) {
	doc := NewDocument(16, 16, Pico8Palette)

	cRed := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	cBlue := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	// 1. Initial State
	if len(doc.Frames) != 1 {
		t.Fatalf("expected 1 initial frame, got %d", len(doc.Frames))
	}
	if doc.GetPixel(0, 0) != (color.RGBA{}) {
		t.Fatalf("expected blank pixel at (0,0)")
	}

	// 2. Modify with Snapshot
	doc.SaveSnapshot()
	doc.SetPixel(0, 0, cRed)
	if doc.GetPixel(0, 0) != cRed {
		t.Fatalf("expected pixel at (0,0) to be red")
	}

	// 3. Undo
	if !doc.Undo() {
		t.Fatalf("expected Undo to succeed")
	}
	if doc.GetPixel(0, 0) != (color.RGBA{}) {
		t.Fatalf("expected pixel at (0,0) to revert to blank after undo")
	}

	// 4. Redo
	if !doc.Redo() {
		t.Fatalf("expected Redo to succeed")
	}
	if doc.GetPixel(0, 0) != cRed {
		t.Fatalf("expected pixel at (0,0) to be red after redo")
	}

	// 5. Frame Management
	doc.AddFrame()
	if len(doc.Frames) != 2 || doc.ActiveFrameIndex != 1 {
		t.Fatalf("expected 2 frames with active index 1, got %d frames and index %d", len(doc.Frames), doc.ActiveFrameIndex)
	}

	doc.SetPixel(5, 5, cBlue)
	if doc.GetPixel(5, 5) != cBlue {
		t.Fatalf("expected blue pixel on frame 2")
	}

	// Duplicate
	doc.DuplicateFrame()
	if len(doc.Frames) != 3 || doc.ActiveFrameIndex != 2 {
		t.Fatalf("expected 3 frames after duplicate, got %d", len(doc.Frames))
	}
	if doc.GetPixel(5, 5) != cBlue {
		t.Fatalf("expected duplicated frame to retain blue pixel")
	}

	// Delete
	if !doc.DeleteFrame() {
		t.Fatalf("expected DeleteFrame to succeed")
	}
	if len(doc.Frames) != 2 {
		t.Fatalf("expected 2 frames after delete, got %d", len(doc.Frames))
	}
}

func TestBresenhamAndFloodFill(t *testing.T) {
	doc := NewDocument(16, 16, Pico8Palette)
	cGreen := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	// Draw diagonal line from (2, 2) to (5, 5)
	DrawStrokeLine(doc, 2, 2, 5, 5, cGreen)
	for i := 2; i <= 5; i++ {
		if doc.GetPixel(i, i) != cGreen {
			t.Fatalf("expected pixel at (%d,%d) on line to be green", i, i)
		}
	}

	// Flood Fill
	cYellow := color.RGBA{R: 255, G: 255, B: 0, A: 255}
	// Create a closed 3x3 box
	DrawRectStroke(doc, 7, 7, 10, 10, cGreen, false)
	// Fill inside at (8, 8)
	FloodFill(doc, 8, 8, cYellow)

	if doc.GetPixel(8, 8) != cYellow || doc.GetPixel(9, 9) != cYellow {
		t.Fatalf("expected inner pixels to be filled with yellow")
	}
	// Border should still be green
	if doc.GetPixel(7, 7) != cGreen {
		t.Fatalf("expected border to remain green")
	}
	// Outside at (0, 0) should remain empty
	if doc.GetPixel(0, 0) != (color.RGBA{}) {
		t.Fatalf("expected outside pixel (0,0) to remain blank")
	}
}
