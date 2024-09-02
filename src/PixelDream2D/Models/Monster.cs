using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Models;

public abstract class Monster : Entity
{
    public Texture2D Texture;
    public int Health;
    public int Damage;

    public Monster(Vector2 position, bool visible = true) : base(position, visible){

    }

    protected internal override void Draw(GameTime gameTime, SpriteBatch spriteBatch)
    {
        spriteBatch.Draw(Texture, Position, Color.Red);

        base.Draw(gameTime, spriteBatch);
    }   
}