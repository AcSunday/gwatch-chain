package utils

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	bin "github.com/gagliardetto/binary"
	"reflect"
)

// for solana
const ProgramDataPrefix = "Program data: "

// UnmarshalBorsh for solana
func UnmarshalBorsh(dataBytes []byte, obj any) error {
	err := bin.UnmarshalBorsh(obj, dataBytes[8:])
	if err != nil {
		return err
	}
	if obj == nil {
		return fmt.Errorf("object is nil, dataBytes: %v", dataBytes)
	}

	return nil
}

// SigHash for solana
func SigHash(instructName string) [8]byte {
	sign := []byte("global:" + instructName)
	hash := sha256.Sum256(sign)
	return [8]byte(hash[:8])
}

// SigHashEvent for solana
func SigHashEvent(instructName string) [8]byte {
	sign := []byte("event:" + instructName)
	hash := sha256.Sum256(sign)
	return [8]byte(hash[:8])
}

func GetMethodHash(obj any) []byte {
	methodHash := SigHashEvent(reflect.TypeOf(obj).Name())
	return methodHash[:]
}

func GetMethodHashString(obj any) string {
	methodHash := SigHashEvent(reflect.TypeOf(obj).Name())
	buffer := bytes.NewBuffer(methodHash[:])
	return buffer.String()
}
