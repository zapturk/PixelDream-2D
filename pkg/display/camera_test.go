package display

import (
	"math"
	"testing"
)

func TestCameraCenterWorldToScreen(t *testing.T) {
	cam := NewCamera(320, 180)
	cam.X = 100
	cam.Y = 200

	// The camera center in world space (100, 200) should project to the viewport center (160, 90)
	sx, sy := cam.WorldToScreen(100, 200)
	if sx != 160 || sy != 90 {
		t.Fatalf("expected (160, 90), got (%f, %f)", sx, sy)
	}

	// Inverse should project back to (100, 200)
	wx, wy := cam.ScreenToWorld(160, 90)
	if math.Abs(wx-100) > 0.001 || math.Abs(wy-200) > 0.001 {
		t.Fatalf("expected world (100, 200), got (%f, %f)", wx, wy)
	}
}

func TestCameraPixelSnapping(t *testing.T) {
	cam := NewCamera(320, 180)
	cam.PixelSnap = true
	cam.X = 100.4
	cam.Y = 200.7

	// With pixel snap: camX rounds to 100, camY rounds to 201
	// World (100, 201) -> screen center (160, 90)
	sx, sy := cam.WorldToScreen(100, 201)
	if sx != 160 || sy != 90 {
		t.Fatalf("expected snapped center at (160, 90), got (%f, %f)", sx, sy)
	}
}

func TestCameraBounds(t *testing.T) {
	cam := NewCamera(320, 180)
	cam.SetBounds(0, 0, 1000, 1000)

	// Half viewport: 160, 90
	cam.Follow(-50, -50, 1.0)
	if cam.X != 160 || cam.Y != 90 {
		t.Fatalf("expected clamped at (160, 90), got (%f, %f)", cam.X, cam.Y)
	}

	cam.Follow(2000, 2000, 1.0)
	if cam.X != 840 || cam.Y != 910 {
		t.Fatalf("expected clamped at (840, 910), got (%f, %f)", cam.X, cam.Y)
	}
}
