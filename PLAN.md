# PixelDream-2D: Architecture & Development Plan

## 1. Vision & Core Philosophy

**PixelDream-2D** is a dedicated 2D game engine built in Go specifically designed for pixel art games. It pairs a high-performance, pixel-perfect runtime with an integrated visual editor for tilemaps and pixel art.

### Guiding Principles
- **Pixel-Perfect by Default**: Built-in virtual canvas with integer scaling, sub-pixel positioning handling, and aspect ratio management (letterboxing/pillarboxing).
- **Integrated Tooling**: Unified workflow where pixel art editing, map design, and game playtesting happen in one cohesive in-engine window.
- **Pure Go Performance**: Idiomatic Go architecture with zero external scripting overhead, compile-time safety, and high frame rates.

---

## 2. Technology Stack & Key Decisions

| Component | Selected Approach | Rationale |
| :--- | :--- | :--- |
| **Engine Core & Rendering** | **Ebitengine (`github.com/hajimehoshi/ebiten/v2`)** | Battle-tested, highly optimized for 2D pixel art, cross-platform (macOS, Windows, Linux, Web/WASM), robust audio and input. |
| **Frontend Editor UI** | **Integrated In-Engine UI** (Dear ImGui / EbitenUI) | Runs directly in the engine viewport for seamless instant playtesting, zero IPC/serialization overhead between editor and runtime. |
| **Game Logic** | **Pure Go Engine API** | Clean, type-safe Go packages (`pkg/engine`, `pkg/tilemap`, `pkg/sprite`, etc.), allowing games to be compiled directly with maximum performance. |

---

## 3. High-Level Architecture

```
┌───────────────────────────────────────────────────────────────┐
│                    PixelDream Editor Frontend                 │
│  ┌──────────────────┐  ┌──────────────────┐  ┌─────────────┐  │
│  │ Pixel Art Editor │  │  Tilemap Editor  │  │ Scene/Props │  │
│  │ (Canvas/Palette) │  │  (Layers/Brushes)│  │  Inspector  │  │
│  └─────────┬────────┘  └────────┬─────────┘  └──────┬──────┘  │
│            └────────────────────┼───────────────────┘         │
│                                 ▼                             │
│                      Integrated Editor GUI                    │
├───────────────────────────────────────────────────────────────┤
│                      PixelDream Engine Core                   │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────────┐  │
│  │ Virtual Canvas│  │  Scene / ECS  │  │   Tilemap Engine  │  │
│  │ & Pixel Render│  │  Game Objects │  │  & Collision Grid │  │
│  └───────────────┘  └───────────────┘  └───────────────────┘  │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────────┐  │
│  │ Input Manager │  │ Audio Engine  │  │ Resource/Assets   │  │
│  └───────────────┘  └───────────────┘  └───────────────────┘  │
├───────────────────────────────────────────────────────────────┤
│                     Ebitengine Backend                        │
│            (OpenGL / Metal / DirectX, Audio, Input)           │
└───────────────────────────────────────────────────────────────┘
```

---

## 4. Engine Core Modules

### 4.1 Virtual Resolution & Display Pipeline
- **Fixed Virtual Canvas**: Games render internally to a low-resolution target (e.g., $320 \times 180$, $256 \times 240$, $640 \times 360$).
- **Scaling Modes**:
  - Integer Scaling (Crisp pixels, zero shimmering, aspect ratio preservation).
  - Fit-to-window with letterboxing/pillarboxing.
- **Camera System**:
  - Pixel-snapping camera coordinates to eliminate sub-pixel jitter.
  - Smooth target tracking, deadzones, camera bounds, and screen shake.

### 4.2 Pixel Art & Animation Pipeline
- **Sprite Sheet Management**: Texture atlas loading, automatic slicing, and origin/anchor handling.
- **Animation System**: Frame-by-frame animation controller with custom frame durations, loop modes, and event dispatchers (e.g., footstep sounds, attack hit frames).
- **Palettes & Shaders**:
  - Indexed color palette swapping.
  - Pixel art-friendly lighting and tinting shaders.

### 4.3 Tilemap System
- **Layering**: Multi-layer tilemaps (Background, Collision/Solid, Foreground, Foreground Overlay).
- **Tile Properties**: Passable/solid flags, tile tags, slope definitions, and animated tiles.
- **Autotiling**: Bitmask-based autotiling (16-pipe and 47-pipe) for automated borders and connections.
- **Spatial Collision Grid**: Fast spatial lookup for pixel-perfect and AABB tile collision.

### 4.4 Entity & Scene Architecture
- **Composition-based Entity Model**: Lightweight game object / component model in idiomatic Go.
- **Physics & Kinematics**:
  - 2D Kinematic character controller with tile collision.
  - Slopes, one-way platforms, moving platforms, and raycasting.

### 4.5 Audio & Input
- **Input Actions**: Unified action mapping (`MoveLeft`, `Jump`, `Attack`) supporting keyboard, mouse, and gamepads.
- **Audio Manager**: Sound effect pooling, spatial panning (2D stereo), and seamless BGM streaming with loop points.

---

## 5. Frontend Editor Suite

### 5.1 Pixel Art / Sprite Editor
- **Drawing Canvas**:
  - Grid overlay with dynamic zoom (1x to 32x), pan, and onion skinning.
  - Tools: Pencil (1px, 2px brush), Eraser, Color Picker, Paint Bucket, Selection/Marquee.
- **Palette Manager**:
  - Color swatches, import/export (GPL, HEX, PNG).
- **Animation Timeline**:
  - Playback preview at target FPS.
  - Frame sequencing, duplication, and onion-skin previews.
- **Hitbox / Hurtbox Tool**:
  - Define custom collision and attack boxes directly on sprite frames.

### 5.2 Tile & Map Editor
- **Tileset Panel**:
  - Load tilesheets, configure grid size ($8 \times 8$, $16 \times 16$, $32 \times 32$), and assign autotile rules.
- **Map Viewport**:
  - Multi-layer editing with visibility/lock toggles.
  - Brushes: Single tile, Rectangle fill, Bucket fill, Random variation, and Autotile brush.
- **Entity & Trigger Placement**:
  - Visually place spawn points, trigger zones, NPC markers, and dynamic objects.

### 5.3 Scene & Project Manager
- **Asset Hierarchy**: Visual explorer for sprites, maps, audio, and scenes.
- **Instant Playtesting**: Switch between "Edit Mode" and "Play Mode" with a single keypress or button directly within the engine window.

---

## 6. Phased Implementation Roadmap

### Phase 1: Engine Foundation & Core Runtime (Completed)
- [x] Initialize Go module (`github.com/zapturk/PixelDream-2D`).
- [x] Establish project directory structure (`cmd/`, `pkg/display/`, `pkg/gfx/`, `pkg/input/`, `pkg/engine/`).
- [x] Implement virtual resolution canvas with integer scaling, letterboxing/pillarboxing, and screen-to-canvas coordinate mapping (`pkg/display/canvas.go`).
- [x] Implement pixel-snapping 2D camera with smoothing, bounds clamping, and screen shake (`pkg/display/camera.go`).
- [x] Implement texture loader, sub-image slicing, and pivot-aware sprite rendering (`pkg/gfx/`).
- [x] Implement unified action mapping for keyboard, mouse, and gamepad (`pkg/input/`).
- [x] Implement engine runner loop and lifecycle (`pkg/engine/`).
- [x] Create runnable Phase 1 interactive demo (`cmd/demo/main.go`).

### Phase 2: Sprite Animation & Tile Engine (Completed)
- [x] Implement frame-by-frame sprite animation controller with loop modes (`LoopForever`, `LoopOnce`, `LoopPingPong`) and frame event callbacks (`pkg/gfx/animation.go`).
- [x] Implement multi-layer tilemap data structures, grid slicing, and camera viewport culling (`pkg/tilemap/`).
- [x] Implement 4-bit autotiling bitmask helper (`pkg/tilemap/autotile.go`).
- [x] Implement 2D AABB and kinematic collision controller with gravity, jumping, wall/ceiling detection, and ground sliding (`pkg/physics/`).
- [x] Implement debug collider boundary visualizer for solid tiles and entity AABB boxes (`pkg/physics/debug.go`).
- [x] Comprehensive unit tests across animation, tilemaps, and physics (`pkg/gfx/animation_test.go`, `pkg/tilemap/tilemap_test.go`, `pkg/physics/kinematic_test.go`).
- [x] Update interactive demo (`cmd/demo/main.go`) showcasing animated player running, jumping, and colliding with a multi-layered tilemap.

### Phase 3: Integrated In-Engine Editor Shell & Map Editor (Completed)
- [x] Implement JSON map serialization, file saving/loading, and layer serialization (`pkg/tilemap/serialize.go`).
- [x] Build lightweight, high-performance in-engine UI toolkit and pixel dark theme (`pkg/editor/ui.go`, `pkg/editor/theme.go`).
- [x] Implement editing tools: Pencil, Eraser, Rect Fill, Bucket (flood fill), and Eyedropper/Picker (`pkg/editor/tools.go`).
- [x] Implement in-engine Editor shell with Toolbar, Tileset Palette, Layer manager, status bar, and camera panning (`pkg/editor/editor.go`).
- [x] Implement instant Playtest Mode toggle (`F5`) to seamlessly test game physics, player kinematics, and collisions directly within the editor.
- [x] Build standalone editor application (`cmd/editor/main.go` -> `bin/editor`).
- [x] Unit tests for map serialization and editor tools (`pkg/tilemap/serialize_test.go`, `pkg/editor/tools_test.go`).

### Phase 4: Integrated Pixel Art & Animation Editor (Completed)
- [x] Implement multi-frame pixel art document model with snapshot-based Undo/Redo history (`pkg/pixeleditor/document.go`).
- [x] Implement retro color palette manager (PICO-8, DawnBringer) with indexed swatches (`pkg/pixeleditor/document.go`).
- [x] Implement pixel drawing tools: Pencil (Bresenham line smoothing), Eraser, Eyedropper/Picker, Flood Fill (Bucket), Line tool, and Rectangle tool (`pkg/pixeleditor/tools.go`).
- [x] Implement pixel art editor canvas with zoom ($20\times$), transparency checkerboard, pixel grid overlay, and Onion Skinning (`pkg/pixeleditor/editor.go`).
- [x] Implement animation timeline with frame thumbnails, frame duplication/deletion, and real-time looping preview (`pkg/pixeleditor/editor.go`).
- [x] Integrated dual-mode editor shell: toggle between **Map Editor** and **Pixel Art Editor** using the top tab bar or `Tab` key (`pkg/editor/editor.go`, `cmd/editor/main.go`).
- [x] Comprehensive unit tests for pixel tools, flood fill, Bresenham lines, and undo/redo stacks (`pkg/pixeleditor/tools_test.go`).

### Phase 5: Engine Polish & Sample Game
- Integrate audio engine (SFX pooling, BGM loops).
- Add pixel particle system (dust, weather, sparks).
- Build a sample 2D platformer / top-down game to validate all systems and ergonomics.
