# Go Projects

This repository showcases a collection of projects/scripts that I've created during my journey of learning Go programming language.

## Version Control System

It's a CLI application that can be used to track file changes. You can stage/track files, commit changes, log commits and checkout to a previous commit. It's a simple implementation of Git. It uses `os`, `bufio` and `io` packages to create, read and write files and `crypto/sha256` package for hashing. By building this project, I learned about the internal working of a version control system like Git.
<br />

**How to use:** check out the [usage instructions](./version-control-system/README.md).\
**Project details:** click [here](https://hyperskill.org/projects/420).

## Simple Bank Manager

It's a simple credit card system that can generate valid credit card numbers by using Luhn algorithm. You can use the generated card and PIN to login and do transactions such as add money, transfer money to another card in the database, close account and log out. It uses GORM package in Go to interact with a SQLite database and store the card details and PIN. By building this project, I learned about models, CRUD operations, transactions and migrations in GORM.
<br />

**How to use:** check out the [usage instructions](./simple-bank-manager/README.md).\
**Project details:** click [here](https://hyperskill.org/projects/413).

## Cinema Room Manager

It's a script that helps manage a cinema theatre: sell tickets, check available seats, display seat map and sales statistics. It covers the fundamentals of Go.
<br />

**Project details:** click [here](https://hyperskill.org/projects/399).
