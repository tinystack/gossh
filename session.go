// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gossh

import (
    "bytes"

    "golang.org/x/crypto/ssh"
)

type Session struct {
    session     *ssh.Session
    stdout      bytes.Buffer
    stderr      bytes.Buffer
}

func (s *Session) RunCmd(cmd string) error {
    s.session.Stdout = &s.stdout
    s.session.Stderr = &s.stderr

    if err := s.session.Run(cmd); err != nil {
        return err
    }
    return nil
}

func (c *Session) Stdout() []byte {
    return c.stdout.Bytes()
}

func (c *Session) Stderr() []byte {
    return c.stderr.Bytes()
}

func (s *Session) Close() {
    if s.session != nil {
        s.session.Close()
    }
}

