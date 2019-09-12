// Copyright 2019 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compile

import (
	"bytes"
	"os"
	"os/exec"
)

// NativeCodeUnit represents compiled native code.
type NativeCodeUnit interface {
	Invoke(stack, locals, globals *[]uint64, mem *[]byte) JITExitSignal
}

func debugPrintAsm(asm []byte) {
	cmd := exec.Command("ndisasm", "-b64", "-")
	cmd.Stdin = bytes.NewReader(asm)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// CompletionStatus describes the final status of a native execution.
type CompletionStatus uint64

// Valid completion statuses.
const (
	CompletionOK CompletionStatus = iota
	CompletionBadBounds
	CompletionUnreachable
	CompletionFatalInternalError
)

// JITExitSignal is the value returned from the execution of a native section.
// The bits of this packed 64bit value is encoded as follows:
// [00:04] Completion Status
// [04:08] Reserved
// [08:32] Index of the WASM instruction where the exit occurred.
// [32:64] Status-specific 32bit value.
type JITExitSignal uint64
