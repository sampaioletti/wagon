// +build !amd64

package exec

import (
	"github.com/go-interpreter/wagon/exec/internal/compile"
)

// nativeCompiler represents a backend for native code generation + execution.
type nativeCompiler struct {
	Scanner   sequenceScanner
	Builder   instructionBuilder
	allocator pageAllocator
}

func (c *nativeCompiler) Close() error {
	return c.allocator.Close()
}
func nativeBackend() (bool, *nativeCompiler) {
	return false, nil
}

func (vm *VM) tryNativeCompile() error {
	return nil
}

// nativeCodeInvocation calls into one of the assembled code blocks.
// Assembled code blocks expect the following two pieces of
// information on the stack:
// [fp:fp+pointerSize]: sliceHeader for the stack.
// [fp+pointerSize:fp+pointerSize*2]: sliceHeader for locals variables.
func (vm *VM) nativeCodeInvocation(asmIndex uint32) {
}

// pageAllocator is responsible for the efficient allocation of
// executable, aligned regions of executable memory.
type pageAllocator interface {
	AllocateExec(asm []byte) (compile.NativeCodeUnit, error)
	Close() error
}

// sequenceScanner is responsible for detecting runs of supported opcodes
// that could benefit from compilation into native instructions.
type sequenceScanner interface {
	// ScanFunc returns an ordered, non-overlapping set of
	// sequences to compile into native code.
	ScanFunc(bytecode []byte, meta *compile.BytecodeMetadata) ([]compile.CompilationCandidate, error)
}

// instructionBuilder is responsible for compiling wasm opcodes into
// native instructions.
type instructionBuilder interface {
	// Build compiles the specified bytecode into native instructions.
	Build(candidate compile.CompilationCandidate, code []byte, meta *compile.BytecodeMetadata) ([]byte, error)
}
