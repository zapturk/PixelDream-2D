using PixelDream2D.Models;

namespace PixelDream2D.Properties;

public class Rotation : Property
{
    /// <summary>
    /// The value of the rotation
    /// </summary>
    private float value;

    public Rotation(Entity entity, float rotation) : base(entity)
    {
        value = rotation;
        name = "Rotation";
    }

    public float getRotation()
    {
        return value;
    }

    public void setRotation(float rotation)
    {
        value = rotation % 360;
    }
}