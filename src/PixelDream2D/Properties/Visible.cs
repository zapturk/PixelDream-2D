using PixelDream2D.Models;

namespace PixelDream2D.Properties;

public class Visible : Property
{
    /// <summary>
    /// Indicates if the entity is visible or not
    /// </summary>
    private bool _value;

    public Visible(Entity entity, bool visible) : base(entity) {
        _value = visible;
        Name = "Visible";
    }

    public void Hide(){
        _value = false;
    }

    public void Show(){
        _value = true;
    }

    public bool IsVisible(){
        return _value;
    }
}