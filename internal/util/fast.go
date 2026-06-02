package util

import (
	"strconv"
	"strings"
	"sync"
)

var (
	percentCache = make([]string, 101)
	dec1Cache    = make([]string, 1000)
	intCache     = make([]string, 10000)
)

func init() {
	for i := range percentCache {
		percentCache[i] = strconv.Itoa(i) + "%"
	}
	for i := range dec1Cache {
		dec1Cache[i] = formatFloat1(float64(i) / 10)
	}
	for i := range intCache {
		intCache[i] = strconv.Itoa(i)
	}
}

func formatFloat1(val float64) string {
	return strconv.FormatFloat(val, 'f', 1, 64)
}

func formatFloat2(val float64) string {
	return strconv.FormatFloat(val, 'f', 2, 64)
}

func FastPercent(val float64) string {
	if val >= 0 && val <= 100 && val == float64(int(val)) {
		return percentCache[int(val)]
	}
	return formatFloat1(val) + "%"
}

func FastPercent1(val float64) string {
	if val >= 0 && val < 1000 {
		idx := int(val * 10)
		if idx >= 0 && idx < 1000 && float64(idx)/10 == val {
			return dec1Cache[idx]
		}
	}
	return formatFloat1(val) + "%"
}

func FastFloat2(val float64) string {
	return formatFloat2(val)
}

func FastFloat1(val float64) string {
	if val >= 0 && val < 1000 {
		idx := int(val * 10)
		if idx >= 0 && idx < 1000 && float64(idx)/10 == val {
			return dec1Cache[idx]
		}
	}
	return formatFloat1(val)
}

func FastInt(val int) string {
	if val >= 0 && val < 10000 {
		return intCache[val]
	}
	return strconv.Itoa(val)
}

func FastInt64(val int64) string {
	return strconv.FormatInt(val, 10)
}

func FastUint64(val uint64) string {
	return strconv.FormatUint(val, 10)
}

func FastCpu(val float64) string {
	return FastPercent1(val)
}

func FastMem(val float64) string {
	return FastPercent1(val)
}

func FastMbPerSec(val float64) string {
	return formatFloat2(val) + " Mb/s"
}

func FastTemp(val float64) string {
	if val >= 0 && val < 1000 {
		idx := int(val)
		if float64(idx) == val {
			return strconv.Itoa(idx) + "°C"
		}
	}
	return formatFloat1(val) + "°C"
}

func FastWatts(val float64) string {
	return formatFloat1(val) + "W"
}

func FastMhz(val int) string {
	return FastInt(val) + " MHz"
}

// percentBuilderMu guards a single reusable strings.Builder. The previous
// implementation shared the builder across FastPctWithPrefix and FastMbUsed
// without synchronization, which raced under -race as soon as two callers
// (e.g. parallel renders) tried to Reset+WriteString at the same time.
var percentBuilderMu sync.Mutex
var percentBuilder strings.Builder

func FastPctWithPrefix(prefix string, val float64) string {
	percentBuilderMu.Lock()
	defer percentBuilderMu.Unlock()
	percentBuilder.Reset()
	percentBuilder.WriteString(prefix)
	percentBuilder.WriteString(": ")
	percentBuilder.WriteString(FastPercent1(val))
	return percentBuilder.String()
}

func FastMbUsed(total, used string) string {
	percentBuilderMu.Lock()
	defer percentBuilderMu.Unlock()
	percentBuilder.Reset()
	percentBuilder.WriteString(used)
	percentBuilder.WriteString(" / ")
	percentBuilder.WriteString(total)
	percentBuilder.WriteString(" MB")
	return percentBuilder.String()
}

func FastMb(value uint64) string {
	return FastInt(int(value)) + " MB"
}

func FastGb(value uint64) string {
	gb := float64(value) / 1024 / 1024 / 1024
	if gb >= 1 {
		return formatFloat1(gb) + " GB"
	}
	mb := float64(value) / 1024 / 1024
	return formatFloat1(mb) + " MB"
}

func FastBytes(value uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	switch {
	case value >= TB:
		return formatFloat1(float64(value)/float64(TB)) + " TB"
	case value >= GB:
		return formatFloat1(float64(value)/float64(GB)) + " GB"
	case value >= MB:
		return formatFloat1(float64(value)/float64(MB)) + " MB"
	case value >= KB:
		return formatFloat1(float64(value)/float64(KB)) + " KB"
	default:
		return FastUint64(value) + " B"
	}
}
