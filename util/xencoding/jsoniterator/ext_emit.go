package jsoniterator

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"reflect"
	"strings"
)

//排除struct中omitempty标签（此功能对于json中有omitempty的可使用）
type EmitDefaultExtension struct {
	jsoniter.DummyExtension
}

// UpdateStructDescriptor ...
func (ed EmitDefaultExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {

	for _, field := range structDescriptor.Fields {
		//field.Field =&myfield{oldField}
		//没有标签的不更改
		if field.Field.Tag() == "" {
			continue
		}
		oldField := field.Field
		field.Field = &myField{oldField}
	}
}

type myField struct{ reflect2.StructField }

func (mf *myField) Tag() reflect.StructTag {
	return reflect.StructTag(strings.Replace(string(mf.StructField.Tag()), ",omitempty", "", -1))
}

/*var fieldCache sync.Map
//备用带缓存的方式
type EmitDefaultExtensionWithCache struct {
	jsoniter.DummyExtension
}

func (ed EmitDefaultExtensionWithCache) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, field := range structDescriptor.Fields {
		//field.Field =&myfield{oldField}
		//没有标签的不更改
		if field.Field.Tag() == "" {
			continue
		}
		key:=field

		//key := fKey.String()
		if v, ok := fieldCache.Load(key); ok {
			if f, ok := v.(reflect2.StructField); ok {
				field.Field = f
				continue
			}
		}
		fmt.Println(key)
		oldField := field.Field
		field.Field = &myField{oldField}
		fieldCache.Store(key, field.Field)
	}
}*/
