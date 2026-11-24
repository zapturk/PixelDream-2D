package gfx

import (
	"image"
	"time"

	"github.com/zapturk/PixelDream-2D/pkg/display"
)

// LoopMode determines how an animation behaves when it reaches the end.
type LoopMode int

const (
	// LoopForever restarts from the first frame indefinitely.
	LoopForever LoopMode = iota
	// LoopOnce plays through once and stops on the last frame.
	LoopOnce
	// LoopPingPong alternates playing forward and backward.
	LoopPingPong
)

// Frame represents a single frame in an animation.
type Frame struct {
	Rect     image.Rectangle
	Duration float64 // in seconds
}

// Animation holds an ordered collection of frames.
type Animation struct {
	Name     string
	Frames   []Frame
	Loop     LoopMode
	OnFinish func()
	// Callbacks triggered on reaching a specific frame index
	FrameEvents map[int]func()
}

// NewAnimation creates a new animation with the specified name and loop mode.
func NewAnimation(name string, loop LoopMode) *Animation {
	return &Animation{
		Name:        name,
		Loop:        loop,
		FrameEvents: make(map[int]func()),
	}
}

// AddFrame adds a frame with the given rectangle and duration (in seconds).
func (a *Animation) AddFrame(rect image.Rectangle, duration float64) *Animation {
	a.Frames = append(a.Frames, Frame{
		Rect:     rect,
		Duration: duration,
	})
	return a
}

// AddUniformFrames slices a row or column of uniform width and height frames from a texture sheet.
func (a *Animation) AddUniformFrames(startX, startY, frameWidth, frameHeight, count int, duration float64) *Animation {
	for i := 0; i < count; i++ {
		x := startX + i*frameWidth
		rect := image.Rect(x, startY, x+frameWidth, startY+frameHeight)
		a.AddFrame(rect, duration)
	}
	return a
}

// SetFrameEvent registers a callback when the animation reaches a specific frame index.
func (a *Animation) SetFrameEvent(frameIndex int, callback func()) *Animation {
	a.FrameEvents[frameIndex] = callback
	return a
}

// AnimationPlayer handles playback of animations.
type AnimationPlayer struct {
	animations   map[string]*Animation
	current      *Animation
	currentFrame int
	elapsedTime  float64
	isPlaying    bool
	isReverse    bool // used for PingPong
}

// NewAnimationPlayer creates a new AnimationPlayer.
func NewAnimationPlayer() *AnimationPlayer {
	return &AnimationPlayer{
		animations: make(map[string]*Animation),
		isPlaying:  false,
	}
}

// AddAnimation registers an animation with the player.
func (ap *AnimationPlayer) AddAnimation(anim *Animation) {
	if anim != nil {
		ap.animations[anim.Name] = anim
	}
}

// Play starts playing the animation with the given name. If already playing, it continues unless restart is true.
func (ap *AnimationPlayer) Play(name string, restart bool) {
	anim, ok := ap.animations[name]
	if !ok || len(anim.Frames) == 0 {
		return
	}

	if ap.current == anim && ap.isPlaying && !restart {
		return
	}

	ap.current = anim
	ap.currentFrame = 0
	ap.elapsedTime = 0
	ap.isPlaying = true
	ap.isReverse = false

	// Trigger event on initial frame if registered
	if ev, ok := anim.FrameEvents[0]; ok && ev != nil {
		ev()
	}
}

// Pause pauses animation playback.
func (ap *AnimationPlayer) Pause() {
	ap.isPlaying = false
}

// Resume unpauses animation playback.
func (ap *AnimationPlayer) Resume() {
	if ap.current != nil && len(ap.current.Frames) > 0 {
		ap.isPlaying = true
	}
}

// Stop stops the animation and resets to the initial frame.
func (ap *AnimationPlayer) Stop() {
	ap.isPlaying = false
	ap.currentFrame = 0
	ap.elapsedTime = 0
	ap.isReverse = false
}

// CurrentAnimationName returns the name of the active animation.
func (ap *AnimationPlayer) CurrentAnimationName() string {
	if ap.current == nil {
		return ""
	}
	return ap.current.Name
}

// CurrentFrameIndex returns the index of the active frame.
func (ap *AnimationPlayer) CurrentFrameIndex() int {
	return ap.currentFrame
}

// IsPlaying returns true if an animation is actively progressing.
func (ap *AnimationPlayer) IsPlaying() bool {
	return ap.isPlaying
}

// CurrentSourceRect returns the source rectangle of the current frame.
func (ap *AnimationPlayer) CurrentSourceRect() image.Rectangle {
	if ap.current == nil || len(ap.current.Frames) == 0 {
		return image.Rectangle{}
	}
	return ap.current.Frames[ap.currentFrame].Rect
}

// Update advances the animation by delta time (dt in seconds).
func (ap *AnimationPlayer) Update(dt float64) {
	if !ap.isPlaying || ap.current == nil || len(ap.current.Frames) <= 1 {
		return
	}

	ap.elapsedTime += dt
	frameDuration := ap.current.Frames[ap.currentFrame].Duration
	if frameDuration <= 0 {
		frameDuration = 1.0 / 60.0 // fallback default
	}

	for ap.elapsedTime >= frameDuration {
		ap.elapsedTime -= frameDuration
		ap.advanceFrame()
		if !ap.isPlaying {
			break
		}
		frameDuration = ap.current.Frames[ap.currentFrame].Duration
	}
}

func (ap *AnimationPlayer) advanceFrame() {
	anim := ap.current
	numFrames := len(anim.Frames)

	switch anim.Loop {
	case LoopForever:
		ap.currentFrame = (ap.currentFrame + 1) % numFrames
	case LoopOnce:
		if ap.currentFrame < numFrames-1 {
			ap.currentFrame++
		} else {
			ap.isPlaying = false
			if anim.OnFinish != nil {
				anim.OnFinish()
			}
			return
		}
	case LoopPingPong:
		if !ap.isReverse {
			if ap.currentFrame < numFrames-1 {
				ap.currentFrame++
			} else {
				ap.isReverse = true
				ap.currentFrame--
			}
		} else {
			if ap.currentFrame > 0 {
				ap.currentFrame--
			} else {
				ap.isReverse = false
				ap.currentFrame++
			}
		}
	}

	// Trigger frame event
	if ev, ok := anim.FrameEvents[ap.currentFrame]; ok && ev != nil {
		ev()
	}
}

// AnimatedSprite combines a Sprite with an AnimationPlayer.
type AnimatedSprite struct {
	*Sprite
	Player *AnimationPlayer
}

// NewAnimatedSprite creates a sprite backed by an animation player.
func NewAnimatedSprite(texture *Texture) *AnimatedSprite {
	s := NewSprite(texture)
	return &AnimatedSprite{
		Sprite: s,
		Player: NewAnimationPlayer(),
	}
}

// Update updates the animation and syncs the sprite's SourceRect.
func (as *AnimatedSprite) Update(dt float64) {
	as.Player.Update(dt)
	rect := as.Player.CurrentSourceRect()
	if !rect.Empty() {
		as.SourceRect = rect
	}
}

// Draw draws the animated sprite.
func (as *AnimatedSprite) Draw(canvas *display.Canvas, cam *display.Camera) {
	as.Sprite.Draw(canvas, cam)
}

// Convenience helper to construct duration from milliseconds
func Ms(millis int) float64 {
	return float64(millis) / 1000.0
}

// Convenience helper to construct duration from time.Duration
func Duration(d time.Duration) float64 {
	return d.Seconds()
}
