package gameserver

import (
	//"log"
	//"sort"
	"staticfunc"
	//"strings"
	"fmt"
	"lib"
	"math"
	"time"
)

var LKPY_TYPE []int = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26}
var ACTIVE_LKPY float64 = 60 //! 多长时间不活跃后退出房间

type Rec_Fishing_Info_LKPY struct {
	GameType int                           `json:"gametype"`
	Time     int64                         `json:"time"`
	Info     []Son_Rec_Fishing_Person_LKPY `json:"info"`
}
type Son_Rec_Fishing_Person_LKPY struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Score int    `json:"score"`
}

type Game_Fishing_LKPY struct {
	MzSpace   float64                      `json:"mzspace"`   //! 命中波动间隔
	MzPersist float64                      `json:"mzpersist"` //! 命中波动持续时长
	Wave      float64                      `json:"wave"`      //! 波动值
	MzUp      bool                         `json:"mzup"`      //! 命中上升还是下降
	PersonMgr [4]*Game_Fishing_LKPY_Person `json:"personmgr"`
	Key       int                          `json:"key"`      //! 鱼的动态标识列
	FishPool  []*Game_Fish_LKPY            `json:"fishpool"` //! 鱼的对象池
	CurFish   []*Game_Fish_LKPY            `json:"curfish"`  //! 出现在屏幕中的鱼
	//!------------------路径--------------------------
	CheShi     int                       `json:"cheshi"`
	PathMap    map[int][]*Game_Fish_LKPY `json:pathmap`      //! 路径map
	FishMap    map[int]int64             `json:"fishmap"`    //! 鱼map
	AllPath    []Game_Path_LKPY          `json:"allpath"`    //!　所有路径
	DjFishTime float64                   `json:"djfishtime"` //! 道具鱼刷新时间

	TopPaths   [][]int `json:"toppaths"` //! 上平行路径
	TopPath    []int   `json:"toppath"`
	DownPaths  [][]int `json:"downpaths"`
	DownPath   []int   `json:"downpath"`
	RightPaths [][]int `json:"rightpaths"`
	RightPath  []int   `json:"rightpath"`
	LeftPaths  [][]int `json:"leftpaths"`
	LeftPath   []int   `json:"leftPath"`
	BossPath   []int   `json:"bosspath"` //! BOOS的路径
	NextPath   int     `json:"nextpath"` //! 下一个出鱼的方向

	room *Room
}

func NewGame_Fishing_LKPY() *Game_Fishing_LKPY {
	game := new(Game_Fishing_LKPY)
	for i := 0; i < len(game.PersonMgr); i++ {
		game.PersonMgr[i] = nil
	}

	game.Key = 1

	game.NextPath = 1

	//!----- 处理路径
	game.PathMap = make(map[int][]*Game_Fish_LKPY, 0)
	pathId := make(map[int]int)
	for _, value := range staticfunc.GetCsvMgr().Data["lkpypath"] {

		var path Game_Path_LKPY
		path.Id = lib.HF_Atoi(value["id"])
		game.PathMap[path.Id] = make([]*Game_Fish_LKPY, 0)
		path.Path = lib.HF_Atoi(value["path"])
		path.Time = lib.HF_Atof64(value["time"])
		path.CurTime = 0
		game.AllPath = append(game.AllPath, path)
		pathId[path.Path/100]++
	}

	for key, value := range pathId {
		if value == 1 {
			for i := 0; i < len(game.AllPath); i++ {
				if key == game.AllPath[i].Path/100 {
					if game.AllPath[i].Path/100000 == 1 {
						game.LeftPath = append(game.LeftPath, game.AllPath[i].Path)
					} else if game.AllPath[i].Path/100000 == 2 {
						game.RightPath = append(game.RightPath, game.AllPath[i].Path)
					} else if game.AllPath[i].Path/100000 == 3 {
						game.TopPath = append(game.TopPath, game.AllPath[i].Path)
					} else if game.AllPath[i].Path/100000 == 4 {
						game.DownPath = append(game.DownPath, game.AllPath[i].Path)
					} else if game.AllPath[i].Path/100000 == 5 {
						game.BossPath = append(game.BossPath, game.AllPath[i].Path)
					}
				}
			}
		}
		if value > 1 { //!----添加到平行路径
			paths := make([]int, 0)
			for i := 0; i < len(game.AllPath); i++ {
				if key == game.AllPath[i].Path/100 {
					paths = append(paths, game.AllPath[i].Path)
				}
			}
			if len(paths) > 1 {
				if paths[0]/100000 == 1 {
					game.LeftPaths = append(game.LeftPaths, paths)
				} else if paths[0]/100000 == 2 {
					game.RightPaths = append(game.RightPaths, paths)
				} else if paths[0]/100000 == 3 {
					game.TopPaths = append(game.TopPaths, paths)
				} else if paths[0]/100000 == 4 {
					game.DownPaths = append(game.DownPaths, paths)
				}
			}
		}
	}

	//lib.GetLogMgr().Output(lib.LOG_DEBUG, " len(LeftPath) : ", len(game.LeftPath))
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, " len(RightPath) : ", len(game.RightPath))
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, " len(TopPath) : ", len(game.TopPath))
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, " len(DownPath) : ", len(game.DownPath))
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, " len(LeftPaths) : ", len(game.LeftPaths))
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, " len(RightPaths) : ", len(game.RightPaths))
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, " len(TopPaths) : ", len(game.TopPaths))
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, " len(DownPaths) : ", len(game.DownPaths))

	game.FishMap = make(map[int]int64)

	game.MzSpace = -999
	game.MzPersist = -999
	game.Wave = -999

	return game
}

type Msg_GameFishing_Exit_LKPY struct {
	Uid  int64 `jsson:"uid"`
	Seat int   `json:"seat"`
}
type Game_Fishing_LKPY_Person struct {
	Uid         int64             `json:"uid"`
	Gold        int               `json:"gold"`
	Total       int               `json:"total"`
	Cannon      *Game_Cannon_LKPY `json:"cannon"`
	Seat        int               `json:"seat"`
	Name        string            `json:"name"`
	Head        string            `json:"head"`
	Address     string            `json:"address"`
	IP          string            `json:"ip"`
	Sex         int               `json:"sex"`
	Expend      int               `json"expend`
	ZZD         float64           `json:"zzd"`
	BDD         float64           `json:"bdd"`
	DjNum       []int             `json:"djnum"`
	ACTIVE_LKPY float64           `json:"ACTIVE_LKPY"`
	CannonInfo  []int             `json:"cannoninfo"`
	PlayerMz    float64           `json:"playermz"`
	MzUp        int               `json:"mzup"` //!　0－没有buff 1-命中增加 2-命中减小
}

type Game_LKPY_GoOut struct {
	Key int `json:"key"`
}

type Game_Cannon_LKPY struct { //! 火炮
	Id     int     `json:"id"`
	Expend int     `json:"expend"` //! 每颗子弹的消耗
	Rad    int     `json:"rad"`
	PL     float64 `json:"pl"`
	Bulspd int     `json:"Bulspd"`
	Space  int     `json:"space"`
}

type Game_Fish_LKPY struct { //! 鱼
	Key  int            `json:"key"`
	Type int            `json:"type"` //! 鱼的种类
	Dod  float64        `json:"dod"`  //! 命中率
	Win  int            `json:"win"`  //! 奖金
	Num  int            `json:"num"`
	Max  int            `json:"max"`
	Path Game_Path_LKPY `json:"path"` //!　路径信息
}

type Game_Path_LKPY struct {
	Id      int     `json:"id"`
	Path    int     `json:"path"`
	CurTime float64 `json:"curtime"`
	Time    float64 `json:"time"`
}

func (self *Game_Fishing_LKPY) GetPathById(id int) *Game_Path_LKPY {
	for i := 0; i < len(self.AllPath); i++ {
		if self.AllPath[i].Path == id {
			return &self.AllPath[i]
		}
	}
	return nil
}

type Msg_GameFishing_State_LKPY struct {
	Type  int     `json:"type"` //! 鱼的type
	Key   int     `json:"key"`  //! 鱼的key
	Path  int     `json:"path"` //! 路径id
	Time  float64 `json:"time"`
	State int     `json:"state"` //! 0-游进屏幕中 1-游出屏幕外
}

type Msg_GameFishing_Win_LKPY struct {
	Uid     int64           `json:"uid"`
	Total   int             `json:"total"`
	Name    string          `json:"name"`
	Head    string          `json:"head"`
	Address string          `json:"address"`
	IP      string          `json:"ip"`
	Sex     int             `json:"sex"`
	Win     int             `json:"win"` //! 赢了多少
	Fish    []Son_Fish_LKPY `json:"fish"`
}
type Son_Fish_LKPY struct {
	Key int `json:"key"`
	Win int `json:"win"`
}

type Msg_GameFishing_Info_LKPY struct {
	PersonMgr [4]Son_GameFishing_Info_LKPY `json:"personmgr"`
	CurFish   []Son_GameFish_LKPY          `json:"curfish"` //! 出现在屏幕中的鱼
}

type Msg_GameFishing_Help_LKPY struct {
	FishInfo []Son_GameFishing_Help_LKPY `json:"fishinfo"`
}

type Son_GameFishing_Help_LKPY struct {
	Type int `json:"type"`
	Win  int `json:"win"`
}

type Son_GameFish_LKPY struct {
	Path    int     `json:"path"` //! 路线id
	Key     int     `json:"key"`  //! 鱼的key
	Type    int     `json:"type"` //! 鱼的type
	Time    float64 `json:"time"`
	CurTime float64 `json:"curtime"` //! 当前时间
}

type Son_GameFishing_Info_LKPY struct {
	Uid       int64     `json:"uid"`
	DjNum     []int     `json:"djnum"`
	Total     int       `json:"total"`  //! 总金币
	Cannon    int       `json:"cannon"` //! 火炮信息
	CanExpent int       `json:"canexpent"`
	BulSpd    int       `json:"bulspd"`
	Space     int       `json:"space"`
	Rad       int       `json:"rad"`
	Seat      int       `json:"seat"` //! 座位
	Name      string    `json:"name"`
	Head      string    `json:"head"`
	Address   string    `json:"address"`
	IP        string    `json:"ip"`
	Sex       int       `json:"sex"`
	Time      []float64 `json:"time"` //! 有效时间  !=0道具开始使用  ==0道具结束
}

type Msg_GameFishing_Fire_LKPY struct {
	Uid     int64   `json:"uid"`
	Angle   float64 `json:"angle"`
	FishKey int     `json:"fishkey"`
	Rad     int     `json:"rad"`
	BulSpd  int     `json:"bulspd"`
	Seat    int     `json:"seat"`
	Total   int     `json:"total"`
	Name    string  `json:"name"`
	Head    string  `json:"head"`
	Address string  `json:"address"`
	IP      string  `json:"ip"`
	Sex     int     `json:"sex"`
}

type Msg_GameFishing_Total_LKPY struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

func (self *Game_Fishing_LKPY) SendTotal(uid int64, total int) {
	var msg Msg_GameFishing_Total_LKPY
	msg.Uid = uid
	msg.Total = total

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	self.room.broadCastMsg("gamegoldtotal", &msg)

}

func (self *Game_Fishing_LKPY) getinfo(uid int64) *Msg_GameFishing_Info_LKPY {
	var msg Msg_GameFishing_Info_LKPY
	for i := 0; i < len(self.CurFish); i++ {
		var fish Son_GameFish_LKPY
		fish.Path = self.CurFish[i].Path.Path
		fish.Time = self.CurFish[i].Path.Time
		fish.Key = self.CurFish[i].Key
		fish.Type = self.CurFish[i].Type
		fish.CurTime = self.CurFish[i].Path.CurTime
		msg.CurFish = append(msg.CurFish, fish)
	}
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i] != nil {
			msg.PersonMgr[i].Uid = self.PersonMgr[i].Uid
			msg.PersonMgr[i].Total = self.PersonMgr[i].Total
			msg.PersonMgr[i].Sex = self.PersonMgr[i].Sex
			msg.PersonMgr[i].Seat = self.PersonMgr[i].Seat
			msg.PersonMgr[i].Name = self.PersonMgr[i].Name
			msg.PersonMgr[i].IP = self.PersonMgr[i].IP
			msg.PersonMgr[i].DjNum = self.PersonMgr[i].DjNum
			msg.PersonMgr[i].Head = self.PersonMgr[i].Head
			msg.PersonMgr[i].Cannon = self.PersonMgr[i].Cannon.Id
			msg.PersonMgr[i].CanExpent = self.PersonMgr[i].Cannon.Expend
			msg.PersonMgr[i].BulSpd = self.PersonMgr[i].Cannon.Bulspd
			msg.PersonMgr[i].Space = self.PersonMgr[i].Cannon.Space
			msg.PersonMgr[i].Rad = self.PersonMgr[i].Cannon.Rad
			msg.PersonMgr[i].Time = make([]float64, 0)
			msg.PersonMgr[i].Time = append(msg.PersonMgr[i].Time, self.PersonMgr[i].ZZD)
			msg.PersonMgr[i].Time = append(msg.PersonMgr[i].Time, self.PersonMgr[i].BDD)
		}
	}
	return &msg
}

type Msg_GameFishing_HitFish_LKPY struct {
	Uid      int64 `json:"uid"`
	FishKey  []int `json:"fishkey"` //! 打中鱼的key
	Rad      int   `json:"rad"`     //! 网的半径
	CannonId int   `json:"cannonid"`
}

type Msg_GameFishing_SetCannon_LKPY struct {
	Uid      int64 `json:"uid"`
	CannonId int   `json:"cannonid"`
}

type Msg_GameFishing_OpenFire_LKPY struct {
	Uid     int64   `json:"uid"`
	Angle   float64 `json:"angle"`
	FishKey int     `json:"fishkey"`
}

type Msg_GameFishing_DJ_LKPY struct {
	//	Uid int64 `json:"uid"`
	//	DjNum []int   `json:"djnum"`
	//	Seat  int     `json:"seat"`
	//	Type  int     `json:"type"` //! 1-追踪弹 2-冰冻弹
	Time float64 `json:"time"` //! 有效时间  !=0道具开始使用  ==0道具结束
}

type Msg_GameFish_GetDJ_LKPY struct {
	Uid   int64 `json:"uid"`
	Key   int   `json:"key"`
	Type  int   `json:"type"`
	DjNum []int `json:"djnum"`
}

func (self *Game_Fishing_LKPY) GetPerson(uid int64) *Game_Fishing_LKPY_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i] == nil {
			continue
		}
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}
	return nil
}

func (self *Game_Fishing_LKPY) GetTypeByKey(key int) int {
	for i := 0; i < len(self.CurFish); i++ {
		if self.CurFish[i].Key == key {
			return self.CurFish[i].Type
		}
	}
	return -1
}

//!　获取单行路径
func (self *Game_Fishing_LKPY) GetPath(_type int) int {
	path := 0
	if _type == 0 {
		switch self.NextPath {
		case 1:
			if len(self.LeftPath) > 0 && lib.HF_GetRandom(100) < 60 || len(self.LeftPaths) == 0 {
				if len(self.LeftPath) == 0 {
					return 0
				}
				path = self.LeftPath[lib.HF_GetRandom(len(self.LeftPath))]
			} else {
				paths := self.LeftPaths[lib.HF_GetRandom(len(self.LeftPaths))]
				path = paths[lib.HF_GetRandom(len(paths))]
			}
		case 2:
			if len(self.RightPath) > 0 && lib.HF_GetRandom(100) < 60 || len(self.RightPaths) == 0 {
				if len(self.RightPath) == 0 {
					return 0
				}
				path = self.RightPath[lib.HF_GetRandom(len(self.RightPath))]
			} else {
				paths := self.RightPaths[lib.HF_GetRandom(len(self.RightPaths))]
				path = paths[lib.HF_GetRandom(len(paths))]
			}
		case 3:
			if len(self.TopPath) > 0 && lib.HF_GetRandom(100) < 60 || len(self.TopPaths) == 0 {
				if len(self.TopPath) == 0 {
					return 0
				}
				path = self.TopPath[lib.HF_GetRandom(len(self.TopPath))]
			} else {
				paths := self.TopPaths[lib.HF_GetRandom(len(self.TopPaths))]
				path = paths[lib.HF_GetRandom(len(paths))]
			}
		case 4:
			if len(self.DownPath) > 0 || lib.HF_GetRandom(100) < 60 || len(self.DownPaths) == 0 {
				if len(self.DownPath) == 0 {
					return 0
				}
				path = self.DownPath[lib.HF_GetRandom(len(self.DownPath))]
			} else {
				paths := self.DownPaths[lib.HF_GetRandom(len(self.DownPaths))]
				path = paths[lib.HF_GetRandom(len(paths))]
			}
		}
		if path == 0 && len(self.AllPath) > 0 {
			path = self.AllPath[lib.HF_GetRandom(len(self.AllPath))].Id
		}

		self.NextPath++
		if self.NextPath > 4 {
			self.NextPath = 1
		}
	} else {
		if len(self.BossPath) > 0 {
			path = self.BossPath[lib.HF_GetRandom(len(self.BossPath))]
		}
	}

	return path
}

//! 获取平行路径
func (self *Game_Fishing_LKPY) GetPaths() []int {
	path := make([]int, 0)
	switch self.NextPath {
	case 1:
		if len(self.LeftPaths) > 0 {
			path = self.LeftPaths[lib.HF_GetRandom(len(self.LeftPaths))]
		}
	case 2:
		if len(self.RightPaths) > 0 {
			path = self.RightPaths[lib.HF_GetRandom(len(self.RightPaths))]
		}
	case 3:
		if len(self.TopPaths) > 0 {
			path = self.TopPaths[lib.HF_GetRandom(len(self.TopPaths))]
		}
	case 4:
		if len(self.DownPath) > 0 {
			path = self.DownPaths[lib.HF_GetRandom(len(self.DownPaths))]
		}
	}

	self.NextPath++
	if self.NextPath > 4 {
		self.NextPath = 1
	}
	return path
}

//!　同步金币
func (self *Game_Fishing_LKPY_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

//! 获取鱼的对象
func (self *Game_Fishing_LKPY) GetFish(id int) {
	fishValue := staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis)
	if len(self.CurFish) > fishValue.MaxFish {
		return
	}

	var fish *Game_Fish_LKPY
	if len(self.FishPool) != 0 {
		fish = self.FishPool[0]
		copy(self.FishPool[0:], self.FishPool[1:])
		self.FishPool = self.FishPool[:len(self.FishPool)-1]
	} else {
		fish = new(Game_Fish_LKPY)
	}
	fishid := LKPY_TYPE[lib.HF_GetRandom(len(LKPY_TYPE))]
	//	maxGold := 0
	bigFish := 0
	mediumFish := 0
	smallFish := 0

	if id != -1 {
		fishid = id
		find := false
		for i := 0; i < len(self.CurFish); i++ {
			if self.CurFish[i].Type == fishid {
				find = true
			}
		}
		if find {
			return
		}
	} else {
		for j := 0; j < len(self.CurFish); j++ {
			if self.CurFish[j].Type > 13 && self.CurFish[j].Type <= 26 {
				bigFish++
			}
			if self.CurFish[j].Type > 7 && self.CurFish[j].Type <= 13 {
				mediumFish++
			}
			if self.CurFish[j].Type <= 7 {
				smallFish++
			}
		}

		if fishid > 13 && fishid <= 26 && bigFish > fishValue.MaxBigFish {
			return
		}
		if fishid > 7 && fishid <= 13 && mediumFish > fishValue.MaxMediumFish {
			return
		}
		if fishid <= 7 && smallFish > fishValue.MaxSmallFish {
			return
		}
	}

	fishvalue, ok := staticfunc.GetFishMgr().GetLKPYFishProperty(fishid, GetServer().Redis)
	if ok {
		fish.Type = fishvalue.Id
		fish.Dod = fishvalue.Dodge
		fish.Win = fishvalue.Win
		fish.Num = fishvalue.State
		fish.Max = fishvalue.Max

		fishNum := 0
		for j := 0; j < len(self.CurFish); j++ {
			if self.CurFish[j].Type%100 == fish.Type%100 {
				fishNum++
			}
		}

		for _, value := range self.PathMap {
			for j := 0; j < len(value); j++ {
				if value[j].Type == fish.Type {
					fishNum++
				}
			}
		}

		if fishNum >= fishvalue.Max {
			return
		}
		if self.FishMap[fish.Type] != 0 && time.Now().Unix()-self.FishMap[fish.Type] < int64(fishvalue.Timespan) && id != -1 { //!　特殊种类的鱼不用判断刷新时间
			return
		}

		if fish.Num == 1 { //! 刷鱼群
			paths := self.GetPaths()
			//! 李逵和企鹅的路径只能是从左开始
			if fish.Type == 21 || fish.Type == 18 {
				self.NextPath = 1
				paths = self.GetPaths()
			}

			num := lib.HF_GetRandom(5) + 1
			if fishNum+num > fishvalue.Max {
				num = fishvalue.Max - fishNum
			}
			index := 0

			for j := 0; j < num; j++ {
				if len(paths) == 0 {
					break
				}

				path := self.GetPathById(paths[index])
				if path != nil {
					var _fish *Game_Fish_LKPY
					_fish = new(Game_Fish_LKPY)
					_fish.Dod = fish.Dod
					_fish.Type = fish.Type
					_fish.Win = fish.Win
					_fish.Path.Id = path.Id
					_fish.Path.Path = path.Path
					_fish.Path.Time = path.Time
					_fish.Path.CurTime = 0
					_fish.Key = self.Key
					self.Key++
					self.PathMap[paths[index]] = append(self.PathMap[paths[index]], _fish)
					self.FishMap[fish.Type] = time.Now().Unix()
					index++
					if index >= len(paths) {
						index = 0
					}
					continue
				}
			}
			return
		} else { //!单刷
			pathId := 0

			pathId = self.GetPath(0)
			if fish.Type == 21 || fish.Type == 18 {
				self.NextPath = 1
				pathId = self.GetPath(0)
			}

			if fish.Type == 18 || fish.Type == 20 || fish.Type == 21 {
				self.GetPath(1)
			}

			path := self.GetPathById(pathId)
			if path != nil {
				fish.Key = self.Key
				fish.Path.Id = path.Id
				fish.Path.Path = path.Path
				fish.Path.Time = path.Time
				fish.Path.CurTime = 0
				self.Key++
				self.PathMap[pathId] = append(self.PathMap[pathId], fish)
				self.FishMap[fish.Type] = time.Now().Unix()
				return
			}
		}
	}

}

//! 杀死条鱼
func (self *Game_Fishing_LKPY) Kill(uid int64, key []int) {
	per := self.GetPerson(uid)
	if per == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "没找到该玩家！")
		return
	}

	for i := 0; i < len(key); i++ {
		fishType := self.GetTypeByKey(key[i])
		if fishType != 98 && fishType != 99 {
			continue
		}
		property, _ := staticfunc.GetFishMgr().GetLKPYFishProperty(fishType, GetServer().Redis)
		notice := "eed33a#恭喜"
		notice += ("#91ec39#" + per.Name)
		notice += "#eed33a#在"
		if self.room.Type%10 == 0 {
			notice += "#59acbe#捕鱼新手场"
		} else if self.room.Type%10 == 1 {
			notice += "#59acbe#捕鱼中级场"
		} else if self.room.Type%10 == 2 {
			notice += "#59acbe#捕鱼高级场"
		}
		notice += fmt.Sprintf("#eed33a#%d号渔场", self.room.Param1)
		notice += ("#cbc853#打死" + property.Name)

		GetServer().SendNotice(notice)
	}

	var msg Msg_GameFishing_Win_LKPY
	msg.Uid = uid
	msg.Address = per.Address
	msg.Sex = per.Sex
	msg.Name = per.Name
	msg.IP = per.IP
	msg.Head = per.Head
	msg.Fish = make([]Son_Fish_LKPY, 0)

	perWin := 0
	getDJ := false
	//	djKey := -1
	boom := 0

	for j := 0; j < len(key); j++ {
		for i := 0; i < len(self.CurFish); i++ {
			if self.CurFish[i].Key == key[j] {
				fishType := self.GetTypeByKey(key[j])

				if fishType == 24 {
					fishValue, ok := staticfunc.GetFishMgr().GetLKPYFishProperty(fishType, GetServer().Redis)
					if ok {
						boom += fishValue.Win
					}
				}
				if fishType == 23 {
					getDJ = true
				}

				fish := self.CurFish[i]
				copy(self.CurFish[i:], self.CurFish[i+1:])
				self.CurFish = self.CurFish[:len(self.CurFish)-1]

				self.FishPool = append(self.FishPool, fish)

				if fishType == 25 { //! 宝箱
					fishWin := int((float64(lib.HF_GetRandom(fish.Win) + 1)) * per.Cannon.PL)
					perWin += fishWin
					var son Son_Fish_LKPY
					son.Key = fish.Key
					son.Win = fishWin
					msg.Fish = append(msg.Fish, son)
				} else {
					perWin += int(float64(fish.Win) * float64(per.Cannon.Expend))
					var son Son_Fish_LKPY
					son.Key = fish.Key
					son.Win = int(float64(fish.Win) * float64(per.Cannon.Expend))
					msg.Fish = append(msg.Fish, son)
				}

				break
			}
		}
	}

	dealwin := per.Expend - perWin
	if dealwin != 0 {
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
	}
	per.Expend = 0

	GetServer().SetLkpyMoney(self.room.Type%280000, GetServer().LkpyMoney[self.room.Type%280000]-int64(perWin))
	cost := int(math.Ceil(float64(perWin) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
	if perWin-cost > 0 {
		perWin -= cost
		GetServer().SqlAgentGoldLog(uid, cost, self.room.Type)
		GetServer().SqlAgentBillsLog(uid, cost/2, self.room.Type)
	}

	if getDJ {
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i] != nil {
				self.GamePlay(self.PersonMgr[i].Uid, 2)
				break
			}
		}
	}

	per.Total += perWin

	msg.Total = per.Total
	msg.Win = perWin
	self.room.broadCastMsg("gamefishkill", &msg)

	if boom > 0 {
		self.Boom(uid, boom)
	}
}

//爆炸
func (self *Game_Fishing_LKPY) Boom(uid int64, boom int) {
	//	per := self.GetPerson(uid)
	kill := make([]int, 0)
	money := GetServer().LkpyMoney[self.room.Type%280000]
	for i := 0; i < len(self.CurFish); i++ {
		kill = append(kill, self.CurFish[i].Key)
		money -= int64(self.CurFish[i].Win)
	}
	self.Kill(uid, kill)
}

//! 游出屏幕外
func (self *Game_Fishing_LKPY) GoOut(key int) {
	for i := 0; i < len(self.CurFish); i++ {
		if self.CurFish[i].Key == key {
			copy(self.CurFish[i:], self.CurFish[i+1:])
			self.CurFish = self.CurFish[:len(self.CurFish)-1]
			break
		}
	}
}

//! 游入屏幕中
func (self *Game_Fishing_LKPY) ComeIn(fish *Game_Fish_LKPY) {
	fish.Path.CurTime = 0
	self.CurFish = append(self.CurFish, fish)
	var msg Msg_GameFishing_State_LKPY
	msg.Type = fish.Type
	msg.Path = fish.Path.Path
	msg.Time = fish.Path.Time
	msg.Key = fish.Key
	msg.State = 0
	self.room.broadCastMsg("gamefishstate", &msg)

}

//!　设置炮的等级
func (self *Game_Fishing_LKPY) SetCannon(uid int64, id int) { //!　0--　1++
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, "---------------------------------------")
	per := self.GetPerson(uid)
	if per == nil {
		return
	}
	cannon := per.Cannon.Id

	if id == 0 { //! 减一等级
		if cannon == 1 || id == 7 || id == 13 { //!　最小的炮
			if cannon == 13 {
				cannon = 18
			}
			if id == 7 {
				cannon = 12
			}
			if id == 1 {
				cannon = 6
			}
		} else {
			cannon--
		}
	} else { //! 加一等级
		if cannon == 6 || cannon == 12 || cannon == 18 {
			if cannon == 6 {
				cannon = 1
			}
			if cannon == 12 {
				cannon = 7
			}
			if cannon == 18 {
				cannon = 13
			}
		} else {
			cannon++
		}
	}

	fishValue := staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis)
	per.ACTIVE_LKPY = fishValue.Action

	gunvalue, ok := staticfunc.GetFishMgr().GetLKPYGunProperty(cannon, GetServer().Redis)
	if !ok {
		return
	}

	per.Cannon.Id = gunvalue.Id
	//	per.Cannon.Crit = gunvalue.Crit
	//	per.Cannon.Atk = gunvalue.Power
	per.Cannon.Expend = gunvalue.Cost
	per.Cannon.Rad = gunvalue.Radius
	per.Cannon.PL = gunvalue.PL
	per.Cannon.Bulspd = gunvalue.BulSpd
	per.Cannon.Space = gunvalue.Space

	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "id : ", gunvalue.Id, " crit : ", gunvalue.Crit, " atk : ", gunvalue.Power, " expend : ", gunvalue.Cost, " radius : ", gunvalue.Radius, " pl : ", gunvalue.PL)
	var msg Son_GameFishing_Info_LKPY
	msg.Uid = per.Uid
	msg.Total = per.Total
	msg.Sex = per.Sex
	msg.Seat = per.Seat
	msg.Name = per.Name
	msg.IP = per.IP
	msg.Head = per.Head
	msg.Address = per.Address
	msg.Cannon = per.Cannon.Id
	msg.CanExpent = per.Cannon.Expend
	msg.BulSpd = per.Cannon.Bulspd
	msg.Space = per.Cannon.Space
	msg.Rad = per.Cannon.Rad
	self.room.broadCastMsg("gamefishingcannon", &msg)
}

//!　打中鱼
func (self *Game_Fishing_LKPY) HitFish(uid int64, key []int, rad int, canInfo int) {

	per := self.GetPerson(uid)
	if per == nil {
		return
	}
	if rad != per.Cannon.Rad {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------------- rad != per.Cannon.Rad   rad : ", rad, " can.rad : ", per.Cannon.Rad)
		return
	}

	if per.CannonInfo[canInfo-1] <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------------------- per.CannonInfo[canInfo-1] <= 0")
		return
	}
	per.CannonInfo[canInfo-1]--

	if key == nil || len(key) == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "per.CannonInfo[canInfo-1]-- : ", per.CannonInfo[canInfo-1])
		return
	}

	fishValue := staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis)
	per.ACTIVE_LKPY = fishValue.Action

	kill := make([]int, 0)
	for j := 0; j < len(key); j++ {
		find := false
		var fish *Game_Fish_LKPY
		for i := 0; i < len(self.CurFish); i++ {
			if key[j] == self.CurFish[i].Key {
				find = true
				fish = self.CurFish[i]
				break
			}
		}

		if !find {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "目标鱼不在屏幕中 ", key[j])
			var msg Game_LKPY_GoOut
			msg.Key = key[j]
			self.room.SendMsg(uid, "gamelkpygoout", &msg)
			continue
		}
		//! 基础命中
		playerhit := float64(staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis).PlayerMz) / 100.0
		//! 鱼的命中
		fishhit := float64(fish.Dod) / 100.0
		//! buff
		//! 座位加成
		buff := 1.0 + staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis).SeatMz[per.Seat]/100.0
		buff += per.PlayerMz / 100
		//lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------------奖池 : ", GetServer().LkpyMoney[self.room.Type%280000])
		//lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------------收益 : ", int64(float64(fish.Win)*per.Cannon.PL))
		//lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------------JackPotMin : ", lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin)
		//lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------------步长 : ", staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis).Step)
		buff += float64((GetServer().LkpyMoney[self.room.Type%280000]-int64(float64(fish.Win)*per.Cannon.PL)-lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin)/int64(staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis).Step)) * 0.01
		if buff < 0 {
			buff = 0
		}
		//! 最终命中率
		//lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------------基础命中 : ", playerhit)
		//lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------------鱼的命中 : ", fishhit)
		//lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------------座位加成 : ", buff)
		hit := playerhit * fishhit * buff
		if lib.HF_GetRandom(10000) < int(hit*10000) {
			kill = append(kill, key[j])
		}
	}
	if len(kill) > 0 {
		self.Kill(uid, kill)
	}

}

//! 射击
func (self *Game_Fishing_LKPY) Fire(uid int64, angle float64, fishkey int) {
	per := self.GetPerson(uid)
	if per == nil {
		return
	}

	if per.Total < per.Cannon.Expend {
		self.room.SendErr(uid, "您的金币不足，请前往充值!")
		return
	}

	gunValue, ok := staticfunc.GetFishMgr().GetLKPYGunProperty(per.Cannon.Id, GetServer().Redis)

	if ok {
		fishValue := staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis)
		per.ACTIVE_LKPY = fishValue.Action

		per.Total -= per.Cannon.Expend
		per.Expend += per.Cannon.Expend
		bl := float64(lib.GetManyMgr().GetProperty(self.room.Type).DealCost)
		cost := int(math.Ceil(float64(per.Cannon.Expend) * bl / 100.0))
		GetServer().SetLkpyMoney(self.room.Type%280000, GetServer().LkpyMoney[self.room.Type%280000]+int64(per.Cannon.Expend-cost))

		var msg Msg_GameFishing_Fire_LKPY
		msg.Uid = uid
		msg.Seat = per.Seat
		msg.Total = per.Total
		msg.Sex = per.Sex
		msg.Name = per.Name
		msg.IP = per.IP
		msg.BulSpd = gunValue.BulSpd
		msg.Head = per.Head
		msg.Angle = angle
		msg.FishKey = fishkey
		msg.Address = per.Address
		msg.Rad = gunValue.Radius
		self.room.broadCastMsg("gamefishingfire", &msg)

		per.CannonInfo[per.Cannon.Id-1]++
		//lib.GetLogMgr().Output(lib.LOG_DEBUG, "---------------- per.Cannon.Id-1 : ", per.Cannon.Id-1)
		//lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------- per.CannonInfo : ", per.CannonInfo)
	}

}

func (self *Game_Fishing_LKPY) OnBegin() {

}

func (self *Game_Fishing_LKPY) OnEnd() {

}

func (self *Game_Fishing_LKPY) OnInit(room *Room) {
	self.room = room
}

func (self *Game_Fishing_LKPY) OnRobot(robot *lib.Robot) {

}

func (self *Game_Fishing_LKPY) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i] == nil {
			continue
		}
		if self.PersonMgr[i].Uid == person.Uid {
			self.PersonMgr[i].IP = person.ip
			self.PersonMgr[i].Address = person.minfo.Address
			self.PersonMgr[i].Sex = person.Sex
			self.PersonMgr[i].SynchroGold(person.Gold)
			self.room.broadCastMsg("gamefishinginfo", self.getinfo(person.Uid))
			return
		}
	}

	_person := new(Game_Fishing_LKPY_Person)
	_person.Uid = person.Uid
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Address = person.minfo.Address
	_person.Head = person.Imgurl
	_person.IP = person.ip
	_person.Name = person.Name
	_person.Sex = person.Sex

	gunvalue, ok := staticfunc.GetFishMgr().GetLKPYGunProperty(1, GetServer().Redis)

	_person.CannonInfo = make([]int, len(staticfunc.GetFishMgr().GetAllLKPYGunProperty(GetServer().Redis)))
	if !ok {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, " 读取gunvalue出错")
		return
	}
	//	if _person.Total < gunvalue.Cost {
	//		self.room.KickPerson(person.Uid, 99)
	//		return
	//	}
	_person.Cannon = new(Game_Cannon_LKPY)
	_person.Cannon.Id = gunvalue.Id
	//	_person.Cannon.Crit = gunvalue.Crit
	//	_person.Cannon.Atk = gunvalue.Power
	_person.Cannon.Expend = gunvalue.Cost
	_person.Cannon.Rad = gunvalue.Radius
	_person.Cannon.PL = gunvalue.PL
	_person.Cannon.Bulspd = gunvalue.BulSpd
	_person.Cannon.Space = gunvalue.Space
	_person.DjNum = make([]int, 2)
	fishValue := staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis)
	_person.ACTIVE_LKPY = fishValue.Action
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i] == nil {
			_person.Seat = i
			self.PersonMgr[i] = _person
			break
		}
	}
	_person.PlayerMz = staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis).PlayerMz
	self.room.broadCastMsg("gamefishinginfo", self.getinfo(person.Uid))

}

func (self *Game_Fishing_LKPY) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gameopenfire": //!发射子弹
		self.Fire(msg.Uid, msg.V.(*Msg_GameFishing_OpenFire).Angle, msg.V.(*Msg_GameFishing_OpenFire).FishKey)
	case "gamesetcannon": //! 设置炮的等级
		self.SetCannon(msg.Uid, msg.V.(*Msg_GameFishing_SetCannon).CannonId)
	case "gamehitfish": //! 打中鱼
		self.HitFish(msg.Uid, msg.V.(*Msg_GameFishing_HitFish).FishKey, msg.V.(*Msg_GameFishing_HitFish).Rad, msg.V.(*Msg_GameFishing_HitFish).CannonId)
	case "gameplay":
		self.GamePlay(msg.Uid, msg.V.(*Msg_GamePlay).Type)
	case "gamehelp":
		self.GameHelp(msg.Uid)
	}
}

func (self *Game_Fishing_LKPY) GameHelp(uid int64) {
	per := self.GetPerson(uid)
	if per == nil {
		return
	}
	fishValue := staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis)
	per.ACTIVE_LKPY = fishValue.Action

	fishList := staticfunc.GetFishMgr().GetAllLKPYFishProperty(GetServer().Redis)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "------------ help : ", len(fishList), " fishlist : ", fishList)
	//	cannonList := staticfunc.GetFishMgr().GetAllLKPYGunProperty(GetServer().Redis)
	if len(fishList) > 0 {
		var msg Msg_GameFishing_Help_LKPY
		msg.FishInfo = make([]Son_GameFishing_Help_LKPY, 0)
		for i := 0; i < len(fishList); i++ {
			var son Son_GameFishing_Help_LKPY
			son.Type = fishList[i].Id
			son.Win = fishList[i].Win
			msg.FishInfo = append(msg.FishInfo, son)
		}
		self.room.SendMsg(uid, "gamefishinghelp", &msg)
	}

}

//使用道具
func (self *Game_Fishing_LKPY) GamePlay(uid int64, index int) {
	per := self.GetPerson(uid)
	if per == nil {
		return
	}

	if index < 1 && index > 2 {
		return
	}

	fishValue := staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis)
	per.ACTIVE_LKPY = fishValue.Action

	switch index {
	case 1:
		//		if per.DjNum[0] < 1 {
		//			self.room.SendErr(uid, "道具数量不足")
		//			return
		//		}
		//		per.DjNum[0]--
		//		per.ZZD = fishValue.ZZDTime
		//		var msg Msg_GameFishing_DJ_LKPY
		//		msg.Seat = per.Seat
		//		msg.DjNum = per.DjNum
		//		msg.Type = 1
		//		msg.Uid = per.Uid
		//		msg.Time = fishValue.ZZDTime
		//		self.room.broadCastMsg("gamefishingdj", &msg)
	case 2:
		//		if per.DjNum[1] < 1 {
		//			self.room.SendErr(uid, "道具数量不足")
		//			return
		//		}
		//		per.DjNum[1]--
		per.BDD = fishValue.BBDTime
		var msg Msg_GameFishing_DJ_LKPY
		msg.Time = fishValue.BBDTime
		//		msg.DjNum = per.DjNum
		//		msg.Seat = per.Seat
		//		msg.Type = 2
		//		msg.Uid = per.Uid

		self.room.broadCastMsg("gamefishingdj", &msg)
	}
}

func (self *Game_Fishing_LKPY) OnBye() {

}

func (self *Game_Fishing_LKPY) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i] == nil {
			continue
		}
		if self.PersonMgr[i].Uid == uid {

			//! 将庄家赢的钱存入数据库
			dealwin := self.PersonMgr[i].Expend
			if dealwin != 0 {
				GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
			}
			//! 退出房间同步金币
			gold := self.PersonMgr[i].Total - self.PersonMgr[i].Gold
			if gold > 0 {
				GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
			} else if gold < 0 {
				GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
			}
			self.PersonMgr[i].Gold = self.PersonMgr[i].Total

			var msg Msg_GameFishing_Exit_LKPY
			msg.Uid = self.PersonMgr[i].Uid
			msg.Seat = self.PersonMgr[i].Seat
			self.room.broadCastMsg("gamefishingexit", &msg)

			//! 插入战绩
			if gold != 0 {
				var record Rec_Fishing_Info_LKPY
				record.Time = time.Now().Unix()
				record.GameType = self.room.Type
				var rec Son_Rec_Fishing_Person_LKPY
				rec.Uid = uid
				rec.Name = self.PersonMgr[i].Name
				rec.Head = self.PersonMgr[i].Head
				rec.Score = gold
				record.Info = append(record.Info, rec)
				GetServer().InsertRecord(self.room.Type, self.PersonMgr[i].Uid, lib.HF_JtoA(&record), rec.Score)
			}

			self.PersonMgr[i] = nil

			break
		}
	}
}

func (self *Game_Fishing_LKPY) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_Fishing_LKPY) OnBalance() {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i] == nil {
			continue
		}

		//! 将庄家赢的钱存入数据库
		dealwin := self.PersonMgr[i].Expend
		if dealwin != 0 {
			GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
		}

		//! 退出房间同步金币
		gold := self.PersonMgr[i].Total - self.PersonMgr[i].Gold
		if gold > 0 {
			GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		self.PersonMgr[i].Gold = self.PersonMgr[i].Total
	}
}

func (self *Game_Fishing_LKPY) OnIsBets(uid int64) bool {
	return false
}

func (self *Game_Fishing_LKPY) OnTime() {
	bdd := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i] == nil {
			continue
		}
		if self.PersonMgr[i].BDD > 0 {
			bdd = true
			break
		}
	}

	if !bdd { //! 没有人使用冰冻弹，刷鱼，游动
		for i := 0; i < len(self.CurFish); {
			self.CurFish[i].Path.CurTime += 0.1
			if self.CurFish[i].Path.CurTime > self.CurFish[i].Path.Time {
				self.GoOut(self.CurFish[i].Key)
			} else {
				i++
			}
		}

		self.GetFish(-1)
		if len(self.CurFish) < 100 { //!如果屏幕中的鱼少于30条，加快刷新效率
			self.GetFish(-1)
		}

		for key, value := range self.PathMap {
			if len(value) > 0 {
				self.ComeIn(value[0])
				copy(self.PathMap[key][0:], self.PathMap[key][1:])
				self.PathMap[key] = self.PathMap[key][:len(self.PathMap[key])-1]
			}
		}

		if self.DjFishTime <= 0 {
			value, ok := staticfunc.GetFishMgr().GetLKPYFishProperty(90, GetServer().Redis)
			if ok {
				self.GetFish(90)
				self.DjFishTime = value.Timespan
			}
		} else {
			self.DjFishTime -= 0.1
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i] == nil {
			continue
		}
		//		if self.PersonMgr[i].ZZD != 0 {
		//			self.PersonMgr[i].ZZD -= 0.1
		//			if self.PersonMgr[i].ZZD == 0 {
		//				var msg Msg_GameFishing_DJ_LKPY
		//				msg.Uid = self.PersonMgr[i].Uid
		//				msg.Seat = self.PersonMgr[i].Seat
		//				msg.Type = 1
		//				msg.Time = 0
		//				self.room.broadCastMsg("gamefishingzzd", &msg)
		//			}
		//		}
		if self.PersonMgr[i].BDD > 0 {
			self.PersonMgr[i].BDD -= 0.1
			if self.PersonMgr[i].BDD <= 0 {
				find := false //!  是否有其他人使用冰冻弹
				for j := 0; j < len(self.PersonMgr); j++ {
					if self.PersonMgr[j] != nil && self.PersonMgr[j].BDD > 0 {
						find = true
						break
					}
				}
				if !find {
					var msg Msg_GameFishing_DJ_LKPY
					//					msg.Uid = self.PersonMgr[i].Uid
					//					msg.Seat = self.PersonMgr[i].Seat
					//					msg.Type = 2
					msg.Time = 0
					self.room.broadCastMsg("gamefishingzzd", &msg)
				}

			}
		}

	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i] == nil {
			continue
		}
		self.PersonMgr[i].ACTIVE_LKPY -= 0.1
		if self.PersonMgr[i].ACTIVE_LKPY <= 0 {
			self.room.KickPerson(self.PersonMgr[i].Uid, 96)
		}
	}

	if self.MzSpace == -999 { //! 读取命中波动
		self.MzSpace = staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis).MzSpace
		if self.MzSpace <= 0 {
			self.MzSpace = -999
		}
		self.MzUp = lib.HF_GetRandom(100) < 50
	}
	if self.MzSpace == -999 { //! 命中波动没有设置，return
		return
	}
	self.MzSpace -= 0.1 //!　倒计时--

	if self.MzSpace <= 0 {
		self.MzSpace = -999
		if self.MzPersist == -999 { //! 读取命中波动持续时间
			self.MzPersist = staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis).MzPersist
			if self.MzPersist <= 0 {
				self.MzPersist = -999
			}
		}
	}

	if self.MzPersist == -999 { //! 命中持续时长没有设置
		return
	}

	if self.MzPersist > 0 { //! 增减玩家命中
		self.MzPersist -= 0.1
		self.Wave = staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis).Wave

		mzUp := self.MzUp
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i] == nil {
				continue
			}
			if self.PersonMgr[i].MzUp != 0 {
				continue
			}
			if mzUp {
				self.PersonMgr[i].PlayerMz += self.Wave
				self.PersonMgr[i].MzUp = 1
				mzUp = false
			} else {
				self.PersonMgr[i].PlayerMz -= self.Wave
				self.PersonMgr[i].MzUp = 2
				mzUp = true
			}
		}
	} else {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------- 进入else")
		self.MzPersist = -999
		self.Wave = staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis).Wave
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i] == nil {
				continue
			}
			if self.PersonMgr[i].MzUp == 0 {
				continue
			}
			if self.PersonMgr[i].MzUp == 1 {
				self.PersonMgr[i].PlayerMz -= self.Wave
				self.PersonMgr[i].MzUp = 0
			}
			if self.PersonMgr[i].MzUp == 2 {
				self.PersonMgr[i].PlayerMz += self.Wave
				self.PersonMgr[i].MzUp = 0
			}
		}

	}

}
