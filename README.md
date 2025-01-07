# toml-alias
Tool for easy TOML configurable alias creation.

*Recommended platform: __Windows__*  
*Recommended shell: __cmd prompt__ or __cygwin bash__*

The idea is to use common binary under yet another name (alias name) and execute it strictly the way it is configured in dot-config file (TOML). On Linux such functionality like `aliases`
supports easy way to execute much complex custom commands,
but on Windows `doskey` mechanism is pretty harmful and disenchartmentful. Gratefully `Powershell` or `Python` supports
theirs own mechanisms. Here, the tool offers yet another way for
any executable possibly called by command prompt. 

## Configuration

Configuration file consists custom aliases activated by binary call (eg. by adding the binary into `$PATH` directory). By the binary (for Windows), this means, `toml-alias.exe` and `aliaslib.dll`.

!!! Warning
  For binary files search in `Releases` tab.

`toml-alias.exe` file suppose to be renamed to Your name of alias
(whatever alias You needed). Same alias-name should be configured in TOML file.

Depending on current shell prompt localize configuration file in *`$HOME/.config/toml-alias/config.toml`*.

- In case of `command prompt` configuration file suppose to be in
  path like: `C:/Users/<YOUR_USERNAME>/.config/toml-alias/config.toml`
- In case of Cygwin pseudo shell: `C:/<cygwin_directory>/home/<YOUR_USERNAME>/.config/toml-alias/config.toml`

For a very abstract example of TOML fle plase see [here](https://github.com/kkuba91/toml-alias/examples/config.toml).

Below there is a list of parameters explained by the example:

```
[toml-alias]            # ALIAS NAME
hide-base-help=false    # ALLOW TO SEE `toml-alias` help on `-h`
hide-base-version=true  # ALLOW TO SEE `toml-alias` version on `-V`
custom-help="Check python version."  # CUSTOM HELP TEXT ON `-h` 
custom-version="0.1.0"  # CUSTOM VERSION TEXT ON `-V` 

[[toml-alias.stage]]    # STAGE (STEP) - CREATE AS MANY AS NEEDED
print-stdout=false      # SHOW COMMANDS OUTPUT (ON `true`)
print-match=true        # SHOW MATCHING RESULT (ON `true`)
allow-fail=false        # HIDE "FAIL" MATCHING IF OCCUR (ON `true`)
match-stdout="(?P<match>\\d\\.\\d+)\\s+\\*"  # MATCH REGEX - NOTE <match> IS HARDCODED SEARCHED GROUP
match-msg="Found python:"  # EXTRA STRING FOR MATCHING
print-on-end="This is the end"  # SOME PRINT STRING AT THE END
print-on-success="[PASS]"  # PRINT STRING WHEN MATCHING IS SUCCESSFULL
print-on-failure="[FAIL]"  # PRINT STRING WHEN MATCHING DOS NOT FOUND
cmd=["cmd", "/k", "py -0|findstr *&echo ✅"]  # COMMAND (USE LIST)
pre-cmd=[]   # EXTRA PRE-COMMAND (USE LIST) - DOES NOT MATCH
post-cmd=[]  # EXTRA POST-COMMAND (USE LIST) - DOES NOT MATCH

[[toml-alias.stage]]    # YET ANOTHER STAGE (ADD AS MANY NEEDED)
print-stdout=false
print-on-end="[style.italic]This is the end[style.reset] - [style.bold]stage 2[style.reset] [style.green]✅"  # SOME EXTRA COLORS
cmd=["echo ✅"]
pre-cmd=[]
post-cmd=[]
```

## Usage

Above example allows to execute configured behavior in the way only for the binary from `Releases` tab renamed to `toml-alias.exe` (so as it is by default). To try new, please add another alias configuration in the same file and rename binary.

__[OPTIONAL]__ To work as an alias (like on Linux), simply add Your binary (with `aliaslib.dll`) to the directory included in user (or system) $PATH value.

!!! Warning
    Note, that executing binary files from release requires "allow" for Windows Defender sometimes (as an exception rule). Please, use __ONLY__ these binaries from github source if th tool is intended to use.

To use custom colors in alias custom strings, please see available styles defined in `stylesMap` in [terminal_print.go](https://github.com/kkuba91/toml-alias/shared_library/terminal_print.go)

## Project

The project of tool includes two binary `go` projects:

 - for executiona binary - it only loads dll content - minimal
 - `aliaslib` shared library with all implemented logic, it is
   located in `/shared_library/` project directory 
