package physics

import (
	"testing"
)

func TestKinematicGravityAndGroundCollision(t *testing.T) {
	body := NewKinematicBody(0, 0, 16, 16)
	body.Gravity = 500.0

	// Tile grid: ground at y=2 (world Y = 32)
	isSolid := func(tx, ty int) bool {
		return ty >= 2
	}

	// Step physics forward 0.1s
	body.ApplyGravity(0.1)
	if body.VelocityY <= 0 {
		t.Fatalf("expected positive falling velocity, got %f", body.VelocityY)
	}

	// Move body towards ground
	body.MoveAndSlide(0.1, isSolid, 16)

	// Keep stepping until body lands
	for i := 0; i < 20; i++ {
		body.ApplyGravity(0.05)
		body.MoveAndSlide(0.05, isSolid, 16)
		if body.OnGround {
			break
		}
	}

	if !body.OnGround {
		t.Fatalf("expected body to be on ground")
	}

	// Box bottom should be exactly at 32 (tile Top)
	if body.Box.Bottom() != 32.0 {
		t.Fatalf("expected box bottom at 32.0, got %f", body.Box.Bottom())
	}
	if body.VelocityY != 0 {
		t.Fatalf("expected velocity Y to be 0 upon landing, got %f", body.VelocityY)
	}
}

func TestKinematicWallCollision(t *testing.T) {
	body := NewKinematicBody(0, 0, 16, 16)
	body.VelocityX = 200.0

	// Wall at tile X=2 (world X = 32)
	isSolid := func(tx, ty int) bool {
		return tx >= 2
	}

	body.MoveAndSlide(0.2, isSolid, 16)

	if !body.OnWall {
		t.Fatalf("expected body to hit wall")
	}
	if body.Box.Right() != 32.0 {
		t.Fatalf("expected body right edge at 32.0, got %f", body.Box.Right())
	}
	if body.VelocityX != 0 {
		t.Fatalf("expected VelocityX to be 0 upon hitting wall, got %f", body.VelocityX)
	}
}
