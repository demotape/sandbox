package sandbox

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	//"strings"
)

type KeyPair struct {
	Pub    []byte
	Secret []byte
}

var (
	keyName = "demotape_key"
)

func GenerateKeys(tempDir string) (*KeyPair, error) {
	var args []string

	keyPath := fmt.Sprintf("%s/%s", tempDir, keyName)

	path, err := exec.LookPath("ssh-keygen")

	if err != nil {
		fmt.Println("ssh-keygen command not found")
		return nil, err
	}

	args = append(args, "-t", "rsa", "-b", "768", "-N", "", "-f", keyPath)
	output, err := exec.Command(path, args...).CombinedOutput()

	if err != nil {
		fmt.Printf("%s", output)
		return nil, err
	}

	secKey, _ := os.Open(keyPath)
	pubKey, _ := os.Open(keyPath + ".pub")

	sec, err := ioutil.ReadAll(secKey)
	pub, err := ioutil.ReadAll(pubKey)

	authorizedKeys := fmt.Sprintf("%s/authorized_keys", tempDir)
	ioutil.WriteFile(authorizedKeys, pub, 0400)

	keyPair := KeyPair{
		Secret: sec,
		Pub:    pub,
	}

	return &keyPair, nil
}

func CreateAuthorizedKeysFile(dir string, pubKey []byte) error {
	fileName := "authorized_keys"

	err := ioutil.WriteFile(dir+"/"+fileName, pubKey, 0644)
	if err != nil {
		return err
	}
	return nil
}

func createTempDir(baseDir string, prefix string) (string, error) {
	dir, err := ioutil.TempDir(baseDir, prefix)

	if err != nil {
		fmt.Println("fail to create tempdir")
		return "", err
	}

	return dir, nil
}

/*
func generateSSHKeys(dir string) error {
	var err error
	keyName := "mykey"

	keyPair, err := sshkey.GenerateKeys(dir + "/" + keyName)
	if err != nil {
		fmt.Println("got error")
		return err
	}

	err = sshkey.CreateAuthorizedKeysFile(dir, keyPair.Pub)
	if err != nil {
		fmt.Println("Fail to create authorized_keys")
		return err
	}

	fmt.Println(dir)
	return nil
}
*/
