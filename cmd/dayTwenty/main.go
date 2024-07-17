package dayTwenty

import (
	"fmt"
	"log/slog"
	"math"
	"strings"

	"github.com/spf13/cobra"
)

/**
cables connect to modules
modules connect to machines
modules communicate with either a high or low pulse
module types:
- flip-flop (%) are on or off and start off
  - high pulse: do nothing
  - low pulse: flip state, if it was off send a high pulse otherwise send a low
- conjunction (&) remember the type of the most recent pulse received from each connected input
  - default to remembering a low pulse for each input
  - on a new pulse, update memory, then send a low pulse if all memory is high, otherwise send high pulse
- broadcast (`broadcaster`)
  - when receiving a pulse send that same pulse to all destinations
- button module sends a low pulse to the broadcaster

never push the button if modules are still processing (so no parallelization)
pulses are processed in order, so it's a queue, not a stack, for pulse processing

**/

type PulseType int

const (
	High PulseType = iota
	Low
	Nil
)

func (p PulseType) String() string {
	switch p {
	case High:
		return "High"
	case Low:
		return "Low"
	case Nil:
		return "Nil"
	default:
		panic("wtf")
	}
}

type Module struct {
	ModuleKind       string   `@ModuleKind?`
	Name             string   `@Ident Pointer`
	Receivers        []string `@Ident (ReceiverSeparator @Ident)*`
	flipFlopState    bool
	conjunctionState map[string]PulseType
}

func (m *Module) String() string {
	switch m.ModuleKind {
	case "":
		return fmt.Sprintf("%s(%v)", m.Name, m.Receivers)
	case "%":
		return fmt.Sprintf("flip-flop(%s, %v, %t)", m.Name, m.Receivers, m.flipFlopState)
	case "&":
		return fmt.Sprintf("conjunction(%s, %v, %v)", m.Name, m.Receivers, m.conjunctionState)
	default:
		panic("wtf")
	}
}

func (m *Module) Process(src *Module, p PulseType) (PulseType, []string) {
	nextPulse := p
	destinations := []string{}
	switch m.ModuleKind {
	case "":
		// Broadcast
		destinations = m.Receivers
		nextPulse = p
	case "%":
		if p == High {
			return Nil, nil
		}
		// It's low, so let's send!
		if m.flipFlopState {
			nextPulse = Low
		} else {
			nextPulse = High
		}
		m.flipFlopState = !m.flipFlopState
		destinations = m.Receivers
	case "&":
		/**
		- conjunction (&) remember the type of the most recent pulse received from each connected input
		  - default to remembering a low pulse for each input
		  - on a new pulse, update memory, then send a low pulse if all memory is high, otherwise send high pulse
		  **/
		m.conjunctionState[src.Name] = p
		destinations = m.Receivers
		nextPulse = Low
		for _, v := range m.conjunctionState {
			if v == Low {
				nextPulse = High
				break
			}
		}
	}

	return nextPulse, destinations
}

type Pulse struct {
	src   *Module
	pType PulseType
	dsts  []string
}

func (p *Pulse) String() string {
	return fmt.Sprintf("Pulse(%s, %v, %v)", p.src.Name, p.pType, p.dsts)
}

func Push(ms map[string]*Module, pushes int) (int, int) {
	lows, highs := 0, 0

	for i := 0; i < pushes; i++ {
		pulses := []*Pulse{{nil, Low, []string{"broadcaster"}}}
		for len(pulses) > 0 {
			n := pulses[0]
			pulses = pulses[1:]
			slog.Debug("Processing", "n", n)
			if n.pType == Low {
				lows += len(n.dsts)
			} else if n.pType == High {
				highs += len(n.dsts)
			}

			for _, d := range n.dsts {
				m := ms[d]
				if m == nil {
					// Just a testing destination
					continue
				}
				nextPulse, dsts := m.Process(n.src, n.pType)
				if nextPulse != Nil {
					pulses = append(pulses, &Pulse{m, nextPulse, dsts})
				}
			}
		}
		slog.Debug("Pushed", "i", i, "lows", lows, "highs", highs, "ms", ms)
	}
	return lows, highs
}

func MinimumForRx(ms map[string]*Module) int {
	lows, highs := 0, 0

	found := map[string][]int{"ph": []int{}, "nz": []int{}, "dd": []int{}, "tx": []int{}}
	for i := 1.0; i < math.MaxFloat64; i++ {
		pulses := []*Pulse{{nil, Low, []string{"broadcaster"}}}
		for len(pulses) > 0 {
			n := pulses[0]
			pulses = pulses[1:]
			slog.Debug("Processing", "n", n)
			if n.pType == Low {
				lows += len(n.dsts)
			} else if n.pType == High {
				highs += len(n.dsts)
			}

			for _, d := range n.dsts {
				m := ms[d]
				if d == "rx" && n.pType == Low {
					return lows + highs
				}
				if m == nil {
					// Just a testing destination
					continue
				}
				// Following some Reddit advice, I'm looking for numbers when one of my four inputs
				// that lead to the "final trail" get a high pulse
				allSeenThrice := true
				for _, in := range []string{"ph", "nz", "dd", "tx"} {
					if n.src != nil && n.src.Name == in && n.pType == High {
						found[in] = append(found[in], int(i))

						slog.Info("Found", "i", i, "seen", found[in], "in", in)
					}
					if len(found[in]) < 3 {
						allSeenThrice = false
					}
				}

				if allSeenThrice {
					cycle := 1
					for _, v := range found {
						cycle = cycle * v[0]
					}
					return cycle
				}
				nextPulse, dsts := m.Process(n.src, n.pType)
				if nextPulse != Nil {
					pulses = append(pulses, &Pulse{m, nextPulse, dsts})
				}
			}
		}
		if int(i)%1000000 == 0 {
			slog.Info("Pushed", "i", i, "%", fmt.Sprintf("%.10f", i/math.MaxFloat64*100), "lows", lows, "highs", highs)
		}
	}

	panic("wtf")
}

func partOne(puzzleFile string, pushCount int) {
	slog.Info("Day Twenty part one", "puzzle file", puzzleFile)

	modules := ParseModules(puzzleFile)
	slog.Debug("Parsed modules", "modules", modules)
	lowPulses, highPulses := Push(modules, pushCount)

	slog.Info("Parsing modules", "low pulses", lowPulses, "high pulses", highPulses, "product", lowPulses*highPulses)
}

func PrintGraph(f string) {
	ms := ParseModules(f)
	fmt.Println("flowchart TD")
	for _, m := range ms {
		name := strings.ToUpper(m.Name)
		switch m.ModuleKind {
		case "":
			fmt.Println(fmt.Sprintf("    %s(%s)", name, m.Name))
		case "%":
			fmt.Println(fmt.Sprintf("    %s{%s%s}", name, m.ModuleKind, m.Name))
		case "&":
			fmt.Println(fmt.Sprintf("    %s[%s%s]", name, m.ModuleKind, m.Name))
		}
		for _, r := range m.Receivers {
			fmt.Println(fmt.Sprintf("    %s --> %s", name, strings.ToUpper(r)))
		}
	}
}

func partTwo(puzzleFile string) {
	slog.Info("Day Twenty part two", "puzzle file", puzzleFile)

	modules := ParseModules(puzzleFile)
	slog.Debug("Parsed modules", "modules", modules)
	minimumPulses := MinimumForRx(modules)

	slog.Info("Minimum pulses for rx", "minimum pulses", minimumPulses)
}

var Cmd = &cobra.Command{
	Use: "dayTwenty",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		if cmd.Flag("print-graph").Changed {
			PrintGraph(puzzleInput)
		}
		pushCount, _ := cmd.Flags().GetInt("push-count")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput, pushCount)
		} else {
			partTwo(puzzleInput)
		}
	},
}

func init() {
	Cmd.Flags().Bool("part-two", false, "Whether to run part two of the day's challenge")
	Cmd.Flags().Bool("print-graph", false, "print our graph")
	Cmd.Flags().Int("push-count", 1000, "Push count")
}
