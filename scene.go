package heretic

// scene is currently just a way for the engine to render multiple meshes.
type scene struct {
	Meshes []*Mesh
}

func NewScene() *scene {
	return &scene{
		Meshes: make([]*Mesh, 0),
	}
}

func (s scene) Background() *Background {
	for _, mesh := range s.Meshes {
		if mesh.Background != nil {
			return mesh.Background
		}
	}
	return nil
}
