using PixelDream2D.Models;

namespace PixelDream2D.Properties;

public class Rotation : Property
{
    /// <summary>
    /// The value of the rotation
    /// </summary>
    private float _value;

    public Rotation(Entity entity, float rotation) : base(entity)
    {
        _value = rotation;
        Name = "Rotation";
    }

    public float GetRotation()
    {
        return _value;
    }

    public void SetRotation(float rotation)
    {
        _value = rotation % 360;
    }
}