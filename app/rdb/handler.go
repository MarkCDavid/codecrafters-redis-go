package rdb

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

func HandleMetadata(
	rdbFile *os.File,
) (string, string, error) {
	metadataKey, err := ReadEncodedString(rdbFile)
	if err != nil {
		return "", "", fmt.Errorf("Failure during metadata reading - %s\n", err)
	}
	metadataValue, err := ReadEncodedString(rdbFile)
	if err != nil {
		return "", "", fmt.Errorf("Failure during metadata reading - %s\n", err)
	}

	return metadataKey, metadataValue, nil
}

func HandleDatabase(
	rdbFile *os.File,
	rdb *RDB,
) error {
	dbIndex, err := ReadSize(rdbFile)
	if err != nil {
		return err
	}

	rdb.NewDB(dbIndex)

	indicator, err := ReadSingleByte(rdbFile)
	if err != nil {
		return err
	}

	if indicator != 0xFB {
		return fmt.Errorf("Expected indicator 0xFB, got 0x%x", indicator)
	}

	kvSize, err := ReadSize(rdbFile)
	if err != nil {
		return err
	}
	expirySize, err := ReadSize(rdbFile)
	if err != nil {
		return err
	}

	for index := 0; index < kvSize+expirySize; index++ {
		indicator, err = ReadSingleByte(rdbFile)
		if err != nil {
			return err
		}

		var expiresAt *time.Time

		if indicator == 0xFC {
			expiryBytes, err := ReadBytes(rdbFile, 8)
			if err != nil {
				return fmt.Errorf("Failure during metadata reading - %s\n", err)
			}

			expiryUnix := int64(binary.LittleEndian.Uint64(expiryBytes))
			expiry := time.UnixMilli(expiryUnix)
			expiresAt = &expiry
			indicator, err = ReadSingleByte(rdbFile)
			if err != nil {
				return err
			}
		}
		if indicator == 0xFD {
			expiryBytes, err := ReadBytes(rdbFile, 4)
			if err != nil {
				return fmt.Errorf("Failure during metadata reading - %s\n", err)
			}

			expiryUnix := int64(binary.LittleEndian.Uint32(expiryBytes))
			expiry := time.Unix(expiryUnix, 0)
			expiresAt = &expiry
			indicator, err = ReadSingleByte(rdbFile)
			if err != nil {
				return err
			}
		}

		if indicator == 0x00 {
			key, err := ReadEncodedString(rdbFile)
			if err != nil {
				return fmt.Errorf("Failure during metadata reading - %s\n", err)
			}
			value, err := ReadEncodedString(rdbFile)
			if err != nil {
				return fmt.Errorf("Failure during metadata reading - %s\n", err)
			}

			rdb.Databases[dbIndex][key] = storage.Entry{
				Value:     value,
				Type:      "string",
				ExpiresAt: expiresAt,
			}
		}
	}

	return nil
}
