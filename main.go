package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Por favor, insira seu CEP:")
	var cep string
	_, err := fmt.Scanln(&cep)
	if err != nil {
		fmt.Println(err)
		return
	}

	brasilApiChannel := make(chan []byte)
	viaCepChannel := make(chan []byte)

	go makeRequestForChannel(brasilApiChannel, fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep))
	go makeRequestForChannel(viaCepChannel, fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep))

	select {
	case data := <-brasilApiChannel:
		var parsedData BrasilApiResponse
		err = json.Unmarshal(data, &parsedData)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fmt.Sprintf("Seu CEP foi encontrado em %s. %s - %s - %s - %s - %s", "BrasilAPI", parsedData.Cep, parsedData.City,
			parsedData.State, parsedData.Street, parsedData.Neighborhood))
	case data := <-viaCepChannel:
		var parsedData ViaCepResponse
		err = json.Unmarshal(data, &parsedData)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fmt.Sprintf("Seu CEP foi encontrado em %s. %s - %s - %s - %s - %s", "ViaCep", parsedData.Cep, parsedData.Localidade,
			parsedData.Uf, parsedData.Logradouro, parsedData.Bairro))
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout")
	}

}

func makeRequestForChannel(ch chan<- []byte, requestUrl string) {
	resp, err := http.Get(requestUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch <- data
}

type BrasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type ViaCepResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}
