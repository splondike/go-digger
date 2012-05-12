package main

import (
   "net/http"
   "strconv"
   "strings"
   "io/ioutil"
)

type PlayerApi interface {
   View() string;
   Drop() int;
   Grab() int;
   Carrying() int;
   Next() bool;
   North() bool;
   South() bool;
   East() bool;
   West() bool;
}

type WebPlayerApi struct {
   host string
   port string
   // The password or secret name of the player, used to identify them
   password string
}

func (p WebPlayerApi) getUrl() string {
   return "http://" + p.host + ":" + p.port + "/golddigger/digger/" + p.password + "/"
}

func (p WebPlayerApi) strCmd(path string) string {
   resp,_ := http.Get(p.getUrl() + path)
   rtn,_ := ioutil.ReadAll(resp.Body)
   return strings.TrimRight(string(rtn), "\n")
}

func (p WebPlayerApi) intCmd(path string) int {
   rtnInt, _ := strconv.Atoi(p.strCmd(path))
   return rtnInt
}

func (p WebPlayerApi) boolCmd(path string) (success bool) {
   rtn := p.strCmd(path)

   if rtn == "OK" {
      success = true
   } else {
      success = false
   }

   return
}

func (p WebPlayerApi) View() string {
   return p.strCmd("view")
}

func (p WebPlayerApi) Drop() int {
   return p.intCmd("drop")
}

func (p WebPlayerApi) Grab() int {
   return p.intCmd("grab")
}

func (p WebPlayerApi) Carrying() int {
   return p.intCmd("carrying")
}

func (p WebPlayerApi) Next() bool {
   return p.boolCmd("next")
}

func (p WebPlayerApi) North() bool {
   return p.boolCmd("move/north")
}

func (p WebPlayerApi) South() bool {
   return p.boolCmd("move/south")
}

func (p WebPlayerApi) East() bool {
   return p.boolCmd("move/east")
}

func (p WebPlayerApi) West() bool {
   return p.boolCmd("move/west")
}

type DummyPlayerApi bool
func (p DummyPlayerApi) View() string {return ""}
func (p DummyPlayerApi) Drop() int {return 0}
func (p DummyPlayerApi) Grab() int {return 0}
func (p DummyPlayerApi) Carrying() int {return 0}
func (p DummyPlayerApi) Next() bool {return true}
func (p DummyPlayerApi) North() bool {return true}
func (p DummyPlayerApi) South() bool {return true}
func (p DummyPlayerApi) East() bool {return true}
func (p DummyPlayerApi) West() bool {return true}
