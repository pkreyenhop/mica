# Miranda

This directory vendors the Miranda compiler/runtime source and support files.
It builds the `mira` executable plus the `miralib/` runtime library.

## Upstream

Active maintenance is at:

- `https://codeberg.org/DATurner/miranda`

This copy is vendored into the `mica` repository. Make C/source changes for this
project in this local `miranda/` directory and commit them to `mica`; do not
push project-specific edits to the upstream Codeberg Miranda repository.

Related:

- KRC: `https://codeberg.org/DATurner/KRC`
- Mailing list: `miranda@groups.io` (`http://groups.io/g/miranda`)

## Directory Layout

- `Makefile`: build/install/release targets for `mira`
- `*.c`, `*.h`, `rules.y`: Miranda implementation sources
- `miralib/`: standard library, examples, manuals, and helpers
- `issues/`: issue triage/build helper scripts
- `README`: original legacy text documentation
- `mira.1`, `mira.man.ms`: manpage sources

## Build

Requirements:

- `make`
- C compiler (`cc`, commonly `gcc` or `clang`)
- `byacc` only if you regenerate parser files from `rules.y`

Build in this directory:

```sh
make
```

Run without installing:

```sh
./mira
```

## Install

Default install paths in `Makefile`:

- `BIN=/usr/local/bin`
- `LIB=/usr/local/lib`
- `MAN=/usr/local/share/man/man1`

Install:

```sh
make install
```

If you want different paths or compiler flags, edit the top of `Makefile` before building/installing.

## Notes

- The garbage collector can be sensitive to aggressive optimization. If runtime stability issues appear, reduce/remove optimization flags in `CFLAGS`.
- GNU `bison` is not a drop-in replacement here; use `byacc` for parser regeneration.

## Developer Docs

- `CODEMAP.md`: architecture and file-by-file orientation for the core implementation.
- `README`: original legacy build and porting notes from the upstream distribution.
