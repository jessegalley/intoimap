package main

import (
	"fmt"
	"log"
	"os"

	intoimapcli "github.com/jessegalley/intoimap/internal"
	intoimap "github.com/jessegalley/intoimap/pkg"
	flag "github.com/spf13/pflag"
)

const (
  semVer = "0.1.0"
  progName = "intoimap"
)

func init() {
  // TODO: set up slog logger for proper output
  // intoimapcli.SetupCLIAgs()
  flag.StringP("file", "f", "", "file input instead of stdin")
  flag.StringP("mailbox", "m", "INBOX", "mailbox to which message is appended")
  flag.IntP("port", "p", 993, "imap tls port")
  flag.BoolP("verbose", "v", false, "verbose output")
  flag.BoolP("version", "V", false, "print version and exit")
  flag.Bool("debug", false, "print debug output")
  flag.CommandLine.MarkHidden("debug")
  flag.CommandLine.SortFlags = false
  flag.Usage =  func() {
    fmt.Fprintf(os.Stderr, "usage: %s [OPTS] <host> <user> <pass>\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "\n")
    flag.PrintDefaults()
  }
  flag.Parse()

  verflag, _ := flag.CommandLine.GetBool("version")
  if verflag {
    fmt.Println(progName)
    fmt.Println("v"+semVer)
    os.Exit(1)
  }

  expectedArgs := 3
  if(len(flag.Args()) != expectedArgs) {
    flag.Usage()
    os.Exit(2)
  }

}

func main() {
  // setup args and opts
  inputstr := ""
  if f, _ := flag.CommandLine.GetString("file"); f != "" {
    // fmt.Println("getting message content from file")
    filestr, err := intoimapcli.ReadFileIntoString(f)
    if err != nil {
      log.Fatal(err)
    }
    inputstr = filestr
  } else {
    stdinstr, err := intoimapcli.ReadInputToString() 
    if err != nil {
      log.Fatal(err)
    }
    inputstr = stdinstr
  }
  // fmt.Println(inputstr)

  // setup logger/output 
  // parse stdin
  // connect imap 
  // fmt.Println(flag.Args())
  host := flag.Arg(0)
  port, _ := flag.CommandLine.GetInt("port") //TODO: does this err need to happen with defaults?
  user := flag.Arg(1)
  pass := flag.Arg(2)

  c, err := intoimap.StartSession(host, port, user, pass)
  if err != nil {
    log.Fatal(err)
  }

  defer c.Close()
  // append imap
  // msg := generateTestMsg()
  msg := inputstr 
  // TODO: parse rfc5322 message with net/mail before attempting append
  appenderr := intoimap.AppendMsg(c, "INBOX", msg)
  if appenderr != nil {
    log.Fatal("err:", appenderr)
  }

  fmt.Println("done")
}



func generateTestMsg() string { 
  headers := make(map[string]string)
  headers["From"] = "Jesse Galley <jgalley@example.com>"
  headers["To"] = "qatest-jg1 <qatest-jg1@example.net>"
  headers["Subject"] = "Test Email, Please Ignore"

  body := "This is a system test, please ignore it and go about your day. Thank You."

  var headerString string
  for header, content := range headers {
    headerString += header + ": " + content + "\r\n"
  }

  messageString := headerString + "\r\n" + body + "\r\n"
 
  return messageString
}


