package wasm

import (
	"fmt"
	"io"
	"reflect"

	"github.com/go-interpreter/wagon/wasm/leb128"
)

func NewBuilder() *Builder {
	b := &Builder{}
	return b
}

// // EmptyBuilder is before the module is loaded
// type EmptyBuilder interface{
// 	Base()ModuleBuilder
// 	Module(M *Module)ModuleBuilder
// 	Decode(r io.Reader)ModuleBuilder
// }
// // ModuleBuilder is after the module has been initialized
// type ModuleBuilder interface{

// }

type Builder struct {
	M   *Module
	Err error
}

func (b *Builder) NewModule() *Builder {
	if b.Err != nil {
		return b
	}
	b.M = NewModule()
	return b
}

func (b *Builder) Decode(r io.Reader) *Builder {
	if b.Err != nil {
		return b
	}
	M, err := DecodeModule(r)
	if err != nil {
		b.Err = err
		return b
	}
	b.M = M
	return b
}
func (b *Builder) InitLinearMemory(size int) *Builder {
	if b.Err != nil {
		return b
	}
	b.M.LinearMemoryIndexSpace = [][]byte{make([]byte, size)}
	return b
}
func (b *Builder) SetLinearMemory(buf []byte) *Builder {
	if b.Err != nil {
		return b
	}
	b.M.LinearMemoryIndexSpace = [][]byte{buf}
	return b
}
func (b *Builder) InitTableMemory(size int) *Builder {
	if b.Err != nil {
		return b
	}
	b.M.TableIndexSpace = [][]uint32{make([]uint32, size)}
	return b
}
func (b *Builder) SetTableMemory(buf []uint32) *Builder {
	if b.Err != nil {
		return b
	}
	b.M.TableIndexSpace = [][]uint32{buf}
	return b
}

func (b *Builder) ResolveImports(r ResolveFunc) *Builder {
	if b.Err != nil {
		return b
	}
	if b.M.Import != nil && r != nil {
		if b.M.Code == nil {
			b.M.Code = &SectionCode{}
		}
		err := b.M.resolveImports(r)
		if err != nil {
			b.Err = err
			return b
		}

	}
	return b
}

func (b *Builder) AddGlobal(init []byte, typ ValueType, name string, export bool) *Builder {
	if b.Err != nil {
		return b
	}
	globalEntry := GlobalEntry{
		Type: GlobalVar{
			Type: ValueTypeI32,
		},
		Init: init,
	}
	b.M.GlobalIndexSpace = append(b.M.GlobalIndexSpace, globalEntry)
	if !export {
		return b
	}
	exportEntry := ExportEntry{
		FieldStr: name,
		Kind:     ExternalGlobal,
		Index:    uint32(len(b.M.GlobalIndexSpace) - 1),
	}
	b.M.Export.Entries[name] = exportEntry
	return b
}
func (b *Builder) AddGlobalVal(name string, export bool, val interface{}) *Builder {
	if b.Err != nil {
		return b
	}
	//init expression currently only supports 32 bit values
	switch t := val.(type) {
	case int32, int64:
		buf := []byte{i32Const}
		buf = leb128.AppendSleb128(buf, t.(int64))
		buf = append(buf, end)
		b.AddGlobal(buf, ValueTypeI32, name, export)
	// case int64:
	// 	buf := []byte{i64Const}
	// 	buf = leb128.AppendSleb128(buf, t)
	// 	buf = append(buf, end)
	// 	b.AddGlobal(buf, ValueTypeI64, name, export)
	case uint32, uint64:
		buf := []byte{i32Const}
		buf = leb128.AppendUleb128(buf, t.(uint64))
		buf = append(buf, end)
		b.AddGlobal(buf, ValueTypeI32, name, export)
	// case uint64:
	// buf := []byte{i64Const}
	// buf = leb128.AppendUleb128(buf, t)
	// buf = append(buf, end)
	// b.AddGlobal(buf, ValueTypeI64, name, export)
	case float32, float64:
		b.Err = fmt.Errorf("float globals not implemented yet")
	default:
		b.Err = fmt.Errorf("Invalid Global Type %T must be int32, int64, float32 or float64", val)
	}
	return b
}

func (b *Builder) AddFunction(name string, export bool, sig FunctionSig, v interface{}) *Builder {
	if b.Err != nil {
		return b
	}

	f := Function{
		Sig:  &sig,
		Host: reflect.ValueOf(v),
		Body: &FunctionBody{},
		Name: name,
	}
	b.M.FunctionIndexSpace = append(b.M.FunctionIndexSpace, f)
	if !export {
		return b
	}
	entry := ExportEntry{
		FieldStr: name,
		Kind:     ExternalFunction,
		Index:    uint32(len(b.M.FunctionIndexSpace) - 1),
	}
	b.M.Export.Entries[entry.FieldStr] = entry

	return b
}

func (b *Builder) AddMemory(name string, export bool, initial uint32, max uint32) *Builder {
	if b.Err != nil {
		return b
	}
	def := Memory{
		Limits: ResizableLimits{Initial: initial},
	}
	if max > initial {
		def.Limits.Maximum = max
		def.Limits.Flags = 1
	}
	b.M.Memory.Entries = append(b.M.Memory.Entries, def)
	entry := ExportEntry{
		FieldStr: name,
		Kind:     ExternalMemory,
		Index:    uint32(len(b.M.Memory.Entries) - 1),
	}
	if !export {
		return b
	}
	b.M.Export.Entries[name] = entry
	return b
}

func (b *Builder) AddTable(name string, export bool, initial uint32, max uint32) *Builder {
	if b.Err != nil {
		return b
	}
	def := Table{
		Limits: ResizableLimits{Initial: initial},
	}
	if max > initial {
		def.Limits.Maximum = max
		def.Limits.Flags = 1
	}
	b.M.Table.Entries = append(b.M.Table.Entries, def)
	entry := ExportEntry{
		FieldStr: name,
		Kind:     ExternalTable,
		Index:    uint32(len(b.M.Memory.Entries) - 1),
	}
	if !export {
		return b
	}
	b.M.Export.Entries[name] = entry
	return b
}
func (b *Builder) Populate() *Builder {
	if b.Err != nil {
		return b
	}
	for _, fn := range []func() error{
		b.M.populateGlobals,
		b.M.populateFunctions,
		b.M.populateTables,
		b.M.populateLinearMemory,
	} {
		if err := fn(); err != nil {
			b.Err = err
			return b
		}
	}
	logger.Printf("There are %d entries in the function index space.", len(b.M.FunctionIndexSpace))
	return b
}
