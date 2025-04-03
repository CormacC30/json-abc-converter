# json-abc-converter
Command line application that converts JSON to ABC Notation

Intended for use on JSON files containing ABC headers and notation, converts to standard ABC format

## Usage

- Convert a file to multiple files: one file per tune:

```
./abc_converter -input <filename>.json -output ./abc_files
```

This will add all the abc tune files to a directory called `abc_files`

- Convert a file to one single combined file

```
./abc_converter -input <filename>.json -single -outfile <outputfile>.abc
```
