# fsize

This is a project of mine, basically to end my problem with
the damn stat command that is unreadable.

Also because I had free time and a desire to do something different
that does not focus on Microsoft Windows.

The project is focused to be simple and easy to read,
so don't expect very advanced functions, but I will probably
add more according to my needs.

[![asciicast](https://asciinema.org/a/661870.svg)](https://asciinema.org/a/661870)

## Install

You can install it running

```bash
go install -v github.com/Tom5521/fsize@latest
```

However, I recommend this method more:

Basically because go takes a long time to automatically detect the newest tags.
And also this installs the completions

```bash
git clone https://github.com/Tom5521/fsize.git
just linux-install # or "sudo just linux-install" which will copy it to /usr/bin
```

## Documentation

```
Displays the file/folder properties.

Usage:
  fsize [flags]

Flags:
  -c, --config strings        Configure the variables used for preferences
                              Example: "fsize --config 'AlwaysShowProgress=true,AlwaysPrintOnWalk=false'".

                              To see the available variables and their values run "fsize --print-settings".
      --gen-bash-completion   Generate a completion file for bash
                              if any, the first argument will be taken as output file.
      --gen-fish-completion   Generate a completion file for fish
                              if any, the first argument will be taken as output file.
      --gen-zsh-completion    Generate a completion file for zsh
                              if any, the first argument will be taken as output file.
  -h, --help                  help for fsize
      --no-walk               Skips walking inside the directories.
      --no-warns              Hide possible warnings.
      --print-on-walk         Prints the name of the file being walked if a directory has been selected.
      --print-settings        Prints the current configuration values.
  -p, --progress              Displays a file count and progress bar when counting and summing file sizes. (default true)
  -v, --version               version for fsize
```

Or by copying one of the
[binaries](https://github.com/Tom5521/fsize/releases/latest) to your system to PATH
