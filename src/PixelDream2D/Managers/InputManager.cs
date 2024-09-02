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
        if(inputManager.FindNewGamepads()){
            inputManager.SetActivePlayers(new List<int> { 0 });
            inputManager.GetCapabilities(0);
        }
    }

    public static void Update(GameTime gameTime)
    {
        inputManager.Update(gameTime);
    }
}