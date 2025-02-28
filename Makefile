# Makefile para acbr-cep-api
# Suporta desenvolvimento no Windows e produção no Linux Alpine com Docker

# Variáveis
APP_NAME = acbr-cep-api
PORT = 8080

# Detecta o sistema operacional
OS := $(shell go env GOOS)
ifeq ($(OS),windows)
    EXECUTABLE = $(APP_NAME).exe
else
    EXECUTABLE = $(APP_NAME)
endif



# Regra padrão: build
all: build

# Compila o binário
build:
	go build -o $(EXECUTABLE)

# Executa o binário localmente
run: build
	./$(EXECUTABLE)

# Testa a API após iniciar
test:
	@echo "Testando a API em http://localhost:$(PORT)/consultarCEP/01001000"
	@curl -s http://localhost:$(PORT)/consultarCEP/01001000 || echo "Erro ao testar a API. Verifique se ela está rodando."



# Limpa arquivos gerados
clean:
	rm -f $(EXECUTABLE)
	go clean -cache -modcache -i -r


# Instala dependências
deps:
	go mod tidy
	go get golang.org/x/sys/windows
	go get golang.org/x/sys/unix
