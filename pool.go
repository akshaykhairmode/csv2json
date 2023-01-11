package main

import "sync"

var anyPool = &sync.Pool{
	New: func() any {
		return new(any)
	},
}

var mapAnyPool = &sync.Pool{
	New: func() any {
		if len(keepCols) > 0 {
			return make(map[string]any, len(keepCols))
		}
		return make(map[string]any, len(headers))
	},
}
