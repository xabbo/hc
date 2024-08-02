package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var opts struct {
	vl64  bool
	vl64e bool
	b64   bool
	b64e  bool
	value bool
}

func main() {
	flag.BoolVar(&opts.vl64, "vl64", false, "Decode VL64")
	flag.BoolVar(&opts.vl64e, "vl64e", false, "Encode VL64")
	flag.BoolVar(&opts.b64, "b64", false, "Decode B64")
	flag.BoolVar(&opts.b64e, "b64e", false, "Encode B64")
	flag.BoolVar(&opts.value, "values", false, "Print only values when decoding VL64")
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "hc: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	n := 0
	for _, b := range []bool{opts.vl64, opts.vl64e, opts.b64, opts.b64e} {
		if b {
			n++
		}
	}

	if n != 1 {
		return errors.New("must provide one of -vl64(e) or -b64(e)")
	}

	args := flag.Args()
	if len(args) < 1 {
		return errors.New("no argument provided")
	}

	for _, arg := range args {
		if opts.vl64e || opts.b64e {
			var value int
			value, err = strconv.Atoi(arg)
			if err != nil {
				return fmt.Errorf("not an integer: %q", arg)
			}
			switch {
			case opts.b64e:
				buf := [2]byte{}
				B64Encode(buf[:], value)
				fmt.Println(encodeBytes(buf[:]))
			case opts.vl64e:
				n := VL64EncodeLen(value)
				buf := make([]byte, n)
				VL64Encode(buf, value)
				fmt.Println(encodeBytes(buf))
			}
		} else {
			switch {
			case opts.b64:
				var data []byte
				data, err = parseString(arg)
				if err != nil {
					return
				}
				if len(data) != 2 {
					return fmt.Errorf("b64 must be 2 bytes: %q", arg)
				}
				var v int
				v, err = B64Decode(data)
				if err != nil {
					return
				}
				fmt.Println(v)
			case opts.vl64:
				var data []byte
				data, err = parseString(arg)
				if err != nil {
					return
				}
				for len(data) > 0 {
					var v, n int
					v, n, err = VL64Decode(data)
					if err != nil {
						return
					}
					if opts.value {
						fmt.Println(v)
					} else {
						fmt.Printf("%s: %d\n", string(data[:n]), v)
					}
					data = data[n:]
				}
			}
		}
	}
	return
}

func encodeBytes(b []byte) string {
	var sb strings.Builder
	for _, v := range b {
		escape := false
		if v <= 0x20 || v >= 0x7f {
			escape = true
		} else {
			switch v {
			case '[', ']', '{', '}':
				escape = true
			}
		}
		if escape {
			sb.WriteRune('[')
			sb.WriteString(strconv.Itoa(int(v)))
			sb.WriteRune(']')
		} else {
			sb.WriteRune(rune(v))
		}
	}
	return sb.String()
}

func parseString(s string) (buf []byte, err error) {
	buf = make([]byte, 0, len(s))
	escaping := false
	escapeBuf := []rune{}
	for _, r := range s {
		switch r {
		case '[':
			if escaping {
				err = fmt.Errorf("character '[' inside escape sequence")
				return
			}
			escaping = true
		case ']':
			if !escaping {
				err = fmt.Errorf("character ']' outside escape sequence")
				return
			}
			escaping = false
			var n int
			n, err = strconv.Atoi(string(escapeBuf))
			if err != nil {
				err = fmt.Errorf("invalid escape code: %q", string(escapeBuf))
				return
			}
			if n < 0 || n > 255 {
				err = fmt.Errorf("byte value out of range: %d", n)
				return
			}
			buf = append(buf, byte(n))
			escapeBuf = []rune{}
		default:
			if escaping {
				if r < '0' || r > '9' {
					err = fmt.Errorf("invalid character inside escape sequence")
					return
				}
				escapeBuf = append(escapeBuf, r)
			} else {
				buf = append(buf, byte(r))
			}
		}
	}
	if escaping {
		err = fmt.Errorf("unterminated escape sequence")
	}
	return
}
