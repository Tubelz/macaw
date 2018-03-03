package system

import (
	"log"
	"github.com/tubelz/macaw/entity"
	"github.com/tubelz/macaw/math"
	"github.com/veandco/go-sdl2/sdl"
)

// CollisionSystem is the system responsible to handle collisions
type CollisionSystem struct {
	Entities []entity.Entitier
	Name string
	Subject
}

// Assign assign entities with this system
func (c *CollisionSystem) Assign(entities []entity.Entitier) {
	c.Entities = entities
}

// Update check for collision and notify observers
func (c *CollisionSystem) Update() {
	var component interface{}
	var ok bool

	for _, obj := range c.Entities {
		if component, ok = obj.GetComponent("position"); !ok {
			continue
		}
		position := component.(*entity.PositionComponent)
		if component, ok = obj.GetComponent("collision"); !ok {
			continue
		}
		collision := component.(*entity.CollisionComponent)

		for _, obj2 := range c.Entities {
			if obj == obj2 {
				continue
			}

			if component, ok = obj2.GetComponent("position"); !ok {
				continue
			}
			position2 := component.(*entity.PositionComponent)
			if component, ok = obj2.GetComponent("collision"); !ok {
				continue
			}

			collision2 := component.(*entity.CollisionComponent)
			rect1 := &sdl.Rect{position.Pos.X, position.Pos.Y, collision.Size.X, collision.Size.Y}
			rect2 := &sdl.Rect{position2.Pos.X, position2.Pos.Y, collision2.Size.X, collision2.Size.Y}
			if (rect1.HasIntersection(rect2)) {
				c.NotifyEvent(&CollisionEvent{Ent: obj.(*entity.Entity), With: obj2.(*entity.Entity)})
			}
		}
	}
}

// CollisionEvent has the entity (Ent) that produced the collision and the entity that got collided (With)
type CollisionEvent struct {
	Ent *entity.Entity
	With *entity.Entity
}

// Name returns the collision event name
func (c *CollisionEvent) Name() string {
	return "collision event"
}

/*
	----
	Util functions for handling collision events
	----
*/

// InvertVel invert the vel of the collided object.
func InvertVel(event Event) {
	collision := event.(*CollisionEvent)
	log.Printf("Inverting pos and mov of obj %d", collision.Ent.GetID())

	component, ok := collision.Ent.GetComponent("position")
	if !ok {
		return
	}
	position := component.(*entity.PositionComponent)

	component, ok = collision.Ent.GetComponent("physics")
	if !ok {
		return
	}
	physics := component.(*entity.PhysicsComponent)

	intersectRect := intersection(collision.Ent, collision.With)
	displacementPos := &sdl.Point{intersectRect.W, intersectRect.H}

	// TODO: Clean this a little bit...
	if displacementPos.X < displacementPos.Y {
		physics.Vel.X *= -1
		physics.Acc.X *= -1
		if physics.Vel.X > 0 {
			position.Pos.X = position.Pos.X + displacementPos.X
		} else if physics.Vel.X < 0 {
			position.Pos.X = position.Pos.X - displacementPos.X
		}
		physics.FuturePos.X = float32(position.Pos.X) + physics.Vel.X
	} else if displacementPos.Y < displacementPos.X {
		physics.Vel.Y *= -1
		physics.Acc.Y *= -1
		if physics.Vel.Y > 0 {
			position.Pos.Y = position.Pos.Y + displacementPos.Y
		} else if physics.Vel.Y < 0 {
			position.Pos.Y = position.Pos.Y - displacementPos.Y
		}
		physics.FuturePos.Y = float32(position.Pos.Y) + physics.Vel.Y
	} else {
		physics.Vel = math.MulFPointWithFloat(physics.Vel, -1)
		physics.Acc = math.MulFPointWithFloat(physics.Acc, -1)
		if physics.Vel.X > 0 {
			position.Pos.X = position.Pos.X + displacementPos.X
		} else if physics.Vel.X < 0{
			position.Pos.X = position.Pos.X - displacementPos.X
		}
		if physics.Vel.Y > 0 {
			position.Pos.Y = position.Pos.Y + displacementPos.Y
		} else if physics.Vel.Y < 0 {
			position.Pos.Y = position.Pos.Y - displacementPos.Y
		}
		physics.FuturePos = math.ConvertPointToFPoint(math.SumPointWithFPoint(position.Pos, physics.Vel))
	}
}

// intersection get the intersection rectangle between two objects
func intersection(obj1, obj2 *entity.Entity) sdl.Rect {
	c, _ := obj1.GetComponent("position")
	position := c.(*entity.PositionComponent)
	c, _ = obj2.GetComponent("position")
	position2 := c.(*entity.PositionComponent)
	c, _ = obj1.GetComponent("collision")
	collision := c.(*entity.CollisionComponent)
	c, _ = obj2.GetComponent("collision")
	collision2 := c.(*entity.CollisionComponent)

	rect1 := &sdl.Rect{position.Pos.X, position.Pos.Y, collision.Size.X, collision.Size.Y}
	rect2 := &sdl.Rect{position2.Pos.X, position2.Pos.Y, collision2.Size.X, collision2.Size.Y}

	displacement, _ := rect1.Intersect(rect2)
	return displacement
}
