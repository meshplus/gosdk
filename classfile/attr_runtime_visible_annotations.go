package classfile

type RuntimeVisibleAnnotationsAttribute struct {
	cp          *ConstantPool
	annotations []Annotation
}

func (runtimeVisibleAnnotationsAttribute *RuntimeVisibleAnnotationsAttribute) Annotations() []Annotation {
	return runtimeVisibleAnnotationsAttribute.annotations
}

type Annotation struct {
	cp                      *ConstantPool
	type_index              uint16
	num_element_value_pairs uint16
	element_value_pairs     []ElementValuePair
}

func (annotation *Annotation) Type() string {
	return annotation.cp.getUtf8(annotation.type_index)
}

type ElementValuePair struct {
	cp                 *ConstantPool
	element_name_index uint16
	element_value      ElementValue
}

type ElementValue struct {
	tag uint8

	cp                *ConstantPool
	const_value_index uint16
	enum_const_value  EnumConstValue
	class_info_index  uint16
	annotation_value  Annotation
	array_value       ArrayValue
}

type EnumConstValue struct {
	cp               *ConstantPool
	type_name_index  uint16
	const_name_index uint16
}

type ArrayValue struct {
	num_values uint16
	values     []ElementValue
}

func (attribute *RuntimeVisibleAnnotationsAttribute) readInfo(reader *ClassReader) {
	numAnnotations := reader.readUint16()
	attribute.annotations = make([]Annotation, numAnnotations)
	for i := range attribute.annotations {
		attribute.annotations[i] = readAnnotation(reader, attribute.cp)
	}
}

func readAnnotation(reader *ClassReader, cp *ConstantPool) Annotation {
	annotation := Annotation{}
	annotation.cp = cp
	annotation.type_index = reader.readUint16()
	annotation.num_element_value_pairs = reader.readUint16()
	annotation.element_value_pairs = make([]ElementValuePair, annotation.num_element_value_pairs)
	for i := range annotation.element_value_pairs {
		annotation.element_value_pairs[i] = readElementValuePair(reader, cp)
	}
	return annotation
}

func readElementValuePair(reader *ClassReader, cp *ConstantPool) ElementValuePair {
	elementValuePair := ElementValuePair{}
	elementValuePair.cp = cp
	elementValuePair.element_name_index = reader.readUint16()
	elementValuePair.element_value = readElementValue(reader, cp)
	return elementValuePair
}

func readElementValue(reader *ClassReader, cp *ConstantPool) ElementValue {
	tag := reader.readUint8()

	elementValue := ElementValue{}
	elementValue.cp = cp
	elementValue.tag = tag
	switch tag {
	case byte('B'), byte('C'), byte('D'), byte('F'), byte('I'), byte('J'), byte('S'), byte('Z'), byte('s'):
		elementValue.const_value_index = reader.readUint16()
	case byte('e'):
		elementValue.enum_const_value = readEnumConstValue(reader, cp)
	case byte('c'):
		elementValue.class_info_index = reader.readUint16()
	case byte('@'):
		elementValue.annotation_value = readAnnotation(reader, cp)
	case byte('['):
		elementValue.array_value = readArrayValue(reader, cp)
	default:
		panic("unregornized tag")
	}
	return elementValue
}

func readEnumConstValue(reader *ClassReader, cp *ConstantPool) EnumConstValue {
	enumConstValue := EnumConstValue{}
	enumConstValue.cp = cp
	enumConstValue.type_name_index = reader.readUint16()
	enumConstValue.const_name_index = reader.readUint16()
	return enumConstValue
}

func readArrayValue(reader *ClassReader, cp *ConstantPool) ArrayValue {
	num_values := reader.readUint16()
	arrayValue := ArrayValue{}
	arrayValue.num_values = num_values
	arrayValue.values = make([]ElementValue, num_values)
	for i := range arrayValue.values {
		arrayValue.values[i] = readElementValue(reader, cp)
	}
	return arrayValue
}
