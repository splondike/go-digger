package main

import (
   "math"
   "strconv"
)

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

/*************************************
 * Different move finding algorithms *
 *************************************/

// Spew out moves to the actionPipe after signaled by needMovesSignal until we need to wait for api input
// then block and wait for needMovesSignal to trigger again

type moveGeneratorFunc func (*World) ([]Action, nextAction)

// Stuff actions into actionPipe
func actionsToPipe (actions []Action, actionPipe chan Action) {
   for _, action := range(actions) {
      actionPipe <- action
   }
}

// Update a world based on the actions without sending requests to the server
func applyToWorld (world *World, actions []Action) {
   for _, action := range(actions) {
      action.Do(DummyPlayerApi{}, world)
   }
}

// The possible actions the generator can pass to the generator template
type nextAction int
const (
   CONTINUE nextAction = iota
   WAIT_INPUT
   NEXT_LEVEL
)

func moveGeneratorTemplate(player *Player, actionPipe chan Action, needMovesSignal chan bool, generator moveGeneratorFunc) {
   for {
      // Wait for the main routine to signal it needs more moves generated
      wait: <-needMovesSignal

      internalWorld := CloneWorld(player.CurrentWorld)
      for {
         actions, nextAction := generator(internalWorld)
         actionsToPipe(actions, actionPipe)
         switch nextAction {
            case CONTINUE:
               applyToWorld(internalWorld, actions)
            case WAIT_INPUT:
               actionPipe <- &signalMoreMoves{needMovesSignal}
               goto wait
            case NEXT_LEVEL:
               actionPipe <- &initializeNewWorld{player}
               actionPipe <- &signalMoreMoves{needMovesSignal}
               goto wait
         }
      }
   }
}

// A move generator which prioritizes getting gold over exploring
func GoldThenUnknownGenerator(player *Player, actionPipe chan Action, needMovesSignal chan bool) {
   moveGeneratorTemplate(player, actionPipe, needMovesSignal, goldThenUnknownGenerator)
}
func goldThenUnknownGenerator(world *World) ([]Action, nextAction) {
   println(world.String())
   fullOfGold := world.Gold == MAX_GOLD
   movesToGold := findNextGold(world)
   movesToUnknown := findNextUnknown(world)
   baseMoves := findNextBase(world)
   haveLastGold := world.Gold > 0 && movesToGold == nil

   // Dumping gold is first priority if we can find a base
   if baseMoves != nil && (fullOfGold || (haveLastGold && movesToUnknown == nil)) {
      println("dumping gold")
      actions := append(baseMoves, &Drop{})
      return actions, CONTINUE
   }

   // Next is finding more gold
   if !fullOfGold && movesToGold != nil {
      println("finding gold")
      actions := append(movesToGold, &Grab{})
      return actions, CONTINUE
   }

   // Last is looking for unexplored territory
   if movesToUnknown != nil {
      println("exploring")
      actions := append(movesToUnknown, &View{})
      return actions, WAIT_INPUT
   }

   // Failing all that, it's time for a new level
   println("next level")
   return []Action{&Next{}}, NEXT_LEVEL
}

// A move generator which prioritizes getting gold over exploring
func UnknownThenGoldGenerator(player *Player, actionPipe chan Action, needMovesSignal chan bool) {
   moveGeneratorTemplate(player, actionPipe, needMovesSignal, unknownThenGoldGenerator)
}
func unknownThenGoldGenerator(world *World) ([]Action, nextAction) {
   println(world.String())
   movesToUnknown := findNextUnknown(world)

   // First priority is looking for unexplored territory
   if movesToUnknown != nil {
      println("exploring")
      actions := append(movesToUnknown, &View{})
      return actions, WAIT_INPUT
   }

   fullOfGold := world.Gold == MAX_GOLD
   movesToGold := findNextGold(world)
   baseMoves := findNextBase(world)
   haveLastGold := world.Gold > 0 && movesToGold == nil
   // Dumping gold is next if we can find a base
   if baseMoves != nil && (fullOfGold || (haveLastGold && movesToUnknown == nil)) {
      println("dumping gold")
      actions := append(baseMoves, &Drop{})
      return actions, CONTINUE
   }

   // Last is finding more gold
   if !fullOfGold && movesToGold != nil {
      println("finding gold")
      actions := append(movesToGold, &Grab{})
      return actions, CONTINUE
   }

   // Failing all that, it's time for a new level
   println("next level")
   return []Action{&Next{}}, NEXT_LEVEL
}
