// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gossh

import (
    "net"
    "time"

    "golang.org/x/crypto/ssh"
    "github.com/pkg/sftp"
)

type Conn struct {
    Addr        string
    User        string
    Password    string
    KeyBytes    []byte
    sshClient   *ssh.Client
}

func (c *Conn) Connect() error {
    var (
        auth    []ssh.AuthMethod
        err     error
        client  *ssh.Client
    )
    auth, err = c.makeAuth()
    if err != nil {
        return err
    }
    clientConfig := &ssh.ClientConfig{
        User: c.User,
        Auth: auth,
        Timeout: 30 * time.Second,
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
    }
    if client, err = ssh.Dial("tcp", c.Addr, clientConfig); err != nil {
        return err
    }
    c.sshClient = client

    return nil
}

func (c Conn) NewSession() (*Session, error) {
    var (
        session *ssh.Session
        err     error
    )
    if session, err = c.sshClient.NewSession(); err != nil {
        return nil, err
    }
    return &Session{
        session: session,
    }, nil
}

func (c Conn) NewSftpClient() (*SftpClient, error) {
    sftpClient, err := sftp.NewClient(c.sshClient);
    if err != nil {
        return nil, err
    }
    return &SftpClient{
        client: sftpClient,
    }, nil
}

func (c Conn) makeAuth() ([]ssh.AuthMethod, error) {
    var (
        signer  ssh.Signer
        err     error
    )
    auth := make([]ssh.AuthMethod, 0)
    if len(c.KeyBytes) == 0 {
        auth = append(auth, ssh.Password(c.Password))
    } else {
        if c.Password == "" {
            signer, err = ssh.ParsePrivateKey(c.KeyBytes)
        } else {
            signer, err = ssh.ParsePrivateKeyWithPassphrase(c.KeyBytes, []byte(c.Password))
        }
        if err != nil {
            return nil, err
        }
        auth = append(auth, ssh.PublicKeys(signer))
    }
    return auth, nil
}

