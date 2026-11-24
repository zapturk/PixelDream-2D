package gfx

import (
	"image"
	"testing"
)

func TestAnimationLoopModes(t *testing.T) {
	// 3 frames, 0.1s each
	anim := NewAnimation("walk", LoopForever)
	anim.AddFrame(image.Rect(0, 0, 16, 16), 0.1)
	anim.AddFrame(image.Rect(16, 0, 32, 16), 0.1)
	anim.AddFrame(image.Rect(32, 0, 48, 16), 0.1)

	player := NewAnimationPlayer()
	player.AddAnimation(anim)
	player.Play("walk", false)

	if player.CurrentFrameIndex() != 0 {
		t.Fatalf("expected frame 0, got %d", player.CurrentFrameIndex())
	}

	player.Update(0.1)
	if player.CurrentFrameIndex() != 1 {
		t.Fatalf("expected frame 1, got %d", player.CurrentFrameIndex())
	}

	player.Update(0.1)
	if player.CurrentFrameIndex() != 2 {
		t.Fatalf("expected frame 2, got %d", player.CurrentFrameIndex())
	}

	// Should loop back to 0
	player.Update(0.1)
	if player.CurrentFrameIndex() != 0 {
		t.Fatalf("expected loop back to frame 0, got %d", player.CurrentFrameIndex())
	}
}

func TestAnimationLoopOnce(t *testing.T) {
	finished := false
	anim := NewAnimation("attack", LoopOnce)
	anim.AddFrame(image.Rect(0, 0, 16, 16), 0.1)
	anim.AddFrame(image.Rect(16, 0, 32, 16), 0.1)
	anim.OnFinish = func() {
		finished = true
	}

	player := NewAnimationPlayer()
	player.AddAnimation(anim)
	player.Play("attack", false)

	player.Update(0.1)
	if player.CurrentFrameIndex() != 1 {
		t.Fatalf("expected frame 1, got %d", player.CurrentFrameIndex())
	}

	player.Update(0.1)
	if player.IsPlaying() {
		t.Fatalf("expected animation to stop")
	}
	if !finished {
		t.Fatalf("expected OnFinish callback to be triggered")
	}
}

func TestAnimationLoopPingPong(t *testing.T) {
	anim := NewAnimation("hover", LoopPingPong)
	anim.AddFrame(image.Rect(0, 0, 16, 16), 0.1)
	anim.AddFrame(image.Rect(16, 0, 32, 16), 0.1)
	anim.AddFrame(image.Rect(32, 0, 48, 16), 0.1)

	player := NewAnimationPlayer()
	player.AddAnimation(anim)
	player.Play("hover", false)

	// Forward: 0 -> 1 -> 2
	player.Update(0.1) // -> 1
	player.Update(0.1) // -> 2
	if player.CurrentFrameIndex() != 2 {
		t.Fatalf("expected frame 2, got %d", player.CurrentFrameIndex())
	}

	// Ping-pong reverse: 2 -> 1 -> 0
	player.Update(0.1) // -> 1
	if player.CurrentFrameIndex() != 1 {
		t.Fatalf("expected frame 1 on bounce, got %d", player.CurrentFrameIndex())
	}
	player.Update(0.1) // -> 0
	if player.CurrentFrameIndex() != 0 {
		t.Fatalf("expected frame 0 on bounce, got %d", player.CurrentFrameIndex())
	}
}

func TestAnimationFrameEvent(t *testing.T) {
	hitTriggered := false
	anim := NewAnimation("punch", LoopOnce)
	anim.AddFrame(image.Rect(0, 0, 16, 16), 0.1)
	anim.AddFrame(image.Rect(16, 0, 32, 16), 0.1)
	anim.SetFrameEvent(1, func() {
		hitTriggered = true
	})

	player := NewAnimationPlayer()
	player.AddAnimation(anim)
	player.Play("punch", false)

	if hitTriggered {
		t.Fatalf("event should not trigger on frame 0")
	}

	player.Update(0.1)
	if !hitTriggered {
		t.Fatalf("event should trigger on frame 1")
	}
}
