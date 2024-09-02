using System;
using System.Collections.Generic;
using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Content;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Models;

public abstract class Scene : Disposable
{
    protected Game game { get; private set; }
    public bool HasInitialized { get; private set; }
    private uint entityIds;
    internal Dictionary<uint, Entity> Entities { get; set; }

    public Scene(Game gameState)
    {
        game = gameState;

        Entities = new Dictionary<uint, Entity>();
    }

    protected internal virtual void Init()
    {
        HasInitialized = true;
    }

    public abstract void LoadContent(ContentManager content);

    protected internal virtual void Update(GameTime gameTime){

        foreach (Entity entity in Entities.Values) {
            entity.Update(gameTime);
        }
    }

    protected internal virtual void Draw(GameTime gameTime, SpriteBatch spriteBatch)
    {
        foreach (var entity in Entities.Values) {
            entity.Draw(gameTime, spriteBatch);
        }
    }

    public abstract void UnloadContent();

    public Entity GetEntity(uint id) {
        return Entities[id];
    }

    public void AddEntity(Entity entity) {
        if (Entities.ContainsValue(entity)) {
            Console.WriteLine($"The entity with the id: [{entity.Id}] is already present in the Scene!");
            return;
        }

        Entities.Add(entityIds, entity);

        entityIds++;
    }

    protected override void Dispose(bool disposing) {
        if (disposing) {
            foreach (Entity entity in Entities.Values) {
                entity.Dispose();
            }
            
            entityIds = 0;
        }
    }
}