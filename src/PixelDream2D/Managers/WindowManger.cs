using PixelDream2D.Models;
using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Managers;

/// <summary>
/// Manges the window with resizing and full screen
/// </summary>
public static class WindowManger
{
    private static GraphicsDeviceManager graphics;
    private static GameBase gameBase;
    private static Canvas canvas;

    public static void Init(GameBase game, GraphicsDeviceManager graphics, Canvas canvas)
    {
        gameBase = game;
        WindowManger.graphics = graphics;
        WindowManger.canvas = canvas;
        SetResolution(gameBase.Settings.Width, gameBase.Settings.Height);
    }

    public static void Update()
    {
        if (!gameBase.Window.IsBorderless)
            SetResolution(gameBase.GraphicsDevice.Viewport.Width, gameBase.GraphicsDevice.Viewport.Height);
    }

    public static void SetResolution(int height, int width)
    {
        graphics.PreferredBackBufferWidth = height;
        graphics.PreferredBackBufferHeight = width;
        gameBase.Window.IsBorderless = false;
        graphics.ApplyChanges();
        canvas.SetDestinationRectangle();
    }

    private static void SetFullScreen()
    {
        graphics.PreferredBackBufferWidth = GraphicsAdapter.DefaultAdapter.CurrentDisplayMode.Width;
        graphics.PreferredBackBufferHeight = GraphicsAdapter.DefaultAdapter.CurrentDisplayMode.Height;
        gameBase.Window.IsBorderless = true;
        graphics.ApplyChanges();
        canvas.SetDestinationRectangle();
    }

    /// <summary>
    /// Toggles between full screen and window
    /// </summary>
    public static void ToggleFullScreen()
    {
        if (gameBase.Window.IsBorderless)
        {
            SetResolution(gameBase.Settings.Width, gameBase.Settings.Height);
        }
        else
        {
            SetFullScreen();
        }
    }
}