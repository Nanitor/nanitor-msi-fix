# nanitor-msi-fix
This tool is used to clean up registry keys associated with problems in Windows Installer that sometimes cause Nanitor Agent installations to fail on Windows.

NOTE: This should be used as a last resort and only if installs are failing.

## Usage
Simply run the binary and it will go over and remove registry keys in the Windows Installer that might be causing problems.
```
> nanitor-msi-fix.exe
```

The tool will print out a report of what keys it went through and what was removed.

Once the tool has been run, the Nanitor agent should install without problems.


