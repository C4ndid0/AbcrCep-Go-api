package api

import (
	"acbr-Gocep-api/internal/cep"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	cepLib *cep.CEP
}

func NewHandler(cepLib *cep.CEP) *Handler {
	return &Handler{cepLib: cepLib}
}

type CEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento,omitempty"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge,omitempty"`
	Erro        string `json:"erro,omitempty"`
}

func (h *Handler) ConsultarCEP(w http.ResponseWriter, r *http.Request) {
	cep := strings.TrimPrefix(r.URL.Path, "/consultarCEP/")
	if cep == "" {
		log.Println("CEP não informado")
		h.sendError(w, r, "CEP não informado", http.StatusBadRequest)
		return
	}

	log.Printf("Consultando CEP: %s", cep)
	resultado, err := h.cepLib.BuscarPorCep(cep)
	if err != nil {
		log.Printf("Erro ao consultar CEP %s: %v", cep, err)
		h.sendError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Resultado do CEP %s: %s", cep, resultado)
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		h.sendJSON(w, resultado)
	} else {
		h.sendText(w, resultado)
	}
}

func (h *Handler) sendText(w http.ResponseWriter, resultado string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(resultado))
}

func (h *Handler) sendJSON(w http.ResponseWriter, resultado string) {
	fields := strings.Split(resultado, "|")
	resp := CEPResponse{}
	if len(fields) >= 7 {
		resp.CEP = fields[0]
		resp.Logradouro = fields[1]
		resp.Complemento = fields[2]
		resp.Bairro = fields[3]
		resp.Localidade = fields[4]
		resp.UF = fields[5]
		resp.IBGE = fields[6]
	} else {
		log.Printf("Formato de resposta inválido: %s", resultado)
		resp.Erro = "Formato de resposta inválido"
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Erro ao gerar JSON: %v", err)
		h.sendError(w, nil, "Erro ao gerar JSON", http.StatusInternalServerError)
	}
}

func (h *Handler) sendError(w http.ResponseWriter, r *http.Request, msg string, status int) {
	var accept string
	if r != nil {
		accept = r.Header.Get("Accept")
	}
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(CEPResponse{Erro: msg})
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(status)
		w.Write([]byte(msg))
	}
}
