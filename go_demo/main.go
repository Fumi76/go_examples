package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	out, err := exec.Command("aws", os.Args[1], os.Args[2], "--generate-cli-skeleton").CombinedOutput()
	if err != nil {
		if out != nil {
			fmt.Println(string(out))
		}
		log.Fatal(err)
	}

	replaced1 := strings.Replace(os.Args[1], "-", "_", -1)
	replaced2 := strings.Replace(os.Args[2], "-", "_", -1)

	dirName := replaced1 + "_" + replaced2

	if err := os.MkdirAll(dirName, 0777); err != nil {
		fmt.Println(err)
	}

	writeBytes(dirName+"/skeleton.json", out)

	bat := "\naws " + os.Args[1] + " " + os.Args[2] + " --cli-input-json file://%1\n"

	writeBytes(dirName+"/apply.bat", []byte(bat))

	bat2 := "\naws " + os.Args[1] + " " + os.Args[2] + " --generate-cli-skeleton > skeleton.json\n"

	writeBytes(dirName+"/generate_skeleton.bat", []byte(bat2))

}

func writeBytes(filename string, bytes []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
