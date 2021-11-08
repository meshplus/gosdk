package kvsql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/meshplus/gosdk/kvsql/mysql"
	"github.com/meshplus/gosdk/kvsql/util/hack"
	"math"
	"unsafe"
)

type buffer struct {
	pos int
	buf []byte
}

func (b *buffer) readInt1() byte {
	res := b.buf[b.pos]
	b.pos++
	return res
}

func (b *buffer) readInt2() uint16 {
	res := *(*uint16)(unsafe.Pointer(&b.buf[b.pos]))
	b.pos += 2
	return res
}

func (b *buffer) readInt3() uint32 {
	b1 := uint32(b.readInt1())
	b2 := uint32(b.readInt1())
	b3 := uint32(b.readInt1())
	return (b1 & 0xff) | ((b2 & 0xff) << 8) | ((b3 & 0xff) << 16)
}

func (b *buffer) readInt4() uint32 {
	res := *(*uint32)(unsafe.Pointer(&b.buf[b.pos]))
	b.pos += 4
	return res
}

func (b *buffer) readInt8() uint64 {
	res := *(*uint64)(unsafe.Pointer(&b.buf[b.pos]))
	b.pos += 8
	return res
}

func (b *buffer) readLENENC() uint64 {
	sw := b.readInt1()
	switch sw {
	case 251:
		return math.MaxUint64 // represents a NULL in a ProtocolText::ResultsetRow
	case 252:
		return uint64(b.readInt2())
	case 253:
		return uint64(b.readInt3())
	case 254:
		return b.readInt8()
	default:
		return uint64(sw)

	}
}

func (b *buffer) getBytes(len int) []byte {
	res := make([]byte, len)
	copy(res, b.buf[b.pos:b.pos+len])
	b.pos += len
	return res
}

func (b *buffer) readLenByteArray() []byte {
	length := b.readLENENC()
	if length == math.MaxUint64 {
		return nil
	}
	if length == 0 {
		return make([]byte, 0)
	}
	return b.getBytes(int(length))
}

type Field struct {
	collationIndex     uint16
	colType            uint8
	colFlag            uint16
	colDecimals        uint8
	colLength          uint32
	tableName          string
	originalTableName  string
	columnName         string
	originalColumnName string
}

// resultChunk is used to append chunks to chunk.
type resultChunk struct {
	chk *Chunk
}
type ResultSet struct {
	columnCount  uint32
	lastInsertID uint64
	updatedCount uint32

	rowNumber  int // the number of rows
	resChunk   *resultChunk
	columnInfo []*Field
}

func (c ResultSet) ColumnCount() uint32 {
	return c.columnCount
}

func (c ResultSet) LastInsertID() uint64 {
	return c.lastInsertID
}

func (c ResultSet) UpdatedCount() uint32 {
	return c.updatedCount
}

func (c ResultSet) Columns() []string {
	var names []string
	for _, v := range c.columnInfo {
		names = append(names, v.columnName)
	}
	return names
}

func (c ResultSet) RowNumber() int {
	return c.rowNumber
}

func (rs *ResultSet) ToExecuteResult() (*ExecuteResult, error) {
	answer := &ExecuteResult{
		ColumnCount:  rs.columnCount,
		LastInsertID: rs.lastInsertID,
		UpdatedCount: rs.updatedCount,
		RowNumber:    rs.rowNumber,
		Columns:      rs.Columns(),
		Result:       nil,
	}
	var result [][]interface{}
	for rowIndex := 0; rowIndex < rs.rowNumber; rowIndex++ {
		var data []interface{}
		for columnIndex := 0; columnIndex < int(rs.columnCount); columnIndex++ {
			if v, err := rs.getValue(rowIndex, columnIndex); err != nil {
				return nil, err
			} else {
				if _, ok := v.([]byte); ok {
					data = append(data, fmt.Sprintf("%v", v))
				} else {
					data = append(data, v)
				}
			}
		}
		result = append(result, data)
	}
	answer.Result = result
	return answer, nil
}

type ExecuteResult struct {
	ColumnCount  uint32          `json:"column_count"`
	LastInsertID uint64          `json:"last_insert_id"`
	UpdatedCount uint32          `json:"updated_count"`
	Columns      []string        `json:"columns"`
	RowNumber    int             `json:"row_number"`
	Result       [][]interface{} `json:"result"`
}

func (rs ExecuteResult) String() string {
	res, err := json.Marshal(rs)
	if err != nil {
		return ""
	}
	return string(res)
}

const (
	decodeVersion1 = iota
)

// DecodeRecordSet use this method to decode call ret to a ResultSet Object
func DecodeRecordSet(buf []byte) *ResultSet {
	buff := &buffer{
		buf: buf,
		pos: 0,
	}
	switch buff.readInt1() {
	case decodeVersion1:
		return decodeV1(buff)
	default:
		return nil
	}
}

func decodeV1(buff *buffer) *ResultSet {
	columnCount := buff.readInt4()
	if columnCount == 0 {
		return &ResultSet{
			columnCount:  columnCount,
			updatedCount: buff.readInt4(),
			lastInsertID: buff.readInt8(),
		}
	}

	return &ResultSet{
		columnCount:  columnCount,
		lastInsertID: 0,
		updatedCount: 0,
		columnInfo:   decodeColumnInfo(buff, int(columnCount)),
		rowNumber:    int(buff.readInt4()),
		resChunk:     decodeChunk(buff, int(columnCount)),
	}
}

func decodeColumnInfo(buff *buffer, columnCount int) []*Field {
	fields := make([]*Field, 0, columnCount)
	for i := 0; i < columnCount; i++ {
		f := &Field{
			tableName:          string(buff.readLenByteArray()),
			originalTableName:  string(buff.readLenByteArray()),
			columnName:         string(buff.readLenByteArray()),
			originalColumnName: string(buff.readLenByteArray()),
			collationIndex:     buff.readInt2(),
			colLength:          buff.readInt4(),
			colType:            buff.readInt1(),
			colFlag:            buff.readInt2(),
			colDecimals:        buff.readInt1(),
		}
		fields = append(fields, f)
	}
	return fields
}

func decodeChunk(buff *buffer, columnCount int) *resultChunk {
	columns := make([]*Column, 0, columnCount)
	for i := 0; i < columnCount; i++ {
		data := buff.readLenByteArray()
		nullBitMap := buff.readLenByteArray()
		length := buff.readLENENC()
		offsets := make([]int64, length)
		for j := uint64(0); j < length; j++ {
			offsets[j] = int64(buff.readInt4())
		}
		col := NewColumnWithData(data, nullBitMap, offsets)
		columns = append(columns, col)
	}
	return &resultChunk{chk: NewChunkWithColumn(columns)}
}

// NewChunkWithColumn create a empty chunk to resultChunk
func NewChunkWithColumn(columns []*Column) *Chunk {
	return &Chunk{
		columns: columns,
	}
}

func (rs *ResultSet) GetRow(rowIndex int) Row {
	if rowIndex >= rs.rowNumber || rowIndex < 0 {
		panic(fmt.Sprintf("row index out of range, need [0:%d] get %d", rs.rowNumber, rowIndex))
	}
	return rs.resChunk.chk.GetRow(rowIndex)
}

func (rs *ResultSet) getValue(rowIndex, columnIndex int) (interface{}, error) {
	row := rs.resChunk.chk.GetRow(rowIndex)
	field := rs.columnInfo[columnIndex]
	switch field.colType {
	case mysql.TypeDatetime, mysql.TypeTimestamp:
		if row.IsNull(columnIndex) {
			return "", nil
		}
		return row.GetTime(columnIndex).DateFormat("%Y-%m-%d %H:%i:%s")
	case mysql.TypeDate:
		if row.IsNull(columnIndex) {
			return "", nil
		}
		return row.GetTime(columnIndex).DateFormat("%Y-%m-%d")
	case mysql.TypeDuration:
		if row.IsNull(columnIndex) {
			return "", nil
		}
		return row.GetDuration(columnIndex, 0).String(), nil
	case mysql.TypeTiny:
		if row.IsNull(columnIndex) {
			return 0, nil
		}
		if (field.colFlag & 32) > 0 {
			return row.GetUint8(columnIndex), nil
		} else {
			return row.GetInt8(columnIndex), nil
		}
	case mysql.TypeYear:
		if row.IsNull(columnIndex) {
			return 0, nil
		}
		return row.GetUint16(columnIndex), nil
	case mysql.TypeShort:
		if row.IsNull(columnIndex) {
			return 0, nil
		}
		if (field.colFlag & 32) > 0 {
			return row.GetUint16(columnIndex), nil
		} else {
			return row.GetInt16(columnIndex), nil
		}
	case mysql.TypeInt24:
		if row.IsNull(columnIndex) {
			return 0, nil
		}
		return row.GetInt32(columnIndex), nil
	case mysql.TypeLong:
		if row.IsNull(columnIndex) {
			return 0, nil
		}
		if (field.colFlag & 32) > 0 {
			return row.GetUint32(columnIndex), nil
		} else {
			return row.GetInt32(columnIndex), nil
		}
	case mysql.TypeLonglong:
		if row.IsNull(columnIndex) {
			return 0, nil
		}
		if (field.colFlag & 32) > 0 {
			return row.GetUint64(columnIndex), nil
		} else {
			return row.GetInt64(columnIndex), nil
		}
	case mysql.TypeFloat:
		if row.IsNull(columnIndex) {
			return 0.0, nil
		}
		return row.GetFloat32(columnIndex), nil
	case mysql.TypeDouble:
		if row.IsNull(columnIndex) {
			return 0.0, nil
		}
		return row.GetFloat64(columnIndex), nil
	case mysql.TypeVarString, mysql.TypeVarchar, mysql.TypeString, mysql.TypeSet, mysql.TypeEnum,
		mysql.TypeNewDecimal, mysql.TypeJSON, mysql.TypeBlob, mysql.TypeLongBlob, mysql.TypeMediumBlob,
		mysql.TypeTinyBlob, mysql.TypeBit:
		if row.IsNull(columnIndex) {
			return "", nil
		}
		// note: will split by zero byte
		data := row.GetBytes(columnIndex)
		if i := bytes.IndexByte(data, 0); i != -1 {
			return hack.String(data[:i]), nil
		}
		return row.GetString(columnIndex), nil
	//case mysql.TypeBit:
	//	if row.IsNull(columnIndex) {
	//		return []byte(nil), nil
	//	}
	//	return row.GetBytes(columnIndex), nil
	default:
		return nil, fmt.Errorf("not support type %d", field.colType)
	}
}
