﻿using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Input;

namespace PixelDream2D.Test;

public class Game1 : GameBase
{

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

        base.Draw(gameTime);
    }
}
