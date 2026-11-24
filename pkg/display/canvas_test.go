package display

import (
	"testing"
)

func TestCanvasIntegerScaling(t *testing.T) {
	canvas := NewCanvas(320, 180)
	canvas.SetScaleMode(ScaleModeInteger)

	// Window: 1280x720 (exactly 4x)
	canvas.UpdateLayout(1280, 720)
	sx, sy, ox, oy := canvas.LayoutMetrics()

	if sx != 4 || sy != 4 {
		t.Fatalf("expected scale 4x, got sx=%f sy=%f", sx, sy)
	}
	if ox != 0 || oy != 0 {
		t.Fatalf("expected 0 offset, got ox=%f oy=%f", ox, oy)
	}

	// Window with excess width (1400x720) -> should still be 4x with pillarbox
	canvas.UpdateLayout(1400, 720)
	sx, sy, ox, oy = canvas.LayoutMetrics()
	if sx != 4 || sy != 4 {
		t.Fatalf("expected scale 4x, got sx=%f sy=%f", sx, sy)
	}
	expectedOx := (1400.0 - 320.0*4.0) / 2.0 // (1400 - 1280) / 2 = 60
	if ox != expectedOx || oy != 0 {
		t.Fatalf("expected offset (%f, 0), got (%f, %f)", expectedOx, ox, oy)
	}
}

func TestCoordinateConversion(t *testing.T) {
	canvas := NewCanvas(320, 180)
	canvas.SetScaleMode(ScaleModeInteger)
	canvas.UpdateLayout(1400, 720) // scale=4, ox=60, oy=0

	// Screen position right at canvas (0,0) -> screen (60, 0)
	cx, cy, inBounds := canvas.ScreenToCanvas(60, 0)
	if !inBounds || cx != 0 || cy != 0 {
		t.Fatalf("expected canvas (0,0), got (%d,%d), inBounds=%v", cx, cy, inBounds)
	}

	// Screen position outside canvas (pillarbox area)
	_, _, inBounds = canvas.ScreenToCanvas(10, 50)
	if inBounds {
		t.Fatalf("expected out of bounds for pillarbox area")
	}

	// Canvas center (160, 90) -> screen coords
	sx, sy := canvas.CanvasToScreen(160, 90)
	expectedSx := 60 + 160*4
	expectedSy := 0 + 90*4
	if sx != expectedSx || sy != expectedSy {
		t.Fatalf("expected screen (%d,%d), got (%d,%d)", expectedSx, expectedSy, sx, sy)
	}
}
