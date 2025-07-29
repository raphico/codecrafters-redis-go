package persistence

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

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
	var ttl *time.Duration
	for {
		pb, err := reader.Peek(1)
		if err != nil {
			return err
		}

		if pb[0] == OpcodeEOF || pb[0] == OpcodeSelectDB {
			break
		}

		if pb[0] == OpcodeExpiryMillis || pb[0] == OpcodeExpirySecs {
			b, err := reader.ReadByte()
			if err != nil {
				return err
			}

			if b == OpcodeExpiryMillis {
				//The expire timestamp is expressed in Unix time,
				// stored as an 8-byte unsigned long, in little-endian
				expiryBytes := make([]byte, 8)
				if _, err := io.ReadFull(reader, expiryBytes); err != nil {
					return err
				}

				expiryMillis := binary.LittleEndian.Uint64(expiryBytes)
				ttl = computeTTL(int64(expiryMillis), true)
			} else {
				// The expire timestamp, expressed in Unix time,
				// stored as an 4-byte unsigned integer, in little-endian
				expiryBytes := make([]byte, 4)
				if _, err := io.ReadFull(reader, expiryBytes); err != nil {
					return err
				}

				expirySeconds := binary.LittleEndian.Uint32(expiryBytes)
				ttl = computeTTL(int64(expirySeconds), false)
			}
		}

		b, err := reader.ReadByte()
		if err != nil {
			return err
		}

		if b != ValueTypeString {
			return fmt.Errorf("unsupported type: 0x%X", b)
		}

		key, err := readStringEncoded(reader)
		if err != nil {
			return err
		}

		value, err := readStringEncoded(reader)
		if err != nil {
			return err
		}

		if ttl != nil {
			store.Set(key, value, ttl)
		}
		ttl = nil
	}

	return nil
}

func computeTTL(expiry int64, isMillis bool) *time.Duration {
	if isMillis {
		ttlMillis := expiry - time.Now().UnixMilli()
		if ttlMillis <= 0 {
			return nil
		}

		ttl := time.Duration(ttlMillis) * time.Millisecond
		return &ttl
	}

	ttlSeconds := expiry - time.Now().Unix()
	if ttlSeconds <= 0 {
		return nil
	}

	ttl := time.Duration(ttlSeconds) * time.Second
	return &ttl
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
