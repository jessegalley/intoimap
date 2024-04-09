package intoimapcli

import (
	"bufio"
	"errors"
	"io"
	"os"

)

func ReadFileIntoString (filename string) (string, error) { 
  var out string 
  var contents []byte

  contents, err := os.ReadFile(filename)
  if err != nil {
    return "", errors.Join(errors.New("couldn't read file"), err)
  }

  out = string(contents)

  return out, nil
}


func ReadInputToString () (string, error) {
  rdr := bufio.NewReader(os.Stdin)
  // out := os.Stdout 
  out := ""
  for {
    switch line, err := rdr.ReadString('\n'); err {

      // If the read succeeded (the read `err` is nil),
      // write out out the uppercased line. Check for an
      // error on this write as we do on the read.
    case nil:
      out += line

      // The `EOF` error is expected when we reach the
      // end of input, so exit gracefully in that case.
    case io.EOF:
      return out, nil

      // Otherwise there's a problem; print the
      // error and exit with non-zero status.
    default:
      return "", errors.Join(errors.New("stdin errror"), err)
    }
  }
}
