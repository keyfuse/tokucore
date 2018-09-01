// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"github.com/tokublock/tokucore/xerror"
)

const (
	maxInt32               = 1<<31 - 1
	minInt32               = -1 << 31
	defaultScriptNumMaxLen = 4
)

// ScriptNum --
type ScriptNum int64

// Bytes -- returns the number serialized as a little endian with a sign bit
// Example encodings:
//       127 -> [0x7f]
//      -127 -> [0xff]
//       128 -> [0x80 0x00]
//      -128 -> [0x80 0x80]
//       129 -> [0x81 0x00]
//      -129 -> [0x81 0x80]
//       256 -> [0x00 0x01]
//      -256 -> [0x00 0x81]
//     32767 -> [0xff 0x7f]
//    -32767 -> [0xff 0xff]
//     32768 -> [0x00 0x80 0x00]
//    -32768 -> [0x00 0x80 0x80]
func (n ScriptNum) Bytes() []byte {
	// Zero encodes as an empty byte slice.
	if n == 0 {
		return nil
	}
	neg := n < 0
	absvalue := n
	if neg {
		absvalue = -n
	}

	// Encode to little endian.
	result := make([]byte, 0, 9)
	for absvalue > 0 {
		result = append(result, byte(absvalue&0xff))
		absvalue >>= 8
	}

	// - If the most significant byte is >= 0x80 and the value is positive, push a
	// new zero-byte to make the significant byte < 0x80 again.

	// - If the most significant byte is >= 0x80 and the value is negative, push a
	// new 0x80 byte that will be popped off when converting to an integral.

	// - If the most significant byte is < 0x80 and the value is negative, add
	// 0x80 to it, since it will be subtracted and interpreted as a negative when
	// converting to an integral.
	last := result[len(result)-1]
	if (last & 0x80) > 0 {
		if neg {
			result = append(result, 0x80)
		} else {
			result = append(result, 0x00)
		}
	} else if neg {
		result[len(result)-1] |= 0x80
	}
	return result
}

// Int32 -- returns the script number clamped to a valid int32.
func (n ScriptNum) Int32() int32 {
	if n > maxInt32 {
		return maxInt32
	}
	if n < minInt32 {
		return minInt32
	}
	return int32(n)
}

// checkMinimalDataEncoding --
// Returns whether or not the passed byte array adheres to the minimal encoding requirements.
func checkMinimalDataEncoding(v []byte) error {
	if len(v) == 0 {
		return nil
	}

	// Check that the number is encoded with the minimum possible
	// number of bytes.
	//
	// If the most-significant-byte - excluding the sign bit - is zero
	// then we're not minimal.  Note how this test also rejects the
	// negative-zero encoding, [0x80]
	if v[len(v)-1]&0x7f == 0 {
		// One exception: if there's more than one byte and the most
		// significant bit of the second-most-significant-byte is set
		// it would conflict with the sign bit.  An example of this case
		// is +-255, which encode to 0xff00 and 0xff80 respectively.
		// (big-endian)
		if len(v) == 1 || v[len(v)-2]&0x80 == 0 {
			return xerror.NewError(Errors, ER_SCRIPTNUM_MINIMAL_DATA, v)
		}
	}
	return nil
}

// MakeScriptNum -- convert the byte to script num.
func MakeScriptNum(v []byte, scriptNumLen int) (ScriptNum, error) {
	// Interpreting data requires that it is not larger than
	// the the passed scriptNumLen value
	if len(v) > scriptNumLen {
		return 0, xerror.NewError(Errors, ER_SCRIPTNUM_TOO_BIG, v, len(v), scriptNumLen)
	}

	// Enforce minimal encoded if requested.
	if err := checkMinimalDataEncoding(v); err != nil {
		return 0, err
	}

	// Zero.
	if len(v) == 0 {
		return 0, nil
	}

	// Decode from little endian
	var result int64
	for i, val := range v {
		result |= int64(val) << uint8(8*i)
	}

	last := v[len(v)-1]
	if (last & 0x80) > 0 {
		result &= ^(int64(0x80) << uint8(8*(len(v)-1)))
		return ScriptNum(-result), nil
	}
	return ScriptNum(result), nil
}
