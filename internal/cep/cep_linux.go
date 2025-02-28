//go:build linux

package cep

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>
#include <string.h>

// Definindo funções conforme a biblioteca real (nm output)
typedef int (*CEP_InicializarFunc)(const char* eArqConfig, const char* eChaveCrypt);
typedef int (*CEP_FinalizarFunc)(void);
typedef int (*CEP_ConfigGravarValorFunc)(const char* eSessao, const char* eChave, const char* eValor);
typedef int (*CEP_BuscarPorCEPFunc)(const char* eCEP, char* buffer, int* bufferSize);
typedef int (*CEP_UltimoRetornoFunc)(char* buffer, int* bufferSize);

static void* loadLib(const char* path) {
    return dlopen(path, RTLD_LAZY);
}

static void* loadSymbol(void* handle, const char* symbol) {
    return dlsym(handle, symbol);
}

static int call_Inicializar(void* fn, const char* eArqConfig, const char* eChaveCrypt) {
    return ((CEP_InicializarFunc)fn)(eArqConfig, eChaveCrypt);
}
static int call_Finalizar(void* fn) {
    return ((CEP_FinalizarFunc)fn)();
}
static int call_ConfigGravarValor(void* fn, const char* eSessao, const char* eChave, const char* eValor) {
    return ((CEP_ConfigGravarValorFunc)fn)(eSessao, eChave, eValor);
}
static int call_BuscarPorCEP(void* fn, const char* eCEP, char* buffer, int* bufferSize) {
    return ((CEP_BuscarPorCEPFunc)fn)(eCEP, buffer, bufferSize);
}
static int call_UltimoRetorno(void* fn, char* buffer, int* bufferSize) {
    return ((CEP_UltimoRetornoFunc)fn)(buffer, bufferSize);
}

static char* makeCString(const char* s) {
    return strdup(s);
}

static void freeCString(char* s) {
    free(s);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const BUFFER_SIZE = 4096

type CEP struct {
	handle                unsafe.Pointer
	cep_Inicializar       unsafe.Pointer
	cep_Finalizar         unsafe.Pointer
	cep_ConfigGravarValor unsafe.Pointer
	cep_BuscarPorCEP      unsafe.Pointer
	cep_UltimoRetorno     unsafe.Pointer
}

func NewCEP(configFile, cryptKey string) (*CEP, error) {
	libPath := C.CString("./lib/libacbrcep64.so")
	defer C.free(unsafe.Pointer(libPath))

	handle := C.loadLib(libPath)
	if handle == nil {
		return nil, fmt.Errorf("failed to load libacbrcep64.so: %s", C.GoString(C.dlerror()))
	}

	cep := &CEP{handle: handle}
	if err := cep.loadFunctions(); err != nil {
		C.dlclose(handle)
		return nil, err
	}

	cConfigFile := C.makeCString(C.CString(configFile))
	cCryptKey := C.makeCString(C.CString(cryptKey))
	defer C.freeCString(cConfigFile)
	defer C.freeCString(cCryptKey)

	ret := C.call_Inicializar(cep.cep_Inicializar, cConfigFile, cCryptKey)
	if ret != 0 {
		errMsg := cep.getLastError(ret)
		C.dlclose(handle)
		return nil, fmt.Errorf("initialization error (code %d): %s", ret, errMsg)
	}

	return cep, nil
}

func (c *CEP) loadFunctions() error {
	c.cep_Inicializar = C.loadSymbol(c.handle, C.CString("CEP_Inicializar"))
	if c.cep_Inicializar == nil {
		return fmt.Errorf("failed to load CEP_Inicializar: %s", C.GoString(C.dlerror()))
	}
	c.cep_Finalizar = C.loadSymbol(c.handle, C.CString("CEP_Finalizar"))
	if c.cep_Finalizar == nil {
		return fmt.Errorf("failed to load CEP_Finalizar: %s", C.GoString(C.dlerror()))
	}
	c.cep_ConfigGravarValor = C.loadSymbol(c.handle, C.CString("CEP_ConfigGravarValor"))
	if c.cep_ConfigGravarValor == nil {
		return fmt.Errorf("failed to load CEP_ConfigGravarValor: %s", C.GoString(C.dlerror()))
	}
	c.cep_BuscarPorCEP = C.loadSymbol(c.handle, C.CString("CEP_BuscarPorCEP"))
	if c.cep_BuscarPorCEP == nil {
		return fmt.Errorf("failed to load CEP_BuscarPorCEP: %s", C.GoString(C.dlerror()))
	}
	c.cep_UltimoRetorno = C.loadSymbol(c.handle, C.CString("CEP_UltimoRetorno"))
	if c.cep_UltimoRetorno == nil {
		return fmt.Errorf("failed to load CEP_UltimoRetorno: %s", C.GoString(C.dlerror()))
	}
	return nil
}

func (c *CEP) Close() {
	if c.cep_Finalizar != nil {
		C.call_Finalizar(c.cep_Finalizar)
	}
	if c.handle != nil {
		C.dlclose(c.handle)
	}
}

func (c *CEP) SetConfig(section, key, value string) error {
	cSection := C.makeCString(C.CString(section))
	cKey := C.makeCString(C.CString(key))
	cValue := C.makeCString(C.CString(value))
	defer C.freeCString(cSection)
	defer C.freeCString(cKey)
	defer C.freeCString(cValue)

	ret := C.call_ConfigGravarValor(c.cep_ConfigGravarValor, cSection, cKey, cValue)
	if ret != 0 {
		return fmt.Errorf("error setting config %s.%s (code %d): %s", section, key, ret, c.getLastError(ret))
	}
	return nil
}

func (c *CEP) BuscarPorCep(cep string) (string, error) {
	cCEP := C.makeCString(C.CString(cep))
	defer C.freeCString(cCEP)

	buffer := make([]byte, BUFFER_SIZE)
	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	bufferSize := C.int(BUFFER_SIZE)

	ret := C.call_BuscarPorCEP(c.cep_BuscarPorCEP, cCEP, cBuffer, &bufferSize)
	if ret != 0 {
		return "", fmt.Errorf("error querying CEP %s (code %d): %s", cep, ret, c.getLastError(ret))
	}
	return C.GoStringN(cBuffer, bufferSize), nil
}

func (c *CEP) getLastError(retCode C.int) string {
	buffer := make([]byte, BUFFER_SIZE)
	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	bufferSize := C.int(BUFFER_SIZE)

	C.call_UltimoRetorno(c.cep_UltimoRetorno, cBuffer, &bufferSize)
	return C.GoStringN(cBuffer, bufferSize)
}
