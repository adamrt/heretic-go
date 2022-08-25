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

##### Notes:

- NDC == Normalized Device Coordinates, AKA Image Space
- The value in the perspective projection matrix
  (MatrixMakePersProj()) at m[3][2] is 1.0. During multiplication this
  will store the original z value of the Vec4 (z*1.0==z). Then we can
  use the z value later to handle perspective divide in
  m.MulVec4Proj().
- The DOT product of the face normal and the light direction gives a
  float representing alignment. Then that float between 0-1 can be
  multiplied by the original color.
- The normal must be normalized for this to work.
