package iphelper

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
)

type IpStore struct {
	bodyLength   uint32
	metaLength   uint32
	headerBuffer []byte
	bodyBuffer   []byte
	metaBuffer   []byte
	IpTableMap   map[uint64]*IpRow
	IpTable      []*IpRow // ip信息表 按范围自增
	metaTable    map[string][]string
}

func NewIpStore(filename string) *IpStore {
	store := IpStore{
		headerBuffer: make([]byte, HEADER_LENGTH),
		metaTable:    make(map[string][]string),
		IpTableMap:   make(map[uint64]*IpRow),
	}
	store.parseStore(filename)
	return &store
}

// 获取ip的区域信息列表
func (this *IpStore) GetMetaTable() map[string][]string {
	return this.metaTable
}

func (this *IpStore) GetGeoByGeocode(areacode uint64) (*IpRow, error) {
	var row, ok = this.IpTableMap[areacode]
	if !ok {
		return nil, errors.New("fail to find")
	}
	return row, nil
}

// 获取ip所在ip段的信息
func (this *IpStore) GetGeoByIp(ipSearch string) (row *IpRow, err error) {
	search := uint32(IP2Num(ipSearch))
	// fmt.Println(search)
	var start uint32 = 0
	var end uint32 = uint32(len(this.IpTable) - 1)
	var offset uint32 = 0
	for start <= end {
		mid := uint32(math.Floor(float64((end - start) / 2)))
		offset = start + mid
		IpRow := this.IpTable[offset]
		// fmt.Println(IpRow)
		if search >= IpRow.start {
			if search <= IpRow.end {
				return IpRow, nil
			} else {
				start = offset + 1
				continue
			}
		} else {
			end = offset - 1
			continue
		}
	}
	return row, errors.New("fail to find")
}

func (this *IpStore) parseStore(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic("error opening file: %v\n" + err.Error())
	}
	defer file.Close()
	fmt.Println("open file: ", filename)
	var buf [HEADER_LENGTH]byte

	if _, err := file.Read(buf[0:4]); err != nil {
		panic("error read header" + err.Error())
	}

	this.bodyLength = binary.BigEndian.Uint32(buf[0:4])
	fmt.Println("body length is: ", this.bodyLength)
	if _, err := file.Read(buf[0:4]); err != nil {
		panic("error read header" + err.Error())
	}
	this.metaLength = binary.BigEndian.Uint32(buf[0:4])
	fmt.Println("meta length is: ", this.metaLength)
	if err := this.paseBody(file); err != nil {
		panic("parse body  failed:" + err.Error())
	}

	if err := this.parseMeta(file); err != nil {
		panic("pase meta failed" + err.Error())
	}
}

func (this *IpStore) paseBody(file *os.File) error {
	this.bodyBuffer = make([]byte, this.bodyLength)
	if _, err := file.ReadAt(this.bodyBuffer, HEADER_LENGTH+HEADER_LENGTH); err != nil {
		panic("read body error")
	}
	buf := bytes.NewBuffer(this.bodyBuffer)
	var offset uint32 = 0
	for offset < this.bodyLength {
		line := buf.Next(BODYLINE_LENGTH)
		row, err := this.parseBodyLine(line)
		if err != nil {
			return err
		}
		this.IpTableMap[row.Code()] = row
		this.IpTable = append(this.IpTable, row)
		offset += BODYLINE_LENGTH
	}
	return nil
}

func (this *IpStore) parseMeta(file *os.File) (err error) {
	this.metaBuffer = make([]byte, this.metaLength)
	if _, err := file.ReadAt(this.metaBuffer, int64(HEADER_LENGTH+HEADER_LENGTH+this.bodyLength)); err != nil {
		panic("read meta error")
	}
	return json.Unmarshal(this.metaBuffer, &this.metaTable)
}

// @TODO Parse by Reflect IpRow
func (this *IpStore) parseBodyLine(buffer []byte) (row *IpRow, err error) {
	row = &IpRow{}
	buf := bytes.NewBuffer(buffer)
	if err = binary.Read(buf, binary.BigEndian, &row.start); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.end); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.Country); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.Province); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.City); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.Zone); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.Location); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.Operator); err != nil {
		goto fail
	}
fail:
	row.store = this
	return row, err
}
