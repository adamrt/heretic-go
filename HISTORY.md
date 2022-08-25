# A history of Heretic

### Aug 22 2022

- [x] Create an SDL Window with color
- [x] Add color buffer
- [x] Draw pixel
- [x] Draw grid
- [x] Draw rectangle
- [x] Draw flat cube
- [x] Draw perspective corrected cube
- [x] Draw rotating cube (multiaxis)

Everything was good and understandable today. There is no matrix math
yet, so we are rotating with simple functions so far.

##### Trig reminder:

```
            /|
Hypotenous / |
          /  |
         /   | Opposite
        /    |
       /     |
      /______|
      Adjacent
```

- Sin(a) = Opposite/Hypotenuse (**s=o/h**)
- Cos(a) = Adjacent/Hypotenuse (**c=a/h**)
- Tan(a) = Opposite/Adjacent (**t=o/a**)

### Aug 23 2022

- [x] Add simple timestep. Requires improvement
- [x] Add Triangle

### Aug 24 2022

- [x] Draw Line
- [x] Draw Triangle
- [x] Draw Cube from Triangles
- [x] Add Mesh
- [x] Load Mesh from OBJ file
- [x] Vector Math
- [x] Add Backface culling
- [x] Draw Filled Triangle
- [x] Add controls to manage culling and rendering modes
- [x] Add Face (Vec3 + Color) to represent pre-projected triangle
- [x] Add painters algorithm
- [x] Matrix Math
- [x] Add Scale Matrix
- [x] Add Translation Matrix
- [x] Add Rotation Maticies
- [x] Add World Matrix

##### Notes:

- Normals are for lighting, but also for backface culling.
- Identity matrix in matrix multiplication acts as a one real number multiplication
- World matrix is a combined matrix of scale, rotation, translation
- The order of matrix multiplation matters. Scale then rotation then translation

### Aug 25 2022

- [x] Add perspective projection matrix
- [x] Add flat light shading
- [x] Add texture mapping
- [x] Add texture perspective correction
- [x] Add png textures
- [x] Add OBJ vt parsing
- [x] Render arbitrarily complex textured models
- [x] Add View Matrix and Camera

##### Notes:

After adding textures I realized that lighting wasnt working
anymore. This is because we calculate light intensity on the Triangle
Color, but then if we are using textures, we use DrawTexel that uses
the image and not the modified Triangle.color. I was able to add the
lighting back in to textures by passing the lightIntensity along with
the triangle struct and then pass it to DrawTexturedTriangle() and in
turn DrawTexel(), which would get the texel color, then apply the
light intensity. I removed this for now as my lighting gets more
complex, I will probably have to refactor anyway.

- NDC == Normalized Device Coordinates, AKA Image Space
- The value in the perspective projection matrix
  (MatrixMakePersProj()) at m[3][2] is 1.0. During multiplication this
  will store the original z value of the Vec4 (z*1.0==z). Then we can
  use the z value later to handle perspective divide in
  m.MulVec4Proj().
- The DOT product of the face normal and the light direction gives a
  float representing alignment. Then that float between 0-1 can be
  multiplied by the original color. The normal must be normalized for
  this to work.
- Modelspace -> WorldSpace -> View/Camera Space -> Screen Space
  verts       * worldMatrix * viewMatrix         * projMatrix
