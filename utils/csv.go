package utils

import (
	"encoding/csv"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type CloudflareIPData struct {
	ip            net.IPAddr
	pingCount     int
	pingReceived  int
	recvRate      float32
	downloadSpeed float32
	pingTime      time.Duration
}

func (cf *CloudflareIPData) getRecvRate() float32 {
	if cf.recvRate == 0 {
		pingLost := cf.pingCount - cf.pingReceived
		cf.recvRate = float32(pingLost) / float32(cf.pingCount)
	}
	return cf.recvRate
}

func ExportCsv(filePath string, data []CloudflareIPData) {
	fp, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("创建文件[%s]失败：%v", filePath, err)
		return
	}
	defer fp.Close()
	w := csv.NewWriter(fp) //创建一个新的写入文件流
	w.Write([]string{"IP 地址", "已发送", "已接收", "丢包率", "平均延迟", "下载速度 (MB/s)"})
	w.WriteAll(convertToString(data))
	w.Flush()
}

func (cf *CloudflareIPData) toString() []string {
	result := make([]string, 6)
	result[0] = cf.ip.String()
	result[1] = strconv.Itoa(cf.pingCount)
	result[2] = strconv.Itoa(cf.pingReceived)
	result[3] = strconv.FormatFloat(float64(cf.getRecvRate()), 'f', 2, 32)
	result[4] = cf.pingTime.String()
	result[5] = strconv.FormatFloat(float64(cf.downloadSpeed)/1024/1024, 'f', 2, 32)
	return result
}

func convertToString(data []CloudflareIPData) [][]string {
	result := make([][]string, 0)
	for _, v := range data {
		result = append(result, v.toString())
	}
	return result
}
