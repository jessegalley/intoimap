package intoimap

import (
	"errors"
	"fmt"

	"github.com/emersion/go-imap/v2/imapclient"
)

type ImapURI struct {
  host string
  port int
  user string
  pass string
}

// Connect takes in the host and port of an imap server, and a user and password,
// returning a pointer to an imapclient.Client 
func StartSession(host string, port int, user string, pass string) (*imapclient.Client, error) {
  // TODO: since there is not a clear option in v2/imapclient for plaintext
  // we will instead explicitly return an error on 143. Later this can be 
  // refactored to allow for a plaintext option.
  if port == 143 {
    return nil, errors.New("plaintext/143 unsupported, please use 993")
  }

  c, err := imapclient.DialTLS(fmt.Sprintf("%s:%d", host, port), nil) //TODO: add option for starttls/plain
  if err != nil {
    return nil, errors.Join(errors.New("Failed to connect"), err)
  }

  loginerr := c.Login(user, pass).Wait()
  if loginerr != nil {
    return nil, errors.Join(errors.New("failed to LOGIN"), err)
  }

  return c, nil
}

func AppendMsg(c *imapclient.Client, mailbox string, msg string) error {
  messageBuff := []byte(msg)
  messageLen := len(messageBuff)

  cmd := c.Append(mailbox, int64(messageLen), nil)
  if _, err := cmd.Write(messageBuff); err != nil {
    return errors.Join(errors.New("couldn't append message"), err)
  }

  if err := cmd.Close(); err != nil {
    return errors.Join(errors.New("couldn't append/close message"), err)
  }

  if _, err := cmd.Wait(); err != nil {
    return errors.Join(errors.New("append command failed!"), err)
  }

  return nil
}



