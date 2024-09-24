
using System.Reflection;

namespace PixelDream2D;

public struct GameSettings()
{
   public string Title { get; init; } = Assembly.GetEntryAssembly()?.GetName().Name ?? "Pixel Dream 2D";
   public int Width { get; init; } = 1280;
   public int Height { get; init; } = 720;
   public int TargetFps { get; init; } = 60;
   public string ContentRoot { get; init; } = "Content";
}