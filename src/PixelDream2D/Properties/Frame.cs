using PixelDream2D.Enums;
using PixelDream2D.Models;

namespace PixelDream2D.Properties;

public class Frame : Property
{
    private int _value = 0;
    private readonly int _maxFrame = 0;
    private FrameState _state = FrameState.Stop;
    private FrameDirection _direction = FrameDirection.Forward;

    public Frame(Entity entity, int value, int maxFrame = 0) : base(entity)
    {
        _value = value;
        Name = "Frame";
        _maxFrame = maxFrame;
    }

    protected internal override void Update()
    {
        base.Update();
        // do nothing if state is stopped
        if (_state == FrameState.Stop) return;
        
        switch (_direction)
        {
            case FrameDirection.Forward when _value >= _maxFrame:
                _value = 0;
                break;
            case FrameDirection.Forward:
                _value += 1;
                break;
            case FrameDirection.Backward when _value <= 0:
                _value = _maxFrame;
                break;
            case FrameDirection.Backward:
                _value -= 1;
                break;
        }
    }

    public void Play()
    {
        _direction = FrameDirection.Forward;
        _state = FrameState.Play;
    }

    public void PlayBackward()
    {
        _direction = FrameDirection.Backward;
        _state = FrameState.Play;
    }

    /// <summary>
    /// pauses the frame count on the current frame
    /// </summary>
    public void Pause()
    {
        _state = FrameState.Stop;
    }

    /// <summary>
    /// Stop the frame count and resets to the first frame
    /// </summary>
    public void Stop()
    {
        Pause();
        Reset();
    }

    public int GetCurrentFrame()
    {
        return _value;
    }

    public void SetCurrentFrame(int frame)
    {
        _value = frame;
    }

    public void Reset()
    {
        _value = _direction switch
        {
            FrameDirection.Forward => 0,
            FrameDirection.Backward => _maxFrame,
            _ => _value
        };
    }
}