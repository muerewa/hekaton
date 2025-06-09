package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/muerewa/hekaton/structs"
	"gopkg.in/yaml.v2"
)

func LoadConfig(path string) ([]structs.Monitor, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var monitors []structs.Monitor
	err = yaml.Unmarshal(data, &monitors)
	return monitors, err
}

func RunBashCommand(cmd string) (string, error) {
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func Compare(result string, comp *structs.Compare) (bool, error) {
	// Попытка численного сравнения
	if numRes, err := strconv.Atoi(result); err == nil {
		if numComp, ok := comp.Value.(int); ok {
			switch comp.Operator {
			case ">":
				return numRes > numComp, nil
			case ">=":
				return numRes >= numComp, nil
			case "<":
				return numRes < numComp, nil
			case "<=":
				return numRes <= numComp, nil
			}
		}
	}

	// Строковое сравнение
	strComp := fmt.Sprintf("%v", comp.Value)
	switch comp.Operator {
	case "==":
		return result == strComp, nil
	case "!=":
		return result != strComp, nil
	default:
		return false, fmt.Errorf("unsupported operator")
	}
}
