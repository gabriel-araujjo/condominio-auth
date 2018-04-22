package json_patcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Tailor interface {
	Add(obj interface{}, path string, value interface{}) error
	Remove(obj interface{}, path string) error
	Move(obj interface{}, path string, from uint64, to uint64) error
	Replace(obj interface{}, path string, value interface{}) error
}

type Operation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	From  string      `json:"from"`
	Value interface{} `json:"value"`
}

type Patch = []Operation

func Mend(tailor Tailor, patch Patch, obj interface{}) error {
	if obj == nil {
		return errors.New("json_patcher: nil obj")
	}

	if tailor == nil {
		return errors.New("json_patcher: nil tailor")
	}

	for _, op := range patch {
		switch op.Op {
		case "add":
			if err := tailor.Add(obj, op.Path, op.Value); err != nil {
				return err
			}
		case "replace":
			if err := tailor.Replace(obj, op.Path, op.Value); err != nil {
				return err
			}
		case "remove":
			if err := tailor.Remove(obj, op.Path); err != nil {
				return err
			}
		case "move":
			{
				splitIndex := strings.LastIndex(op.From, "/")
				if splitIndex < 0 || len(op.From) < splitIndex+1 {
					return fmt.Errorf("json_patcher: invalid move source %q", op.From)
				}
				if len(op.Path) < splitIndex+1 {
					return fmt.Errorf("json_patcher: invalid move dest %q", op.Path)
				}
				base := op.From[:splitIndex]
				if base != op.Path[:splitIndex] {
					return fmt.Errorf("json_patcher: move source %q and dest %q don't match", op.From[:splitIndex], op.Path[:splitIndex])
				}
				var from, to uint64
				var err error

				if from, err = strconv.ParseUint(op.From[splitIndex+1:], 10, 64); err != nil {
					return err
				}

				if to, err = strconv.ParseUint(op.Path[splitIndex+1:], 10, 64); err != nil {
					return err
				}

				if err := tailor.Move(obj, base, from, to); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("json_patcher: unsuported operation: %q", op.Op)
		}
	}
	return nil
}

func NewPatch(buf []byte) (Patch, error) {
	var patch []Operation
	err := json.Unmarshal(buf, &patch)
	if err != nil {
		return nil, err
	}
	return patch, nil
}
