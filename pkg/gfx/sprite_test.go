package gfx

import (
	"image"
	"image/color"
	"testing"
)

func TestSpriteDimensionsAndOrigins(t *testing.T) {
	tex := NewSolidColorTexture(32, 64, color.White)
	sprite := NewSprite(tex)

	if sprite.FrameWidth() != 32 || sprite.FrameHeight() != 64 {
		t.Fatalf("expected (32, 64), got (%d, %d)", sprite.FrameWidth(), sprite.FrameHeight())
	}

	sprite.SetOriginPreset(OriginCenter)
	ox, oy := sprite.Origin()
	if ox != 16 || oy != 32 {
		t.Fatalf("expected center origin (16, 32), got (%f, %f)", ox, oy)
	}

	sprite.SetOriginPreset(OriginBottomCenter)
	ox, oy = sprite.Origin()
	if ox != 16 || oy != 64 {
		t.Fatalf("expected bottom-center origin (16, 64), got (%f, %f)", ox, oy)
	}

	// When SourceRect changes to a smaller 16x16 animation frame, Origin() must dynamically adapt
	sprite.SourceRect = image.Rect(0, 0, 16, 16)
	if sprite.FrameWidth() != 16 || sprite.FrameHeight() != 16 {
		t.Fatalf("expected sub-rect (16, 16), got (%d, %d)", sprite.FrameWidth(), sprite.FrameHeight())
	}
	ox, oy = sprite.Origin()
	if ox != 8 || oy != 16 {
		t.Fatalf("expected dynamic bottom-center origin (8, 16), got (%f, %f)", ox, oy)
	}
}
