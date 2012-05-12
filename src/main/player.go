package main

import (
   "strconv"
   "math"
   "os"
)

type Player struct {
   CurrentWorld *World
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

// Extra sneaky actions for the use of generateMoves
type signalMoreMoves struct {
   moreMovesSignal chan bool
}

func (m *signalMoreMoves) Do(api PlayerApi, world *World) {
   m.moreMovesSignal <- true
}

func (m *signalMoreMoves) Verify() bool {
   return true
}

func (m *signalMoreMoves) String() string {
   return "signalMoreMoves"
}

type initializeNewWorld struct {
   player *Player
}

func (m *initializeNewWorld) Do(api PlayerApi, world *World) {
   m.player.CurrentWorld = NewWorld()
   v := View{}
   v.Do(api, m.player.CurrentWorld)
   m.player.CurrentWorld.Gold = api.Carrying()
}

func (m *initializeNewWorld) Verify() bool {
   return true
}

func (m *initializeNewWorld) String() string {
   return "initializeNewWorld"
}

// Spew out moves to the actionPipe after signaled by needMovesSignal until we need to wait for api input
// then block and wait for needMovesSignal to trigger again
func generateMoves(player *Player, actionPipe chan Action, needMovesSignal chan bool) {
   // Stuff actions into actionPipe
   actionsToPipe := func (actions []Action) {
      for _, action := range(actions) {
         actionPipe <- action
      }
   }

   // Update a world based on the actions without sending requests to the server
   applyToDummyWorld := func (world *World, actions []Action) {
      for _, action := range(actions) {
         action.Do(DummyPlayerApi{}, world)
      }
   }

   for {
      // Wait for the main routine to signal it needs more moves generated
      <-needMovesSignal

      dummyWorld := CloneWorld(player.CurrentWorld)
      for {
         println(dummyWorld.String())
         fullOfGold := dummyWorld.Gold == MAX_GOLD
         movesToGold := findNextGold(dummyWorld)
         movesToUnknown := findNextUnknown(dummyWorld)
         baseMoves := findNextBase(dummyWorld)
         haveLastGold := dummyWorld.Gold > 0 && movesToGold == nil

         // Dumping gold is first priority if we can find a base
         if baseMoves != nil && (fullOfGold || (haveLastGold && movesToUnknown == nil)) {
            println("dumping gold")
            actions := append(baseMoves, &Drop{})

            actionsToPipe(actions)
            applyToDummyWorld(dummyWorld, actions)
            continue
         }

         // Next is finding more gold
         if !fullOfGold && movesToGold != nil {
            println("finding gold")
            actions := append(movesToGold, &Grab{})

            actionsToPipe(actions)
            applyToDummyWorld(dummyWorld, actions)
            continue
         }

         // Last is looking for unexplored territory
         if movesToUnknown != nil {
            println("exploring")
            actionsToPipe(append(movesToUnknown, &View{}))

            // We need to wait for the signal that more info's come in
            actionPipe <- &signalMoreMoves{needMovesSignal}
            break
         }

         // Failing all that, it's time for a new level
         println("next level")
         actionPipe <- &Next{}
         actionPipe <- &initializeNewWorld{player}
         // We need to wait for the signal that more info's come in
         actionPipe <- &signalMoreMoves{needMovesSignal}
         break
      }
   }
}

func main() {
   if len(os.Args) != 4 {
      println("Usage: host port password")
      println("e.g. localhost 8066 password")
      os.Exit(1)
   }
   var pa = WebPlayerApi{os.Args[1], os.Args[2], os.Args[3]}

   player := NewPlayer()

   actionPipe := make(chan Action, 100)
   needMovesSignal := make(chan bool)
   go generateMoves(player, actionPipe, needMovesSignal)

   // Set up the world and start the moves coming
   actionPipe <- &initializeNewWorld{player}
   actionPipe <- &signalMoreMoves{needMovesSignal}

   // Play out the generated moves to the server
   for {
      move := <-actionPipe
      move.Do(pa, player.CurrentWorld)
      //println(move.String(), player.CurrentWorld.String())
   }
}
