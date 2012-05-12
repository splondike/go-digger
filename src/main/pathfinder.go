package main

import (
   "container/list"
   "math"
)

// Return nil when we can't find a path
// Uses the A* algorithm (taken from Wikipedia)
func FindPath(w World, src Coord, dest Coord) []*Move {
   closedset := make(map[Coord]bool)
   openset := map[Coord]bool{src: true}
   came_from := make(map[Coord]Coord)

   g_score := map[Coord]int{src: 0}
   h_score := map[Coord]int{src: costEstimate(src, dest)}
   f_score := map[Coord]int{src: g_score[src] + h_score[src]}

   for len(openset) != 0 {
      // Find the lowest f_score coord in openset
      var current Coord
      prev_lowest_score := math.MaxInt32 // int can be 32 or 64 bit, this is the safest option
      for c, _ := range(openset) {
         if f_score[c] < prev_lowest_score {
            current = c
            prev_lowest_score = f_score[c]
         }
      }
      if current == dest {
         return reconstructPath(came_from, dest)
      }

      delete(openset, current)
      closedset[current] = true
      for e := getOpenNeighbours(w, current, dest).Front(); e != nil; e = e.Next() {
         neighbour := e.Value.(Coord)
         if _, present := closedset[neighbour]; present {
            continue
         }
         tentative_g_score := g_score[current] + 1

         tentative_is_better := false
         if _, present := openset[neighbour]; !present {
            openset[neighbour] = true
            h_score[neighbour] = costEstimate(neighbour, dest)
            tentative_is_better = true
         } else if tentative_g_score < g_score[neighbour] {
            tentative_is_better = true
         }

         if tentative_is_better {
            came_from[neighbour] = current
            g_score[neighbour] = tentative_g_score
            f_score[neighbour] = g_score[neighbour] + h_score[neighbour]
         }
      }
   }

   return nil
}

func costEstimate(src Coord, dest Coord) int {
   rtn := math.Abs(float64(src.X - dest.X)) + math.Abs(float64(src.Y - dest.Y))
   return int(rtn)
}

func reconstructPath(came_from map[Coord]Coord, current_node Coord) []*Move {
   list := list.New()
   for present := true ; present ; _, present = came_from[current_node] {
      next_node := came_from[current_node]
      list.PushFront(findMove(current_node, next_node))

      current_node = next_node
   }

   // Convert to a list and return
   rtn := make([]*Move, list.Len())
   i := 0
   for e := list.Front(); e != nil; e = e.Next() {
      move := e.Value.(Move)
      rtn[i] = &move
      i++
   }
   return rtn
}

func getOpenNeighbours(w World, startingCoord Coord, targetCoord Coord) *list.List {
   l := list.New()
   possibleCoords := []Coord {
      Coord{startingCoord.X, startingCoord.Y - 1},
      Coord{startingCoord.X - 1, startingCoord.Y},
      Coord{startingCoord.X + 1, startingCoord.Y},
      Coord{startingCoord.X, startingCoord.Y + 1},
   }
   for _, c := range(possibleCoords) {
      v := w.GetCoord(c)
      if c != targetCoord && (v == "w" || v == "") {
         continue
      }

      l.PushBack(c)
   }

   return l
}

func findMove(current_node Coord, next_node Coord) (rtn Move) {
   if current_node.X < next_node.X {
      rtn = NewMove(West)
   } else if current_node.X > next_node.X {
      rtn = NewMove(East)
   } else if current_node.Y < next_node.Y {
      rtn = NewMove(North)
   } else if current_node.Y > next_node.Y {
      rtn = NewMove(South)
   }

   return
}

func FindPathLinear(w World, src Coord, dest Coord) []Move {
   // For now we'll just do a straight line pathfind, ignoring the world
   list := list.New()

   var xMove Move
   if src.X < dest.X {
      xMove = NewMove(East)
   } else {
      xMove = NewMove(West)
   }
   for i:= 0.0; i < math.Abs(float64(dest.X - src.X)); i++ {
      list.PushBack(xMove)
   }

   var yMove Move
   if src.Y < dest.Y {
      yMove = NewMove(South)
   } else {
      yMove = NewMove(North)
   }
   for i:= 0.0; i < math.Abs(float64(dest.Y - src.Y)); i++ {
      list.PushBack(yMove)
   }

   // Convert to a list and return
   rtn := make([]Move, list.Len())
   i := 0
   for e := list.Front(); e != nil; e = e.Next() {
      rtn[i] = e.Value.(Move)
      i++
   }
   return rtn
}
