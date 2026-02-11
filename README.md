# certinfo

This project is a Go library and command-line tool for retrieving information about the X.509 certificates located in the IDeTRUST credential repository.

## Project Structure

The project is split into a reusable library and a simple command-line wrapper.

```text
/
├── go.mod
├── /certinfo/           <-- The core, reusable library (package certinfo)
    └── /testdata/       <-- The certificate files used for testing
└── /cmd/
    └── /certinfo-cli/   <-- The command-line application (package main)
```

## How to Build

1.  Clone the repository:
    ```shell
    git clone https://github.com/Gandook/certinfo.git
    cd certinfo
    ```

2.  Build the `certinfo-cli` executable. This command compiles the application and places the binary in your current directory.
    ```shell
    go build ./cmd/certinfo-cli/
    ```
    This will create an executable file named `certinfo-cli` (or `certinfo-cli.exe` on Windows).

## How to Run

As of now, the application has only one command: `certinfo`.

### Retrieve Certificate Information

Use the `certinfo` command to see useful information about a certain DigSig X.509 certificate.

**Flags:**
* `-daid <string>`: The certificate's DAID.
* `-cid <string>`: The certificate's CID.
* `-f <string>`: The path to the certificate file.

**Example: Giving a DAID and a CID (Linux/macOS/PowerShell):**
```shell
./certinfo-cli certinfo -daid QCDEMO -cid 3
```

**Example: Giving a DAID and a CID (Windows Command Prompt):**
```shell
.\certinfo-cli.exe certinfo -daid QCDEMO -cid 3
```

**Example: Giving a file path (Linux/macOS/PowerShell):**
```shell
./certinfo-cli certinfo -f path/to/cert.pem
```

**Example: Giving a file path (Windows Command Prompt):**
```shell
.\certinfo-cli.exe certinfo -f path\to\cert.pem
```

### Running Tests

To run the included unit tests for the library:
```shell
go test ./certinfo/
```