package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func (v ViaCep) GetUrl(cep string) string {
	// ctx := context.WithValue(context.Background(), "trim", "true")
	// cep = sanitizeInput(ctx, cep)
	return "https://viacep.com.br/ws/" + cep + "/json/"
}

type ApiCep struct {
	Status   int    `json:"status"`
	Code     string `json:"code"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

func (a ApiCep) GetUrl(cep string) string {
	// ctx := context.WithValue(context.Background(), "trim", "false")
	// cep = sanitizeInput(ctx, cep)
	return "https://cdn.apicep.com/file/apicep/" + cep + ".json"
}

type ApiInterface interface {
	ViaCep | ApiCep
}

// func sanitizeInput(ctx context.Context, s string) string {
// 	if strings.Contains(s, "-") {
// 		if ctx.Value("trim") == "true" {
// 			return strings.Replace(s, "-", "", 1)
// 		}
// 		return s
// 	} else {
// 		if ctx.Value("trim") == "false" {
// 			return s[:5] + "-" + s[5:]
// 		}
// 		return s
// 	}
// }

// TODO: get user input in cli?
func main() {
	c1 := make(chan ApiCep)
	c2 := make(chan ViaCep)

	go GetCep(ApiCep{}.GetUrl(""), c1)
	go GetCep(ViaCep{}.GetUrl(""), c2)

	select {
	case res := <-c1:
		fmt.Printf("ApiCEP responded first:\n CEP: %s,\n Estado: %s,\n Cidade: %s,\n Distrito: %s,\n Endereço: %s\n", res.Code, res.State, res.City, res.District, res.Address)
	case res := <-c2:
		fmt.Printf("ViaCEP responded first:\n CEP: %s,\n Logradouro: %s,\n Complemento: %s,\n Bairro: %s,\n Localidade: %s,\n UF: %s\n", res.Cep, res.Logradouro, res.Complemento, res.Bairro, res.Localidade, res.Uf)
	case <-time.After(time.Second):
		fmt.Println("Failed: timeout reached.")
	}
}

// TODO: error handling
func GetCep[T ApiInterface](url string, ch chan T) {
	req, err := http.Get(url)

	if err != nil {
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
	}

	var c T
	err = json.Unmarshal(body, &c)
	if err != nil {
	}

	ch <- c
}
