package rdb

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

func ReadRdb(storage storage.Store) {
	directory, ok := storage.GetConfig("dir")
	if !ok {
		fmt.Printf("RDB directory is not set.\n")
		return
	}
	filename, ok := storage.GetConfig("dbfilename")
	if !ok {
		fmt.Printf("RDB filename is not set.\n")
		return
	}

	rdbFilePath := path.Join(directory.Value, filename.Value)

	rdbFile, err := os.Open(rdbFilePath)
	if err != nil {
		fmt.Printf("Could not open %s - %s\n", rdbFilePath, err)
		return
	}

	defer rdbFile.Close()
	rdb := NewRDB()

	rdbHeader, err := ReadHeader(rdbFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	rdb.Header = rdbHeader

	for {
		indicator, err := ReadSingleByte(rdbFile)
		if errors.Is(err, io.EOF) {
			fmt.Printf("Finished reading %s\n", rdbFilePath)
			break
		}

		if err != nil {
			fmt.Printf("Failure during reading of %s - %s\n", rdbFilePath, err)
			return
		}

		switch indicator {
		case 0xFA:
			metadataKey, metadataValue, err := HandleMetadata(rdbFile)
			if err != nil {
				fmt.Println(err)
				return
			}
			rdb.Metadata[metadataKey] = metadataValue
		case 0xFE:
			err := HandleDatabase(rdbFile, &rdb)

			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				return
			}
		case 0xFF:

			for databaseIndex := range rdb.Databases {
				database, ok := rdb.Databases[databaseIndex]
				if !ok {
					return
				}
				for key := range database {
					entry, ok := database[key]
					if !ok {
						return
					}
					storage.Set(key, entry.Value, entry.ExpiresAt)
				}
			}

			return
		default:
			fmt.Println(rdb)
			fmt.Printf("Indicator %x handler not implemented.\n", indicator)
			return
		}
	}
}

func ReadEncodedString(
	rdbFile *os.File,
) (string, error) {
	stringLength, err := ReadSize(rdbFile)
	if err != nil {
		return "", err
	}

	if stringLength > 0 {
		stringValue, err := ReadBytes(rdbFile, int64(stringLength))
		if err != nil {
			return "", err
		}

		return string(stringValue), nil
	}

	switch stringLength {
	case -1:
		integerString, err := ReadBytes(rdbFile, 1)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", int(integerString[0])), nil
	case -2:
		integerString, err := ReadBytes(rdbFile, 2)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", int(integerString[1])<<8|int(integerString[0])), nil
	case -3:
		integerString, err := ReadBytes(rdbFile, 4)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", int(integerString[3])<<24|int(integerString[2])<<16|int(integerString[1])<<8|int(integerString[0])), nil

	default:
		return "", fmt.Errorf("Unknown string encoding")
	}
}

func ReadSize(
	rdbFile *os.File,
) (int, error) {
	_size0, err := ReadSingleByte(rdbFile)
	if err != nil {
		return 0, err
	}

	switch _size0 & 0xC0 {
	case 0x00:
		return int(_size0 & 0x3f), nil

	case 0x40:
		_size1, err := ReadSingleByte(rdbFile)
		if err != nil {
			return 0, nil
		}
		return (int(_size0&0x3F) << 8) | int(_size1), nil

	case 0x80:
		bytes, err := ReadBytes(rdbFile, 3)
		if err != nil {
			return 0, nil
		}

		return int(_size0)<<24 | int(bytes[0])<<16 | int(bytes[1])<<8 | int(bytes[2]), nil

	case 0xC0:
		return -int(_size0&0x3f) - 1, nil

	default:
		return 0, fmt.Errorf("Unexpected case in ReadSize")
	}
}

func ReadSingleByte(
	rdbFile *os.File,
) (byte, error) {
	_bytes, err := ReadBytes(rdbFile, 1)
	return _bytes[0], err
}

func ReadBytes(
	rdbFile *os.File,
	size int64,
) ([]byte, error) {
	_bytes := make([]byte, size)
	_, err := rdbFile.Read(_bytes)
	return _bytes, err
}

func ReadHeader(
	rdbFile *os.File,
) (string, error) {
	header := make([]byte, 9)
	_, err := rdbFile.Read(header)
	return string(header), err
}
