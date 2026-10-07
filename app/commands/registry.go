package commands

import "strings"

// Registry is built once at startup and read-only afterwards,
// so lookups need no mutex.
type Registry struct {
	specs map[string]Spec
}

func NewRegistry(specs ...Spec) *Registry {
	r := &Registry{specs: make(map[string]Spec, len(specs))}
	for _, s := range specs {
		key := strings.ToUpper(s.Name)
		if _, dup := r.specs[key]; dup {
			panic("duplicate command: " + key)
		}
		r.specs[key] = s
	}
	return r
}

func (r *Registry) Lookup(name string) (Spec, bool) {
	s, ok := r.specs[strings.ToUpper(name)]
	return s, ok
}

func DefaultRegistry() *Registry {
	return NewRegistry(
		Spec{Name: "BLPOP", MinArgs: 2, MaxArgs: 2, Flags: FlagWrite, Parse: parseBLPop},
		Spec{Name: "DISCARD", MinArgs: 0, MaxArgs: 0, Parse: parseDiscard},
		Spec{Name: "ECHO", MinArgs: 1, MaxArgs: 1, Parse: parseEcho},
		Spec{Name: "EXEC", MinArgs: 0, MaxArgs: 0, Flags: FlagWrite, Parse: parseExec},
		Spec{Name: "GET", MinArgs: 1, MaxArgs: 1, Parse: parseGet},
		Spec{Name: "INCR", MinArgs: 1, MaxArgs: 1, Flags: FlagWrite, Parse: parseIncr},
		Spec{Name: "INFO", MinArgs: 1, MaxArgs: 1, Parse: parseInfo},
		Spec{Name: "LLEN", MinArgs: 1, MaxArgs: 1, Parse: parseLLen},
		Spec{Name: "LPOP", MinArgs: 1, MaxArgs: 2, Flags: FlagWrite, Parse: parseLPop},
		Spec{Name: "LPUSH", MinArgs: 2, MaxArgs: -1, Flags: FlagWrite, Parse: parseLPush},
		Spec{Name: "LRANGE", MinArgs: 3, MaxArgs: 3, Parse: parseLRange},
		Spec{Name: "MULTI", MinArgs: 0, MaxArgs: 0, Parse: parseMulti},
		Spec{Name: "PING", MinArgs: 0, MaxArgs: 0, Flags: FlagNotPropagate, Parse: parsePing},
		Spec{Name: "PSYNC", MinArgs: 0, MaxArgs: 2, Flags: FlagNotPropagate, Parse: parsePSync},
		Spec{Name: "REPLCONF", MinArgs: 1, MaxArgs: -1, Flags: FlagNotPropagate, Parse: parseReplconf},
		Spec{Name: "RPUSH", MinArgs: 1, MaxArgs: -1, Flags: FlagWrite, Parse: parseRPush},
		Spec{Name: "SET", MinArgs: 2, MaxArgs: 4, Flags: FlagWrite, Parse: parseSet},
		Spec{Name: "TYPE", MinArgs: 1, MaxArgs: 1, Parse: parseType},
		Spec{Name: "UNWATCH", MaxArgs: 0, Parse: parseUnwatch},
		Spec{Name: "WATCH", MaxArgs: -1, Parse: parseWatch},
		Spec{Name: "XADD", MinArgs: 2, MaxArgs: -1, Flags: FlagWrite, Parse: parseXAdd},
		Spec{Name: "XRANGE", MinArgs: 3, MaxArgs: 3, Parse: parseXRange},
		Spec{Name: "XREAD", MinArgs: 1, MaxArgs: -1, Parse: parseXRead},
		Spec{Name: "WAIT", MinArgs: 0, MaxArgs: -1, Parse: parseWait},

		Spec{Name: "COMMAND", MaxArgs: -1, Flags: FlagNotPropagate, Parse: parseCommand},
	)
}
