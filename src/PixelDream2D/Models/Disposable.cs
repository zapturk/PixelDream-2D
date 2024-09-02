namespace PixelDream2D.Models;

public abstract class Disposable : IDisposable {
    
    public bool HasDisposed { get; private set; }
    
    /// <summary>
    /// Disposes of the object, allowing for proper resource cleanup and finalization.
    /// </summary>
    ~Disposable() {
        Dispose();
    }
    
    /// <summary>
    /// Disposes of the object, allowing for proper resource cleanup and finalization.
    /// </summary>
    public void Dispose() {
        if (HasDisposed) {
            Console.WriteLine($"This object of type [{GetType().Name}] has already been disposed.");
            return;
        }
        
        Dispose(true);
        GC.SuppressFinalize(this);
        HasDisposed = true;
    }

    /// <summary>
    /// Disposes of the object and releases associated resources. 
    /// </summary>
    /// <param name="disposing">True if called from user code; false if called from a finalizer.</param>
    protected abstract void Dispose(bool disposing);
    
    /// <summary>
    /// Throws an exception if the object has been disposed, indicating that it is no longer usable.
    /// </summary>
    protected void ThrowIfDisposed() {
        if (HasDisposed) {
            Console.WriteLine(new ObjectDisposedException(GetType().Name));
        }
    }
}