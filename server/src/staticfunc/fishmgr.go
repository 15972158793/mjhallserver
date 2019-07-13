package staticfunc

import (
	"encoding/json"
	"fmt"
	"lib"

	"github.com/garyburd/redigo/redis"
	//"sync"
)

var DefaulFishValue = FishValue{0, 0, 0, 40, 500, 100, 6, 20, 40, 6, 5, 15, 60, 180, []float64{1.0, 1.0, -5.0, 1.0}}

//! 鱼的动态结构
type FishProperty struct {
	Id       int     `json:"id"`       //! id
	Win      int     `json:"win"`      //! 金币
	State    int     `json:"state"`    //! 0普通  1鱼群
	Dodge    float64 `json:"dodge"`    //! 命中率
	Max      int     `json:"max"`      //! 鱼池中最多几条
	Timespan float64 `json:"timespan"` //! 时间间隔
	Name     string  `json:"name"`     //!  名字
}

//! 炮的结构
type GunProperty struct {
	Id     int     `json:"id"`     //! id
	Cost   int     `json:"cost"`   //! 消耗金币
	Radius int     `json:"radius"` //! 半径
	PL     float64 `json:"pl"`     //! 炮的赔率
	BulSpd int     `json:"bulspd"` //! 速度
	Space  int     `json:"space"`  //! 间隔
}

//! 基础设定
type FishValue struct {
	MzSpace       float64   `json:"mzspace"`       //! 命中波动间隔
	MzPersist     float64   `json:"mzpersist"`     //! 命中波动持续时长
	Wave          float64   `json:"wave"`          //! 波动值
	PlayerMz      float64   `json:"playermz"`      //! 玩家命中率
	Step          int       `json:"step"`          //! 增长值，奖池每增加多少金币，命中率+0.01
	MaxFish       int       `json:"maxfish"`       //! 鱼池最大鱼数
	MaxBigFish    int       `json:"maxbigfish"`    //! 大型鱼最大数量
	MaxMediumFish int       `json:"maxmediumfish"` //! 中型鱼最大数量
	MaxSmallFish  int       `json:"maxsmallfish"`  //! 小型鱼数量
	MaxGold       int       `json:"maxgold"`       //! 最大黄金鱼数量
	BBDTime       float64   `json:"bbdtime"`       //! 道具时间
	ZZDTime       float64   `json:"zzdtime"`
	BossTime      float64   `json:"bosstime"` //!　boss刷新间隔
	Action        float64   `json:"action"`   //!活跃时间
	SeatMz        []float64 `json:"seatmz"`   //! 座位命中
}

//! 捕鱼管理者
type FishMgr struct {
}

var fishmgrSingleton *FishMgr = nil

//! 得到服务器指针
func GetFishMgr() *FishMgr {
	if fishmgrSingleton == nil {
		fishmgrSingleton = new(FishMgr)
	}

	return fishmgrSingleton
}

//! 得到属性
func (self *FishMgr) GetFishProperty(id int, _redis *redis.Pool) (FishProperty, bool) {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	var value FishProperty

	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("fish_%d", id)))
	if err == nil {
		json.Unmarshal(v, &value)
		return value, true
	}

	//! 从csv里面获取
	csv, ok := GetCsvMgr().Data["fish"][id]
	if !ok {
		return value, false
	}
	value.Id = lib.HF_Atoi(csv["type"])
	value.Dodge = lib.HF_Atof64(csv["dod"])
	value.Win = lib.HF_Atoi(csv["win"])
	value.State = lib.HF_Atoi(csv["num"])
	value.Timespan = lib.HF_Atof64(csv["timespan"])
	value.Max = lib.HF_Atoi(csv["max"])
	value.Name = csv["name"]
	return value, true
}

//! 设置属性
func (self *FishMgr) SetFishProperty(id int, _redis *redis.Pool, value FishProperty) {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()
	if value.Dodge <= 0 {
		value.Dodge = 1
	}

	c.Do("SET", fmt.Sprintf("fish_%d", id), lib.HF_JtoB(&value))
}

//! 得到所有鱼的属性
func (self *FishMgr) GetAllFishProperty(_redis *redis.Pool) []FishProperty {
	lst := make([]FishProperty, 0)
	for _, value := range GetCsvMgr().Data["fish"] {
		id := lib.HF_Atoi(value["type"])
		value, ok := self.GetFishProperty(id, _redis)
		if !ok {
			continue
		}
		lst = append(lst, value)
	}
	return lst
}

/////////////////////////////////////////////////////////////////////////
//! 得到炮属性
func (self *FishMgr) GetGunProperty(id int, _redis *redis.Pool) (GunProperty, bool) {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	var value GunProperty

	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("gun_%d", id)))
	if err == nil {
		json.Unmarshal(v, &value)
		return value, true
	}

	//! 从csv里面获取
	csv, ok := GetCsvMgr().Data["cannon"][id]
	if !ok {
		return value, false
	}
	value.Id = lib.HF_Atoi(csv["id"])
	value.Cost = lib.HF_Atoi(csv["expend"])
	value.Radius = lib.HF_Atoi(csv["rad"])
	value.PL = lib.HF_Atof64(csv["pl"])
	value.BulSpd = lib.HF_Atoi(csv["bulspd"])
	value.Space = lib.HF_Atoi(csv["space"])
	return value, true
}

//! 设置炮属性
func (self *FishMgr) SetGunProperty(id int, _redis *redis.Pool, value GunProperty) {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	if value.BulSpd <= 0 {
		value.BulSpd = 1000
	}
	if value.Space <= 0 {
		value.Space = 200
	}
	c.Do("SET", fmt.Sprintf("gun_%d", id), lib.HF_JtoB(&value))
}

//! 得到所有炮的属性
func (self *FishMgr) GetAllGunProperty(_redis *redis.Pool) []GunProperty {
	lst := make([]GunProperty, 0)
	for _, value := range GetCsvMgr().Data["cannon"] {
		id := lib.HF_Atoi(value["id"])
		value, ok := self.GetGunProperty(id, _redis)
		if !ok {
			continue
		}
		lst = append(lst, value)
	}
	return lst
}

///////////////////////////
//! 得到炮属性
func (self *FishMgr) GetFishValue(_redis *redis.Pool) FishValue {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	var value FishValue
	v, err := redis.Bytes(c.Do("GET", "fishvalue"))
	if err == nil {
		json.Unmarshal(v, &value)
		self.correct(&value)
		return value
	}

	return DefaulFishValue
}

//! 设置炮属性
func (self *FishMgr) SetFishValue(_redis *redis.Pool, value FishValue) {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	self.correct(&value)

	c.Do("SET", "fishvalue", lib.HF_JtoB(&value))
}

func (self *FishMgr) correct(value *FishValue) {
	if value.Step <= 0 {
		value.Step = 10000
	}
	if value.PlayerMz <= 0 {
		value.PlayerMz = 40
	}
	if value.Wave < 0 {
		value.Wave = 0
	}
	if value.MaxFish <= 0 {
		value.MaxFish = 200
	}
	if value.MaxBigFish <= 0 {
		value.MaxBigFish = 5
	}
	if value.MaxMediumFish <= 0 {
		value.MaxMediumFish = 50
	}
	if value.MaxSmallFish <= 0 {
		value.MaxSmallFish = 150
	}
	if value.BossTime <= 0 {
		value.BossTime = 60
	}
	if value.ZZDTime <= 0 {
		value.ZZDTime = 15
	}
	if value.BBDTime <= 0 {
		value.BBDTime = 10
	}
	if value.Action <= 0 {
		value.Action = 180
	}

	if len(value.SeatMz) != 4 {
		value.SeatMz = []float64{1.0, 1.0, -5.0, 1.0}
	}
}

//////////////////////////////////////////////////////////////////////! 李逵劈鱼

func (self *FishMgr) GetLKPYFishValue(_redis *redis.Pool) FishValue {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	var value FishValue
	v, err := redis.Bytes(c.Do("GET", "lkpyfishvalue"))
	if err == nil {
		json.Unmarshal(v, &value)
		self.correct(&value)
		return value
	}

	return DefaulFishValue
}

func (self *FishMgr) SetLKPYFishValue(_redis *redis.Pool, value FishValue) {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	self.correct(&value)

	c.Do("SET", "lkpyfishvalue", lib.HF_JtoB(&value))
}

//! 得到鱼属性
func (self *FishMgr) GetLKPYFishProperty(id int, _redis *redis.Pool) (FishProperty, bool) {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	var value FishProperty

	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("lkpyfish_%d", id)))
	if err == nil {
		json.Unmarshal(v, &value)
		return value, true
	}

	//! 从csv里面获取
	csv, ok := GetCsvMgr().Data["lkpyfish"][id]
	if !ok {
		return value, false
	}
	value.Id = lib.HF_Atoi(csv["type"])
	value.Dodge = lib.HF_Atof64(csv["dod"])
	value.Win = lib.HF_Atoi(csv["win"])
	value.State = lib.HF_Atoi(csv["num"])
	value.Timespan = lib.HF_Atof64(csv["timespan"])
	value.Max = lib.HF_Atoi(csv["max"])
	value.Name = csv["name"]
	return value, true
}

//! 设置鱼属性
func (self *FishMgr) SetLKPYFishProperty(id int, _redis *redis.Pool, value FishProperty) {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()
	if value.Dodge <= 0 {
		value.Dodge = 1
	}

	c.Do("SET", fmt.Sprintf("lkpyfish_%d", id), lib.HF_JtoB(&value))
}

//! 得到所有鱼的属性
func (self *FishMgr) GetAllLKPYFishProperty(_redis *redis.Pool) []FishProperty {
	lst := make([]FishProperty, 0)
	for _, value := range GetCsvMgr().Data["lkpyfish"] {
		id := lib.HF_Atoi(value["type"])
		value, ok := self.GetLKPYFishProperty(id, _redis)
		if !ok {
			continue
		}
		lst = append(lst, value)
	}
	return lst
}

//! 得到炮属性
func (self *FishMgr) GetLKPYGunProperty(id int, _redis *redis.Pool) (GunProperty, bool) {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	var value GunProperty

	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("lkpygun_%d", id)))
	if err == nil {
		json.Unmarshal(v, &value)
		return value, true
	}

	//! 从csv里面获取
	csv, ok := GetCsvMgr().Data["lkpycannon"][id]
	if !ok {
		return value, false
	}
	value.Id = lib.HF_Atoi(csv["id"])
	value.Cost = lib.HF_Atoi(csv["expend"])
	value.Radius = lib.HF_Atoi(csv["rad"])
	value.PL = lib.HF_Atof64(csv["pl"])
	value.BulSpd = lib.HF_Atoi(csv["bulspd"])
	value.Space = lib.HF_Atoi(csv["space"])
	return value, true
}

//! 设置炮属性
func (self *FishMgr) SetLKPYGunProperty(id int, _redis *redis.Pool, value GunProperty) {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	if value.BulSpd <= 0 {
		value.BulSpd = 1000
	}
	if value.Space <= 0 {
		value.Space = 200
	}
	c.Do("SET", fmt.Sprintf("lkpygun_%d", id), lib.HF_JtoB(&value))
}

//! 得到所有炮的属性
func (self *FishMgr) GetAllLKPYGunProperty(_redis *redis.Pool) []GunProperty {
	lst := make([]GunProperty, 0)
	for _, value := range GetCsvMgr().Data["lkpycannon"] {
		id := lib.HF_Atoi(value["id"])
		value, ok := self.GetLKPYGunProperty(id, _redis)
		if !ok {
			continue
		}
		lst = append(lst, value)
	}
	return lst
}
