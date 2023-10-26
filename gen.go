// Copyright 2021-2022, Offchain Labs, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type HardHatArtifact struct {
	Format       string        `json:"_format"`
	ContractName string        `json:"contractName"`
	SourceName   string        `json:"sourceName"`
	Abi          []interface{} `json:"abi"`
	Bytecode     string        `json:"bytecode"`
}

type moduleInfo struct {
	contractNames []string
	abis          []string
	bytecodes     []string
}

func (m *moduleInfo) addArtifact(artifact HardHatArtifact) {
	abi, err := json.Marshal(artifact.Abi)
	if err != nil {
		log.Fatal(err)
	}
	m.contractNames = append(m.contractNames, artifact.ContractName)
	m.abis = append(m.abis, string(abi))
	m.bytecodes = append(m.bytecodes, artifact.Bytecode)
}

func (m *moduleInfo) exportABIs(dest string) {
	for i, name := range m.contractNames {
		path := filepath.Join(dest, name+".abi")
		abi := m.abis[i] + "\n"

		// #nosec G306
		err := os.WriteFile(path, []byte(abi), 0o644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func capitalizeFirstLetter(input string) string {
	if input == "" {
		return input
	}

	// Split the input into the first letter and the rest of the string
	firstLetter := string(input[0])
	restOfInput := input[1:]

	// Capitalize the first letter and concatenate it with the rest of the string
	capitalizedFirstLetter := strings.ToUpper(firstLetter)
	capitalizedString := capitalizedFirstLetter + restOfInput

	return capitalizedString
}

func uncapitalizeFirstLetter(input string) string {
	if input == "" {
		return input
	}

	// Split the input into the first letter and the rest of the string
	firstLetter := string(input[0])
	restOfInput := input[1:]

	// Capitalize the first letter and concatenate it with the rest of the string
	capitalizedFirstLetter := strings.ToLower(firstLetter)
	capitalizedString := capitalizedFirstLetter + restOfInput

	return capitalizedString
}

type Param struct {
	Name string
	Type string
}

func createTemplate(dirPath string) {
	path := filepath.Join(dirPath, "precompiles", "fheOps.go")
	outPath := filepath.Join(dirPath, "FheOps_gen.sol")
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()
	defer outFile.Close()

	outFile.WriteString(`// SPDX-License-Identifier: BSD-3-Clause-Clear

pragma solidity >=0.4.21 <0.9.0;
			
interface FheOps {
`)

	// Create a bufio scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line in the file
	hadBanner := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=====") {
			hadBanner = true
			continue
		}

		if !hadBanner {
			continue
		}

		Name := ""
		Ret := ""
		var params []Param
		if strings.Contains(line, "func ") {
			chunks := strings.Fields(line)
			inner := strings.Split(chunks[1], "(")
			Name = inner[0]
			param := Param{
				Name: "",
				Type: "",
			}
			Break := false
			if inner[1] != ")" {
				param.Name = inner[1]
				if strings.Contains(chunks[2], ")") {
					param.Type = strings.Replace(chunks[2], ")", "", 1)
					Break = true
				} else {
					param.Type = strings.Replace(chunks[2], ",", "", 1)
				}
				params = append(params, param)

			}

			i := 3
			for ; i < len(chunks); i += 2 {
				if Break {
					break
				}
				param.Name = strings.Replace(chunks[i], ",", "", 1)
				if strings.Contains(chunks[i+1], ")") {
					param.Type = strings.Replace(chunks[i+1], ")", "", 1)
					Break = true
				} else {
					param.Type = strings.Replace(chunks[i+1], ",", "", 1)
				}
				params = append(params, param)
			}

			if strings.Contains(chunks[i], "(") {
				Ret = strings.Replace(chunks[i], "(", "", 1)
				Ret = strings.Replace(Ret, ",", "", 1)
			}

			outLine := "\tfunction " + uncapitalizeFirstLetter(Name) + "("
			for count, param := range params {
				if param.Type == "[]byte" {
					param.Type = "bytes memory"
				}

				outLine += param.Type + " " + param.Name
				if count < len(params)-1 {
					outLine += ", "
				}
			}

			outLine += ") external view"
			if Ret != "" {
				if Ret == "[]byte" {
					Ret = "bytes memory"
				}
				outLine += " returns (" + Ret + ")"
			}

			outLine += ";\n"

			_, err = outFile.WriteString(outLine)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}

		}
	}

	outFile.WriteString("}")

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

}

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("bad path")
	}
	parent := filepath.Dir(filename)

	if os.Args[1] == "1" {
		createTemplate(parent)
		return
	}

	filePaths, err := filepath.Glob(filepath.Join(parent, "precompiles", "artifacts", "contracts", "*", "*.json"))
	fmt.Printf("%s, %+v\n", parent, filePaths)
	if err != nil {
		log.Fatal(err)
	}

	modules := make(map[string]*moduleInfo)

	for _, path := range filePaths {
		if strings.Contains(path, ".dbg.json") {
			continue
		}

		dir, file := filepath.Split(path)
		dir, _ = filepath.Split(dir[:len(dir)-1])
		_, module := filepath.Split(dir[:len(dir)-1])
		module = strings.ReplaceAll(module, "-", "_")
		module += "gen"

		name := file[:len(file)-5]

		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatal("could not read", path, "for contract", name, err)
		}

		artifact := HardHatArtifact{}
		if err := json.Unmarshal(data, &artifact); err != nil {
			log.Fatal("failed to parse contract", name, err)
		}

		if strings.Contains(fmt.Sprintf("%+v", artifact), "BatchPoster") {
			fmt.Printf("File: %v, Module %v, Artifact %+v\n", path, module, artifact)
		}
		modInfo := modules[module]
		if modInfo == nil {
			modInfo = &moduleInfo{}
			modules[module] = modInfo
		}
		modInfo.addArtifact(artifact)
	}

	abi := ""

	for module, info := range modules {

		code, err := bind.Bind(
			info.contractNames,
			info.abis,
			info.bytecodes,
			nil,
			module,
			bind.LangGo,
			nil,
			nil,
		)

		if err != nil {
			log.Fatal(err)
		}

		abi = info.abis[0]

		folder := filepath.Join(parent, "go", module)

		err = os.MkdirAll(folder, 0o755)
		if err != nil {
			log.Fatal(err)
		}

		/*
			#nosec G306
			This file contains no private information so the permissions can be lenient
		*/
		fmt.Printf("%s\n", filepath.Join(folder, module+".go"))
		err = os.WriteFile(filepath.Join(folder, module+".go"), []byte(code), 0o644)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("successfully generated go abi files")

	blockscout := filepath.Join(parent, "nitro-testnode", "blockscout", "init", "data")
	if _, err := os.Stat(blockscout); err != nil {
		fmt.Println("skipping abi export since blockscout is not present")
	} else {
		modules["precompilesgen"].exportABIs(blockscout)
		modules["node_interfacegen"].exportABIs(blockscout)
		fmt.Println("successfully exported abi files")
	}

	var functions []Function

	err = json.Unmarshal([]byte(abi), &functions)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return
	}

	var operations []Operation

	for _, f := range functions {
		parameters := ""
		innerParameters := ""
		for count, arg := range f.Inputs {
			t := arg.Type
			if t == "bytes" {
				t = "[]byte"
			}

			parameters += fmt.Sprintf("%s %s", arg.Name, t)
			innerParameters += arg.Name
			if count < (len(f.Inputs) - 1) {
				parameters += ", "
				innerParameters += ", "
			}
		}

		ret := "void"
		if len(f.Outputs) != 0 {
			ret = f.Outputs[0].Type
			if ret == "bytes" {
				ret = "[]byte"
			}
		}

		operations = append(operations, Operation{
			Name:        capitalizeFirstLetter(f.Name),
			Inputs:      parameters,
			InnerInputs: innerParameters,
			ReturnType:  ret,
		})
	}

	// operations := []Operation{
	// 	{"Add", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Verify", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Reencrypt", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Lte", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Sub", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Mul", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Lt", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"OptReq", "input []byte, inputLen uint32", "input, inputLen", "void"},
	// 	{"Cast", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"TrivialEncrypt", "input []byte", "input", "[]byte"},
	// 	{"And", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Or", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Xor", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Eq", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Gte", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Gt", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Shl", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Shr", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Ne", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Min", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Max", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Neg", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Not", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Decrypt", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// 	{"Div", "input []byte, inputLen uint32", "input, inputLen", "[]byte"},
	// }

	f := "FheOps_gen.go"
	file, err := os.Create(f)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(`package precompiles

import (
	fheos "github.com/fhenixprotocol/fheos/precompiles"
)

type FheOps struct {
	Address addr // 0x80
}
`)
	defer file.Close()
	for _, op := range operations {
		template := GenerateFHEOperationTemplate(op.ReturnType)
		err = template.Execute(file, op)
		if err != nil {
			fmt.Println("Error writing")
			return
		}
	}
}

type Operation struct {
	Name        string
	Inputs      string
	InnerInputs string
	ReturnType  string
}

type Function struct {
	Inputs          []Argument `json:"inputs"`
	Name            string     `json:"name"`
	Outputs         []Argument `json:"outputs"`
	StateMutability string     `json:"stateMutability"`
	Type            string     `json:"type"`
}

// Argument represents a function argument
type Argument struct {
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}

func GenerateFHEOperationTemplate(returnType string) *template.Template {
	templateText := `
func (con FheOps) {{.Name}}(c ctx, evm mech, {{.Inputs}}) ({{.ReturnType}}, error) {
    err := fheos.SetEvmInterpreter(evm.Interpreter())
	if err != nil {
		return nil, err
	}

    return fheos.{{.Name}}({{.InnerInputs}})
} 
`

	if returnType == "void" {
		templateText = `
func (con FheOps) {{.Name}}(c ctx, evm mech, {{.Inputs}}) error {
	err := fheos.SetEvmInterpreter(evm.Interpreter())
	if err != nil {
		return err
	}

	fheos.{{.Name}}({{.InnerInputs}})
	return nil
}
`
	}

	tmpl, err := template.New("functionTemplate").Parse(templateText)
	if err != nil {
		fmt.Printf("Error parsing template: %s\n", err)
		return nil
	}

	return tmpl
}
