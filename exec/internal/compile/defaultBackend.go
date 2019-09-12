// +build !amd64

package compile

import (
	ops "github.com/go-interpreter/wagon/wasm/operators"
	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/obj/x86"
)

type PlatformBackend struct {
}

func (b *PlatformBackend) paramsForMemoryOp(op byte) (size uint, inst obj.As) {
	switch op {
	case ops.I64Load, ops.F64Load:
		return 8, x86.AMOVQ
	case ops.I32Load, ops.F32Load:
		return 4, x86.AMOVL
	case ops.I64Store, ops.F64Store:
		return 8, x86.AMOVQ
	case ops.I32Store, ops.F32Store:
		return 4, x86.AMOVL
	}
	panic("unreachable")
}
