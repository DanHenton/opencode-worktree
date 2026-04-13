package main

func reorderKnownBoolFlags(args []string, knownFlags ...string) []string {
	if len(args) == 0 || len(knownFlags) == 0 {
		return args
	}

	known := make(map[string]struct{}, len(knownFlags))
	for _, flagName := range knownFlags {
		known[flagName] = struct{}{}
	}

	reordered := make([]string, 0, len(args))
	remaining := make([]string, 0, len(args))
	for _, arg := range args {
		if _, ok := known[arg]; ok {
			reordered = append(reordered, arg)
			continue
		}
		remaining = append(remaining, arg)
	}

	return append(reordered, remaining...)
}
