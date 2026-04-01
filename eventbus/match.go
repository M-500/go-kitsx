package eventbus

import (
	"path"
	"strings"
)

func matches(topic, candidate string, mode MatchMode) bool {
	switch mode {
	case MatchExact:
		return topic == candidate
	case MatchPrefix:
		return strings.HasPrefix(candidate, topic)
	case MatchPattern:
		ok, err := path.Match(topic, candidate)
		return err == nil && ok
	default:
		return topic == candidate
	}
}
