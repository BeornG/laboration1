# Project overview

This Project is a command-line tool for parsing and rendering `.diag`files into SVG diagrams.

It includes a custom lexer and parser and a renderer, a rendering engine and a CLI interface.

The project is modular and written in Go.

## Purpose

The goal of the project is to transform text.based diagram descriptions into a visual SVG output, enabling programmitc creation of diagrams.

## Features

- Custom lexer and parser for `.diag` files
- SVG rendering engine with layout management
- Command-line interface for transforming `.diag` files into SVG code
- Example files for demonstrating and testing

## Project structure

    diagra/
    ├── assets                  # static or template files
    ├── cmd                     # CLI entrypoint
    ├── example                 # contains example files
    ├── interpreter             # Interprets
    ├── renderer                # renders
    ├── runtime                 # runtime
    ├── test                    # tests



