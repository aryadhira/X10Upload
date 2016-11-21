package x10upload

import (
	//"os"

	"os"
	"os/exec"
	"strings"

	tk "github.com/eaciit/toolkit"
)

func ConvertPdfToXml(PathFrom string, PathTo string, FName string) error {

	Name := strings.TrimRight(FName, ".pdf")

	tk.Printf("Converting %#v....\n", FName)

	FileName := PathFrom + "/" + FName
	ResultName := PathTo + "/" + Name + ".xml"

	formattedName := strings.Replace(FileName, " ", "\\ ", -1)
	formattedResultName := strings.Replace(ResultName, " ", "\\ ", -1)

	if _, err := os.Stat(FileName); err == nil {
		cmdStr := []string{"pdftohtml", "-xml", formattedName, formattedResultName}
		finalcmd := strings.Join(cmdStr, " ")
		if err := exec.Command("/bin/sh", "-c", finalcmd).Run(); err != nil {
			return err
		} else {
			tk.Println("Converting Success")
		}
	} else {
		tk.Println("File Doesn't Exist")
	}
	return nil
}
