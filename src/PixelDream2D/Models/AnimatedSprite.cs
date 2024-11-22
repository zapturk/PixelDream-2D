using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Content;
using Microsoft.Xna.Framework.Graphics;

namespace PixelDream2D.Models;
/// <summary>
/// A helper class for handling animated sprite.
/// </summary>
public class AnimatedSprite
{
    // Number of frames in the animation.
    private int numOfFrames;
    
    // The animation spritesheet.
    private Texture2D _myTexture;
    
    // The number of frames to draw per second.
    private float _timePerFrame;
    
    // The current frame being drawn.
    private int _currentFrame;
    
    // Total amount of time the animation has been running.
    private float _totalElapsed;
    
    // Is the animation currently running?
    private bool _isPaused;

    // The current rotation, scale and draw depth for the animation.
    public readonly float Rotation;
    public readonly float Scale;
    public readonly float Depth;

    // The origin point of the animated texture.
    public Vector2 Origin;

    public AnimatedSprite(Vector2 origin, float rotation, float scale, float depth)
    {
        Origin = origin;
        Rotation = rotation;
        Scale = scale;
        Depth = depth;
    }

    public void Load(ContentManager content, string asset, int frameCount, int framesPerSec)
    {
        numOfFrames = frameCount;
        _myTexture = content.Load<Texture2D>(asset);
        _timePerFrame = 1.0f / framesPerSec;
        _currentFrame = 0;
        _totalElapsed = 0;
        _isPaused = false;
    }

    public void UpdateFrame(float elapsed)
    {
        // do nothing if it is paused
        if (_isPaused)
            return;
        
        _totalElapsed += elapsed;
        if (_totalElapsed > _timePerFrame)
        {
            _currentFrame++;
            // Keep the Frame between 0 and the total frames, minus one.
            _currentFrame %= numOfFrames;
            _totalElapsed -= _timePerFrame;
        }
    }

    public void DrawFrame(SpriteBatch batch, Vector2 screenPos)
    {
        DrawFrame(batch, _currentFrame, screenPos);
    }

    public void DrawFrame(SpriteBatch batch, int frameToDraw, Vector2 screenPos)
    {
        int frameWidth = _myTexture.Width / numOfFrames;
        
        Rectangle sourcerect = new Rectangle(frameWidth * frameToDraw, 0, frameWidth, _myTexture.Height);
        
        batch.Draw(_myTexture, screenPos, sourcerect, Color.White, Rotation, Origin, Scale, SpriteEffects.None, Depth);
    }

    public bool IsPlaying()
    {
        return !_isPaused;
    } 

    public void Reset()
    {
        _currentFrame = 0;
        _totalElapsed = 0f;
    }

    public void Stop()
    {
        Pause();
        Reset();
    }

    public void Play()
    {
        _isPaused = false;
    }

    public void Pause()
    {
        _isPaused = true;
    }
}
