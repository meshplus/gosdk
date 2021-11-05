package hvm

import (
	"github.com/meshplus/gosdk/classfile"
	"github.com/syndtr/goleveldb/leveldb"
	"strings"
)

const (
	HyperListClassRef  = "cn/hyperchain/core/HyperList"
	HyperListFieldRef  = "Lcn/hyperchain/core/HyperList;"
	HyperMapClassRef   = "cn/hyperchain/core/HyperMap"
	HyperMapFieldRef   = "Lcn/hyperchain/core/HyperMap;"
	HyperTableClassRef = "cn/hyperchain/core/HyperTable"
	HyperTableFieldRef = "Lcn/hyperchain/core/HyperTable;"

	IntegerJ   = "Ljava/lang/Integer;"
	BooleanJ   = "Ljava/lang/Boolean;"
	ByteJ      = "Ljava/lang/Byte;"
	CharacterJ = "Ljava/lang/Character;"
	DoubleJ    = "Ljava/lang/Double;"
	FloatJ     = "Ljava/lang/Float;"
	LongJ      = "Ljava/lang/Long;"
	ShortJ     = "Ljava/lang/Short;"
	StringJ    = "Ljava/lang/String;"

	NULL          = "null"
	storagePrefix = "-storage"
	addrLen       = 20
)

var (
	NameAndType = []string{"name", "type", "value"}
	Contract    *ObjectType
	ClassMap    map[string]*classfile.ClassFile
	MainClass   string
	DBStorage   *leveldb.DB
)

type HvmLog struct {
	Timestamp uint64
	Type      string
	Body      Body
}

type Body struct {
	Type string
	Name string
	Data map[string]map[string][]byte
}

type KVTemplate struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type FieldType interface {
	//Name() string
	//Class() string
	//Value() string
}

type CommonType struct {
	Name    string
	Class   string
	Columns []string
	Value   string
}

type BaseType struct {
	CommonType
}

type ArrayType struct {
	CommonType
	Dimension int
	Fields    []FieldType
}

type ObjectType struct {
	CommonType
	Fields []FieldType
}

type HyperMap struct {
	CommonType
	KFields []FieldType
	VFields []FieldType
}

type HyperTable struct {
	CommonType
	Columns map[string]map[string]bool
	Items   []TableItem
}

type TableItem struct {
	Row          string
	ColumnFamily string
	Column       string
	Value        string
}

type HyperList struct {
	CommonType
	mapping *ListData
	Fields  []FieldType
}

type ListData struct {
	Table []uint64
	Count uint64
}

//NewFieldType create FieldType with class name
func NewFieldType(class string) FieldType {
	commonType := CommonType{}
	if strings.HasPrefix(class, "[") {
		commonType.Class = class
		arrayType := &ArrayType{}
		arrayType.CommonType = commonType
		return arrayType
	} else if strings.Compare(class, HyperListClassRef) == 0 || strings.Compare(class, HyperListFieldRef) == 0 {
		commonType.Class = class[1 : len(class)-1]
		hyperList := &HyperList{}
		hyperList.mapping = &ListData{}
		hyperList.CommonType = commonType
		return hyperList
	} else if strings.Compare(class, HyperMapClassRef) == 0 || strings.Compare(class, HyperMapFieldRef) == 0 {
		commonType.Class = class[1 : len(class)-1]
		hyperMap := &HyperMap{}
		hyperMap.CommonType = commonType
		return hyperMap
	} else if strings.Compare(class, HyperTableClassRef) == 0 || strings.Compare(class, HyperTableFieldRef) == 0 {
		commonType.Class = class[1 : len(class)-1]
		hyperTable := &HyperTable{}
		hyperTable.Columns = make(map[string]map[string]bool)
		hyperTable.CommonType = commonType
		return hyperTable
	} else if IsBaseType(class) {
		if IsPrimitiveType(class) {
			commonType.Class = class
		} else {
			commonType.Class = class[1 : len(class)-1]
		}
		baseType := &BaseType{}
		baseType.CommonType = commonType
		return baseType
	} else {
		commonType.Class = class[1 : len(class)-1]
		objectType := &ObjectType{}
		objectType.CommonType = commonType
		return objectType
	}
}

func IsBaseType(class string) bool {
	switch class {
	case "I", IntegerJ:
		return true
	case "Z", BooleanJ:
		return true
	case "B", ByteJ:
		return true
	case "C", CharacterJ:
		return true
	case "D", DoubleJ:
		return true
	case "F", FloatJ:
		return true
	case "J", LongJ:
		return true
	case "S", ShortJ:
		return true
	case StringJ:
		return true
	default:
		return false
	}
}

func IsPrimitiveType(class string) bool {
	switch class {
	case "I":
		return true
	case "Z":
		return true
	case "B":
		return true
	case "C":
		return true
	case "D":
		return true
	case "F":
		return true
	case "J":
		return true
	case "S":
		return true
	default:
		return false
	}
}
