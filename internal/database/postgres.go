package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"cpf_importer/internal/models"
)

const (
	dbHost     = "postgres"
	dbPort     = 5432
	dbUser     = "user"
	dbPassword = "password"
	dbName     = "mydatabase"
)

func Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	var db *sql.DB
	var err error

	maxTentativas := 10
	for tentativa := 1; tentativa <= maxTentativas; tentativa++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
		}

		err = db.Ping()
		if err == nil {
			fmt.Println("Conexão com o banco de dados PostgreSQL estabelecida com sucesso!")
			return db, nil
		}

		db.Close()
		tempoEspera := time.Duration(tentativa) * time.Second
		fmt.Printf("Falha ao conectar (tentativa %d/%d): %v. Tentando novamente em %s...\n", tentativa, maxTentativas, err, tempoEspera)
		time.Sleep(tempoEspera)
	}
	return nil, fmt.Errorf("falha ao conectar após %d tentativas: %w", maxTentativas, err)
}

func CreateTable(db *sql.DB) error {
	sqlCreate := `
	CREATE TABLE IF NOT EXISTS clients (
		id SERIAL PRIMARY KEY,
		cpf VARCHAR(20) NOT NULL,
		private BOOLEAN,
		incomplete BOOLEAN,
		last_purchase_date DATE,
		avg_ticket NUMERIC(10, 2),
		last_purchase_ticket NUMERIC(10, 2),
		most_frequent_store VARCHAR(255),
		last_purchase_store VARCHAR(255),
		is_cpf_valid BOOLEAN,
		is_most_frequent_store_cnpj_valid BOOLEAN,
		is_last_purchase_store_cnpj_valid BOOLEAN
	);`

	_, err := db.Exec(sqlCreate)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela clients: %w", err)
	}

	fmt.Println("Tabela 'clients' criada/verificada com sucesso.")
	return nil
}

func InsertClientsBatch(db *sql.DB, clients []models.Client) error {
	if len(clients) == 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
	INSERT INTO clients (
		cpf, private, incomplete, last_purchase_date,
		avg_ticket, last_purchase_ticket, most_frequent_store,
		last_purchase_store, is_cpf_valid, is_most_frequent_store_cnpj_valid,
		is_last_purchase_store_cnpj_valid
	) VALUES `)

	valores := make([]interface{}, 0, len(clients)*11)

	for idx, cliente := range clients {
		placeholders := make([]string, 0, 11)
		for i := 1; i <= 11; i++ {
			placeholders = append(placeholders, fmt.Sprintf("$%d", idx*11+i))
		}

		queryBuilder.WriteString("(" + strings.Join(placeholders, ", ") + ")")
		if idx < len(clients)-1 {
			queryBuilder.WriteString(", ")
		}

		valores = append(valores,
			cliente.CPF,
			cliente.Private,
			cliente.Incomplete,
			cliente.LastPurchaseDate,
			cliente.AvgTicket,
			cliente.LastPurchaseTicket,
			cliente.MostFrequentStore,
			cliente.LastPurchaseStore,
			cliente.IsCPFValid,
			cliente.IsMostFrequentStoreCNPJValid,
			cliente.IsLastPurchaseStoreCNPJValid,
		)
	}
	queryBuilder.WriteString(";")

	_, err = tx.Exec(queryBuilder.String(), valores...)
	if err != nil {
		return fmt.Errorf("erro ao executar inserção em lote: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("erro ao commitar transação: %w", err)
	}

	return nil
}
