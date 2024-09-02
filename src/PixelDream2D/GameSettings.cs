
using System.Reflection;

namespace PixelDream2D;

public struct GameSettings
{
   public string Title { get; init; }
   public int Width { get; init; }
   public int Height { get; init; }
   public int TargetFps { get; init; }
   public string ContentRoot { get; init; }

   public GameSettings() {
        Title = Assembly.GetEntryAssembly()?.GetName().Name ?? "Pixel Dream 2D";
        Width = 1280;
        Height = 720;
        TargetFps = 60;
        ContentRoot = "Content";
    }

}