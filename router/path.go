package router

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

type PathValidator func(Path) error

type PathType uint8

const (
	PathTypeInvalid = PathType(0)
	PathTypeRoot
	PathTypeCollection
	PathTypeDocument
)

const (
	UIPathWallet = Path("/wallet")
	UIPath
)

var (
	ErrPathStart     = errors.New("path should start with / character")
	ErrPathEnd       = errors.New("path should not end with / character")
	ErrWhiteSpace    = errors.New("path should not contain whitespace character")
	ErrInvalidPathID = errors.New("invalid pathID")
	ErrParsingPath   = errors.New("could not parse the path")
)

// Path format is of the following pattern
// /collectionID/documentID/collectionID/documentID....
// path ending with /collectionID is collectionPath
// path ending with /documentID is document path
type Path string

func (p Path) URL() (*url.URL, error) {
	return url.Parse(p.String())
}

func (p Path) Validate() (err error) {
	_, err = p.URL()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(p.String(), "/") {
		return ErrPathStart
	}
	if strings.HasSuffix(p.String(), "/") {
		return ErrPathEnd
	}
	// If true, it indicates presence of whitespace character
	if regexp.MustCompile(`\s`).MatchString(p.String()) {
		return ErrWhiteSpace
	}
	return err
}

func (p Path) Type() (pathType PathType) {
	if p.Validate() != nil {
		return PathTypeInvalid
	}
	if p == "/" {
		return PathTypeRoot
	}
	if len(strings.Split(p.String(), "/"))%2 == 0 {
		return PathTypeCollection
	}
	return PathTypeDocument
}

// LastSegmentPathID Returns the last segment of the path as string
func (p Path) LastSegmentPathID() (string, error) {
	var pathID string
	err := p.Validate()
	if err != nil {
		return pathID, err
	}
	lastIndex := strings.LastIndex(p.String(), "/")
	if lastIndex < 0 {
		return pathID, ErrParsingPath
	}
	isValid := regexp.MustCompile(`\s`).MatchString(pathID)
	if !isValid {
		return pathID, ErrInvalidPathID
	}
	pathID = string(p[lastIndex:])
	return pathID, nil
}

func (p Path) String() string {
	return string(p)
}
