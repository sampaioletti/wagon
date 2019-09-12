//go:generate emcc -Os src/puts.c -s SIDE_MODULE=1 -o puts.wasm -s TOTAL_MEMORY=65536 -s TOTAL_STACK=4096

package test
