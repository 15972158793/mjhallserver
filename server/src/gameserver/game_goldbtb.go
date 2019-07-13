package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

type Game_BTB struct {
}

type BTB_Person struct {
}

func (self *Game_BTB) OnBegin() {

}

func (self *Game_BTB) OnEnd() {

}

func (self *Game_BTB) OnInit(room *Room) {

}

func (self *Game_BTB) OnSendInfo(person *Person) {

}

func (self *Game_BTB) OnRobot(robot *lib.Robot) {

}

func (self *Game_BTB) OnMsg(msg *RoomMsg) {
}

func (self *Game_BTB) OnBye() {

}

func (self *Game_BTB) OnExit(uid int64) {

}

func (self *Game_BTB) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_BTB) OnIsBets(uid int64) bool {

}

func (self *Game_BTB) OnBalance() {

}

func (self *Game_BTB) OnTime() {

}
