using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Models;

public class Sprite(Texture2D texture, Vector2 position)
{
    protected Texture2D texture = texture;
    public Vector2 Position { get; set; } = position;

    public virtual void Draw(SpriteBatch spriteBatch)
    {
        spriteBatch.Draw(texture, Position, Color.White);
    }
}