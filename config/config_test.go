package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	temp, _ := ioutil.TempFile("", "config.yml")
	temp.WriteString("port: 1337\ntds:\n")
	defer os.Remove(temp.Name())
	c, err := Load(temp.Name())
	assert.NoError(t, err)
	assert.Equal(t, 1337, c.Port)
}

func TestSave(t *testing.T) {
	temp, _ := ioutil.TempFile("", "config.yml")
	defer os.Remove(temp.Name())
	c, _ := Load(temp.Name())
	c.Port = 1337
	c.Save()
	c2, _ := Load(temp.Name())
	assert.Equal(t, 1337, c2.Port)
}
