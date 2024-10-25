using PixelDream2D.Handlers;
using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Input;

namespace PixelDream2D.Managers;

public static class InputManager
{
    private static InputHandler inputManager;

    public static void Init(Game game)
    {
        inputManager = new InputHandler(game);
        
        // if there are no controllers connected exit
        if (!inputManager.FindNewGamepads()) return;
        
        inputManager.SetActivePlayers([0]);
        inputManager.GetCapabilities(0);
    }

    public static void Update(GameTime gameTime)
    {
        inputManager.Update(gameTime);
    }
}