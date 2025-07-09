# GoTeX

GoTeX is a LaTeX interpreter/compiler written in pure Go. The project is currently in early development stages.

## Current Limitations

- Package system is not implemented
- Font handling requires improvement
- Limited mathematical symbol coverage
- Trigonometric and logarithmic functions are not implemented

## Usage

Compile a TeX file to PDF:

```bash
gotex document.tex -o output.pdf
```

Or without the .pdf extension (automatically added):

```bash
gotex document.tex -o output
```

You can also omit the .tex extension from the input file:

```bash
gotex assignment -o assignment  # Compiles assignment.tex if it exists
```

Scan directory for TeX files:

```bash
gotex -s
```

Or

```bash
gotex --scan
```
