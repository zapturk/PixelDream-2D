package engine

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zapturk/PixelDream-2D/pkg/display"
	"github.com/zapturk/PixelDream-2D/pkg/input"
)

// Game is the interface that client games implement.
type Game interface {
	// Update is called every tick. dt is the delta time in seconds.
	Update(dt float64) error

	// Draw renders onto the virtual pixel canvas.
	Draw(canvas *display.Canvas)
}

// Config specifies window and resolution settings.
type Config struct {
	Title         string
	VirtualWidth  int
	VirtualHeight int
	WindowWidth   int
	WindowHeight  int
	ScaleMode     display.ScaleMode
	TargetFPS     int
	Resizable     bool
}

// DefaultConfig provides sensible defaults for a pixel-art game.
func DefaultConfig() Config {
	return Config{
		Title:         "PixelDream-2D Game",
		VirtualWidth:  320,
		VirtualHeight: 180,
		WindowWidth:   1280,
		WindowHeight:  720,
		ScaleMode:     display.ScaleModeInteger,
		TargetFPS:     60,
		Resizable:     true,
	}
}

// Runner wraps Ebitengine's Game loop and manages the virtual Canvas.
type Runner struct {
	game     Game
	canvas   *display.Canvas
	inputMap *input.ActionMap
	lastTime time.Time
}

// NewRunner creates an engine runner for the given game and configuration.
func NewRunner(game Game, config Config) *Runner {
	canvas := display.NewCanvas(config.VirtualWidth, config.VirtualHeight)
	canvas.SetScaleMode(config.ScaleMode)

	return &Runner{
		game:     game,
		canvas:   canvas,
		inputMap: input.NewActionMap(),
		lastTime: time.Now(),
	}
}

// Canvas returns the engine's virtual canvas.
func (r *Runner) Canvas() *display.Canvas {
	return r.canvas
}

// InputMap returns the engine's action map.
func (r *Runner) InputMap() *input.ActionMap {
	return r.inputMap
}

// Update implements ebiten.Game.
func (r *Runner) Update() error {
	now := time.Now()
	dt := now.Sub(r.lastTime).Seconds()
	r.lastTime = now

	// Cap delta time to prevent physics explosions on hiccups or pause
	if dt > 0.1 {
		dt = 0.1
	}

	return r.game.Update(dt)
}

// Draw implements ebiten.Game.
func (r *Runner) Draw(screen *ebiten.Image) {
	r.canvas.Clear()
	r.game.Draw(r.canvas)
	r.canvas.RenderToScreen(screen)
}

// Layout implements ebiten.Game.
func (r *Runner) Layout(outsideWidth, outsideHeight int) (int, int) {
	r.canvas.UpdateLayout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

// Run starts the engine game loop.
func Run(game Game, config Config) error {
	ebiten.SetWindowTitle(config.Title)
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if config.TargetFPS > 0 {
		ebiten.SetTPS(config.TargetFPS)
	}

	runner := NewRunner(game, config)
	return ebiten.RunGame(runner)
}
