package decode

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type YamlBlock = map[string]interface{}

//将 Unity 由yaml组成的资源文件解析为多个通用yaml类型
func UnmarshalAsset(r io.Reader, fn func(assetType int, fileID string, block YamlBlock, err error) (goon bool)) error {
	return SpliteAssetFile(r, func(assetType int, fileID, block string) bool {
		br := strings.NewReader(block)
		b, uerr := UnmarshalCommon(br)
		return fn(assetType, fileID, b, uerr)
	})
}

var ErrNoYAMLAsset = errors.New("not unity yaml asset")

//将 Unity 由yaml组成的资源文件拆分成多个yaml块
func SpliteAssetFile(r io.Reader, fn func(assetType int, fileID string, block string) (goon bool)) error {
	reader := bufio.NewReader(r)

	var err error
	var line []byte

	line, _, err = reader.ReadLine()
	if err != nil {
		return ErrNoYAMLAsset
	}
	if string(bytes.TrimSpace(line)) != `%YAML 1.1` {
		return ErrNoYAMLAsset
	}

	line, _, err = reader.ReadLine()
	if err != nil {
		return ErrNoYAMLAsset
	}
	if string(bytes.TrimSpace(line)) != `%TAG !u! tag:unity3d.com,2011:` {
		return ErrNoYAMLAsset
	}

	block := make([]byte, 0)
	var fileID string
	var assetType int

	for err == nil {
		line, _, err = reader.ReadLine()
		if err == io.EOF {
			err = nil
			break
		}

		if len(line) == 0 {
			continue
		}

		if string(line[0:3]) == "---" {
			if len(block) > 0 {
				if !fn(assetType, fileID, string(block)) {
					return err
				}
				block = block[:0]
			}

			assetType, fileID, err = parseHeader(line)
		} else {
			block = append(block, line...)
			block = append(block, '\n')
		}
	}

	return err
}

//通用解析yaml
func UnmarshalCommon(r io.Reader) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := Unmarshal(r, out)
	return out, err
}

//特定类型解析yaml
func Unmarshal(r io.Reader, out interface{}) error {
	decoder := yaml.NewDecoder(r)
	return decoder.Decode(out)
}

func parseHeader(str []byte) (assetType int, fileID string, err error) {
	finds := regexp.MustCompile(`[0-9]+`).FindAll(str, -1)
	assetType, err = strconv.Atoi(string(finds[0]))
	fileID = string(finds[1])
	return
}
