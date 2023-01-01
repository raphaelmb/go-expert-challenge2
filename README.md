## Instruções

Execute o arquivo principal: `go run main.go`

Um prompt aparecerá pedindo o CEP: `"Digite o CEP desejado. Exemplo: 12345-678 ou 12345678"`

Digite o CEP desejado e pressione `enter`.

A resposta recebida será dos dados do CEP bem como qual API respondeu primeiro.

---

Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.

As duas requisições serão feitas simultaneamente para as seguintes APIs:

https://cdn.apicep.com/file/apicep/" + cep + ".json

http://viacep.com.br/ws/" + cep + "/json/

Os requisitos para este desafio são:

- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.

- O resultado da request deverá ser exibido no command line, bem como qual API a enviou.

- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.
