package main

import (
	"log"
	"os"
	"strings"
)

const (
	NOTHING = -1
	LO      = 0
	HI      = 1
)

type Signal struct {
	sender  string
	signal  int
	outputs []string
}

type Module interface {
	process(*Signal) *Signal
}

type Broadcast struct {
	outputs   []string
	label     string
	numInputs int
}

type FlipFlop struct {
	Broadcast
	on bool
}

type Conjunction struct {
	Broadcast
	memory map[string]int
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need arg")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s", err)
	}
	contents := string(contentsBytes)
	labelToModule := make(map[string]Module)
	modsToOutputs := make(map[string]int)
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			continue
		}
		log.Printf("%s", line)
		lineParts := strings.Split(line, " -> ")
		right := lineParts[1]
		out := make([]string, 0)
		for _, s := range strings.Split(right, ", ") {
			prev, ok := modsToOutputs[s]
			if !ok {
				modsToOutputs[s] = 1
			} else {
				modsToOutputs[s] = prev + 1
			}
			out = append(out, s)
		}
		left := lineParts[0]
		var m Module
		label := left
		if left == "broadcaster" {
			m = &Broadcast{outputs: out, label: label}
			labelToModule[label] = m
		} else if left[0] == '%' {
			label = left[1:]
			m = &FlipFlop{
				Broadcast: Broadcast{outputs: out, label: label},
				on:        false,
			}
		} else if left[0] == '&' {
			label = left[1:]
			m = &Conjunction{
				Broadcast: Broadcast{outputs: out, label: label},
				memory:    make(map[string]int),
			}
		} else {
			log.Fatal(left)
		}
		log.Printf("label: %s, outputs: %v", label, out)
		labelToModule[label] = m
	}

	for label, mod := range labelToModule {
		c, ok := mod.(*Conjunction)
		if !ok {
			continue
		}
		c.numInputs = modsToOutputs[label]
	}

	loCount := 0
	hiCount := 0
	for i := 0; i < 1000; i++ {
		q := make([]*Signal, 0)
		first := &Signal{
			signal:  LO,
			outputs: []string{"broadcaster"},
		}
		q = append(q, first)
		for j := 0; j < len(q); j++ {
			currSignal := q[j]
			if currSignal.signal == NOTHING {
				continue
			}
			if currSignal.signal == LO {
				loCount += len(currSignal.outputs)
			} else if currSignal.signal == HI {
				hiCount += len(currSignal.outputs)
			}
			for _, rec := range currSignal.outputs {
				if currSignal.signal == LO {
					log.Printf("%s -low-> %s", currSignal.sender, rec)
				} else {
					log.Printf("%s -high-> %s", currSignal.sender, rec)
				}
				m, ok := labelToModule[rec]
				if !ok {
					continue
				}
				out := m.process(currSignal)
				q = append(q, out)
			}
		}
	}
	log.Printf("lo: %d, hi: %d", loCount, hiCount)
	result := loCount * hiCount
	log.Printf("%d", result)
}

func (b *Broadcast) process(signal *Signal) (res *Signal) {
	res = &Signal{
		sender:  b.label,
		signal:  LO,
		outputs: b.outputs,
	}
	return
}

func (b *FlipFlop) process(signal *Signal) (res *Signal) {
	if signal.signal == HI {
		res = &Signal{
			sender: b.label,
			signal: NOTHING,
		}
		return
	}
	if b.on {
		res = &Signal{
			sender:  b.label,
			signal:  LO,
			outputs: b.outputs,
		}
	} else {
		res = &Signal{
			sender:  b.label,
			signal:  HI,
			outputs: b.outputs,
		}
	}
	b.on = !b.on
	return
}

func (b *Conjunction) process(signal *Signal) (res *Signal) {
	b.memory[signal.sender] = signal.signal
	isAllHi := true
	if len(b.memory) >= b.numInputs {
		for _, v := range b.memory {
			if v == LO {
				isAllHi = false
				break
			}
		}
		if isAllHi {
			res = &Signal{
				sender:  b.label,
				signal:  LO,
				outputs: b.outputs,
			}
			return
		}
	}
	res = &Signal{
		sender:  b.label,
		signal:  HI,
		outputs: b.outputs,
	}
	return
}
