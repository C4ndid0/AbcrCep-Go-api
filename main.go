package main

import (
	"acbr-Gocep-api/internal/api"
	"acbr-Gocep-api/internal/cep"
	"acbr-Gocep-api/internal/config"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg := config.LoadConfig()

	// Criar a pasta log se ela não existir
	logDir := "log"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.Mkdir(logDir, 0755); err != nil {
			log.Fatalf("Erro ao criar a pasta log: %v", err)
		}
	}

	cepLib, err := cep.NewCEP("ACBrLib.ini", "")
	if err != nil {
		log.Fatalf("Erro ao inicializar ACBrLibCEP: %v", err)
	}
	defer cepLib.Close()

	handler := api.NewHandler(cepLib)
	http.HandleFunc("/consultarCEP/", handler.ConsultarCEP)

	log.Printf("Servidor rodando em %s", cfg.ServerAddr)
	log.Fatal(http.ListenAndServe(cfg.ServerAddr, nil))
}
