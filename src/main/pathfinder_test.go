package main

import "testing"

/**
 * Check a path is found from O to O on this map,:
 * O | . .
 * . | O .
 * . _ _ .
 * . . . .
 */
func Test_FindsShortestPath(t *testing.T) {
   // Given a world
   w := NewWorld()
   // With a certain map
   w.MergeData(".w..\n.w..\n.ww.\n....", Coord{0, 0})
   // And a source and destination coordinate
   s := Coord{0, 0}
   e := Coord{2, 1}

   // When we find the path between the coordinates
   path := FindPath(*w, s, e)

   // Then the calculated path is what we expect
   if pathToString(path) != "ssseeennw" {
      t.Error("Didn't get expected path, got [" + pathToString(path) + "]")
   }
}

/**
 * Check a path is found from O to ? on this map:
 * O ?
 */
func Test_FindsPathToUnknownPoint(t *testing.T) {
   // Given a world
   w := NewWorld()
   // With a certain map
   w.MergeData(".", Coord{0, 0})
   // And a source and destination coordinate
   s := Coord{0, 0}
   e := Coord{1, 0}

   // When we find the path between the coordinates
   path := FindPath(*w, s, e)

   // Then the calculated path is what we expect
   if pathToString(path) != "e" {
      t.Error("Didn't get expected path, got [" + pathToString(path) + "]")
   }
}

/**
 * Check that no path is found from O to O on this map, the function returns nil:
 * O . | O
 * . . | _
 * . . . .
 * . . . .
 */
func Test_FindsNilForNoPath(t *testing.T) {
   // Given a world
   w := NewWorld()
   // With a certain map
   w.MergeData("..w.\n..ww\n....\n....", Coord{0, 0})
   // And a source and destination coordinate
   s := Coord{0, 0}
   e := Coord{3, 0}

   // When we find the path between the coordinates
   path := FindPath(*w, s, e)

   // Then the path is what we expect
   if path != nil {
      t.Error("Didn't get nil for an impossible path, got [" + pathToString(path) + "]")
   }
}

func pathToString(moves []*Move) (rtn string) {
   for _, move := range(moves) {
      switch move.Dir {
         case North:
            rtn += "n"
         case South:
            rtn += "s"
         case East:
            rtn += "e"
         case West:
            rtn += "w"
      }
   }

   return
}
