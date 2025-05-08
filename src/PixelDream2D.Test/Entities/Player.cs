using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;
using PixelDream2D.Models;

namespace PixelDream2D.Test.Entities;

public class Player : Entity
{
    private Vector2 _position;
    private bool _visible;
    private Sprite _sprite;
    
    public Player(Vector2 position, bool visible = true) : base()
{
    _position = position;
    _visible = visible;
    
}

public void Move(int direction)
{
    switch (direction)
    {
        case 0: // Left
            _position.X -= 5; // Adjust by a desired amount
            break;
        case 1: // Right
            _position.X += 5; // Adjust by a desired amount
            break;
        case 2: // Up
            _position.Y -= 5; // Adjust by a desired amount
            break;
        case 3: // Down
            _position.Y += 5; // Adjust by a desired amount
            break;
    }
}

}
