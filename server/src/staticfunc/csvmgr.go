package staticfunc

import (
	"encoding/csv"
	"fmt"
	"io"
	"lib"
	"log"
	"net/http"
	"os"
	"rjmgr"
	"strconv"
	"strings"
)

type CsvNode map[string]string

type CsvMgr struct {
	Data map[string]map[int]CsvNode

	//! 拼天九结构
	MapPTJ map[int]CsvNode
}

var csvmgrsingleton *CsvMgr = nil

//! public
func GetCsvMgr() *CsvMgr {
	if csvmgrsingleton == nil {
		csvmgrsingleton = new(CsvMgr)
		csvmgrsingleton.InitStruct()
	}

	return csvmgrsingleton
}

func (self *CsvMgr) InitStruct() {
	self.Data = make(map[string]map[int]CsvNode)
	self.MapPTJ = make(map[int]CsvNode)
}

func (self *CsvMgr) InitData() {
	self.ReadData("game")
	self.ReadData("ptj")
	self.ReadData("fish")
	self.ReadData("cannon")
	self.ReadData("path")
	self.ReadData("lkpyfish")
	self.ReadData("lkpypath")
	self.ReadData("lkpycannon")
	self.ReadData("sign")
	self.ReadData("enter")
	self.LoadPTJ()
}

func (self *CsvMgr) Reload() {
	self.InitStruct()
	self.InitData()
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "reload csv")
}

func (self *CsvMgr) ReadData(name string) {
	_, ok := self.Data[name]
	if ok {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "重复读入csv:", name)
		return
	}

	file, err := os.Open("../csv/" + name + ".csv")
	if err != nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "csv err1:", name, err)
		return
	}
	defer file.Close()

	buf := make([]byte, 65535)
	n, err := file.Read(buf)
	if err != nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "csv err4:", name, err)
		return
	}
	buf = buf[0:n]
	//! 这里解密
	if LIBDEBUG != 1 {
		buf = buf[3:n]
		key := make([]byte, 0)
		key = append(key, []byte(rjmgr.GetRJMgr().Sign)...)
		key = append(key, []byte(rjmgr.GetRJMgr().Flag)...)
		for i := 0; i < len(buf); i++ {
			buf[i] ^= key[i%len(key)]
		}
	}

	header := make([]string, 0)
	reader := csv.NewReader(strings.NewReader(string(buf)))
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "csv err2:", name, err)
			return
		}

		if len(header) == 0 {
			header = record
		} else {
			id, err := strconv.Atoi(record[0])
			if err != nil {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "csv err3:", name, err)
				return
			}

			_, ok := self.Data[name]
			if !ok {
				self.Data[name] = make(map[int]CsvNode)
			}

			_, ok = self.Data[name][id]
			if !ok {
				self.Data[name][id] = make(CsvNode)
			}

			for i := 0; i < len(record); i++ {
				self.Data[name][id][header[i]] = record[i]
			}
		}
	}
}

func (self *CsvMgr) LoadPTJ() {
	csv := self.Data["ptj"]
	for _, value := range csv {
		self.MapPTJ[lib.HF_Atoi(value["card1"])*100+lib.HF_Atoi(value["card2"])] = value
	}
}

func (self *CsvMgr) GetPTJ(card1 int, card2 int) (CsvNode, bool) {
	value, ok := self.MapPTJ[card1*100+card2]
	return value, ok
}

//! 得到底分
func (self *CsvMgr) GetDF(gametype int) int {
	index := gametype % 100 / 10

	var df []int = []int{50, 100, 200, 300, 500, 1000}
	csv, ok := self.Data["enter"][gametype-index*10]
	if !ok {
		return df[index]
	}

	return lib.HF_Atoi(csv[fmt.Sprintf("df%d", index)])
}

//! 得到准入
func (self *CsvMgr) GetZR(gametype int) int {
	index := gametype % 100 / 10

	var zr []int = []int{1000, 2000, 4000, 6000, 10000, 20000}
	csv, ok := self.Data["enter"][gametype-index*10]
	if !ok {
		return zr[index]
	}

	return lib.HF_Atoi(csv[fmt.Sprintf("zr%d", index)])
}

func (self *CsvMgr) OnInit(w http.ResponseWriter, req *http.Request) {
	log.Fatalln("")
}
