using PixelDream2D.Managers;
using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Content;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Models;

public abstract class Entity : Disposable
{
    public uint Id { get; internal set; }
    public bool HasInitialized { get; private set; }
    internal readonly Dictionary<Type, Property> Properties;

    protected Entity()
    {
        Properties = new Dictionary<Type, Property>();
    }

    protected internal virtual void Init()
    {
        foreach (Property property in Properties.Values)
        {
            if (!property.HasInitialized)
            {
                property.Init();
            }
        }
        HasInitialized = true;
    }

    protected internal virtual void LoadContent(ContentManager contentManager)
    {
        
    }

    protected internal virtual void Update(GameTime gameTime)
    {
        foreach (Property property in Properties.Values) {
            if (property.HasInitialized) {
                property.Update();
            }
        }
    }

    protected internal virtual void Draw(GameTime gameTime, SpriteBatch spriteBatch)
    {
        foreach (Property property in Properties.Values) {
            if (property.HasInitialized) {
                property.Draw(gameTime, spriteBatch);
            }
        }
    }

    public void AddProperty(Property property)
    {
        if (Properties.Any(comp => comp.GetType() == property.GetType()))
        {
            Console.WriteLine($"The property type [{property.GetType().Name}] is already present in the Entity[{Id}]!");
            return;
        }

        if (!property.HasInitialized)
        {
            property.OwningEntity = this;
            property.Init();
        }
        else
        {
            Console.WriteLine($"This property has already been added to a different Entity[{property.OwningEntity.Id}]. Please create a new property to add it to this Entity[{this.Id}].");
            return;
        }

        Properties.Add(property.GetType(), property);
    }

    public T? GetProperty<T>() where T : Property{
        return Properties[typeof(T)] as T;
    }

    public void RemoveProperty(Property property)
    {
        property.Dispose();
    }

    protected override void Dispose(bool disposing)
    {
        if (disposing)
        {
            foreach (Property property in Properties.Values)
            {
                property.Dispose();
            }

            SceneManager.ActiveScene?.Entities.Remove(Id);
        }
    }
}