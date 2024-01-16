package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const lookForNormal = `C:\Program Files\Nanitor\Nanitor Agent\nssm.exe`
const lookForOther = `C:\Program Files (x86)\Nanitor\Nanitor Agent\nssm.exe`
const regInstallFolder = `SOFTWARE\Microsoft\Windows\CurrentVersion\Installer\Folders`
const regInstallProduct = `SOFTWARE\classes\installer\products`

var regFixedEntries = []string{
	`Software\Nanitor`,
	`system\controlset001\services\nanitor agent`,
	`system\currentcontrolset\services\nanitor agent`,
}

var dataFolders = []string{
	`C:\Program Files\Nanitor`,
	`C:\ProgramData\Nanitor`,
}

func checkNanitorInstalled() bool {
	checkForAll := []string{lookForNormal, lookForOther}

	for _, checkFor := range checkForAll {
		if _, err := os.Stat(checkFor); !os.IsNotExist(err) {
			// Path found.
			return true
		}

	}

	return false
}

func delOsFolder() {
	for _, folderName := range dataFolders {
		if _, err := os.Stat(folderName); !os.IsNotExist(err) {
			osErr := os.RemoveAll(folderName)
			if osErr == nil {
				fmt.Println("Successfully removed", folderName)
			} else {
				fmt.Println("Failed to remove", folderName, osErr)
			}
		}
		fmt.Println(folderName, "does not exists")
	}
}

func delRegKey(regEntry string) {
	regKey := `HKLM\` + regEntry
	fmt.Printf("About to delete key: %s \n", regKey)
	cmdObj := exec.Command("reg", "delete", regKey, "/f")
	out, err := cmdObj.CombinedOutput()
	if err != nil {
		fmt.Println("problem with command:", err)
	}
	fmt.Println(string(out))
}

func main() {
	if checkNanitorInstalled() {
		fmt.Printf("Nanitor is already installed, not needing to do any cleanup\n")
		os.Exit(0)
	}
	delOsFolder()
	for _, regEntry := range regFixedEntries {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, regEntry, registry.QUERY_VALUE)
		if err != nil {
			fmt.Printf("Reg path not found: HKLM\\%s - err(%v)\n", regEntry, err)
		} else {
			fmt.Printf("Reg path found: HKLM\\%s \n", regEntry)
			delRegKey(regEntry)
		}
		k.Close()
	}

	keyInstallFolder, err := registry.OpenKey(registry.LOCAL_MACHINE, regInstallFolder, registry.ALL_ACCESS)
	if err != nil {
		fmt.Printf("Reg path not found: HKLM\\%s - err(%v)\n", regInstallFolder, err)
	} else {
		//fmt.Printf("Reg path found: HKLM\\%s \n", regInstallFolder)
		InstallFolderValues, err := keyInstallFolder.ReadValueNames(-1)
		if err != nil {
			fmt.Printf("Value Names not found: HKLM\\%s - err(%v)\n", regInstallFolder, err)
		}
		var FolderCount int = 0
		for _, FolderName := range InstallFolderValues {
			if strings.Contains(FolderName, "Nanitor") {
				FolderCount++
				//fmt.Println(FolderName)
				err = keyInstallFolder.DeleteValue(FolderName)
				if err != nil {
					fmt.Printf("Failed to delete folder value %s - err(%v)\n", FolderName, err)
				} else {
					fmt.Printf("Successfully deleted folder value %s\n", FolderName)
				}
			}
		}
		fmt.Printf("Found %v Nanitor Folders.\n", FolderCount)
		keyInstallFolder.Close()
	}

	KeyInstallProd, err := registry.OpenKey(registry.LOCAL_MACHINE, regInstallProduct, registry.ALL_ACCESS)
	if err != nil {
		fmt.Printf("Reg path not found: HKLM\\%s - err(%v)\n", regInstallProduct, err)
	} else {
		//fmt.Printf("Reg path found: HKLM\\%s \n", regInstallProduct)
		InstallProdKeys, err := KeyInstallProd.ReadSubKeyNames(-1)
		if err != nil {
			fmt.Printf("Subkeys not found: HKLM\\%s - err(%v)\n", regInstallProduct, err)
		}
		KeyInstallProd.Close()
		for _, keyProdID := range InstallProdKeys {
			regSubKey := regInstallProduct + `\` + keyProdID
			KeySubProd, err := registry.OpenKey(registry.LOCAL_MACHINE, regSubKey, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
			if err != nil {
				fmt.Printf("Reg path not found: HKLM\\%s - err(%v)\n", regSubKey, err)
				continue
			}
			ProdNameStr, _, _ := KeySubProd.GetStringValue("ProductName")

			if strings.Contains(ProdNameStr, "Nanitor") {
				delRegKey(regSubKey)
			}
			KeySubProd.Close()
		}
	}
}
