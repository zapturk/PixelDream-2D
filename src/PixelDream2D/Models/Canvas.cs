using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Models;

/// <summary>
/// Modle class to draw on the screen to mainitain aspect ratio
/// </summary>
public class Canvas
{
    private readonly RenderTarget2D target;
    private readonly GraphicsDevice graphicsDevice;
    private Rectangle destinationRectangle;

    public Canvas(GraphicsDevice graphicsDevice, int width, int height)
    {
        this.graphicsDevice = graphicsDevice;
        target = new(this.graphicsDevice, width, height);
    }

    public void SetDestinationRectangle()
    {
        var screenSize = graphicsDevice.PresentationParameters.Bounds;

        float scaleX = (float)screenSize.Width / target.Width;
        float scaleY = (float)screenSize.Height / target.Height;
        float scale = Math.Min(scaleX, scaleY);

        int newWidth = (int)(target.Width * scale);
        int newHeight = (int)(target.Height * scale);

        int posX = (screenSize.Width - newWidth) / 2;
        int posY = (screenSize.Height - newHeight) / 2;

        destinationRectangle = new Rectangle(posX, posY, newWidth, newHeight);
    }

    public void Activate()
    {
        graphicsDevice.SetRenderTarget(target);
        graphicsDevice.Clear(Color.DarkGray);
    }

    public void Draw(SpriteBatch spriteBatch)
    {
        graphicsDevice.SetRenderTarget(null);
        graphicsDevice.Clear(Color.Black);
        spriteBatch.Begin(SpriteSortMode.Deferred, null, SamplerState.PointClamp);
        spriteBatch.Draw(target, destinationRectangle, Color.White);
        spriteBatch.End();
    }
}