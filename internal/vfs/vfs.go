package vfs

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type BlockType string

const (
	BlockTypeSingle BlockType = "single"
	BlockTypeAdd    BlockType = "add"
	BlockTypeMerge  BlockType = "merge"
)

type OctalMode os.FileMode

func (o OctalMode) String() string {
	return strconv.FormatUint(uint64(o), 8)
}

func (o OctalMode) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: o.String(),
	}, nil
}

func (o *OctalMode) UnmarshalYAML(value *yaml.Node) error {
	v, err := strconv.ParseUint(value.Value, 8, 32)
	if err != nil {
		return err
	}

	*o = OctalMode(v)

	return nil
}

type Directory struct {
	Name  string
	Mode  OctalMode
	Files []File      `yaml:",omitempty"`
	Dirs  []Directory `yaml:",omitempty"`
}

type File struct {
	Name   string
	Mode   OctalMode
	Blocks []Block `yaml:",omitempty"`
}

type Block struct {
	Name string
	Type BlockType
	Data []string `yaml:",omitempty"`
}
