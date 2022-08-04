package main

import (
	"fmt"
	"os"

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

	// [HKEY_CLASSES_ROOT\Installer\Products\259B3D0E1047167499E9170C5AAB9FD1]
	// "Clients"=hex(7):3a,00,00,00,00,00
	// "ProductName"="Nanitor Agent (64-bit)"
	// "PackageCode"="BA4F162EDB6C82040A6C89A2F9CC4BE5"
	// "Language"=dword:00000409
	// "Version"=dword:01070002
	// "Assignment"=dword:00000001
	// "AdvertiseFlags"=dword:00000180
	// "ProductIcon"="C:\\WINDOWS\\Installer\\{E0D3B952-7401-4761-999E-71C0A5BAF91D}\\icon.ico"
	// "InstanceType"=dword:00000000
	// "AuthorizedLUAApp"=dword:00000000
	// "DeploymentFlags"=dword:00000001
	// [HKEY_CLASSES_ROOT\Installer\Products\259B3D0E1047167499E9170C5AAB9FD1\SourceList]
	// "PackageName"="nanitor-agent-1.7.2.6301_windows_amd64.msi"
	// "LastUsedSource"="n;1;C:\\WINDOWS\\ccmcache\\n\\"
	`259B3D0E1047167499E9170C5AAB9FD1`,

	// [HKEY_CLASSES_ROOT\Installer\Products\7F22E1E791607F74A8C5660D21D76744]
	// "Clients"=hex(7):3a,00,00,00,00,00
	// "DeploymentFlags"=dword:00000002
	// "AuthorizedLUAApp"=dword:00000000
	// "InstanceType"=dword:00000000
	// "ProductIcon"="C:\\WINDOWS\\Installer\\{7E1E22F7-0619-47F7-8A5C-66D0127D7644}\\icon.ico"
	// "AdvertiseFlags"=dword:00000184
	// "Assignment"=dword:00000001
	// "Version"=dword:02050000
	// "Language"=dword:00000409
	// "PackageCode"="A5FF85443CEECB042A596AE8F2E373F4"
	// "ProductName"="Nanitor Agent (64-bit)"
	// [HKEY_CLASSES_ROOT\Installer\Products\7F22E1E791607F74A8C5660D21D76744\SourceList]
	// "LastUsedSource"=hex(2):6e,00,3b,00,31,00,3b,00,43,00,3a,00,5c,00,55,00,73,00,\
	// 65,00,72,00,73,00,5c,00,73,00,63,00,69,00,62,00,6f,00,6e,00,61,00,64,00,5c,\
	// 00,44,00,6f,00,77,00,6e,00,6c,00,6f,00,61,00,64,00,73,00,5c,00,00,00
	// "PackageName"="nanitor-agent-2.5.0.10219_windows_amd64.msi"
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
