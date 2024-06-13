package util

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

const (
	ATOMIC_UNIT float64 = 1e12
)

// NewPaymentID64 generates a 64 bit payment ID (hex encoded).
// With 64 bit IDs, there is a non-negligible chance of a collision
// if they are randomly generated. It is up to recipients generating
// them to sanity check for uniqueness.
//
// 1 million IDs at 64-bit (simplified): 1,000,000^2 / (2^64 * 2) = ~1/36,893,488 so
// there is a 50% chance a collision happens around 5.06 billion IDs generated.
func NewPaymentID64() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}

// NewPaymentID256 generates a 256 bit payment ID (hex encoded).
func NewPaymentID256() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}

// XMRToDecimal converts a raw atomic XMR balance to a more
// human readable format.
func XMRToDecimal(xmr uint64) string {
	str0 := fmt.Sprintf("%013d", xmr)
	l := len(str0)
	return str0[:l-12] + "." + str0[l-12:]
}

// XMRToFloat64 converts raw atomic XMR to a float64
func XMRToFloat64(xmr uint64) float64 {
	return float64(xmr) / ATOMIC_UNIT
}

// Float64ToXMR converts a float64 to a raw atomic XMR
func Float64ToXMR(xmr float64) uint64 {
	return uint64(xmr * ATOMIC_UNIT)
}

// StringToXMR converts a string to a raw atomic XMR
func StringToXMR(xmr string) (uint64, error) {
	f, err := strconv.ParseFloat(xmr, 64)
	if err != nil {
		return 0, err
	}
	return uint64(f * ATOMIC_UNIT), nil
}

// JSON Mapping Related
func parseJson[R any](data []byte) (*R, error) {
	var result R
	if err := json.Unmarshal(data, &result); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &result, nil
}

func ParseJsonString[R any](str *string) (*R, error) {
	if str == nil {
		return nil, errors.New("util.ParseJsonString: string pointer can't be null")
	}

	return parseJson[R]([]byte(*str))
}

func ParseResponse[R any](body io.Reader) (*R, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return parseJson[R](data)
}
