(module
 (type $FUNCSIG$ii (func (param i32) (result i32)))
 (data (global.get $__memory_base) "Hello")
 (import "env" "__memory_base" (global $__memory_base i32))
 (import "env" "_puts" (func $_puts (param i32) (result i32)))
 (memory $memory 1)
 (global $STACKTOP (mut i32) (i32.const 0))
 (global $STACK_MAX (mut i32) (i32.const 0))
 (export "__post_instantiate" (func $__post_instantiate))
 (export "_main" (func $_main))
 (func $_main (; 1 ;) (; has Stack IR ;) (result i32)
  ;;@ src/puts.c:5:0
  (drop
   (call $_puts
    (global.get $__memory_base)
   )
  )
  ;;@ src/puts.c:6:0
  (i32.const 0)
 )
 (func $__post_instantiate (; 2 ;) (; has Stack IR ;)
  (global.set $STACKTOP
   (i32.add
    (global.get $__memory_base)
    (i32.const 16)
   )
  )
  (global.set $STACK_MAX
   (i32.add
    (global.get $STACKTOP)
    (i32.const 5242880)
   )
  )
 )
)
