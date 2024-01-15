package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const lookForNormal = `C:\Program Files\Nanitor\Nanitor Agent\nssm.exe`
const lookForOther = `C:\Program Files (x86)\Nanitor\Nanitor Agent\nssm.exe`

const regPath = `SOFTWARE\Classes\Installer\Products`

// regEntriesToClean lists possible registry keys under HKEY_CLASSES_ROOT\Installer\Products for previous Nanitor
// installations that might be interfering with the installation.
var regEntriesToClean = []string{
	`BBE9FCC2D1F201D4B8CA64DF69F93571`,
	`467C5BBFB3579BB419A8462DD9C60F47`,
	`7C26439BE1E9C1C4F8F2CA131CD58559`,
	`3A4836739F5FF934C9FEFB4B51CC4731`,
	`4CCDAA5DB9EF1E0499A94A0A2044F4B8`,
	`259B3D0E1047167499E9170C5AAB9FD1`,
	`7F22E1E791607F74A8C5660D21D76744`,
}

func getRegFullPath(regEntry string) string {
	return regPath + `\` + regEntry
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

func main() {
	if checkNanitorInstalled() {
		fmt.Printf("Nanitor is already installed, not needing to do any cleanup\n")
		os.Exit(0)
	}

	var regFixedEntries = []string{
		`Software\Nanitor`,
		`system\controlset001\services\nanitor agent`,
		`system\currentcontrolset\services\nanitor agent`,
		}
	var regInstallFolder = `SOFTWARE\Microsoft\Windows\CurrentVersion\Installer\Folders`
	var regInstallProduct = `SOFTWARE\classes\installer\products`
	for _, regEntry := range regFixedEntries {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, regEntry, registry.QUERY_VALUE)
		if err != nil {
			fmt.Printf("Reg path not found: %s - err(%v)\n", regEntry, err)
			//continue
		}	else {
			fmt.Printf("Reg path found: %s \n", regEntry)
		}
		k.Close()
	}
	keyInstallFolder, err := registry.OpenKey(registry.LOCAL_MACHINE, regInstallFolder, registry.QUERY_VALUE)
	if err != nil {
		fmt.Printf("Reg path not found: %s - err(%v)\n", regInstallFolder, err)
	}	else {
		fmt.Printf("Reg path found: %s \n", regInstallFolder)
		InstallFolderValues, _ := keyInstallFolder.ReadValueNames(-1)
		keyInstallFolder.Close()
		for _, FolderName := range InstallFolderValues{
			if strings.Contains(FolderName,"Nanitor"){
				fmt.Println(FolderName)
				}
		}
	}

	KeyInstallProd, err := registry.OpenKey(registry.LOCAL_MACHINE, regInstallProduct, registry.QUERY_VALUE | registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		fmt.Printf("Reg path not found: %s - err(%v)\n", regInstallProduct, err)
	}	else {
		fmt.Printf("Reg path found: %s \n", regInstallProduct)
		InstallProdKeys, _ := KeyInstallProd.ReadSubKeyNames(-1)
		KeyInstallProd.Close()
		for _, keyProdID := range InstallProdKeys{
			regSubKey := regInstallProduct + `\` + keyProdID
			KeySubProd, err := registry.OpenKey(registry.LOCAL_MACHINE, regSubKey, registry.QUERY_VALUE | registry.ENUMERATE_SUB_KEYS)
			if err != nil {
				fmt.Printf("Reg path not found: %s - err(%v)\n", regInstallProduct, err)
				continue
			}
			ProdNameStr,_,_ := KeySubProd.GetStringValue("ProductName")
			if strings.Contains(ProdNameStr,"Nanitor"){
				fmt.Printf("Deleting %s - %s",keyProdID,ProdNameStr)
				err = registry.DeleteKey(registry.LOCAL_MACHINE, regSubKey)
				if err != nil {
					fmt.Printf("Failed to delete key: %s - err(%v)\n", regSubKey, err)
					continue
				}
			}
			KeySubProd.Close()
		}
	}
}

func old(){

	for _, regEntry := range regEntriesToClean {
		regPathFull := getRegFullPath(regEntry)
		if len(regPathFull) < 10 {
			fmt.Printf("Invalid regPathFull\n")
			continue
		}

		k, err := registry.OpenKey(registry.LOCAL_MACHINE, regPathFull, registry.QUERY_VALUE)
		if err != nil {
			fmt.Printf("Reg path not found: %s - err(%v)\n", regPathFull, err)
			continue
		}
		k.Close()

		curDelete := regPathFull + `\SourceList\Media`
		err = registry.DeleteKey(registry.LOCAL_MACHINE, curDelete)
		if err != nil {
			fmt.Printf("Failed to delete key: %s - err(%v)\n", curDelete, err)
			continue
		}

		curDelete = regPathFull + `\SourceList\Net`
		err = registry.DeleteKey(registry.LOCAL_MACHINE, curDelete)
		if err != nil {
			fmt.Printf("Failed to delete key: %s - err(%v)\n", curDelete, err)
			continue
		}

		curDelete = regPathFull + `\SourceList`
		err = registry.DeleteKey(registry.LOCAL_MACHINE, curDelete)
		if err != nil {
			fmt.Printf("Failed to delete key: %s - err(%v)\n", curDelete, err)
			continue
		}

		curDelete = regPathFull
		err = registry.DeleteKey(registry.LOCAL_MACHINE, curDelete)
		if err != nil {
			fmt.Printf("Failed to delete key: %s - err(%v)\n", curDelete, err)
			continue
		}

		fmt.Printf("Key cleaned: %s\n", regPathFull)
	}

}
