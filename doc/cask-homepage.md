# cask-homepage

Check homepages of the casks and try to fix them.

## Description

The `cask-homepage` script is designed to check homepages of the casks for
availability by checking their response statuses. It's useful in finding casks
that are no longer maintained or available. It also checks the URL for these
rules below (in brackets) and gives appropriate warnings if violated:

- HTTPS is available _(https)_
- Redirect found _(redirect)_
- Missing a trailing slash at the end _(slash)_

> The word in brackets represents the rule itself and can be found a warning in
> CSV or can be used with `-i, --ignore` option.

The script also can automatically fix those warnings when running it in
combination with the `-f, --fix` option.

### Lists of homepages

Since there are a lot of casks in [Homebrew-Cask](https://github.com/caskroom/homebrew-cask)
and [Homebrew-Versions](https://github.com/caskroom/homebrew-versions) some of
them can become unavailable due to various reasons which usually is reflected on
the homepage. To help maintainers find such casks the weekly result of this script
is available here:

- <http://caskroom.victorpopkov.com/homebrew-cask/homepages.csv>
- <http://caskroom.victorpopkov.com/homebrew-versions/homepages.csv>

Please note, that the lists are regenerated **every week at 16:00 (GMT+0) on
Saturday**. Please double check before submitting a PR, to make sure that the
cask hasn't been updated by someone else yet.

### Available options

#### `-i, --ignore <warnings>`

Ignore specified warnings separated by a comma: redirect https and slash.

It's useful when you would like to disable some warnings during the check. For
example to disable all warnings and display only casks that have homepage
errors:

```bash
cask-homepage -i redirect,https,slash
```

#### `-H, --header <header>`

Set browser header.

By default `User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2)
AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.152 Safari/537.36` is
used.

#### `-o, --output <filepath>`

Output the results in CSV format into a file.

For example, to output a CSV list in your home directory:

```bash
cask-homepage -o ~/outdated.csv
```

#### `-f, --fix (<total>)`

Try to fix warnings automatically (optional: total number of casks to fix).

#### `-a, --all`

Show and output all casks even if homepages are good.

By default, only the casks with homepage errors or warnings are shown on
screen and added to the CSV file. This parameter forces to show all the casks
even if the homepage is good.

#### `-d, --delete-branch`

Deletes local and remote branch named 'cask-homepage_fix'.

#### `-v, --version`

Show current script version.

#### `-h, --help`

Show the usage message with options descriptions.

### Available warnings

#### HTTPS is available _(https)_

Since the HTTPS is more preferred over the usual HTTP the script checks if it's
available and gives an appropriate warning message.

Example:

```
Cask name:       010-editor
Cask homepage:   http://www.sweetscape.com/
Status:          warning

                 1. HTTPS is available → https://www.sweetscape.com/
```

#### Redirect found _(redirect)_

Usually when the homepage changes a redirect is added to the old URL to help
users find the new homepage location. However, the usage of old URL is not
recommended since in the future in can become unavailable.

Example:

```
Cask name:       1password
Cask homepage:   https://agilebits.com/onepassword [301]
Status:          warning

                 1. Redirect found → https://1password.com/
```

#### Missing a trailing slash at the end _(slash)_

It's highly recommended to use a trailing slash in a bare domain URL like a
homepage. However, most mainstream browsers "append a trailing slash"
automatically to the request.

From [RFC 2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec5.html):

> Note that the absolute path cannot be empty; if none is present in the
> original URI, it MUST be given as "/" (the server root).

Example:

```
Cask name:       5kplayer
Cask homepage:   https://www.5kplayer.com
Status:          warning

                 1. Missing a trailing slash at the end → https://www.5kplayer.com/
```

## Examples

By default you just have to `cd` into the Casks directory and run the script:

```bash
$ cd ~/path/to/homebrew-cask/Casks
$ cask-homepage
Checking homepages of 3420 casks...
---------------------------------------------------------------------------------------------------------------------------------------------
Cask name:       010-editor
Cask homepage:   http://www.sweetscape.com/
Status:          warning

                 1. HTTPS is available → https://www.sweetscape.com/
---------------------------------------------------------------------------------------------------------------------------------------------
Cask name:       1password
Cask homepage:   https://agilebits.com/onepassword [301]
Status:          warning

                 1. Redirect found → https://1password.com/
---------------------------------------------------------------------------------------------------------------------------------------------
Cask name:       4peaks
Cask homepage:   http://nucleobytes.com/index.php/4peaks [301]
Status:          warning

                 1. Redirect found → http://nucleobytes.com/4peaks/
...
```

### CSV lists

In some cases it's useful to output all the found outdated casks into separate
CSV lists:

- <http://caskroom.victorpopkov.com/homebrew-cask/homepages.csv>
- <http://caskroom.victorpopkov.com/homebrew-versions/homepages.csv>

In order to do that you have to `cd` into the Casks directory and run the script
with `-o` option and `~/path/to/output.csv` argument. For example:

```bash
cd ~/path/to/homebrew-versions/Casks
cask-homepage -u -o '~/path/to/output.csv'
```