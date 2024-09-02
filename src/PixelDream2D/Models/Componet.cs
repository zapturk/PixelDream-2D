using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Models;

public abstract class Component : Disposable {
    
    protected internal Entity owningEntity { get; internal set; }

    protected internal string name { get; internal set; }

    public bool HasInitialized { get; private set; }
    
    public Component(Entity entity) {
        owningEntity = entity;
    }
    
    protected internal virtual void Init() {
        HasInitialized = true;
    }

    protected internal virtual void Update() { }

    protected internal virtual void Draw(GameTime gameTime, SpriteBatch spriteBatch) { }

    protected override void Dispose(bool disposing) {
        if (disposing) {
            owningEntity.Components.Remove(GetType());
        }
    }
}