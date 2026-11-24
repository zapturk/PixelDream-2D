package gfx

import (
	"image"
	"image/color"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// Texture wraps an ebiten.Image with convenient metadata and region creation.
type Texture struct {
	image  *ebiten.Image
	width  int
	height int
}

// NewTextureFromImage wraps an existing *ebiten.Image as a Texture.
func NewTextureFromImage(img *ebiten.Image) *Texture {
	bounds := img.Bounds()
	return &Texture{
		image:  img,
		width:  bounds.Dx(),
		height: bounds.Dy(),
	}
}

// NewSolidColorTexture creates a 1x1 or custom-dimension solid colored texture.
func NewSolidColorTexture(width, height int, col color.Color) *Texture {
	img := ebiten.NewImage(width, height)
	img.Fill(col)
	return NewTextureFromImage(img)
}

// LoadTextureFromFile loads an image from the filesystem into a Texture.
func LoadTextureFromFile(path string) (*Texture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	ebitenImg := ebiten.NewImageFromImage(img)
	return NewTextureFromImage(ebitenImg), nil
}

// Image returns the underlying *ebiten.Image.
func (t *Texture) Image() *ebiten.Image {
	return t.image
}

// Width returns the texture width in pixels.
func (t *Texture) Width() int {
	return t.width
}

// Height returns the texture height in pixels.
func (t *Texture) Height() int {
	return t.height
}

// SubImage returns a sub-region of the texture.
func (t *Texture) SubImage(rect image.Rectangle) *Texture {
	sub := t.image.SubImage(rect).(*ebiten.Image)
	return NewTextureFromImage(sub)
}
