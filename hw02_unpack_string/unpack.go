package hw02unpackstring

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputStr string) (string, error) {
	var (
		sb                                           strings.Builder
		pr                                           rune
		prIsDigit, prIsSlash, prIsSpec, prIsSpecText bool
		err                                          error
	)

	runes := []rune(inputStr)
	length := utf8.RuneCountInString(inputStr)
	if inputStr == "" {
		return "", nil
	}
	for i := 0; i < length; i++ {
		cr := runes[i]
		crIsDigit := unicode.IsDigit(cr)
		crIsSlash := cr == rune(92) //  '\'

		switch {
		case crIsDigit:
			if !prIsSpec && !prIsSpecText {
				err = currentIsDigitNoSpec(&sb, cr, pr, prIsDigit, prIsSlash, i)
			} else {
				err = currentIsDigitWithSpec(&sb, cr, pr, prIsDigit, prIsSpec, prIsSpecText, i)
			}
		case !crIsDigit && !crIsSlash:
			currentIsSymbol(&sb, pr, prIsDigit, prIsSpec, i)
		case crIsSlash:
			currentIsSlash(&sb, cr, pr, prIsSlash, prIsSpec, prIsSpecText)
		}
		if err != nil {
			return "", ErrInvalidString
		}

		prIsSpec, prIsSpecText = checkSpecSymbols(prIsSlash, crIsSlash, crIsDigit, prIsSpec, prIsSpecText)
		pr = cr
		prIsDigit = crIsDigit
		prIsSlash = crIsSlash
	}
	if !prIsDigit {
		sb.WriteRune(pr)
	}
	return sb.String(), nil
}

func currentIsDigitNoSpec(sb *strings.Builder, cr, pr rune, prIsDigit, prIsSlash bool, i int) error {
	switch {
	case (i == 0):
		return ErrInvalidString
	case prIsDigit:
		return ErrInvalidString
	case !prIsSlash:
		if num, err := strconv.Atoi(string(cr)); err == nil {
			sb.WriteString(strings.Repeat(string(pr), num))
		}
	case prIsSlash:
		sb.WriteRune(cr)
	}

	return nil
}

func currentIsDigitWithSpec(sb *strings.Builder, cr, pr rune, prIsDigit, prIsSpec, prIsSpecText bool, i int) error {
	switch {
	case (i == 0):
		return ErrInvalidString

	case prIsDigit || prIsSpec:
		if num, err := strconv.Atoi(string(cr)); err == nil {
			sb.WriteString(strings.Repeat(string(pr), num-1))
		}
	case prIsSpecText:
		if num, err := strconv.Atoi(string(cr)); err == nil {
			writeSlashSymbol(sb, pr, num)
		}
	}
	return nil
}

func currentIsSymbol(sb *strings.Builder, pr rune, prIsDigit, prIsSpec bool, i int) {
	if !prIsDigit && !prIsSpec && i != 0 {
		sb.WriteRune(pr)
	}
}

func currentIsSlash(sb *strings.Builder, cr, pr rune, prIsSlash, prIsSpec, prIsSpecText bool) {
	switch {
	case prIsSlash && !prIsSpec && !prIsSpecText:
		sb.WriteRune(cr)
	case !prIsSpec && !prIsSpecText:
		sb.WriteRune(pr)
	}
}

func writeSlashSymbol(sb *strings.Builder, symbol rune, num int) {
	buffer := bytes.Buffer{}
	buffer.WriteRune(rune(92))
	buffer.WriteRune(symbol)
	sb.WriteString(string(symbol))
	sb.WriteString(strings.Repeat(buffer.String(), num-1))
	buffer.Reset()
}

func checkSpecSymbols(prIsSlash, crIsSlash, crIsDigit, prIsSpec, prIsSpecText bool) (bool, bool) {
	if prIsSlash && (crIsSlash || crIsDigit) && !prIsSpec {
		prIsSpec = true
		prIsSpecText = false
	} else {
		prIsSpec = false
	}
	if prIsSlash && !crIsSlash && !crIsDigit && !prIsSpecText {
		prIsSpecText = true
	} else {
		prIsSpecText = false
	}
	return prIsSpec, prIsSpecText
}
