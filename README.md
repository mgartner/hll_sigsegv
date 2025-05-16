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
PC=0x45b48c m=24 sigcode=1 addr=0xa2aa0343

goroutine 0 gp=0xc000884540 m=24 mp=0xc000880708 [idle]:
runtime.fpTracebackPCs(...)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/tracestack.go:258
runtime.traceStack(0xc000880708?, 0xc000105a78?, 0x2)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/tracestack.go:116 +0x2ac fp=0xc00015bec8 sp=0xc00015ba60 pc=0x45b48c
runtime.traceLocker.stack(...)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/traceevent.go:176
runtime.traceLocker.GoStop({0xc000880708?, 0xc000d02b70?}, 0x2)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/traceruntime.go:480 +0x85 fp=0xc00015bf50 sp=0xc00015bec8 pc=0x45a665
runtime.traceLocker.GoPreempt(...)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/traceruntime.go:475
runtime.goschedImpl(0xc000168e00, 0x1?)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/proc.go:4118 +0x7a fp=0xc00015bfa0 sp=0xc00015bf50 pc=0x43ccda
runtime.gopreempt_m(0xc000168e00?)
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/proc.go:4154 +0x18 fp=0xc00015bfc0 sp=0xc00015bfa0 pc=0x43d018
runtime.mcall()
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/asm_amd64.s:459 +0x4e fp=0xc00015bfd8 sp=0xc00015bfc0 pc=0x46b40e

goroutine 93 gp=0xc000168e00 m=24 mp=0xc000880708 [running]:
runtime.asyncPreempt2()
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/preempt.go:308 +0x39 fp=0xc000b0c9b0 sp=0xc000b0c990 pc=0x433bb9
runtime.asyncPreempt()
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/preempt_amd64.s:53 +0xdb fp=0xc000b0cb38 sp=0xc000b0c9b0 pc=0x46e7db
github.com/axiomhq/hyperloglog.hashFunc({0xc000414c00?, 0x6517005f8c22fad3?, 0xc0008a27c0?})
        /home/marcus_cockroachlabs_com/go/pkg/mod/github.com/axiomhq/hyperloglog@v0.2.5/utils.go:45 +0x41 fp=0xc000b0cb70 sp=0xc000b0cb38 pc=0x4972c1
github.com/axiomhq/hyperloglog.(*Sketch).Insert(0xc0008a27c0, {0xc000414c00?, 0x3e8?, 0x3e8?})
        /home/marcus_cockroachlabs_com/go/pkg/mod/github.com/axiomhq/hyperloglog@v0.2.5/hyperloglog.go:148 +0x2d fp=0xc000b0cb98 sp=0xc000b0cb70 pc=0x49650d
main.hll(0xc00011a040)
        /home/marcus_cockroachlabs_com/workspace/hll_sigsegv/main.go:45 +0x85 fp=0xc000b0cfc8 sp=0xc000b0cb98 pc=0x499005
main.main.gowrap1()
        /home/marcus_cockroachlabs_com/workspace/hll_sigsegv/main.go:26 +0x25 fp=0xc000b0cfe0 sp=0xc000b0cfc8 pc=0x498f45
runtime.goexit({})
        /home/marcus_cockroachlabs_com/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.9.linux-amd64/src/runtime/asm_amd64.s:1700 +0x1 fp=0xc000b0cfe8 sp=0xc000b0cfe0 pc=0x46d2c1
created by main.main in goroutine 1
        /home/marcus_cockroachlabs_com/workspace/hll_sigsegv/main.go:26 +0x6b
```

<pre>
$ objdump --disassemble='github.com/axiomhq/hyperloglog.hashFunc' ./hll_sigsegv

./hll_sigsegv:     file format elf64-x86-64


Disassembly of section .text:

0000000000497280 <github.com/axiomhq/hyperloglog.hashFunc>:
  497280:       49 3b 66 10             cmp    0x10(%r14),%rsp
  497284:       76 41                   jbe    4972c7 <github.com/axiomhq/hyperloglog.hashFunc+0x47>
  497286:       55                      push   %rbp
  497287:       48 89 e5                mov    %rsp,%rbp
  49728a:       48 83 ec 28             sub    $0x28,%rsp
  49728e:       48 89 44 24 38          mov    %rax,0x38(%rsp)
  497293:       48 89 04 24             mov    %rax,(%rsp)
  497297:       48 89 5c 24 08          mov    %rbx,0x8(%rsp)
  49729c:       48 89 4c 24 10          mov    %rcx,0x10(%rsp)
  4972a1:       48 c7 44 24 18 39 05    movq   $0x539,0x18(%rsp)
  4972a8:       00 00
  4972aa:       e8 51 ea ff ff          callq  495d00 <github.com/dgryski/go-metro.Hash64.abi0>
  4972af:       45 0f 57 ff             xorps  %xmm15,%xmm15
  4972b3:       64 4c 8b 34 25 f8 ff    mov    %fs:0xfffffffffffffff8,%r14
  4972ba:       ff ff
  4972bc:       48 8b 44 24 20          mov    0x20(%rsp),%rax
  <b>4972c1:       48 83 c4 28             add    $0x28,%rsp</b>
  4972c5:       5d                      pop    %rbp
  4972c6:       c3                      retq
  4972c7:       48 89 44 24 08          mov    %rax,0x8(%rsp)
  4972cc:       48 89 5c 24 10          mov    %rbx,0x10(%rsp)
  4972d1:       48 89 4c 24 18          mov    %rcx,0x18(%rsp)
  4972d6:       e8 c5 42 fd ff          callq  46b5a0 <runtime.morestack_noctxt.abi0>
  4972db:       48 8b 44 24 08          mov    0x8(%rsp),%rax
  4972e0:       48 8b 5c 24 10          mov    0x10(%rsp),%rbx
  4972e5:       48 8b 4c 24 18          mov    0x18(%rsp),%rcx
  4972ea:       eb 94                   jmp    497280 <github.com/axiomhq/hyperloglog.hashFunc>
</pre>
