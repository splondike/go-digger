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

type Drop struct {}

type Grab struct {}

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
   droppedGold := api.Drop()
   world.Gold -= droppedGold

   ground := world.GetCoord(*world.Pos)
   if ground == "b" {
      return
   }

   groundGold := 0
   if gold, err := strconv.Atoi(ground); err != nil {
      groundGold = gold
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
   grabbedGold := api.Grab()
   world.Gold += grabbedGold

   groundGold, _ := strconv.Atoi(world.GetCoord(*world.Pos))
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
