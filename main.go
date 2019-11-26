package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fire988/utils"
)

const verString = "logserver 1.0"

func init() {
	initConfig()
	initLogger()
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(verString + " 当前时间：" + time.Now().Format("2006-01-02 15:04:05")))
}

func doLogHandler(w http.ResponseWriter, r *http.Request) {
	IP := r.RemoteAddr[:strings.LastIndexByte(r.RemoteAddr, ':')]
	if !isValidIP(IP) {
		return
	}

	r.ParseForm()

	app := r.PostFormValue("app")
	logrecs := r.PostFormValue("logs")

	recs := []string{}
	err := json.Unmarshal([]byte(logrecs), &recs)
	if err != nil {
		logger.Info("json err: %s", err.Error())
		return
	}

	err = writeLinesToFile(app+"/log.txt", recs)
	if err != nil {
		err = makeEmptyLogFile(app + "/log.txt")
		if err != nil {
			logger.Info("make file err: %s", err.Error())
			return
		}
		writeLinesToFile(app+"/log.txt", recs)
		if err != nil {
			logger.Info("write to file err: %s", err.Error())
			return
		}
	}
}

func winLogHandler(w http.ResponseWriter, r *http.Request) {
	IP := r.RemoteAddr[:strings.LastIndexByte(r.RemoteAddr, ':')]
	if !isValidIP(IP) {
		return
	}

	r.ParseForm()

	app := r.PostFormValue("app")
	logrecs := r.PostFormValue("logs")

	recs := []string{}
	err := json.Unmarshal([]byte(logrecs), &recs)
	if err != nil {
		logger.Info("json err: %s", err.Error())
		return
	}

	wes, err := parseWinLog(recs)
	if err != nil {
		logger.Info("parse win log err: %s", err.Error())
		return
	}

	err = writeLinesToFile(app+"/log.txt", wes)
	if err != nil {
		err = makeEmptyLogFile(app + "/log.txt")
		if err != nil {
			logger.Info("make file err: %s", err.Error())
			return
		}
		writeLinesToFile(app+"/log.txt", wes)
		if err != nil {
			logger.Info("write to file err: %s", err.Error())
			return
		}
	}
}

func parseWinLog(recs []string) ([]string, error) {
	wes := make([]string, 0)

	if recs != nil {
		l := len(recs)
		for i := 0; i < l; i++ {
			rec := recs[i]
			if strings.Contains(rec, "{") {
				rec = rec[strings.Index(rec, "{"):]
			}
			wes = append(wes, rec)
		}
	}

	return wes, nil
}

func isValidIP(IP string) bool {
	validIPs := cfg.ValidIPs
	validIPs = strings.TrimSpace(validIPs)
	if len(validIPs) == 0 {
		return true
	}
	return strings.Contains(validIPs, IP)
}

func makeEmptyLogFile(file string) error {
	if utils.IsFileExist(file) {
		return nil
	}

	err := utils.CreateFullDir(file)
	if err != nil {
		return err
	}
	err = utils.WriteToFile([]byte(""), file)
	if err != nil {
		logger.Info("错误:%s", err.Error())
	}
	return err
}

// fileName:文件名字(带全路径)
// lines: 写入的lines
func writeLinesToFile(fileName string, lines []string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		logger.Info("%s create failed. err: %s", fileName, err.Error())
	} else {
		// 查找文件末尾的偏移量
		f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		for i := 0; i < len(lines); i++ {
			fmt.Fprintf(f, "%s", lines[i])
		}
		// _, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}

func setServiceHandler() {
	http.HandleFunc("/test", testHandler)     //
	http.HandleFunc("/dolog", doLogHandler)   //
	http.HandleFunc("/winlog", winLogHandler) //
}

func main() {
	setServiceHandler()

	logger.Info("listenning port: %d", cfg.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil)
	if err != nil {
		logger.Error("Error:%s", err.Error())
	}
}
