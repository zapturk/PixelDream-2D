using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Models;

public class Sprite(Texture2D texture, Vector2 position)
{
    private readonly Texture2D Texture = texture;
    private readonly Vector2 Position = position;

    public virtual void Draw(SpriteBatch spriteBatch)
    {
        spriteBatch.Draw(Texture, Position, Color.White);
    }
}