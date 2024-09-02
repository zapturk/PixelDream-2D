using PixelDream2D.Managers;
using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Models;

public abstract class Entity : Disposable
{
    public uint Id { get; internal set; }
    public bool HasInitialized { get; private set; }
    public float Rotation { get; private set; }
    internal readonly Dictionary<Type, Component> Components;

    public Entity(Vector2 position, bool visible = true)
    {
        Components = new Dictionary<Type, Component>();
        AddComponent(new Position(this, position));
        AddComponent(new Visible(this, visible));
    }

    protected internal virtual void Init()
    {
        foreach (Component component in this.Components.Values)
        {
            if (!component.HasInitialized)
            {
                component.Init();
            }
        }
        HasInitialized = true;
    }

    protected internal virtual void Update(GameTime gameTime)
    {
        foreach (Component component in this.Components.Values) {
            if (component.HasInitialized) {
                component.Update();
            }
        }
    }

    protected internal virtual void Draw(GameTime gameTime, SpriteBatch spriteBatch)
    {
        // if not visible, don't draw
        if(!GetComponent<Visible>().IsVisible()) 
            return;
        
        foreach (Component component in this.Components.Values) {
            if (component.HasInitialized) {
                component.Draw(gameTime, spriteBatch);
            }
        }
    }

    public void AddComponent(Component component)
    {
        if (Components.Any(comp => comp.GetType() == component.GetType()))
        {
            Console.WriteLine($"The component type [{component.GetType().Name}] is already present in the Entity[{this.Id}]!");
            return;
        }

        if (!component.HasInitialized)
        {
            component.owningEntity = this;
            component.Init();
        }
        else
        {
            Console.WriteLine($"This component has already been added to a different Entity[{component.owningEntity.Id}]. Please create a new component to add it to this Entity[{this.Id}].");
            return;
        }

        Components.Add(component.GetType(), component);
    }

    public T? GetComponent<T>() where T : Component{
        return Components[typeof(T)] as T;
    }

    public void RemoveComponent(Component component)
    {
        component.Dispose();
    }

    public void SetRotation(float rotation)
    {
        Rotation = rotation % 360;
    }

    protected override void Dispose(bool disposing)
    {
        if (disposing)
        {
            foreach (Component component in this.Components.Values)
            {
                component.Dispose();
            }

            SceneManager.ActiveScene?.Entities.Remove(Id);
        }
    }
}