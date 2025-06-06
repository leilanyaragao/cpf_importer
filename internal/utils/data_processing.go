package utils

import (
	"strings"
	"unicode"

	"cpf_importer/internal/models"
)

// divide uma linha do arquivo de entrada em colunas separadas por espaços múltiplos.
func SplitLine(line string) []string {
	return strings.Fields(line)
}

// higieniza os campos de um cliente, removendo acentos e caracteres não numéricos.
func CleanseData(client models.Client) models.Client {
	client.MostFrequentStore = sanitizeString(client.MostFrequentStore)
	client.LastPurchaseStore = sanitizeString(client.LastPurchaseStore)

	client.CPF = removeNonDigits(client.CPF)
	client.MostFrequentStore = removeNonDigits(client.MostFrequentStore)
	client.LastPurchaseStore = removeNonDigits(client.LastPurchaseStore)

	return client
}

// remove acentuação e converte a string para maiúsculas.
func sanitizeString(s string) string {
	s = strings.ToUpper(s)
	replacer := strings.NewReplacer(
		"Á", "A", "À", "A", "Ã", "A", "Â", "A",
		"É", "E", "Ê", "E",
		"Í", "I",
		"Ó", "O", "Õ", "O", "Ô", "O",
		"Ú", "U",
		"Ç", "C",
	)
	return replacer.Replace(s)
}

// elimina todos os caracteres que não são dígitos de uma string.
func removeNonDigits(s string) string {
	var builder strings.Builder
	builder.Grow(len(s))
	for _, r := range s {
		if unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

// verifica se o CPF é válido numericamente.
func ValidateCPF(cpf string) bool {
	cpf = removeNonDigits(cpf)
	if len(cpf) != 11 {
		return false
	}

	// Verifica se todos os dígitos são iguais
	allEqual := true
	for i := 1; i < 11; i++ {
		if cpf[i] != cpf[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return false
	}

	// Calcula o primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(cpf[i]-'0') * (10 - i)
	}
	firstDigit := 11 - (sum % 11)
	if firstDigit >= 10 {
		firstDigit = 0
	}
	if int(cpf[9]-'0') != firstDigit {
		return false
	}

	// Calcula o segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(cpf[i]-'0') * (11 - i)
	}
	secondDigit := 11 - (sum % 11)
	if secondDigit >= 10 {
		secondDigit = 0
	}
	return int(cpf[10]-'0') == secondDigit
}

// verifica se o CNPJ é válido numericamente.
func ValidateCNPJ(cnpj string) bool {
	cnpj = removeNonDigits(cnpj)
	if len(cnpj) != 14 {
		return false
	}

	// Verifica se todos os dígitos são iguais
	allEqual := true
	for i := 1; i < 14; i++ {
		if cnpj[i] != cnpj[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return false
	}

	// Pesos para o cálculo dos dígitos verificadores
	weightsFirst := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weightsSecond := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	// Calcula o primeiro dígito verificador
	sum := 0
	for i := 0; i < 12; i++ {
		sum += int(cnpj[i]-'0') * weightsFirst[i]
	}
	firstDigit := 11 - (sum % 11)
	if firstDigit >= 10 {
		firstDigit = 0
	}
	if int(cnpj[12]-'0') != firstDigit {
		return false
	}

	// Calcula o segundo dígito verificador
	sum = 0
	for i := 0; i < 13; i++ {
		sum += int(cnpj[i]-'0') * weightsSecond[i]
	}
	secondDigit := 11 - (sum % 11)
	if secondDigit >= 10 {
		secondDigit = 0
	}

	return int(cnpj[13]-'0') == secondDigit
}
