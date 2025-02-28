# acbr-Gocep-api
# ACBr CEP API

REST API for CEP (Brazilian postal code) lookup using ACBrLibCEP, compatible with both Windows and Linux.

## Purpose
This project provides a containerized REST API for querying Brazilian postal codes (CEPs), leveraging the ACBrLibCEP library. It is based on:
- PHP: https://svn.code.sf.net/p/acbr/code/trunk2/Projetos/ACBrLib/Demos/PHP/ConsultaCEP/
- .NET: https://github.com/OpenAC-Net/OpenAC.Net.CEP
- ACBrLib Documentation: https://acbr.sourceforge.io/ACBrLib/ComoUsar.html

The API supports both Windows (using `ACBrCEP.dll`) and Linux (using `libacbrcep64.so`), with platform-specific implementations abstracted behind a common interface.

## Structure
- **cmd/**: Entry point (`main.go`).
- **internal/cep/**:
  - `cep_linux.go`: Linux implementation using `cgo` with `libacbrcep64.so`.
  - `cep_windows.go`: Windows implementation using `golang.org/x/sys/windows` with `ACBrCEP.dll`.
- **internal/api/**: REST endpoints (`handler.go`).
- **internal/config/**: Configuration handling (`config.go`).
- **lib/**: Contains platform-specific libraries:
  - `ACBrCEP.dll` (Windows).
  - `libacbrcep64.so` (Linux).

## Prerequisites
- **Go**: Version 1.24.0 or higher installed.
- **Windows**:
  - `ACBrCEP.dll` in the `lib/` directory.
- **Linux (WSL Ubuntu)**:
  - `libacbrcep64.so` in the `lib/` directory.
  - Dependencies: `libgtk2.0-0`, `libxml2`, `xvfb`, `xauth`, `ttf-mscorefonts-installer`.
- **Postman**: For testing the API.

## Build and Run Locally

### Windows
1. **Install Go**:
   - Download and install from https://go.dev/dl/.

2. **Navigate to project directory**:
   ```cmd
   cd C:\path\to\acbr-cep-api