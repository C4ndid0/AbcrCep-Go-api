//go:build windows

package cep

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const BUFFER_SIZE = 4096

type CEP struct {
	dll                   *windows.DLL
	cep_Inicializar       *windows.Proc
	cep_Finalizar         *windows.Proc
	cep_ConfigGravarValor *windows.Proc
	cep_BuscarPorCEP      *windows.Proc
	cep_UltimoRetorno     *windows.Proc
}

func NewCEP(configFile, cryptKey string) (*CEP, error) {
	dll, err := windows.LoadDLL("./lib/ACBrCEP.dll")
	if err != nil {
		return nil, fmt.Errorf("failed to load ACBrCEP.dll: %v", err)
	}

	cep := &CEP{dll: dll}
	if err := cep.loadFunctions(); err != nil {
		dll.Release()
		return nil, err
	}

	configFilePtr, err := syscall.BytePtrFromString(configFile)
	if err != nil {
		dll.Release()
		return nil, fmt.Errorf("failed to convert configFile to C string: %v", err)
	}
	cryptKeyPtr, err := syscall.BytePtrFromString(cryptKey)
	if err != nil {
		dll.Release()
		return nil, fmt.Errorf("failed to convert cryptKey to C string: %v", err)
	}

	ret, _, err := cep.cep_Inicializar.Call(uintptr(unsafe.Pointer(configFilePtr)), uintptr(unsafe.Pointer(cryptKeyPtr)))
	if ret != 0 {
		errMsg := cep.getLastError(int(ret))
		dll.Release()
		return nil, fmt.Errorf("initialization error (code %d): %s", ret, errMsg)
	}

	return cep, nil
}

func (c *CEP) loadFunctions() error {
	c.cep_Inicializar = c.dll.MustFindProc("CEP_Inicializar")
	c.cep_Finalizar = c.dll.MustFindProc("CEP_Finalizar")
	c.cep_ConfigGravarValor = c.dll.MustFindProc("CEP_ConfigGravarValor")
	c.cep_BuscarPorCEP = c.dll.MustFindProc("CEP_BuscarPorCEP")
	c.cep_UltimoRetorno = c.dll.MustFindProc("CEP_UltimoRetorno")
	return nil
}

func (c *CEP) Close() {
	c.cep_Finalizar.Call()
	c.dll.Release()
}

func (c *CEP) SetConfig(section, key, value string) error {
	sectionPtr, err := syscall.BytePtrFromString(section)
	if err != nil {
		return fmt.Errorf("failed to convert section to C string: %v", err)
	}
	keyPtr, err := syscall.BytePtrFromString(key)
	if err != nil {
		return fmt.Errorf("failed to convert key to C string: %v", err)
	}
	valuePtr, err := syscall.BytePtrFromString(value)
	if err != nil {
		return fmt.Errorf("failed to convert value to C string: %v", err)
	}

	ret, _, err := c.cep_ConfigGravarValor.Call(
		uintptr(unsafe.Pointer(sectionPtr)),
		uintptr(unsafe.Pointer(keyPtr)),
		uintptr(unsafe.Pointer(valuePtr)),
	)
	if ret != 0 {
		return fmt.Errorf("error setting config %s.%s (code %d): %s", section, key, ret, c.getLastError(int(ret)))
	}
	return nil
}

func (c *CEP) BuscarPorCep(cep string) (string, error) {
	cepPtr, err := syscall.BytePtrFromString(cep)
	if err != nil {
		return "", fmt.Errorf("failed to convert cep to C string: %v", err)
	}

	buffer := make([]byte, BUFFER_SIZE)
	bufferSize := int32(BUFFER_SIZE)

	ret, _, err := c.cep_BuscarPorCEP.Call(
		uintptr(unsafe.Pointer(cepPtr)),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&bufferSize)),
	)
	if ret != 0 {
		return "", fmt.Errorf("error querying CEP %s (code %d): %s", cep, ret, c.getLastError(int(ret)))
	}
	return string(buffer[:bufferSize]), nil
}

func (c *CEP) getLastError(retCode int) string {
	buffer := make([]byte, BUFFER_SIZE)
	bufferSize := int32(BUFFER_SIZE)

	c.cep_UltimoRetorno.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&bufferSize)),
	)
	return string(buffer[:bufferSize])
}
