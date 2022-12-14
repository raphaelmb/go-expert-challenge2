package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
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
	ctx := context.WithValue(context.Background(), "trim", "true")
	cep = ParseInput(ctx, cep)
	return "https://viacep.com.br/ws/" + cep + "/json/"
}

func (v ViaCep) ViaCepOutput(res ViaCep) {
	fmt.Printf("ViaCEP respondeu primeiro:\n CEP: %s,\n Logradouro: %s,\n Complemento: %s,\n Bairro: %s,\n Localidade: %s,\n UF: %s\n", res.Cep, res.Logradouro, res.Complemento, res.Bairro, res.Localidade, res.Uf)
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
	ctx := context.WithValue(context.Background(), "trim", "false")
	cep = ParseInput(ctx, cep)
	return "https://cdn.apicep.com/file/apicep/" + cep + ".json"
}

func (a ApiCep) ApiCepOutput(res ApiCep) {
	fmt.Printf("ApiCEP respondeu primeiro:\n CEP: %s,\n Estado: %s,\n Cidade: %s,\n Distrito: %s,\n Endereço: %s\n", res.Code, res.State, res.City, res.District, res.Address)
}

type ApiInterface interface {
	ViaCep | ApiCep
}

func main() {
	fmt.Println("Digite o CEP desejado. Exemplo: 12345-678 ou 12345678")
	var cep string
	_, err := fmt.Scan(&cep)
	if err != nil {
		log.Println(err)
		return
	}

	err = IsValidCep(cep)
	if err != nil {
		fmt.Println(err)
		return
	}

	apiCepChan := make(chan ApiCep)
	viaCepChan := make(chan ViaCep)

	go GetCep(ApiCep{}.GetUrl(cep), apiCepChan)
	go GetCep(ViaCep{}.GetUrl(cep), viaCepChan)

	select {
	case res := <-apiCepChan:
		ApiCep{}.ApiCepOutput(res)
	case res := <-viaCepChan:
		ViaCep{}.ViaCepOutput(res)
	case <-time.After(time.Second):
		fmt.Println("Failed: timeout reached.")
	}
}

func GetCep[T ApiInterface](url string, ch chan T) {
	req, err := http.Get(url)

	if err != nil {
		return
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	var c T
	err = json.Unmarshal(body, &c)
	if err != nil {
		return
	}

	ch <- c
}

// utils
func IsValidCep(input string) error {
	match, err := regexp.Match(`(^[0-9]{8}$)|(^[0-9]{5}-[0-9]{3}$)`, []byte(input))
	if err != nil {
		return err
	}
	if !match {
		return errors.New("CEP digitado no formato errado")
	}
	return nil
}

func ParseInput(ctx context.Context, s string) string {
	if strings.Contains(s, "-") {
		if ctx.Value("trim") == "true" {
			return strings.Replace(s, "-", "", 1)
		}
		return s
	} else {
		if ctx.Value("trim") == "false" {
			return s[:5] + "-" + s[5:]
		}
		return s
	}
}
