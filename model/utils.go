package model

import (
	"errors"
	"strings"
)

var ErrInvalidDataFormat = errors.New("invalid data format")

func EncodeStringArray(pData []string) (string, error) {
	// check empty
	if len(pData) == 0 {
		return "", nil
	}

	// check data format
	for _, v := range pData {
		if strings.Contains(v, "|") {
			return "", ErrInvalidDataFormat
		}
	}

	rst := ""
	for i := 0; i < len(pData)-1; i++ {
		rst += pData[i] + "|"
	}
	rst += pData[len(pData)-1]
	return rst, nil
}

func DecodeStringArray(pData string) []string {
	if pData == "" {
		return []string{}
	}
	return strings.Split(pData, "|")
}

func AppendStringArray(pArray, pNew string) (string, error) {
	if strings.Contains(pNew, "||") || pNew[0] == '|' || pNew[len(pNew)-1] == '|' {
		return "", ErrInvalidDataFormat
	}
	if pArray == "" {
		return pNew, nil
	}
	if pNew == "" {
		return pArray, nil
	}
	return pArray + "|" + pNew, nil
}

func StringArrayContains(pArray, pElement string) bool {
	s := DecodeStringArray(pArray)
	for _, v := range s {
		if v == pElement {
			return true
		}
	}
	return false
}

func EncodeSsKeyPair(pKey, pPair string) (string, error) {
	if strings.Contains(pKey, "=>") || strings.Contains(pPair, "=>") ||
		pKey == "" || pPair == "" {
		return "", ErrInvalidDataFormat
	}
	return pKey + "=>" + pPair, nil
}

func DecodeSsKeyPair(pPair string) (string, string) {
	r := strings.Split(pPair, "=>")
	if len(r) == 0 {
		return "", ""
	}
	if len(r) == 1 {
		return r[0], ""
	}
	return r[0], r[1]
}
