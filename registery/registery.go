package reg

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

var hives = map[string]registry.Key{
	"HKCR":                registry.CLASSES_ROOT,
	"HKCU":                registry.CURRENT_USER,
	"HKLM":                registry.LOCAL_MACHINE,
	"HKU":                 registry.USERS,
	"HKCC":                registry.CURRENT_CONFIG,
	"HKEY_CLASSES_ROOT":   registry.CLASSES_ROOT,
	"HKEY_CURRENT_USER":   registry.CURRENT_USER,
	"HKEY_LOCAL_MACHINE":  registry.LOCAL_MACHINE,
	"HKEY_USERS":          registry.USERS,
	"HKEY_CURRENT_CONFIG": registry.CURRENT_CONFIG,
}

func RegHandler(command string) string {
	commands := strings.Split(command, " ")
	var output string
	if commands[0] == "enumkey" {
		if len(commands) == 3 {
			if commands[1] == "-k" {
				tokens := strings.SplitN(commands[2], "\\", 2)
				results, err := EnumKey(tokens[0], tokens[1])
				if err != nil {
					output = err.Error()
				} else {
					for _, res := range results {
						output += res
						output += "\n"
					}
				}
			}
		}
	} else if commands[0] == "createkey" {
		if len(commands) == 3 {
			if commands[1] == "-k" {
				last := commands[2][strings.LastIndex(commands[2], "\\")+1:]
				tokens := strings.SplitN(commands[2], "\\", 2)
				err := CreateSubKey(tokens[0], strings.Split(tokens[1], last)[0], last)
				if err != nil {
					output = err.Error()
				} else {
					output = "Operation completed successfully"
				}
			}
		}
	} else if commands[0] == "deletekey" {
		if len(commands) == 3 {
			if commands[1] == "-k" {
				last := commands[2][strings.LastIndex(commands[2], "\\")+1:]
				tokens := strings.SplitN(commands[2], "\\", 2)
				err := DeleteKey(tokens[0], tokens[1], last)
				if err != nil {
					output = err.Error()
				} else {
					output = "Operation completed successfully"
				}
			}
		}
	} else if commands[0] == "setval" {
		if len(commands) == 9 {
			if commands[1] == "-k" {
				if commands[3] == "-v" {
					if commands[5] == "-t" {
						if commands[7] == "-d" {
							tokens := strings.SplitN(commands[2], "\\", 2)
							err := WriteKey(tokens[0], tokens[1], commands[4], commands[6], commands[8])
							if err != nil {
								output = err.Error()
							} else {
								output = "Operation completed successfully"
							}
						}
					}
				}
			}
		}

	} else if commands[0] == "queryval" {
		if len(commands) == 5 {
			if commands[1] == "-k" {
				if commands[3] == "-v" {
					tokens := strings.SplitN(commands[2], "\\", 2)
					output = ReadKey(tokens[0], tokens[1], commands[4])
				}
			}
		}
	} else if commands[0] == "deleteval" {
		if len(commands) == 5 {
			if commands[1] == "-k" {
				if commands[3] == "-v" {
					tokens := strings.SplitN(commands[2], "\\", 2)
					err := DeleteValue(tokens[0], tokens[1], commands[4])
					if err != nil {
						output = err.Error()
					} else {
						output = "Operation completed successfully"
					}
				}
			}
		}
	} else {
		output = "Unsupported operation"
	}
	return output
}

func openKey(hive string, key string, access uint32) (*registry.Key, error) {
	hiveKey := hives[hive]
	k, err := registry.OpenKey(hiveKey, key, access)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func ReadKey(hive string, key string, value string) string {
	var (
		buf    []byte
		result strings.Builder
	)

	k, err := openKey(hive, key, registry.QUERY_VALUE)
	if err != nil {
		return err.Error()
	}

	_, valType, err := k.GetValue(value, buf)
	if err != nil {
		return err.Error()
	}
	switch valType {
	case registry.BINARY:
		val, _, err := k.GetBinaryValue(value)
		if err != nil {
			return err.Error()
		}
		result.WriteString(fmt.Sprintf("Name: %s\n", value))
		result.WriteString("Type: Binary\n")
		result.WriteString(fmt.Sprintf("Data: %08x", val))
	case registry.SZ:
		fallthrough
	case registry.EXPAND_SZ:
		val, _, err := k.GetStringValue(value)
		if err != nil {
			return err.Error()
		}
		result.WriteString(fmt.Sprintf("Name: %s\n", value))
		result.WriteString("Type: SZ\n")
		result.WriteString(fmt.Sprintf("Data: %s", val))
	case registry.DWORD:
		fallthrough
	case registry.QWORD:
		val, _, err := k.GetIntegerValue(value)
		if err != nil {
			return err.Error()
		}
		result.WriteString(fmt.Sprintf("Name: %s\n", value))
		result.WriteString("Type: QWORD\n")
		result.WriteString(fmt.Sprintf("Data: 0x%08x", val))
	case registry.MULTI_SZ:
		val, _, err := k.GetStringsValue(value)
		if err != nil {
			return err.Error()
		}
		result.WriteString(fmt.Sprintf("Name: %s\n", value))
		result.WriteString("Type: MULTI_SZ\n")
		result.WriteString(fmt.Sprintf("Data: %s", strings.Join(val, "\n")))
	default:
		return fmt.Sprintf("unhandled type: %d", valType)
	}
	return result.String()
}

func WriteKey(hive string, key string, value string, _type string, data string) error {
	k, err := openKey(hive, key, registry.QUERY_VALUE|registry.SET_VALUE|registry.WRITE)
	if err != nil {
		return err
	}

	switch _type {
	case "DWORD":
		d, err := strconv.ParseUint(data, 0, 32)
		if err != nil {
			return err
		}
		err = k.SetDWordValue(value, uint32(d))
	case "QWORD":
		d, err := strconv.ParseUint(data, 0, 32)
		if err != nil {
			return err
		}
		err = k.SetQWordValue(value, d)
	case "SZ":
		err = k.SetStringValue(value, data)
	case "BINARY":
		d, err := hex.DecodeString(data)
		if err != nil {
			return err
		}
		err = k.SetBinaryValue(value, d)
	default:
		return fmt.Errorf("unknow type")
	}

	return err
}

func DeleteValue(hive string, key string, value string) (err error) {
	k, err := openKey(hive, key, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return
	}
	k.DeleteValue(value)
	return
}

func EnumKey(hive string, path string) (results []string, err error) {
	k, err := openKey(hive, path, registry.READ|registry.RESOURCE_LIST|registry.FULL_RESOURCE_DESCRIPTOR)
	if err != nil {
		return
	}
	kInfo, err := k.Stat()
	if err != nil {
		return
	}
	if kInfo.SubKeyCount != 0 {
		return k.ReadSubKeyNames(int(kInfo.SubKeyCount))
	} else {
		return k.ReadValueNames(int(kInfo.ValueCount))
	}
}

func CreateSubKey(hive string, key string, name string) error {
	k, err := openKey(hive, key, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	_, _, err = registry.CreateKey(*k, name, registry.ALL_ACCESS)
	return err
}

func DeleteKey(hive string, key string, name string) (err error) {
	k, err := openKey(hive, key, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return
	}
	err = registry.DeleteKey(*k, name)
	return
}
