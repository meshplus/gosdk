package classfile

import "strings"

/*
field_info {
    u2             access_flags;
    u2             name_index;
    u2             descriptor_index;
    u2             attributes_count;
    attribute_info attributes[attributes_count];
}
method_info {
    u2             access_flags;
    u2             name_index;
    u2             descriptor_index;
    u2             attributes_count;
    attribute_info attributes[attributes_count];
}
*/
type MemberInfo struct {
	cp              *ConstantPool
	AccessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	AttributeTable
}

// read field or method table
func readMembers(reader *ClassReader, cp *ConstantPool) []MemberInfo {
	memberCount := reader.readUint16()
	members := make([]MemberInfo, memberCount)
	for i := range members {
		members[i] = MemberInfo{cp: cp}
		members[i].read(reader)
	}
	return members
}

func (self *MemberInfo) read(reader *ClassReader) {
	self.AccessFlags = reader.readUint16()
	self.nameIndex = reader.readUint16()
	self.descriptorIndex = reader.readUint16()
	self.attributes = readAttributes(reader, self.cp)
}

func (self *MemberInfo) Name() string {
	return self.cp.getUtf8(self.nameIndex)
}
func (self *MemberInfo) Descriptor() string {
	return self.cp.getUtf8(self.descriptorIndex)
}
func (self *MemberInfo) Signature() string {
	signatureAttr := self.SignatureAttribute()
	if signatureAttr != nil {
		return signatureAttr.Signature()
	}
	return ""
}

func (self *MemberInfo) GetInvokeMethods(params []string) []string {
	var res []string
	isExist := make(map[string]bool)
	for _, v := range self.cp.Infos {
		if temp, ok := v.(ConstantMemberrefInfo); ok {
			if temp.Tag == CONSTANT_Methodref || temp.Tag == CONSTANT_InterfaceMethodref {
				for _, p := range params {
					if strings.Contains(p, temp.ClassName()) {
						name, _ := temp.NameAndDescriptor()
						if !isExist[name] {
							res = append(res, name)
							isExist[name] = true
						}
					}
				}
			}
		}
	}
	return res
}

func (self *MemberInfo) ArgumentTypes() []string {
	numArgumentTypes := 0
	var (
		begin         = 1
		currentOffset int
	)
	s := self.Descriptor()
	var res []string

	for currentOffset = 1; s[currentOffset] != ')'; numArgumentTypes++ {
		for s[currentOffset] == '[' {
			currentOffset++
		}

		if s[currentOffset] == 'L' {
			for i := currentOffset; i < len(s); i++ {
				if s[i] == ';' {
					currentOffset = i + 1
					break
				}
			}
		} else {
			currentOffset++
		}
		res = append(res, s[begin:currentOffset])
		begin = currentOffset
	}

	return res
}
