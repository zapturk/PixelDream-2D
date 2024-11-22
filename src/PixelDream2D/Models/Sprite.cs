using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Models;

public class Sprite(Texture2D texture, Vector2 position)
{
    private Texture2D Texture = texture;

    public virtual void Draw(SpriteBatch spriteBatch)
    {
        spriteBatch.Draw(Texture, position, Color.White);
    }
}