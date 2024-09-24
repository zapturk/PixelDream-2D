using PixelDream2D.Models;
using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Content;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Managers;

/// <summary>
/// Manages the active scene and switching between scenes
/// </summary>
public static class SceneManager
{
    public static Scene? ActiveScene { get; private set; }

    private static GameBase _gameBase;

    public static void Init(GameBase gameInstance)
    {
        _gameBase = gameInstance;
        if (ActiveScene?.HasInitialized == false)
        {
            ActiveScene?.Init();
            ActiveScene?.LoadContent(_gameBase.Content);
        }
    }

    public static void Init(GameBase gameInstance, Scene scene)
    {
        SetDefaultScene(scene);
        Init(gameInstance);
    }

    internal static void LoadContent(ContentManager content)
    {
        ActiveScene?.LoadContent(content);
    }

    public static void Update(GameTime gameTime)
    {
        ActiveScene?.Update(gameTime);
    }

    public static void Draw(GameTime gameTime, SpriteBatch spriteBatch)
    {
        ActiveScene?.Draw(gameTime, spriteBatch);
    }

    public static void SetDefaultScene(Scene scene)
    {
        ActiveScene = scene;
    }

    public static void SetScene(Scene scene)
    {
        // Unload and load the new scene
        ActiveScene?.UnloadContent();

        // Set the new scene
        ActiveScene = scene;

        // Initialize the new scene
        ActiveScene.Init();

        // Load content for the new scene
        ActiveScene.LoadContent(_gameBase.Content);
    }
}