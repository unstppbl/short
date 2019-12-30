package validator

import "regexp"

const longLinkMaxLength = 200

// LongLink represents format validator for original long link
type LongLink struct {
	uriPattern *regexp.Regexp
}

// IsValid checks whether the given long link has valid format.
func (l LongLink) IsValid(longLink *string) bool {
	if longLink == nil {
		return false
	}

	if *longLink == "" {
		return false
	}

	if len(*longLink) >= longLinkMaxLength {
		return false
	}

	if !l.uriPattern.MatchString(*longLink) {
		return false
	}

	return true
}

// NewLongLink creates long link validator.
func NewLongLink() LongLink {
	uriPattern := regexp.MustCompile(`^[a-zA-Z]+://.+$`)
	return LongLink{
		uriPattern: uriPattern,
	}
}
