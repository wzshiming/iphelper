package iphelper

import (
	"fmt"
	"strconv"
)

// 获取ip段信息中文名称
type IpRowCn struct {
	Country  string
	Province string
	City     string
	Zone     string
	Location string
	Operator string
}

//func (this *IpRowCn) String() string {
//	return this.Country + this.Province + this.City + this.Zone + this.Location + this.Operator
//}

// 获取ip段信息
type IpRow struct {
	store    *IpStore
	start    uint32
	end      uint32
	Country  uint16
	Province uint16
	City     uint16
	Zone     uint16
	Location uint16
	Operator uint16
}

func (this *IpRow) Cn() *IpRowCn {
	return &IpRowCn{
		Country:  this.store.metaTable[AREA_COUNTRY][this.Country],
		Province: this.store.metaTable[AREA_PROVINCE][this.Province],
		City:     this.store.metaTable[AREA_CITY][this.City],
		Zone:     this.store.metaTable[AREA_ZONE][this.Zone],
		Location: this.store.metaTable[AREA_LOCATION][this.Location],
		Operator: this.store.metaTable[AREA_OPERATOR][this.Operator],
	}
}

func (this *IpRow) String() string {
	countryCode := strconv.Itoa(int(this.Country))
	provinceCode := fmt.Sprintf("%04d", this.Province)
	cityCode := fmt.Sprintf("%04d", this.City)
	zoneCode := fmt.Sprintf("%04d", this.Zone)
	provoderCode := fmt.Sprintf("%02d", this.Location)
	OperatorCode := fmt.Sprintf("%02d", this.Operator)
	return countryCode + provinceCode + cityCode + zoneCode + provoderCode + OperatorCode
}

func (this *IpRow) Code() (code uint64) {
	fmt.Sscan(this.String(), &code)
	return
}
