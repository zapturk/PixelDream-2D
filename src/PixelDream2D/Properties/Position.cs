using Microsoft.Xna.Framework;
using PixelDream2D.Models;

namespace PixelDream2D.Properties;

public class Position : Property
{
    /// <summary>
    /// The value of the position
    /// </summary>
    private Vector2 _value;

    public Position(Entity entity, Vector2 pos) : base(entity)
    {
        _value = pos;
        Name = "Position";
    }

    public Vector2 GetPosition()
    {
        return _value;
    }

    public void SetPosition(Vector2 pos)
    {
        _value = pos;
    }
}