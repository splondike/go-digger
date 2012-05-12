package main

import (
   "strconv"
   "math"
)

var pa = WebPlayerApi{"localhost", "8066", "password"}

type Player struct {
   CurrentWorld *World
}

func (p *Player) View() {
   p.CurrentWorld.MergeViewData(pa.View())
}

func NewPlayer() *Player {
   return &Player{NewWorld()}
}

type candidateFunc func(Coord, string) bool

func findNextSquare(world *World, isCandidateFunc candidateFunc) (rtn []Action) {
   bestLen := math.MaxInt32
   var nextMoves []*Move = nil

   // TODO: Add a different iterate which does a back and forth scanner
   // TODO: Allow turning off of 'find best option', slows down revealing the map
   world.Iterate(func(c Coord, val string) bool {
      isCandidate := isCandidateFunc(c, val)

      if isCandidate {
         moves := FindPath(*world, *world.Pos, c)
         if moves != nil && len(moves) < bestLen {
            nextMoves = moves
            bestLen = len(moves)
         }
      }

      return true
   })

   if nextMoves != nil {
      // Can't just cast a []*Move to an []Action
      rtn = make([]Action, len(nextMoves))
      for i, move := range(nextMoves) {
         rtn[i] = move
      }
   }

   return
}

func findNextUnknown(w *World) []Action {
   return findNextSquare(w, func(c Coord, val string) bool {
      b := w.GetBoundingBox()
      if val == "" {
         // It's an unknown square within our bounds
         return true
      } else if val != "w" && (c.X == b.MinX || c.X == b.MaxX || c.Y == b.MinY || c.Y == b.MaxY) {
         // It's a unenclosed square at the edge of our board
         return true
      }

      return false
   })
}

func findNextBase(w *World) []Action {
   return findNextSquare(w, func(c Coord, val string) bool {
      return val == "b"
   })
}

func findNextGold(w *World) []Action {
   return findNextSquare(w, func(c Coord, val string) bool {
      _, err := strconv.Atoi(val)
      return err == nil
   })
}

func main() {
   p := NewPlayer()

   for {
      p.CurrentWorld = NewWorld()
      p.View()
      p.CurrentWorld.Gold = pa.Carrying()

      // Reveal the whole map
      for next := findNextUnknown(p.CurrentWorld);next != nil;next = findNextUnknown(p.CurrentWorld) {
         shouldView := true
         for _, move := range(next) {
            move.Do(pa, p.CurrentWorld)
            shouldView = move.Verify()
            println("move", shouldView, p.CurrentWorld.String())
         }
         if shouldView {
            p.View()
         }
      }

      // Get all the gold
      // TODO: Turn this into a go routine
      for {
         var moves []Action
         if p.CurrentWorld.Gold == MAX_GOLD {
            baseMoves := findNextBase(p.CurrentWorld)
            moves = append(baseMoves, &Drop{})
            println("drop", p.CurrentWorld.String())
         } else {
            // Find gold
            goldMoves := findNextGold(p.CurrentWorld)

            // If gold, get it
            if goldMoves != nil {
               moves = append(goldMoves, &Grab{})
               println("get", p.CurrentWorld.String())
            } else if p.CurrentWorld.Gold == 0 {
               pa.Next()
               println("next", p.CurrentWorld.String())
               break
            } else {
               baseMoves := findNextBase(p.CurrentWorld)
               println("drop2", p.CurrentWorld.String())
               moves = append(baseMoves, &Drop{})
            }
         }

         for _, move := range(moves) {
            move.Do(pa, p.CurrentWorld)
            println("move", p.CurrentWorld.String())
         }
      }
   }
}
