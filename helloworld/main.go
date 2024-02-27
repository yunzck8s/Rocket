package main

//go:generate bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS -target bpfel,bpfeb bpf helloworld.bpf.c -- -I $BPF_HEADERS
import (
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGALRM)

	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("Loading objects: %s,", err)
	}
	defer objs.Close()

	kp, err := link.Tracepoint("syscalls", "sys_enter_execve", objs.BpfProg, nil)
	if err != nil {
		log.Fatalf("opening tracepoint: %s", err)
	}
	defer kp.Close()
	log.Printf("Successfully started! Please run \"sudo cat /sys/kernel/debug/tracing/trace_pipe\" to see output of the BPF programs\n")

	<-stopper
	log.Println("Received signal, exiting program..")
}
