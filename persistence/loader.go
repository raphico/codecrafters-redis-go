package persistence

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/redis-starter-go/session"
	"github.com/codecrafters-io/redis-starter-go/store"
)

const (
	OpcodeSelectDB     = 0xFE // Marks the start of a DB section
	OpcodeResizeDB     = 0xFB // Indicates hash table size info follows
	OpcodeAuxField     = 0xFA // Metadata field (key/value)
	OpcodeEOF          = 0xFF // End of file
	OpcodeExpiryMillis = 0xFC // Expiry in milliseconds
	OpcodeExpirySecs   = 0xFD // Expiry in seconds
	ValueTypeString    = 0x00 // Type marker for plain string key/value
)

func LoadRDB(cfg session.ConfigAccessor, store *store.Store) error {
	file, err := os.Open(cfg.GetRDBPath())
	if err != nil {
		return err
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	// read the header section which is the first 9 bytes
	header := make([]byte, 9)
	if _, err := io.ReadFull(reader, header); err != nil {
		return err
	}

	if string(header[:5]) != "REDIS" {
		return fmt.Errorf("invalid RDB header: %s", header)
	}

	// seek the database section
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("reached the EOF before finding the database section")
			}

			return err
		}

		if b == OpcodeSelectDB {
			break
		}
	}

	// skip over db index
	if _, err := readLengthEncoded(reader); err != nil {
		return err
	}

	// skip over hash table hash information
	b, err := reader.Peek(1)
	if err != nil {
		return err
	}
	if b[0] == OpcodeResizeDB {
		reader.ReadByte()
		if _, err := readLengthEncoded(reader); err != nil {
			return err
		}
		if _, err := readLengthEncoded(reader); err != nil {
			return err
		}
	}

	// read database section
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}

		// EOF / start of a new database subsection
		if b == OpcodeSelectDB || b == OpcodeEOF {
			break
		}

		// only handle strings
		if b != ValueTypeString {
			return fmt.Errorf("unsupported type: 0x%X", b)
		}

		// start of a key value pair
		key, err := readStringEncoded(reader)
		if err != nil {
			return err
		}

		value, err := readStringEncoded(reader)
		if err != nil {
			return err
		}

		store.Set(key, value, nil)
	}

	return nil
}

// redis uses a special variable-length format to encode size instead of raw bytes
// the first 2 bits of the first byte indicate how to interpret the rest
func readLengthEncoded(reader *bufio.Reader) (int, error) {
	first, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// extract the first two bits of the byte
	switch first >> 6 {
	// 2 bits == 00 -> the last 6 bits represent the length
	case 0:
		return int(first & 0x3F), nil

	// 2 bits == 01 -> the remaining 6 bits + 8 bits from second byte
	case 1:
		second, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		return ((int(first) & 0x3F) << 8) | int(second), nil

	// special encoding not supported
	case 2:
		return 0, fmt.Errorf("special encoding not supported")

	// 2 bits == 11 -> remaining 6 bits are ignored
	// Next 4 bytes = big-endian 32-bit unsigned int
	case 3:
		buf := make([]byte, 4)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint32(buf)), nil
	default:
		return 0, fmt.Errorf("invalid length encoding")
	}
}

func readStringEncoded(reader *bufio.Reader) (string, error) {
	length, err := readLengthEncoded(reader)
	if err != nil {
		return "", err
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return "", err
	}

	return string(buf), nil
}
