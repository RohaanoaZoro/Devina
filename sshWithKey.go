package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os/user"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

const DefaultTimeout = 30 * time.Second

var HostKeyCallback = ssh.InsecureIgnoreHostKey()

type Client struct {
	SSHClient *ssh.Client
}

func ExecuteScriptInSSHServer(host string, username string, privatekeyPath string, shCmds string) (string, error) {

	client, err := ConnectWithKeyFile(host, username, privatekeyPath)
	if err != nil {
		return "", err
	}

	defer client.Close()

	// var shCmdsArr []string = strings.Split(strings.Replace(shCmds, "\n", "", -1), ";")
	var finalOutput string = ""
	// for i := 0; i < len(shCmdsArr); i++ {
	out, err := execSSHCmd(client, shCmds)
	if err != nil {
		return "", err
	}
	finalOutput += string(out)
	// }

	return finalOutput, nil
}

// Same as ConnectWithKeyFile but allows a custom timeout. If username is empty simplessh will attempt to get the current user.
func ConnectWithKeyFile(host, username, privKeyPath string) (*ssh.Client, error) {
	return ConnectWithKeyFileTimeout(host, username, privKeyPath, DefaultTimeout)
}

func ConnectWithKeyFileTimeout(host, username, privKeyPath string, timeout time.Duration) (*ssh.Client, error) {
	if privKeyPath == "" {
		currentUser, err := user.Current()
		if err == nil {
			privKeyPath = filepath.Join(currentUser.HomeDir, ".ssh", "id_rsa")
		}
	}

	privKey, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return nil, err
	}

	return ConnectWithKeyTimeout(host, username, string(privKey), timeout)
}

// Connect with a private key with a custom timeout. If username is empty simplessh will attempt to get the current user.
func ConnectWithKeyTimeout(host, username, privKey string, timeout time.Duration) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKey([]byte(privKey))
	if err != nil {
		return nil, err
	}

	authMethod := ssh.PublicKeys(signer)

	return connect(username, host, authMethod, timeout)
}

func connect(username, host string, authMethod ssh.AuthMethod, timeout time.Duration) (*ssh.Client, error) {
	if username == "" {
		user, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("Username wasn't specified and couldn't get current user: %v", err)
		}

		username = user.Username
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: HostKeyCallback,
	}

	host = addPortToHost(host)

	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return nil, err
	}
	sshConn, chans, reqs, err := ssh.NewClientConn(conn, host, config)
	if err != nil {
		return nil, err
	}
	c := ssh.NewClient(sshConn, chans, reqs)

	// c := &Client{SSHClient: client}
	return c, nil
}

func addPortToHost(host string) string {
	_, _, err := net.SplitHostPort(host)

	// We got an error so blindly try to add a port number
	if err != nil {
		return net.JoinHostPort(host, "22")
	}

	return host
}

func execSSHCmd(client *ssh.Client, cmd string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		// return err
	}
	defer session.Close()

	return session.CombinedOutput(cmd)
}
