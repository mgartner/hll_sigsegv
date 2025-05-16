This is a reproduction of a segfault when runtime execution tracing is enabled.
It reproduces fairly regularly for me on amd64. This reproduction uses the
`github.com/dgryski/go-metro` library, which has an assembly implementation for
amd64 - perhaps that is relevant.

At Cockroach Labs, we saw a few occurrences of this segfault in CockroachDB
clusters, with stack traces related to our usage of
`github.com/axiomhq/hyperloglog`, which relies on `go-metro`. I was able to
create this minimized reproduction using just `go-metro`.

```
go build . && ./hll_sigsegv
```

```
SIGSEGV: segmentation violation
PC=0x45acec m=8 sigcode=1 addr=0xa2aa0343

goroutine 0 gp=0xc0005841c0 m=8 mp=0xc000580008 [idle]:
runtime.fpTracebackPCs(...)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/tracestack.go:258
runtime.traceStack(0xc000580008?, 0xc00015d1b8?, 0x6)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/tracestack.go:116 +0x2ac fp=0xc000597ec8 sp=0xc000597a60 pc=0x45acec
runtime.traceLocker.stack(...)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/traceevent.go:176
runtime.traceLocker.GoStop({0xc000580008?, 0xc00041dba8?}, 0x2)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/traceruntime.go:480 +0x85 fp=0xc000597f50 sp=0xc000597ec8 pc=0x459ec5
runtime.traceLocker.GoPreempt(...)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/traceruntime.go:475
runtime.goschedImpl(0xc000168fc0, 0x1?)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/proc.go:4118 +0x7a fp=0xc000597fa0 sp=0xc000597f50 pc=0x43c53a
runtime.gopreempt_m(0xc000168fc0?)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/proc.go:4154 +0x18 fp=0xc000597fc0 sp=0xc000597fa0 pc=0x43c878
runtime.mcall()
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/asm_amd64.s:459 +0x4e fp=0xc000597fd8 sp=0xc000597fc0 pc=0x469f8e

goroutine 96 gp=0xc000168fc0 m=8 mp=0xc000580008 [running]:
runtime.asyncPreempt2()
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/preempt.go:308 +0x39 fp=0xc0009989e8 sp=0xc0009989c8 pc=0x433419
runtime.asyncPreempt()
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/preempt_amd64.s:53 +0xdb fp=0xc000998b70 sp=0xc0009989e8 pc=0x46d33b
main.hash({0xc000a10000?, 0x3e8?, 0x3e8?})
        /home/marcus_cockroachlabs_com/workspace/hll_sigsegv/main.go:54 +0x41 fp=0xc000998ba8 sp=0xc000998b70 pc=0x47b981
main.hll(0xc00011a040)
        /home/marcus_cockroachlabs_com/workspace/hll_sigsegv/main.go:46 +0x6b fp=0xc000998fc8 sp=0xc000998ba8 pc=0x47b86b
main.main.gowrap1()
        /home/marcus_cockroachlabs_com/workspace/hll_sigsegv/main.go:27 +0x25 fp=0xc000998fe0 sp=0xc000998fc8 pc=0x47b7c5
runtime.goexit({})
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/asm_amd64.s:1700 +0x1 fp=0xc000998fe8 sp=0xc000998fe0 pc=0x46be21
created by main.main in goroutine 1
        /home/marcus_cockroachlabs_com/workspace/hll_sigsegv/main.go:27 +0x6b
```

<pre>
$ objdump --disassemble='main.hash' ./hll_sigsegv   w/ marcus_cockroachlabs_com@gceworker-marcus

./hll_sigsegv:     file format elf64-x86-64


Disassembly of section .text:

000000000047b940 <main.hash>:
  47b940:       49 3b 66 10             cmp    0x10(%r14),%rsp
  47b944:       76 41                   jbe    47b987 <main.hash+0x47>
  47b946:       55                      push   %rbp
  47b947:       48 89 e5                mov    %rsp,%rbp
  47b94a:       48 83 ec 28             sub    $0x28,%rsp
  47b94e:       48 89 44 24 38          mov    %rax,0x38(%rsp)
  47b953:       48 89 04 24             mov    %rax,(%rsp)
  47b957:       48 89 5c 24 08          mov    %rbx,0x8(%rsp)
  47b95c:       48 89 4c 24 10          mov    %rcx,0x10(%rsp)
  47b961:       48 c7 44 24 18 39 05    movq   $0x539,0x18(%rsp)
  47b968:       00 00
  47b96a:       e8 51 fa ff ff          callq  47b3c0 <github.com/dgryski/go-metro.Hash64.abi0>
  47b96f:       45 0f 57 ff             xorps  %xmm15,%xmm15
  47b973:       64 4c 8b 34 25 f8 ff    mov    %fs:0xfffffffffffffff8,%r14
  47b97a:       ff ff
  47b97c:       48 8b 44 24 20          mov    0x20(%rsp),%rax
  <b>47b981:       48 83 c4 28             add    $0x28,%rsp</b>
  47b985:       5d                      pop    %rbp
  47b986:       c3                      retq
  47b987:       48 89 44 24 08          mov    %rax,0x8(%rsp)
  47b98c:       48 89 5c 24 10          mov    %rbx,0x10(%rsp)
  47b991:       48 89 4c 24 18          mov    %rcx,0x18(%rsp)
  47b996:       e8 85 e7 fe ff          callq  46a120 <runtime.morestack_noctxt.abi0>
  47b99b:       48 8b 44 24 08          mov    0x8(%rsp),%rax
  47b9a0:       48 8b 5c 24 10          mov    0x10(%rsp),%rbx
  47b9a5:       48 8b 4c 24 18          mov    0x18(%rsp),%rcx
  47b9aa:       eb 94                   jmp    47b940 <main.hash>
</pre>
