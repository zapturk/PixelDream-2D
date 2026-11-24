package input

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestActionMapBinding(t *testing.T) {
	am := NewActionMap()
	am.BindKey("jump", ebiten.KeySpace, ebiten.KeyW)
	am.BindMouseButton("shoot", ebiten.MouseButtonLeft)

	b, ok := am.bindings["jump"]
	if !ok || len(b.Keys) != 2 {
		t.Fatalf("expected 2 keys for 'jump', got %v", b)
	}

	bShoot, ok := am.bindings["shoot"]
	if !ok || len(bShoot.MouseButtons) != 1 {
		t.Fatalf("expected 1 mouse button for 'shoot', got %v", bShoot)
	}
}
