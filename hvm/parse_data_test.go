package hvm

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestParseObjectData(t *testing.T) {
	datas := `{"id":"111","name":"name","age":20}`

	iterator := jsoniter.NewIterator(jsoniter.ConfigDefault)
	iterator.ResetBytes([]byte(datas))
	any := iterator.ReadAny()
	fmt.Println("any", any.Get("id").ToString())

}

func TestParseData(t *testing.T) {
	//mainClassFile := ClassMap[MainClass]
	//
	//Contract = NewFieldType(MainClass).(*ObjectType)
	//Contract.Fields = make([]FieldType, 0)
	//for _, field := range mainClassFile.Fields {
	//	if annotationAttr := field.RuntimeVisibleAnnotationsAttributeData(); annotationAttr == nil {
	//		continue
	//	} else {
	//		flag := true
	//		for _, annotation := range annotationAttr.Annotations() {
	//			if annotation.Type() == "Lcn/hyperchain/annotations/StoreField;" {
	//				flag = false
	//				break
	//			}
	//		}
	//		if flag {
	//			continue
	//		}
	//	}
	//
	//	descriptor := field.Descriptor()
	//	Contract.Fields = append(Contract.Fields, NewFieldType(descriptor))
	//	Contract.Columns = append(Contract.Columns, field.Name())
	//
	//	// Process data
	//	fmt.Println("field.Name****:", field.Name())
	//	fmt.Println("field.Class****:", field.Descriptor())
	//
	//}
	//
	//fmt.Println("************HyperMap************")
	//parseContractData(Contract.Fields[0], mainClassFile.Fields[0], hypermap_data...)
	//fmt.Println()
	//
	//fmt.Println("************Student************")
	//parseContractData(Contract.Fields[1], mainClassFile.Fields[1], student_data)
	//fmt.Println()
	//
	//fmt.Println("************HyperList************")
	//parseContractData(Contract.Fields[2], mainClassFile.Fields[2], hyperlist_data...)
	//fmt.Println()
	//
	//fmt.Println("************HyperList************")
	//parseContractData(Contract.Fields[3], mainClassFile.Fields[3], hypertable_data...)
	//fmt.Println()
	//
	//fmt.Println("************Array************")
	//parseContractData(Contract.Fields[4], mainClassFile.Fields[4], array_data)
	//fmt.Println()
	//
	//for key, _ := range ClassMap {
	//	fmt.Println(key)
	//}
}

func TestParseContractData(t *testing.T) {
	t.Skip()
	bytes, err := DecompressJar("test-jar/hvm-bench-test-2.0-SNAPSHOT-student.jar")
	if err != nil {
		fmt.Println(err)
		return
	}
	InitContract(bytes)
	ParseContractData(Contract, ClassMap[MainClass], "eef907d9b8a4a61f353d6759055f47ec7e09ee61", true, "/Users/dong/go/src/github.com/meshplus/hvmd/storage/runtime/data", nil)
}

func TestParseContractData2(t *testing.T) {
	t.Skip("no hvmd")
	addr := "1f69751671f4aa9c3753959be9d79ce546989bea"
	path := "/Users/dong/go/src/github.com/meshplus/hvmd/storage/runtime/data"
	bytes, err := DecompressJar("test-jar/hvm-bench-test-2.0-SNAPSHOT-student.jar")
	if err != nil {
		fmt.Println(err)
		return
	}
	InitContract(bytes)
	ParseContractData(Contract, ClassMap[MainClass], addr, true, path, nil)
	fmt.Println("Contract:", Contract.Class)
	fmt.Println("Field0_stduents:", Contract.Fields[0])
	fmt.Println("Field1_student:", Contract.Fields[1])
	fmt.Println("Field2_studentList:", Contract.Fields[2])
	fmt.Println("Field3_studentArray:", Contract.Fields[3])
	fmt.Println(Contract.Fields[3].(*ArrayType).Fields[0])
	fmt.Println(Contract.Fields[3].(*ArrayType).Fields[1])
}

func TestProcessDatas(t *testing.T) {
	t.Skip()
	//path := "/Users/dong/go/src/github.com/meshplus/hvmd/storage/runtime/data"
	//DBStorage = OpenDB(path)
	//addr := "eef907d9b8a4a61f353d6759055f47ec7e09ee61"
	//datas := ProcessDatas(addr, "students")
	//fmt.Println("***************************")
	//fmt.Println(len(datas))
	//fmt.Println([]byte(`"`))
	//for _, data := range datas {
	//	fmt.Println("data:", data)
	//	fmt.Println(string(data))
	//}
}

func TestData(t *testing.T) {
	t.Skip()
	datasStr := `7b226b6579223a2273747564656e74732d5f5f73697a655f5f222c2276616c7565223a2233227d,7b226b6579223a2273747564656e7473405c226964315c22222c2276616c7565223a227b5c2269645c223a5c226964315c222c5c226e616d655c223a5c226e616d65315c222c5c226167655c223a32307d227d,7b226b6579223a2273747564656e7473405c226964325c22222c2276616c7565223a227b5c2269645c223a5c226964325c222c5c226e616d655c223a5c226e616d65325c222c5c226167655c223a32307d227d,7b226b6579223a2273747564656e7473405c226964335c22222c2276616c7565223a227b5c2269645c223a5c226964335c222c5c226e616d655c223a5c226e616d65335c222c5c226167655c223a32307d227d`
	datas := DataStringToBytes(datasStr, ",")

	bytes, err := DecompressJar("test-jar/hvm-bench-test-2.0-SNAPSHOT-student.jar")
	if err != nil {
		fmt.Println(err)
		return
	}
	InitContract(bytes)

	mainClassFile := ClassMap[MainClass]
	Contract = NewFieldType(MainClass).(*ObjectType)
	Contract.Fields = make([]FieldType, 0)
	for _, field := range mainClassFile.Fields {
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
		Contract.Fields = append(Contract.Fields, NewFieldType(descriptor))
		Contract.Columns = append(Contract.Columns, field.Name())
	}

	parseContractData(Contract.Fields[0], mainClassFile.Fields[0], datas...)
}

func TestHyperList(t *testing.T) {
	t.Skip()
	decompressJar, err := DecompressJar("/Users/dong/go/src/github.com/meshplus/gosdk/hvm/test-jar/J2sTest-1.0-SNAPSHOT-Hyperlist.jar")
	if err != nil {
		panic(err)
	}
	//fmt.Println("****")
	InitContract(decompressJar)
	//fmt.Println("****")
	ParseContractData(Contract, ClassMap[MainClass], "96a7f4f37b1d349b534891ac8734211a88ee2998", true, "/Users/dong/go/src/github.com/meshplus/hyperchain/build/node1/namespaces/global/data/leveldb/blockchain", nil)
	InitContract(decompressJar)
	ParseContractData(Contract, ClassMap[MainClass], "96a7f4f37b1d349b534891ac8734211a88ee2998", true, "/Users/dong/go/src/github.com/meshplus/hyperchain/build/node1/namespaces/global/data/leveldb/blockchain", nil)

}

var (
	incrementData1 string = `{"timestamp":1578294887729063000,"type":"MQHvm_2c21f1e17","body":{"type":"MQHvm_2c21f1e17","name":"Hvm","data":{"abfc54ef2479c207930d8b0f803ceeb7a1a72617":{"map@4":"NAE=","map@7":"NwE=","map@9":"OQE=","map@2":"MgE=","map-__size__":"MTAB","map@0":"MAE=","map@6":"NgE=","map@8":"OAE=","map@5":"NQE=","map@1":"MQE=","map@3":"MwE="}}}}`
)

func TestParseIncrementData(t *testing.T) {
	hvmLog := &HvmLog{}
	err := jsoniter.UnmarshalFromString(incrementData1, hvmLog)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%v\n", hvmLog)

	preDatas := hvmLog.Body.Data
	for addr, kvs := range preDatas {
		fmt.Println(addr)
		fmt.Println(kvs)
		for k, v := range kvs {
			fmt.Println("key:", k)
			fmt.Println("value:", string(v))
		}
	}
}

func TestParseIncrementData2(t *testing.T) {
	Reset()
	jarCode, err := DecompressJar("test-jar/hypermap-1.0-hypermap.jar")
	if err != nil {
		t.Error(err)
		return
	}

	InitContract(jarCode)
	ParseIncrementData(Contract, ClassMap[MainClass], "abfc54ef2479c207930d8b0f803ceeb7a1a72617", incrementData1)

	fmt.Printf("%v\n", Contract.Fields[0])
	fmt.Printf("%v\n", Contract.Fields[0].(*HyperMap).KFields[0])
	fmt.Printf("%v\n", Contract.Fields[0].(*HyperMap).VFields[0])
}
