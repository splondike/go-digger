Intro
=====

This program plays the GoldDigger programming game found here: https://github.com/ellnestam/GoldDigger

I wrote it to learn something about Google's Go language (http://golang.org).

Probably won't do anything else with it now that I've done that.

Usage
=====
Compile using Google Go:

    cd src/main
    go build -o digger *.go (or go run *.go ..args.. for dev)

Now run your digger executable and it should explain how to use it.

Code
==============

The code isn't especially well documented, but there isn't much going on so I'll explain things here from the perspective of the program just having been run:

1. player.go creates an api.go by connecting to the server.
2. player.go spawns a new moveFinder instance using the appropriate algorithm.
3. moveFinder uses pathfinder and its own logic to work out what action.go(s) it wants to do and pushes them onto player.go's actionPipe. When it needs to wait for more data from the server it pushes a special action to the pipe which tells it to unstall.
4. Meanwhile player.go is loading action.go objects off the move channel and executing them as fast as the server will let it. Some map directly to server commands, some are special (like the unstall action mentioned previously).

intmath.go provides functions from stdlib's math except for the int data type.

world.go is an object representing a map and related data (e.g. gold count).
