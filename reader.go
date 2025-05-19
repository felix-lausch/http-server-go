package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (reader rot13Reader) Read(p []byte) (n int, err error) {
	n, err = reader.r.Read(p)

	for i := range n {
		val := p[i]

		if val >= 97 && val <= 122 {
			zeroed := val - 97

			zeroed = (zeroed + 13) % 26

			p[i] = zeroed + 97
		}
	}

	return n, err
}

func main2() {
	s := strings.NewReader("lbh penpxrq gur pbqr!")
	r := rot13Reader{s}

	io.Copy(os.Stdout, &r)
}

// func (reader rot13Reader) FuckMe() string {
// 	return "FUCK ME HARDER"
// }
// func (reader MyReader) Read(bytes []byte) (int, error) {
// 	for i := range bytes {
// 		bytes[i] = 'A'
// 	}

// 	return 1, nil
// }

// type BrainRotError int

// func (e BrainRotError) Error() string {
// 	return "OMG WHAT TH FUCKK"
// }
