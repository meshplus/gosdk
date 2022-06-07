package hvm

import (
	"bytes"
	"encoding/hex"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/meshplus/gosdk/classfile"
	"github.com/opentracing/opentracing-go/log"
	"github.com/syndtr/goleveldb/leveldb"
	iterator2 "github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func InitContract(contractCode []byte) {
	ClassMap = make(map[string]*classfile.ClassFile)
	parseContractJar(contractCode)

	Contract = NewFieldType(" " + MainClass + " ").(*ObjectType)
	Contract.Fields = make([]FieldType, 0)
}

func parseContractJar(contractJar []byte) {
	mainClassLen := BytesToInt32(contractJar[0:2])
	MainClass = string(contractJar[2 : 2+mainClassLen])
	start := 2 + int(mainClassLen)

	for start < len(contractJar) {
		classLen := int(BytesToInt32(contractJar[start : start+4]))
		start += 4
		classNameLen := int(BytesToInt32(contractJar[start : start+2]))
		start += 2

		class := contractJar[start : start+classLen]
		start += classLen
		className := string(contractJar[start : start+classNameLen])
		start += classNameLen

		classFile, err := classfile.Parse(class)
		if err != nil {
			panic(err)
		}

		ClassMap[className] = classFile
	}
}

func ParseHistoryData(contract *ObjectType, classFile *classfile.ClassFile, contractAddress string, dbPath string) {
	ParseContractData(contract, classFile, contractAddress, true, dbPath, nil)
}

func ParseIncrementData(contract *ObjectType, classFile *classfile.ClassFile, contractAddress, mqLog string) {
	hvmLog := &HvmLog{}
	err := jsoniter.UnmarshalFromString(mqLog, hvmLog)
	if err != nil {
		fmt.Println(err)
		return
	}

	ParseContractData(contract, classFile, contractAddress, false, "", hvmLog.Body.Data)
}

func ParseContractData(contract *ObjectType, classFile *classfile.ClassFile, contractAddress string, db bool, dbPath string, source map[string]map[string][]byte) {
	contract.Fields = make([]FieldType, 0)
	contract.Columns = make([]string, 0)
	//nolint
	datas := make([][]byte, 0)
	var kvs []*KVTemplate
	var storageCnt = 0

	contractAddress = strings.TrimPrefix(contractAddress, "0x")

	if db {
		if DBStorage == nil {
			DBStorage = OpenDB(dbPath)
		}
	} else {
		kvs = processIncrementedDatas(contractAddress, source)
	}

	for _, field := range classFile.Fields {
		if annotationAttr := field.RuntimeVisibleAnnotationsAttributeData(); annotationAttr == nil {
			continue
		} else {
			flag := true
			for _, annotation := range annotationAttr.Annotations() {
				if annotation.Type() == "Lcn/hyperchain/annotations/StoreField;" {
					flag = false
					break
				}
			}
			if flag {
				continue
			}
		}

		descriptor := field.Descriptor()
		contract.Fields = append(contract.Fields, NewFieldType(descriptor))
		contract.Columns = append(contract.Columns, field.Name())

		// Process data
		if db {
			datas = processDatas(contractAddress, field.Name(), descriptor)
		} else {
			datas = processNewDatas(field.Name(), kvs, descriptor)
		}

		if len(datas) != 0 {
			parseContractData(contract.Fields[storageCnt], field, datas...)
		}
		storageCnt++
	}
}

func parseContractData(field FieldType, memberInfo classfile.MemberInfo, datas ...[]byte) {
	fieldName := memberInfo.Name()
	iterator := jsoniter.NewIterator(jsoniter.ConfigDefault)
	iterator.ResetBytes([]byte(datas[0]))
	any := iterator.ReadAny()

	switch field.(type) {
	case *ObjectType:
		objectType := field.(*ObjectType)
		objectType.Name = fieldName
		data := any.Get("value").ToString()
		if newClassFile, ok := ClassMap[objectType.Class]; ok {
			parseObjectData(objectType, newClassFile, data)
		} else {
			objectType.Value = data
		}
	case *ArrayType:
		arrayType := field.(*ArrayType)
		arrayType.Name = fieldName
		data := any.Get("value").ToString()
		parseArrayData(arrayType, data)
	case *BaseType:
		baseType := field.(*BaseType)
		baseType.Name = fieldName
		baseType.Columns = NameAndType
		data := any.Get("value").ToString()
		baseType.Value = data
	case *HyperList:
		hyperList := field.(*HyperList)
		hyperList.Name = fieldName
		//fmt.Println("开始解析HyperList", hyperList.Name)
		parseHyperListData(hyperList, memberInfo, datas)
		//fmt.Println("结束解析HyperList", hyperList.Name)
	case *HyperMap:
		hyperMap := field.(*HyperMap)
		hyperMap.Name = fieldName
		//fmt.Println("开始解析HyperMap", hyperMap.Name)
		parseHyperMapData(hyperMap, memberInfo, datas)
		//fmt.Println("结束解析HyperMap", hyperMap.Name)
	case *HyperTable:
		hyperTable := field.(*HyperTable)
		hyperTable.Name = fieldName
		parseHyperTableData(hyperTable, datas)
	default:
		fmt.Println("error type")
		return
	}
}

func parseObjectData(object *ObjectType, classFile *classfile.ClassFile, data string) {
	if data == "" || data == NULL {
		object.Value = NULL
		return
	}

	object.Value = data
	iterator := jsoniter.NewIterator(jsoniter.ConfigDefault)
	iterator.ResetBytes([]byte(data))
	any := iterator.ReadAny()

	for i, field := range classFile.Fields {
		descriptor := field.Descriptor()
		object.Fields = append(object.Fields, NewFieldType(descriptor))
		object.Columns = append(object.Columns, field.Name())

		switch object.Fields[i].(type) {
		case *ObjectType:
			objectType := object.Fields[i].(*ObjectType)
			objectType.Name = field.Name()
			innerData := any.Get(objectType.Name).ToString()
			if newClassFile, ok := ClassMap[objectType.Class]; ok {
				parseObjectData(objectType, newClassFile, innerData)
			} else {
				objectType.Value = innerData
			}
		case *ArrayType:
			arrayType := object.Fields[i].(*ArrayType)
			arrayType.Name = field.Name()
			innerData := any.Get(arrayType.Name).ToString()
			parseArrayData(arrayType, innerData)
		case *BaseType:
			baseType := object.Fields[i].(*BaseType)
			baseType.Name = field.Name()
			innerData := any.Get(baseType.Name).ToString()
			baseType.Value = innerData
		default:
			fmt.Println("error type")
			return
		}
	}
}

func parseArrayData(array *ArrayType, data string) {
	if data == "" {
		array.Value = NULL
		return
	}

	iterator := jsoniter.NewIterator(jsoniter.ConfigDefault)
	iterator.ResetBytes([]byte(data))
	any := iterator.ReadAny()

	array.Value = data
	array.Dimension = strings.Count(array.Class, "[")
	array.Fields = make([]FieldType, 0)
	array.Columns = make([]string, 0)

	componentClass := array.Class[1:len(array.Class)]

	realLen := any.Size()
	for i := 0; i < realLen; i++ {
		array.Fields = append(array.Fields, NewFieldType(componentClass))
		array.Columns = append(array.Columns, strconv.Itoa(i))
	}

	if array.Dimension > 1 {
		for i, field := range array.Fields {
			arrayType := field.(*ArrayType)
			arrayType.Value = any.Get(i).ToString()
		}
		return
	}

	for i, field := range array.Fields {
		componentData := any.Get(i).ToString()
		switch array.Fields[i].(type) {
		case *ObjectType:
			objectType := field.(*ObjectType)
			if newClassFile, ok := ClassMap[objectType.Class]; ok {
				parseObjectData(objectType, newClassFile, componentData)
			} else {
				objectType.Value = componentData
			}
		case *BaseType:
			baseType := field.(*BaseType)
			baseType.Value = componentData
			baseType.Columns = NameAndType
		default:
			fmt.Println("error type")
			return
		}
	}
}

func parseHyperListData(hyperList *HyperList, memberInfo classfile.MemberInfo, datas [][]byte) {
	if len(datas) <= 1 {
		return
	}
	signature := memberInfo.Signature()
	hyperList.Value = string(bytes.Join(datas, []byte(",")))
	hyperList.Fields = make([]FieldType, 0)
	var generics string

	if signature != "" {
		re, err := regexp.Compile("<.*>")
		if err != nil {
			fmt.Println(err)
			return
		}
		generics = re.FindString(signature)
		generics = generics[1 : len(generics)-1]
		switch NewFieldType(generics).(type) {
		case *ObjectType:
			newClassFile, ok := ClassMap[generics]
			if ok {
				hyperList.Columns = make([]string, len(newClassFile.Fields)+1)
				for i, field := range newClassFile.Fields {
					hyperList.Columns[i+1] = field.Name()
				}
				hyperList.Columns[0] = "index"
			}
		case *BaseType:
			columns := make([]string, 0)
			columns = append(columns, "index")
			columns = append(columns, NameAndType...)
			hyperList.Columns = columns
		}
	}

	iterator := jsoniter.NewIterator(jsoniter.ConfigDefault)
	iterator.ResetBytes(datas[0])
	any := iterator.ReadAny()
	listData := any.Get("value").ToString()
	err := jsoniter.Unmarshal([]byte(listData), hyperList.mapping)
	if err != nil {
		fmt.Println(err)
	}

	datas = datas[1:] // table, kv

	mapData := make(map[string]string)
	for _, data := range datas {
		iterator.ResetBytes(data)
		any = iterator.ReadAny()
		key := any.Get("key").ToString()
		value := any.Get("value").ToString()
		mapData[key] = value
	}

	for i := range datas {
		hyperList.Fields = append(hyperList.Fields, NewFieldType(generics))
		innerData := mapData[hyperList.Name+"@"+strconv.FormatUint(hyperList.mapping.Table[i], 10)]

		switch hyperList.Fields[i].(type) {
		case *ObjectType:
			objectType := hyperList.Fields[i].(*ObjectType)
			if newClassFile, ok := ClassMap[objectType.Class]; ok {
				parseObjectData(objectType, newClassFile, innerData)
			} else {
				objectType.Value = innerData
			}
		case *BaseType:
			baseType := hyperList.Fields[i].(*BaseType)
			baseType.Value = innerData
			baseType.Columns = NameAndType
		default:
			fmt.Println("error type")
			return
		}
	}

}

// parseHyperMapData datas format: [{key:"", value:""}]
func parseHyperMapData(hyperMap *HyperMap, memberInfo classfile.MemberInfo, datas [][]byte) {
	if len(datas) == 0 {
		return
	}

	hyperMap.Value = string(bytes.Join(datas, []byte(",")))
	signature := memberInfo.Signature()
	var kclass string
	var vclass string
	var columns []string

	if signature != "" {
		re, err := regexp.Compile("<.*>")
		if err != nil {
			fmt.Println(err)
			return
		}
		generics := strings.Split(re.FindString(signature), ";")
		//fmt.Println("generics:", generics)
		kclass = generics[0][1:] + ";"
		vclass = generics[1][:] + ";"
		columns = make([]string, 0)
		columns = append(columns, "key")

		switch NewFieldType(vclass).(type) {
		case *ObjectType:
			objectType := NewFieldType(vclass).(*ObjectType)
			if newClassFile, ok := ClassMap[objectType.Class]; ok {
				for _, field := range newClassFile.Fields {
					columns = append(columns, field.Name())
				}
				hyperMap.Columns = columns
			}
		case *BaseType:
			baseType := NewFieldType(vclass).(*BaseType)
			columns = append(columns, NameAndType...)
			hyperMap.Columns = columns
			baseType.Columns = NameAndType
		default:
			fmt.Println("error type")
			return
		}
	}

	iterator := jsoniter.NewIterator(jsoniter.ConfigDefault)

	for i, data := range datas {
		hyperMap.KFields = append(hyperMap.KFields, NewFieldType(kclass))
		hyperMap.VFields = append(hyperMap.VFields, NewFieldType(vclass))

		iterator.ResetBytes(data)
		any := iterator.ReadAny()
		key := any.Get("key").ToString()
		value := any.Get("value").ToString()

		switch hyperMap.KFields[i].(type) {
		case *ObjectType:
			objectType := hyperMap.KFields[i].(*ObjectType)
			objectType.Name = memberInfo.Name()
			objectType.Value = key
		case *BaseType:
			baseType := hyperMap.KFields[i].(*BaseType)
			baseType.Name = memberInfo.Name()
			baseType.Value = key
			baseType.Columns = NameAndType
		default:
			fmt.Println("error type")
			return
		}

		switch hyperMap.VFields[i].(type) {
		case *ObjectType:
			objectType := hyperMap.VFields[i].(*ObjectType)
			objectType.Name = memberInfo.Name()
			if newClassFile, ok := ClassMap[objectType.Class]; ok {
				parseObjectData(objectType, newClassFile, value)
			} else {
				objectType.Value = value
			}
		case *BaseType:
			baseType := hyperMap.VFields[i].(*BaseType)
			baseType.Name = memberInfo.Name()
			baseType.Value = value
			baseType.Columns = NameAndType
		default:
			fmt.Println("error type")
			return
		}
	}
}

func parseHyperTableData(hyperTable *HyperTable, datas [][]byte) {
	if len(datas) == 0 {
		return
	}

	hyperTable.Value = string(bytes.Join(datas, []byte(",")))
	hyperTable.Items = make([]TableItem, len(datas))
	hyperTable.Value = string(bytes.Join(datas, []byte(",")))
	//fmt.Println("实际value:", hyperTable.Value)
	iterator := jsoniter.NewIterator(jsoniter.ConfigDefault)

	for i, data := range datas {
		iterator.ResetBytes(data)
		any := iterator.ReadAny()
		key := any.Get("key").ToString()
		value := any.Get("value").ToString()
		items := strings.Split(key, "@")
		//fmt.Println("items:", items)
		// table | row | colf | col
		hyperTable.Items[i].Row = items[1]
		hyperTable.Items[i].ColumnFamily = items[2]
		hyperTable.Items[i].Column = items[3]
		hyperTable.Items[i].Value = value
		if columns, ok := hyperTable.Columns[items[2]]; ok {
			if _, ok := columns[items[3]]; !ok {
				columns[items[3]] = true
				hyperTable.Columns[items[2]] = columns
			}
		} else {
			hyperTable.Columns[items[2]] = make(map[string]bool)
			hyperTable.Columns[items[2]][items[3]] = true
		}
	}

	//fmt.Println("Columns:", hyperTable.Columns)
}

func OpenDB(path string) *leveldb.DB {
	newDB, err := leveldb.OpenFile(path, GetLdbConfig())
	if err != nil {
		log.Error(err)
		return nil
	}
	return newDB
}

func processDatas(contractAddr string, fieldName string, class string) [][]byte {
	addr, err := hex.DecodeString(contractAddr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	prefix := append([]byte(storagePrefix), addr...)
	prefix = append(prefix, []byte(fieldName)...)
	storageLen := len([]byte(storagePrefix))
	result := make([][]byte, 0)

	var iterator iterator2.Iterator
	switch NewFieldType(class).(type) {
	case *BaseType, *ArrayType, *ObjectType:
		rge := &util.Range{Start: prefix, Limit: append(prefix, []byte("-")...)}
		iterator = DBStorage.NewIterator(rge, nil)
	case *HyperList:
		key := append(prefix, []byte("-__table__")...)
		value, err := DBStorage.Get(key, nil)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		kv := &KVTemplate{Key: string(key[storageLen+addrLen:]), Value: string(value)}
		kvBytes, err := jsoniter.Marshal(kv)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		result = append(result, kvBytes)
		rge := util.BytesPrefix(append(prefix, []byte("@")...))
		iterator = DBStorage.NewIterator(rge, nil)
	case *HyperMap, *HyperTable:
		rge := util.BytesPrefix(append(prefix, []byte("@")...))
		iterator = DBStorage.NewIterator(rge, nil)
	default:
		fmt.Println("error type")
		return nil
	}
	//iterator = DBStorage.NewIterator(util.BytesPrefix(prefix), nil)
	for iterator.Next() {
		key := string(iterator.Key()[storageLen+addrLen:])
		var value string
		switch NewFieldType(class).(type) {
		case *HyperTable:
			value = string(iterator.Value())
		default:
			value = string(iterator.Value()[:len(iterator.Value())-1])
		}

		kv := &KVTemplate{Key: key, Value: value}
		kvBytes, err := jsoniter.Marshal(kv)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		result = append(result, kvBytes)
	}
	//fmt.Println("field name:", fieldName)
	//fmt.Println("result length", len(result))
	//fmt.Println("result:", string(bytes.Join(result, []byte(","))))
	return result
}

func processIncrementedDatas(contractAddr string, source map[string]map[string][]byte) []*KVTemplate {
	datas := source[contractAddr]
	keys := make([]string, 0)
	sortDatas := make([]*KVTemplate, 0)

	for field := range datas {
		keys = append(keys, field)
	}

	sort.Strings(keys)
	for _, k := range keys {
		key := k
		value := string(datas[k])
		kv := &KVTemplate{Key: key, Value: value}
		sortDatas = append(sortDatas, kv)
	}
	return sortDatas
}

func processNewDatas(fieldName string, sortDatas []*KVTemplate, class string) [][]byte {
	result := make([][]byte, 0)
	var version = true

	switch NewFieldType(class).(type) {
	case *BaseType, *ArrayType, *ObjectType:
		for _, kv := range sortDatas {
			if kv.Key == fieldName {
				kv.Value = string([]byte(kv.Value)[:len([]byte(kv.Value))-1])
				kvBytes, err := jsoniter.Marshal(kv)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				result = append(result, kvBytes)
				return result
			}
		}
	case *HyperList:
		for _, kv := range sortDatas {
			if kv.Key == fieldName+"-__table__" {
				kv.Value = string([]byte(kv.Value)[:len([]byte(kv.Value))-1])
				kvBytes, err := jsoniter.Marshal(kv)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				result = append(result, kvBytes)
				break
			}
		}
	case *HyperMap:
	case *HyperTable:
		version = false
	default:
		fmt.Println("error type")
		return nil
	}

	limit := util.BytesPrefix([]byte(fieldName + "@"))
	for _, kv := range sortDatas {
		if bytes.Compare([]byte(kv.Key), limit.Limit) > 0 {
			break
		}
		if bytes.Compare([]byte(kv.Key), limit.Start) < 0 {
			continue
		}
		if version {
			kv.Value = string([]byte(kv.Value)[:len([]byte(kv.Value))-1])
		}
		kvBytes, err := jsoniter.Marshal(kv)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		result = append(result, kvBytes)
	}
	return result
}

func GetLdbDefaultConfig() *opt.Options {
	return &opt.Options{}
}

func GetLdbConfig() *opt.Options {
	opt := GetLdbDefaultConfig()

	opt.CompactionTableSize = int((1 << 20) * 8)
	opt.BlockSize = int((1 << 10) * 4)
	opt.BlockCacheCapacity = int((1 << 20) * 8)
	opt.WriteBuffer = int((1 << 20) * 4)
	opt.WriteL0PauseTrigger = 12
	opt.WriteL0SlowdownTrigger = 8

	return opt
}

func Reset() {
	if DBStorage != nil {
		DBStorage.Close()
		DBStorage = nil
	}

	MainClass = ""
	ClassMap = make(map[string]*classfile.ClassFile)
	Contract = nil
}
