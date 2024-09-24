using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Models;

public abstract class Property : Disposable {

    protected internal Entity OwningEntity { get; internal set; }

    protected internal string Name { get; internal set; }

    public bool HasInitialized { get; private set; }

    protected Property(Entity entity) {
        OwningEntity = entity;
    }

    protected internal virtual void Init() {
        HasInitialized = true;
    }

    protected internal virtual void Update() { }

    protected internal virtual void Draw(GameTime gameTime, SpriteBatch spriteBatch) { }

    protected override void Dispose(bool disposing) {
        if (disposing) {
            OwningEntity.Properties.Remove(GetType());
        }
    }
}