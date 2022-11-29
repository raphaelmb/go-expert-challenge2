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
	return "https://cdn.apicep.com/file/apicep/" + cep + ".json"
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
	return "https://viacep.com.br/ws/" + cep + "/json/"
}

type ApiInterface interface {
	ViaCep | ApiCep
}

// TODO: get user input in cli?
func main() {
	c1 := make(chan ApiCep)
	c2 := make(chan ViaCep)

	go GetCep(ViaCep{}.GetUrl(""), c1)
	go GetCep(ApiCep{}.GetUrl(""), c2)

	// TODO: better formatted output
	select {
	case res := <-c1:
		fmt.Println("ApiCEP responded first:", res)
	case res := <-c2:
		fmt.Println("ViaCEP responded first:", res)
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
