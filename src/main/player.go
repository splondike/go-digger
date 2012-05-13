package main

import "os"

type Player struct {
   CurrentWorld *World
}

func NewPlayer() *Player {
   return &Player{NewWorld()}
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
   go GoldThenUnknownGenerator(player, actionPipe, needMovesSignal)

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
