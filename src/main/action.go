package main

import "strconv"

type Direction int

const (
   North Direction = iota
   South
   East
   West
)

type Action interface {
   Do(PlayerApi, *World)
   Verify() bool
   String() string
}

type Move struct {
   Dir Direction
   rtn bool
}

type Drop struct{}

type Grab struct{}

type Next struct {
   rtn bool
}

type View struct{}

func NewMove(dir Direction) Move {
   return Move{dir,false}
}

func (m *Move) Do(api PlayerApi, world *World) {
   switch m.Dir {
      case North:
         m.rtn = api.North()
         if m.Verify() {
            world.Pos.Y--
         } else {
            world.MergeData("w", Coord{world.Pos.X, world.Pos.Y - 1})
         }
      case South:
         m.rtn = api.South()
         if m.Verify() {
            world.Pos.Y++
         } else {
            world.MergeData("w", Coord{world.Pos.X, world.Pos.Y + 1})
         }
      case East:
         m.rtn = api.East()
         if m.Verify() {
            world.Pos.X++
         } else {
            world.MergeData("w", Coord{world.Pos.X + 1, world.Pos.Y})
         }
      case West:
         m.rtn = api.West()
         if m.Verify() {
            world.Pos.X--
         } else {
            world.MergeData("w", Coord{world.Pos.X - 1, world.Pos.Y})
         }
      default:
         m.rtn = false
   }
}

func (m *Move) Verify() bool {
   return m.rtn
}

func (m *Move) String() string {
   switch m.Dir {
      case North:
         return "Move[North]"
      case South:
         return "Move[South]"
      case East:
         return "Move[East]"
      case West:
         return "Move[West]"
   }

   return "Move[Unknown]"
}

func (m *Drop) Do(api PlayerApi, world *World) {
   // Update the gold count
   ground := world.GetCoord(*world.Pos)
   groundGold := 0
   if gold, err := strconv.Atoi(ground); err == nil {
      groundGold = gold
   }
   droppedGold := min(world.Gold, 9 - groundGold)
   world.Gold -= droppedGold

   // Notify the server
   api.Drop()

   // Update the board
   if ground == "b" {
      return
   }
   remaining := groundGold + droppedGold
   remainingStr := "."
   if remaining != 0 {
      remainingStr = strconv.Itoa(remaining)
   }
   world.MergeData(remainingStr, *world.Pos)
}

func (m *Drop) Verify() bool {
   return true
}

func (m *Drop) String() string {
   return "Drop"
}

func (m *Grab) Do(api PlayerApi, world *World) {
   ground := world.GetCoord(*world.Pos)
   if ground == "b" {
      return
   }

   // Update the gold count
   groundGold := 0
   if gold, err := strconv.Atoi(ground); err == nil {
      groundGold = gold
   }
   grabbedGold := min(3 - world.Gold, groundGold)
   world.Gold += grabbedGold

   // Notify the server
   api.Grab()

   // Update the board
   remaining := groundGold - grabbedGold
   remainingStr := "."
   if remaining != 0 {
      remainingStr = strconv.Itoa(remaining)
   }
   world.MergeData(remainingStr, *world.Pos)
}

func (m *Grab) Verify() bool {
   return true
}

func (m *Grab) String() string {
   return "Grab"
}

func (m *Next) Do(api PlayerApi, world *World) {
   m.rtn = api.Next()
}

func (m *Next) Verify() bool {
   return m.rtn
}

func (m *Next) String() string {
   return "Next"
}

func (m *View) Do(api PlayerApi, world *World) {
   world.MergeViewData(api.View())
}

func (m *View) Verify() bool {
   return true
}

func (m *View) String() string {
   return "View"
}
