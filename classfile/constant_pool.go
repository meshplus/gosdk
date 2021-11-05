package classfile

import (
	"fmt"
)

type ConstantPool struct {
	cf    *ClassFile
	Infos []ConstantInfo
}

func (self *ConstantPool) read(reader *ClassReader) {
	cpCount := int(reader.readUint16())
	self.Infos = make([]ConstantInfo, cpCount)

	// The constant_pool table is indexed from 1 to constant_pool_count - 1.
	for i := 1; i < cpCount; i++ {
		self.Infos[i] = readConstantInfo(reader, self)
		switch self.Infos[i].(type) {
		case int64, float64:
			i++
		}
	}
}

func (self *ConstantPool) getConstantInfo(index uint16) ConstantInfo {
	cpInfo := self.Infos[index]
	if cpInfo == nil {
		panic(fmt.Errorf("Bad constant pool index: %v!", index))
	}

	return cpInfo
}

func (self *ConstantPool) getNameAndType(index uint16) (name, _type string) {
	ntInfo := self.getConstantInfo(index).(ConstantNameAndTypeInfo)
	name = self.getUtf8(ntInfo.nameIndex)
	_type = self.getUtf8(ntInfo.descriptorIndex)
	return
}

func (self *ConstantPool) getClassName(index uint16) string {
	classInfo := self.getConstantInfo(index).(ConstantClassInfo)
	return self.getUtf8(classInfo.nameIndex)
}

func (self *ConstantPool) getUtf8(index uint16) string {
	return self.getConstantInfo(index).(string)
}
