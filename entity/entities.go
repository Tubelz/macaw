package entity

import (
	"github.com/tubelz/macaw/math"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Component is the abstract type for each component
type Component interface{}

// Entitier has all the behaviours entities should have
type Entitier interface {
	GetID() uint16
	GetComponents() map[string]Component
	AddComponent(string, Component)
	DelComponent(key string)
	GetComponent(componentName string) (Component, bool)
}

// Entity is the struct that contains the components. Right now the id's are not being used
type Entity struct {
	id         uint16
	etype      string // type of the entity
	components map[string]Component
}

// GetID returns the id of the entity
func (e *Entity) GetID() uint16 {
	return e.id
}

// GetType returns the type of the entity
func (e *Entity) GetType() string {
	return e.etype
}

// GetComponents returns a list of all the components of the entity
func (e *Entity) GetComponents() map[string]Component {
	return e.components
}

// GetComponent returns the given component
func (e *Entity) GetComponent(componentName string) (Component, bool) {
	val, ok := e.components[componentName]
	return val, ok
}

// AddComponent adds a component to the component map
func (e *Entity) AddComponent(name string, c Component) {
	e.components[name] = c
}

// DelComponent removes the given component
func (e *Entity) DelComponent(key string) {
	delete(e.components, key)
}

// Manager is the struct responsible to manage the entities in your game
type Manager struct {
	// count entities allocated
	counter        uint16
	entities       []*Entity
	availableSlots []uint16
}

// Create creates a new entity and returns it
func (m *Manager) Create(etype string) *Entity {
	var i uint16
	entity := new(Entity)
	entity.etype = etype

	// check if we can use an empty slot of our array, or if we have to add a new position
	if len(m.availableSlots) > 0 {
		// pop the first id that was deleted if there is any deleted element. we use FIFO here
		i = m.availableSlots[0]
		m.availableSlots = append(m.availableSlots[:0], m.availableSlots[1:]...)
		entity.id = i
		m.entities[i] = entity
	} else {
		entity.id = m.counter
		m.counter++
		m.entities = append(m.entities, entity)
	}
	entity.components = make(map[string]Component)

	return entity
}

// binarySearch returns the index where either the index of the element that we found,
// or where we should insert the new element that doesn't exist
func binarySearch(arr []uint16, low int, high int, val uint16) int {
	var mid int
	for low <= high {
		mid = (low + high) >> 1
		if arr[mid] <= val {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	if arr[mid] < val {
		return mid + 1
	}
	return mid
}

// Delete removes an entity associated to the given id
func (m *Manager) Delete(id uint16) bool {
	if m.entities[id] == nil {
		return false
	}
	var entity *Entity
	entity = nil
	m.entities[id] = entity
	// insert element at i
	arrSize := len(m.availableSlots) - 1
	if arrSize < 1 {
		m.availableSlots = append(m.availableSlots, id)
	} else {
		i := binarySearch(m.availableSlots, 0, arrSize, id)
		// first we we make sure the array has enough capacity
		m.availableSlots = append(m.availableSlots, 0)
		// shift existing elements (thus, overwrite our 0 that we just inserted)
		copy(m.availableSlots[i+1:], m.availableSlots[i:])
		// set our element to the proper index
		m.availableSlots[i] = id
	}
	return true
}

// Get gets an entity from the array of entities given an id
func (m *Manager) Get(id uint16) *Entity {
	if id < m.counter {
		return m.entities[id]
	}
	var entity *Entity
	return entity
}

// GetAll gets all entities
func (m *Manager) GetAll() []*Entity {
	return m.entities
}

// IterAvailable creates an iterator for the available entities
func (m *Manager) IterAvailable() func() (*Entity, bool) {
	i := -1
	entitySize := len(m.entities)
	return func() (*Entity, bool) {
		for i++; i < entitySize; i++ {
			if m.entities[i] != nil {
				return m.entities[i], true
			}
		}
		return nil, false
	}
}

/////////////////////////////////////////////////
/// ... Basic Components ...
/////////////////////////////////////////////////

//PositionComponent is responsible for the position of the entity
type PositionComponent struct {
	Pos *sdl.Point
}

// PhysicsComponent is responsible for some of the physics
type PhysicsComponent struct {
	FuturePos *math.FPoint // TODO: move this to PositionComponent
	// velocity
	Vel *math.FPoint
	// acceleration
	Acc *math.FPoint
}

// RenderComponent is responsible for the rendering of the entity
type RenderComponent struct {
	Renderer *sdl.Renderer
	Texture  *sdl.Texture
	Crop     *sdl.Rect // part of the texture which will be displayed
	Scroll   bool
	Angle    float64
	Center   *sdl.Point
	Flip     sdl.RendererFlip
}

// CameraComponent is responsible to render only the content of the viewport
type CameraComponent struct {
	ViewportSize sdl.Point
	WorldSize    sdl.Point
	IsActive     bool
}

// AnimationComponent is responsible for animate the entity
type AnimationComponent struct {
	InitialPos     sdl.Point // frame reference
	AnimationSpeed uint8     //animations per second
	PreviousTime   uint32    //last animation time
	Current        int
	Frames         int // total sprites
	RowLength      int // number of sprites per row
	SpriteMap      map[string]int
}

// FontComponent holds the font and text information
type FontComponent struct {
	Font     *ttf.Font
	Text     string
	Modified bool
	Color    *sdl.Color
}

// RectangleComponent has the information to draw a rectangle
type RectangleComponent struct {
	Size   *sdl.Point
	Color  *sdl.Color
	Filled bool
}

// CollisionComponent makes the entity notify if it hits something else
// TODO: Add other type of information such as Shape, Density, Friction etc...
type CollisionComponent struct {
	// CollisionAreas contains the rectangles that will be checked.
	// The position is relative to the upper left corner of the renderer
	CollisionAreas []sdl.Rect
}

// GridComponent is used for debugging
type GridComponent struct {
	Size  *sdl.Point
	Color *sdl.Color
}
