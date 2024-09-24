using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;
using Microsoft.Xna.Framework.Input;

namespace PixelDream2D;

public class GameBase : Game
{
    // private Canvas canvas;
    private GraphicsDeviceManager _graphics;
    public SpriteBatch spriteBatch;
    public readonly GameSettings Settings;

    public GameBase(GameSettings settings)
    {
        _graphics = new GraphicsDeviceManager(this);
        Settings = settings;
        Content.RootDirectory = Settings.ContentRoot;
        IsMouseVisible = true;
        Window.AllowAltF4 = true;
    }

    protected override void Initialize()
    {
        base.Initialize();
    }

    protected override void LoadContent()
    {
        spriteBatch = new SpriteBatch(GraphicsDevice);
    }

    protected override void Update(GameTime gameTime)
    {
        base.Update(gameTime);
    }

    protected override void Draw(GameTime gameTime)
    {
        base.Draw(gameTime);
    }
}
