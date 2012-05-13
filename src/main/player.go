package main

import "os"

type Player struct {
   CurrentWorld *World
}

func NewPlayer() *Player {
   return &Player{NewWorld()}
}

func main() {
   if len(os.Args) < 4 {
      println("Usage: host port password")
      println("e.g. localhost 8066 password [goldfirst (default),explorefirst]")
      os.Exit(1)
   }
   var pa = WebPlayerApi{os.Args[1], os.Args[2], os.Args[3]}

   var moveGenerator func(*Player, chan Action, chan bool) = GoldThenUnknownGenerator
   if len(os.Args) >= 5 {
      switch os.Args[4] {
         case "explorefirst":
            moveGenerator = UnknownThenGoldGenerator
      }
   }

   player := NewPlayer()

   actionPipe := make(chan Action, 100)
   needMovesSignal := make(chan bool)
   go moveGenerator(player, actionPipe, needMovesSignal)

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
