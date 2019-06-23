# github-rev-parse (perl edition)

## Synopsis

```
[usage]
  $ perl github-rev-parse <org> <repo> <key (sha1, branch, tag)>
  options:
    --token=github-token : pass the token of GitHub
```

## Note

```
.
├── github-rev-parse # <= fatpacked code, stand-alone executable
└── src
    └── github-rev-parse.pl # <= original perl code
```

## For developers

### How to fatpack the original perl code

```
$ make installdeps
$ make fatpack
```

# License

Copyright (C) moznion.

This library is free software; you can redistribute it and/or modify
it under the same terms as Perl itself.

# Author

moznion (<moznion@gmail.com>)

