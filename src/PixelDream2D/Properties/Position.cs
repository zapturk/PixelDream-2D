using Microsoft.Xna.Framework;
using PixelDream2D.Models;

namespace PixelDream2D.Properties;

public class Position : Property
{
    /// <summary>
    /// The value of the position
    /// </summary>
    private Vector2 value;

    public Position(Entity entity, Vector2 pos) : base(entity)
    {
        value = pos;
        name = "Position";
    }

    public Vector2 getPosition()
    {
        return value;
    }

    public void setPosition(Vector2 pos)
    {
        value = pos;
    }

}