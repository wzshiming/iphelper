package iphelper

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type datFile struct {
	err error
	*bytes.Buffer
	headerLength int
	bodyLength   int
	geoMap       map[string]map[string]uint16
	geoSlice     map[string][]string
	operator     map[string]int
	writer       io.Writer
}

func NewDatFile(w io.Writer) *datFile {
	m := map[string]map[string]uint16{
		AREA_COUNTRY:  make(map[string]uint16),
		AREA_PROVINCE: make(map[string]uint16),
		AREA_CITY:     make(map[string]uint16),
		AREA_ZONE:     make(map[string]uint16),
		AREA_LOCATION: make(map[string]uint16),
		AREA_OPERATOR: make(map[string]uint16),
	}
	return &datFile{
		Buffer:   bytes.NewBuffer(nil),
		geoMap:   m,
		geoSlice: make(map[string][]string),
		writer:   bufio.NewWriter(w),
	}
}

// get area code by typ
func (d *datFile) getCode(typ string, area string) uint16 {
	var code uint16
	code, ok := d.geoMap[typ][area]
	if !ok {
		code = uint16(len(d.geoMap[typ]))
		d.geoMap[typ][area] = code
		d.geoSlice[typ] = append(d.geoSlice[typ], area)
	}
	return code
}

// @TODO parse fields by reflect the ip row
func (d *datFile) writeBody(fields []string) error {
	if d.err != nil {
		return d.err
	}
	start, _ := strconv.ParseUint(fields[0], 10, 32)
	end, _ := strconv.ParseUint(fields[1], 10, 32)
	binary.Write(d, binary.BigEndian, uint32(start))
	binary.Write(d, binary.BigEndian, uint32(end))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_COUNTRY, fields[2]))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_PROVINCE, fields[3]))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_CITY, fields[4]))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_ZONE, fields[5]))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_LOCATION, fields[6]))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_OPERATOR, fields[7]))
	return d.err
}

// bodylength|body|metalength|meta
func (d *datFile) writeFile() error {
	if d.err != nil {
		return d.err
	}

	bodyLength := d.Buffer.Len()
	meta, err := json.Marshal(d.geoSlice)
	if err != nil {
		d.err = err
		return d.err
	}
	metaLength := len(meta)

	binary.Write(d.writer, binary.BigEndian, uint32(bodyLength))
	binary.Write(d.writer, binary.BigEndian, uint32(metaLength))
	d.writer.Write(d.Buffer.Bytes())
	d.writer.Write(meta)

	fmt.Println("meta length is: ", metaLength)
	fmt.Println("body length is: ", bodyLength)
	return err
}

func MakeDat(infile, outfile string) error {
	in, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 755)
	if err != nil {
		return err
	}
	defer out.Close()
	output := NewDatFile(out)
	r := bufio.NewReader(in)
	count := 0
	for {
		count++
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if len(line) != 0 {
			fields := strings.Fields(line)
			if len(fields) != 8 {
				return errors.New("invalid input file invalid line string")
			}
			if err := output.writeBody(fields); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
	}
	if err := output.writeFile(); err != nil {
		return err
	}
	fmt.Println("amount ip range from ip source: ", count)
	return nil
}
