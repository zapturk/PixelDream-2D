using PixelDream2D.Models;
using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Managers;

/// <summary>
/// Manges the window with resizing and full screen
/// </summary>
public static class WindowManger
{
    private static GraphicsDeviceManager _graphics;
    private static GameBase _gameBase;
    private static Canvas _canvas;

    public static void Init(GameBase game, GraphicsDeviceManager graphics, Canvas canvas)
    {
        _gameBase = game;
        _graphics = graphics;
        _canvas = canvas;
        SetResolution(_gameBase.Settings.Width, _gameBase.Settings.Height);
    }

    public static void Update()
    {
        if (!_gameBase.Window.IsBorderless)
            SetResolution(_gameBase.GraphicsDevice.Viewport.Width, _gameBase.GraphicsDevice.Viewport.Height);
    }

    public static void SetResolution(int height, int width)
    {
        _graphics.PreferredBackBufferWidth = height;
        _graphics.PreferredBackBufferHeight = width;
        _gameBase.Window.IsBorderless = false;
        _graphics.ApplyChanges();
        _canvas.SetDestinationRectangle();
    }

    private static void SetFullScreen()
    {
        _graphics.PreferredBackBufferWidth = GraphicsAdapter.DefaultAdapter.CurrentDisplayMode.Width;
        _graphics.PreferredBackBufferHeight = GraphicsAdapter.DefaultAdapter.CurrentDisplayMode.Height;
        _gameBase.Window.IsBorderless = true;
        _graphics.ApplyChanges();
        _canvas.SetDestinationRectangle();
    }

    /// <summary>
    /// Toggles between full screen and window
    /// </summary>
    public static void ToggleFullScreen()
    {
        if (_gameBase.Window.IsBorderless)
        {
            SetResolution(_gameBase.Settings.Width, _gameBase.Settings.Height);
        }
        else
        {
            SetFullScreen();
        }
    }
}