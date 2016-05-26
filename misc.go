package iphelper

import (
	"strconv"
	"strings"
)

const (
	HEADER_LENGTH   = 4
	BODYLINE_LENGTH = 20
)

const (
	AREA_COUNTRY  = "country"
	AREA_PROVINCE = "province"
	AREA_CITY     = "city"
	AREA_ZONE     = "zone"
	AREA_LOCATION = "location"
	AREA_OPERATOR = "operator"
)

func IP2Num(requestip string) uint64 {
	//获取客户端地址的long
	nowip := strings.Split(requestip, ".")
	if len(nowip) != 4 {
		return 0
	}
	a, _ := strconv.ParseUint(nowip[0], 10, 64)
	b, _ := strconv.ParseUint(nowip[1], 10, 64)
	c, _ := strconv.ParseUint(nowip[2], 10, 64)
	d, _ := strconv.ParseUint(nowip[3], 10, 64)
	ipNum := a<<24 | b<<16 | c<<8 | d
	return ipNum
}

func Num2IP(ipnum uint64) string {
	byte1 := ipnum & 0xff
	byte2 := (ipnum & 0xff00)
	byte2 >>= 8
	byte3 := (ipnum & 0xff0000)
	byte3 >>= 16
	byte4 := (ipnum & 0xff000000)
	byte4 >>= 24
	result := strconv.FormatUint(byte4, 10) + "." +
		strconv.FormatUint(byte3, 10) + "." +
		strconv.FormatUint(byte2, 10) + "." +
		strconv.FormatUint(byte1, 10)
	return result
}
