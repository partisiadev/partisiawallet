package router

import (
	"errors"
	"regexp"
	"strings"
)

const DefaultPathParamPattern = `/:[.]+[^/]`

var (
	ErrPathTypeRoot          = errors.New("path type root not allowed")
	ErrPathAlreadyRegistered = errors.New("path already registered")
)

type Config struct {
	// If Path has wildcards, then pattern is required
	Path
	// Pattern is called with regexp.MatchString(pattern, Path.String())
	// and is used to derive the concrete path
	Pattern  string
	OnActive func(concretePath Path)
	Tag      string
}

// MatchesPath the returned string is the concrete path,
//
//	it's unreliable if there's an error. (Mostly it will be empty if error)
func (c Config) MatchesPath(pathToMatch Path) (bool, error) {
	cPath := c.Path
	cPathStr := c.Path.String()
	if cPath.Type() == PathTypeRoot || pathToMatch.Type() == PathTypeRoot {
		return false, ErrPathTypeRoot
	}
	err := cPath.Validate()
	if err != nil {
		return false, err
	}
	err = pathToMatch.Validate()
	if err != nil {
		return false, err
	}
	cPathArr := strings.Split(cPathStr, "/")
	matchFound := true
	pathToMatchArr := strings.Split(pathToMatch.String(), "/")
	minLen := len(cPathArr)
	if minLen > len(pathToMatchArr) {
		minLen = len(pathToMatchArr)
	}
	for i := 0; i < minLen; i++ {
		eachCPathSeg := cPathArr[i]
		eachPathToMatchSeg := pathToMatchArr[i]
		isMatch := eachCPathSeg == eachPathToMatchSeg
		if !isMatch {
			isPattern, err := regexp.MatchString(c.Pattern, eachCPathSeg)
			if err != nil {
				return false, err
			}
			if isPattern {
				isMatch, err = regexp.MatchString(c.Pattern, eachPathToMatchSeg)
			}
		}
		if !isMatch {
			matchFound = false
		}
	}
	return matchFound, nil
}
