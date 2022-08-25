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

Normals are for lighting, but also for backface culling.

Line equation `y = mx + c`.

- `m` is the slope.
- `c` is the y-intercept.

`rise`/`run`

Identity matrix in matrix multiplication acts as a one real number multiplication
