package middlewares

import (
	"bitbucket.org/HeilaSystems/log"
	"context"
	"sort"
)

// LoggerIncomingContextExtractor creates a context extractor for logger
// This is useful if you want to add different fields from gRPC incoming metadata.MD
func LoggerIncomingContextExtractor(extractedHeaders []string) log.ContextExtractor {

	sort.Slice(extractedHeaders, func(i, j int) bool {
		return len(extractedHeaders[i]) < len(extractedHeaders[j])
	})
	return headerPrefixes(extractedHeaders).Extract
}

type headerPrefixes []string // if this slice will be very large it's better to build a trie map

func (h headerPrefixes) Extract(ctx context.Context) map[string]interface{} {
	var output = make(map[string]interface{})
	for _, s := range h {
		output[s] = ctx.Value(s)
	}
	return output
}
