package emlibc

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-interpreter/wagon/exec"

	"github.com/go-interpreter/wagon/wasm"
)

func ResolveEnv(name string) (*wasm.Module, error) {
	if name == "env" {
		return GetEnv(), nil
	}
	fmt.Println("tried resolve", name)
	return nil, errors.New("Not Found")
}
func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}
func GetEnv() *wasm.Module {

	m := wasm.NewModule()
	print := func(proc *exec.Process, v int32) int32 {
		fmt.Printf("result = %v\n", v)
		return 0
	}
	puts := func(proc *exec.Process, v int32) int32 {

		buf := []byte{}
		temp := make([]byte, 1)
		for i := int(v); i < proc.MemSize(); i++ {
			_, err := proc.ReadAt(temp, int64(i))
			if err != nil {
				fmt.Println(err)
			}
			if temp[0] == 0 {
				break
			}
			buf = append(buf, temp[0])
		}
		fmt.Println(string(buf))
		return 0
	}
	m.Types = &wasm.SectionTypes{
		Entries: []wasm.FunctionSig{
			{
				Form:        0, // value for the 'func' type constructor
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
		},
	}
	m.GlobalIndexSpace = []wasm.GlobalEntry{
		{
			Type: wasm.GlobalVar{
				Type: wasm.ValueTypeI32,
			},
			Init: []byte{65, 0, 11},
		},
	}
	// m.LinearMemoryIndexSpace = [][]byte{make([]byte, 256)}
	m.Memory = &wasm.SectionMemories{
		Entries: []wasm.Memory{
			{
				Limits: wasm.ResizableLimits{Initial: 1},
			},
		},
	}
	m.FunctionIndexSpace = []wasm.Function{
		{
			Sig:  &m.Types.Entries[0],
			Host: reflect.ValueOf(print),
			Body: &wasm.FunctionBody{}, // create a dummy wasm body (the actual value will be taken from Host.)
		},
		{
			Sig:  &m.Types.Entries[0],
			Host: reflect.ValueOf(puts),
			Body: &wasm.FunctionBody{}, // create a dummy wasm body (the actual value will be taken from Host.)
		},
	}
	m.Export = &wasm.SectionExports{
		Entries: map[string]wasm.ExportEntry{
			"print": {
				FieldStr: "print",
				Kind:     wasm.ExternalFunction,
				Index:    0,
			},
			"_puts": {
				FieldStr: "_puts",
				Kind:     wasm.ExternalFunction,
				Index:    1,
			},
			"__memory_base": {
				FieldStr: "__memory_base",
				Kind:     wasm.ExternalGlobal,
				Index:    0,
			},
			"memory": {
				FieldStr: "memory",
				Kind:     wasm.ExternalMemory,
				Index:    0,
			},
		},
	}
	return m
}
