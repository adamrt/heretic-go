package heretic

import (
	"log"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS = 60
	// Number of milliseconds per frame
	TargetFrameTime = (1000 / FPS)
)

type CullMode int

const (
	CullModeNone     CullMode = 0
	CullModeBackFace CullMode = 1
)

type RenderMode int

const (
	RenderModeWire        RenderMode = 1
	RenderModeWireVertex  RenderMode = 2
	RenderModeWireFill    RenderMode = 3
	RenderModeFill        RenderMode = 4
	RenderModeTexture     RenderMode = 5
	RenderModeTextureWire RenderMode = 6
)

var leftButtonDown bool = false

func NewEngine(window *Window, framebuffer *Framebuffer) *Engine {
	// Projection matrix. We only need this calculate this once.
	aspectX := float64(window.width) / float64(window.height)
	aspectY := float64(window.height) / float64(window.width)
	fovY := math.Pi / 3.0 // Same as 180/3 or 60deg
	fovX := math.Atan(math.Tan(fovY/2.0)*aspectX) * 2.0
	znear := 0.3
	zfar := 100.0

	return &Engine{
		window:      window,
		framebuffer: framebuffer,
		IsRunning:   true,

		ambientLight: DirectionalLight{Direction: Vec3{0, 0, 1}},

		projMatrix: MatrixMakePerspective(fovY, aspectY, znear, zfar),
		frustum:    NewFrustum(fovX, fovY, znear, zfar),
		camera:     NewCamera(Vec3{-1.0, 1.0, -1.0}, Vec3{0.0, 0.0, 0.0}, Vec3{0.0, 1.0, 0.0}, window.width, window.height),
		cullMode:   CullModeBackFace,
		renderMode: RenderModeWireFill,

		scene: NewScene(),
		// Rotation is set so if the user presses spacebar they get some
		// rotation. But autoRotation is off by default. Use
		// SetAutoRotation() to override the rotation value.
		autoRotation: false,
		rotation:     Vec3{0, 0.5, 0},
	}
}

// meshReader is a temporary interface to avoid circular imports with the fft
// package. It will be removed once the project is better organized.
type meshReader interface {
	ReadMesh(mapNum int) Mesh
}

// Engine is the top level object that contains windows, renderers, etc.
// It has a basic game loop and in the typical for is run like so:
//
// engine.Setup()
// for engine.IsRunning() {
//   engine.ProcessInput()
//   engine.Update()
//   engine.Render()
// }
//
type Engine struct {
	window      *Window
	framebuffer *Framebuffer

	IsRunning bool

	// Timing
	previous  uint32
	deltaTime float64

	// Rendering
	cullMode   CullMode
	renderMode RenderMode
	projMatrix Matrix
	camera     *Camera
	frustum    Frustum

	// Model
	scene *scene

	MeshReader meshReader
	currentMap int

	ambientLight DirectionalLight

	// These two control the mesh rotating on its own.
	// The amount can be set by SetAutoRotation().
	autoRotation bool
	rotation     Vec3
}

func (e *Engine) Setup() {
	if len(e.scene.Meshes) == 0 {
		log.Fatalln("no mesh specified")
	}

	// If there is any texture on any mesh, show it.
	for _, mesh := range e.scene.Meshes {
		if len(mesh.Texture.data) != 0 {
			e.renderMode = RenderModeTexture
		}
	}
	e.previous = sdl.GetTicks()
}

func (e *Engine) ProcessInput() {
	// The other mouse/keyboard functionality can be handled via polling.
	// Moving the keys into keyboard state polling (above) will appear to be
	// multiple presses in a row when we just want one.
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			e.IsRunning = false
			break
		case *sdl.MouseWheelEvent:
			for i, mesh := range e.scene.Meshes {
				if t.PreciseY > 0 {
					e.scene.Meshes[i].Scale = mesh.Scale.Mul(1.5)
				} else {
					e.scene.Meshes[i].Scale = mesh.Scale.Div(1.5)
				}
			}
		case *sdl.MouseButtonEvent:
			if t.Button == sdl.BUTTON_LEFT {
				leftButtonDown = t.Type == sdl.MOUSEBUTTONDOWN
			}
		case *sdl.MouseMotionEvent:
			if leftButtonDown {
				e.camera.ProcessMouseMovement(float64(t.XRel), float64(t.YRel), e.deltaTime)
			}
		case *sdl.KeyboardEvent:
			if t.Type != sdl.KEYDOWN {
				continue
			}
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				e.IsRunning = false
				break

			case sdl.K_1:
				e.renderMode = RenderModeWire
			case sdl.K_2:
				e.renderMode = RenderModeWireVertex
			case sdl.K_3:
				e.renderMode = RenderModeWireFill
			case sdl.K_4:
				e.renderMode = RenderModeFill
			case sdl.K_5:
				e.renderMode = RenderModeTexture
			case sdl.K_6:
				e.renderMode = RenderModeTextureWire
			case sdl.K_c:
				e.cullMode = CullModeNone
			case sdl.K_b:
				e.cullMode = CullModeBackFace
			case sdl.K_k:
				e.NextMap()
			case sdl.K_j:
				e.PrevMap()
			case sdl.K_SPACE:
				e.autoRotation = !e.autoRotation
			}
		}
	}
}

func (e *Engine) Update() {
	// Target the specified FPS
	wait := TargetFrameTime - (sdl.GetTicks() - e.previous)
	if wait > 0 && wait <= TargetFrameTime {
		sdl.Delay(wait)
	}

	// Using a deltaTime for transformations keeps animation speed
	// consistent regardless of FPS. It basically changes the engine
	// transforms-per-frame to tranforms-per-second.
	e.deltaTime = float64(sdl.GetTicks()-e.previous) / 1000.0
	e.previous = sdl.GetTicks()

	for _, mesh := range e.scene.Meshes {
		// Apply the engine's rotation vector. This is for automatic rotation.
		if e.autoRotation {
			mesh.Rotation = mesh.Rotation.Add(e.rotation.Mul(e.deltaTime))
		}

		// World matrix. Combination of scale, rotation and translation.
		worldMatrix := MatrixIdentity()
		worldMatrix = worldMatrix.Mul(NewScaleMatrix(mesh.Scale))
		worldMatrix = worldMatrix.Mul(NewRotationMatrix(mesh.Rotation))
		worldMatrix = worldMatrix.Mul(NewTranslationMatrix(mesh.Translation))

		// Project each into 2D
		for _, triangle := range mesh.Triangles {
			triangle.Projected = make([]Vec4, 3)

			// Transformation
			for i := 0; i < 3; i++ {
				transformedVertex := triangle.Points[i].Vec4()
				// Transform to world space
				transformedVertex = worldMatrix.MulVec4(transformedVertex)
				// Transform to view space
				transformedVertex = e.camera.ViewMatrix().MulVec4(transformedVertex)
				triangle.Projected[i] = transformedVertex
			}

			// Backface Culling
			//
			// 1. Find the vector between a point in the triangle and the camera origin.
			// 2. Determine the alignment of the ray and the normal
			//
			// Origin is always {0,0,0} since the camera is at in
			// view space. It shouldn't be eye/position.
			if e.cullMode == CullModeBackFace {
				origin := Vec3{0, 0, 0}
				cameraRay := origin.Sub(triangle.Projected[0].Vec3())
				visibility := triangle.Normal().Dot(cameraRay)
				if visibility < 0 {
					continue
				}
			}

			// Currently unused until we improve our lighting.
			// triangle.LightIntensity = -triangle.Normal().Dot(e.ambientLight.Direction)

			// Clip Polygons against the frustum
			clippedTriangles := e.frustum.Clip(triangle)

			// Projection
			for _, triangleToRender := range clippedTriangles {
				for i, point := range triangleToRender.Projected {
					// Multiply the original projection matrix by the vector
					projected := e.projMatrix.MulVec4(point)

					// Perspective Divide with original z value (result.w).  The result.w is
					// populated during MulVec4() because of the projection matrix 3/2==1.
					if projected.W != 0.0 {
						projected.X /= projected.W
						projected.Y /= projected.W
						projected.Z /= projected.W
					}

					// FIXME: Invert Y to deal with obj coordinates
					// system.  I'd like to get rid of this but its
					// more complex than it seems. I think it has to
					// do with the handedness rules.
					projected.Y *= -1

					// Scale into view (tiny otherwise)
					projected.X *= (float64(e.window.width) / 2.0)
					projected.Y *= (float64(e.window.height) / 2.0)

					// Translate the projected points to the
					// middle of the screen.  FIXME: If this
					// is removed, the viewport is the top
					// left only.  I understand the model
					// would be in top left, but I don't
					// understand why the viewport/frustum
					// is changed.
					projected.X += (float64(e.window.width) / 2.0)
					projected.Y += (float64(e.window.height) / 2.0)

					triangleToRender.Projected[i] = projected
				}

				mesh.trianglesToRender = append(mesh.trianglesToRender, triangleToRender)
			}
		}
	}
}

func (e *Engine) Render() {
	// Draw a nice gradient background if we have one (typically from a FFT
	// Map) or fallback to just a black background with a grid.
	if e.scene.Background() != nil {
		e.framebuffer.DrawBackground(*e.scene.Background())
	} else {
		e.framebuffer.Clear(ColorBlack)
		e.framebuffer.DrawGrid(ColorGrey)
	}

	e.framebuffer.ClearDepth()

	for _, mesh := range e.scene.Meshes {
		for _, triangle := range mesh.trianglesToRender {
			if e.renderMode == RenderModeTexture || e.renderMode == RenderModeTextureWire {
				if triangle.HasTexture() {
					e.framebuffer.DrawTexturedTriangle(triangle, mesh.Texture)
				} else {
					e.framebuffer.DrawFilledTriangle(triangle, triangle.Color)
				}

			}

			if e.renderMode == RenderModeFill || e.renderMode == RenderModeWireFill {
				e.framebuffer.DrawFilledTriangle(triangle, triangle.Color)
			}

			if e.renderMode == RenderModeWire || e.renderMode == RenderModeWireVertex || e.renderMode == RenderModeWireFill || e.renderMode == RenderModeTextureWire {
				e.framebuffer.DrawTriangle(triangle, ColorWhite)
			}

			if e.renderMode == RenderModeWireVertex {
				for _, point := range triangle.Projected {
					e.framebuffer.DrawRectangle(int(point.X-2), int(point.Y-2), 4, 4, ColorRed)
				}
			}
		}

		// Clear the slice while retaining capacity so we don't have to
		// keep allocating each frame. The number of triangles can
		// change due to frustum clipping and backface culling, but
		// keeping the capacity at the maximum seems reasonable.
		mesh.trianglesToRender = mesh.trianglesToRender[:0]
	}

	// Render ColorBuffer
	e.window.Update(e.framebuffer)
}

func (e *Engine) SetMesh(mesh Mesh) {
	e.scene.Meshes = []*Mesh{&mesh}
}

func (e *Engine) AppendMesh(mesh Mesh) {
	e.scene.Meshes = append(e.scene.Meshes, &mesh)
}

func (e *Engine) SetAutoRotation(v Vec3) {
	e.rotation = v
	e.autoRotation = true
}

// Move to the next FFT map. This is pretty hacky.
func (e *Engine) NextMap() {
	if e.currentMap < 125 {
		e.currentMap++
		mesh := e.MeshReader.ReadMesh(e.currentMap)
		e.SetMesh(mesh)
		e.Setup()
	}
}

// Move to the previous FFT map. This is pretty hacky.
func (e *Engine) PrevMap() {
	if e.currentMap > 1 {
		e.currentMap--
		mesh := e.MeshReader.ReadMesh(e.currentMap)
		e.SetMesh(mesh)
		e.Setup()
	}
}
