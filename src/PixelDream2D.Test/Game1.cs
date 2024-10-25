using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;
using Microsoft.Xna.Framework.Input;

namespace PixelDream2D.Test;

public class Game1 : GameBase
{
    Texture2D texture2D;

    public Game1(GameSettings settings) : base(settings)
    {
    }

    protected override void Initialize()
    {
        // TODO: Add your initialization logic here
        base.Initialize();
    }

    protected override void LoadContent()
    {
        base.LoadContent();

        texture2D = Content.Load<Texture2D>("test");

        // TODO: use this.Content to load your game content here
    }

    protected override void Update(GameTime gameTime)
    {
         if (GamePad.GetState(PlayerIndex.One).Buttons.Back == ButtonState.Pressed || Keyboard.GetState().IsKeyDown(Keys.Escape))
            Exit();
         // TODO: Add your update logic here
         
         

         base.Update(gameTime);
    }

    protected override void Draw(GameTime gameTime)
    {
        GraphicsDevice.Clear(Color.CornflowerBlue);
        // TODO: Add your drawing code here 

        spriteBatch.Begin();
        spriteBatch.Draw(texture2D, new Vector2(100, 100), Color.White);
        spriteBatch.End();
        base.Draw(gameTime);
    }
}
