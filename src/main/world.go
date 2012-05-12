package main

import (
   "strings"
   "strconv"
)

type Coord struct {
   X, Y int
}

func (c Coord) String() string {
   return "(" + strconv.Itoa(c.X) + "," + strconv.Itoa(c.Y) + ")"
}

type Box struct {
   MinX, MinY, MaxX, MaxY int
}

func (b Box) String() string {
   return "MinX: " + strconv.Itoa(b.MinX) + ", MinY: " + strconv.Itoa(b.MinY) +
         ", MaxX: " + strconv.Itoa(b.MaxX) + ", MaxY: " + strconv.Itoa(b.MaxY)
}

// A player's map, and where they are in it
type World struct {
   board map[Coord]string
   // Board dimensions
   border Box
   // Current pos of the digger, starting pos is 0,0, can get negative
   Pos *Coord
   // Current gold of the digger
   Gold int
}

func NewWorld() *World {
   return &World{board: make(map[Coord]string), Pos: &Coord{0,0}}
}

func CloneWorld(w *World) *World {
   cloneBoard := make(map[Coord]string)
   for k, v := range w.board {
      cloneBoard[k] = v
   }

   return &World{board: cloneBoard, border: w.border, Pos: &Coord{w.Pos.X,w.Pos.Y}, Gold: w.Gold}
}

type WorldIterator func(Coord, string) bool

const MAX_GOLD = 3

// Iterates the board top to bottom, left to right, calling the given iterator for each cell
func (w World) Iterate(iterator WorldIterator) {
   box := w.GetBoundingBox()
   for y := box.MinY;y <= box.MaxY;y++ {
      for x := box.MinX;x <= box.MaxX;x++ {
         c := Coord{x,y}
         cont := iterator(c, w.GetCoord(c))
         if !cont {
            return
         }
      }
   }
}

// Mutates the given world to include the new view data
func (w *World) MergeViewData (viewInfo string) {
   // Assume the view info is square
   rows := strings.Split(viewInfo, "\n")
   viewDelta := (len(rows) - 1) / 2
   w.MergeData(viewInfo, Coord{w.Pos.X - viewDelta, w.Pos.Y - viewDelta})
}

// Mutates the given world to include the new view data
func (w *World) MergeData (viewInfo string, topLeft Coord) {
   rows := strings.Split(viewInfo, "\n")
   height := len(rows)
   width := len(rows[0])

   w.border.MinX = min(w.border.MinX, topLeft.X)
   w.border.MinY = min(w.border.MinY, topLeft.Y)
   w.border.MaxX = max(w.border.MaxX, topLeft.X + width - 1)
   w.border.MaxY = max(w.border.MaxY, topLeft.Y + height - 1)

   for y, rv := range(rows) {
      for x, v := range([]byte(rv)) {
         loc := Coord{x + topLeft.X, y + topLeft.Y}
         w.board[loc] = string(v)
      }
   }
}

func (w World) GetCoord (coord Coord) string {
   return w.board[coord]
}

func (w World) GetBoundingBox () Box {
   return w.border
}

func (w World) String() string {
   rtn := "Pos: " + w.Pos.String() + ", Gold: " + strconv.Itoa(w.Gold) + "\n"
   box := w.GetBoundingBox()
   w.Iterate(func (c Coord, val string) bool {
      if c == *w.Pos {
         rtn = rtn + "o"
      } else if val == "" {
         rtn = rtn + "?"
      } else {
         rtn = rtn + val
      }

      if c.X == box.MaxX {
         rtn = rtn + "\n"
      }

      return true
   })

   return rtn
}
