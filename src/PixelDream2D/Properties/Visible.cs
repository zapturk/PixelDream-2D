using PixelDream2D.Models;

namespace PixelDream2D.Properties;

public class Visible : Property
{
    /// <summary>
    /// Indicates if the entity is visible or not
    /// </summary>
    private bool value = true;

    public Visible(Entity entity, bool visible) : base(entity) {
        value = visible;
        name = "Visible";
    }

    public void Hide(){
        value = false;
    }

    public void Show(){
        value = true;
    }

    public bool IsVisible(){
        return value;
    }
}