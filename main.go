package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"cpf_importer/internal/database"
	"cpf_importer/internal/models"
	"cpf_importer/internal/utils"
)

const (
	filePath      = "./base_teste.txt"
	batchSize     = 1000  
	numProcessors = 4    
)

func main() {
	start := time.Now()
	fmt.Println("Iniciando o serviço de importação e higienização de dados...")

	db, err := database.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao conectar ao banco de dados: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err = database.CreateTable(db); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar tabela: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao abrir o arquivo %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() { 
		fmt.Fprintln(os.Stderr, "Arquivo vazio ou com problema no cabeçalho")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	clientChan := make(chan models.Client, batchSize*2) 

	for i := 0; i < numProcessors; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			batch := make([]models.Client, 0, batchSize)
			for client := range clientChan {
				processed := utils.CleanseData(client)
				processed.IsCPFValid = utils.ValidateCPF(processed.CPF)
				processed.IsMostFrequentStoreCNPJValid = utils.ValidateCNPJ(processed.MostFrequentStore)
				processed.IsLastPurchaseStoreCNPJValid = utils.ValidateCNPJ(processed.LastPurchaseStore)

				batch = append(batch, processed)
				if len(batch) >= batchSize {
					if err := database.InsertClientsBatch(db, batch); err != nil {
						fmt.Fprintf(os.Stderr, "Worker %d: erro ao inserir lote: %v\n", workerID, err)
					}
					batch = batch[:0] 
				}
			}
			if len(batch) > 0 {
				if err := database.InsertClientsBatch(db, batch); err != nil {
					fmt.Fprintf(os.Stderr, "Worker %d: erro ao inserir lote final: %v\n", workerID, err)
				}
			}
		}(i + 1)
	}

	lineNumber := 1
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		parts := utils.SplitLine(line)

		if len(parts) != 8 {
			fmt.Fprintf(os.Stderr, "Linha %d com formato inválido, pulando: %s\n", lineNumber, line)
			continue
		}

		client, err := models.NewClientFromParts(parts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao parsear linha %d: %v, pulando: %s\n", lineNumber, err, line)
			continue
		}

		clientChan <- client
	}

	close(clientChan) 
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro na leitura do arquivo: %v\n", err)
	}

	fmt.Printf("Importação e higienização concluídas em %s\n", time.Since(start))
}
