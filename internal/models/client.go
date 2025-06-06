package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	CPF                          string
	Private                      bool
	Incomplete                   bool
	LastPurchaseDate             *time.Time
	AvgTicket                    *float64
	LastPurchaseTicket           *float64
	MostFrequentStore            string
	LastPurchaseStore            string
	IsCPFValid                   bool
	IsMostFrequentStoreCNPJValid bool
	IsLastPurchaseStoreCNPJValid bool
}

func parseBoolFlagFromString(flagStr string) (bool, error) {
	flagStr = strings.TrimSpace(flagStr)
	flagInt, err := strconv.Atoi(flagStr)
	if err != nil {
		return false, err
	}
	return flagInt == 1, nil
}

func parseNullableFloatFromString(floatStr string) (*float64, error) {
	floatStr = strings.TrimSpace(floatStr)
	if floatStr == "" || floatStr == "NULL" {
		return nil, nil
	}
	floatStr = strings.ReplaceAll(floatStr, ",", ".")
	parsedFloat, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		return nil, err
	}
	return &parsedFloat, nil
}

func parseNullableDateFromString(dateStr string) (*time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" || dateStr == "NULL" {
		return nil, nil
	}
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &parsedDate, nil
}

func NewClientFromParts(parts []string) (Client, error) {
	if len(parts) != 8 {
		return Client{}, fmt.Errorf("esperado 8 campos para cliente, recebido %d", len(parts))
	}

	cpf := strings.TrimSpace(parts[0])

	isPrivateClient, err := parseBoolFlagFromString(parts[1])
	if err != nil {
		return Client{}, fmt.Errorf("erro ao converter campo PRIVATE '%s': %w", parts[1], err)
	}

	hasIncompleteData, err := parseBoolFlagFromString(parts[2])
	if err != nil {
		return Client{}, fmt.Errorf("erro ao converter campo INCOMPLETE '%s': %w", parts[2], err)
	}

	lastPurchaseDate, err := parseNullableDateFromString(parts[3])
	if err != nil {
		return Client{}, fmt.Errorf("erro ao converter DATA DA ÚLTIMA COMPRA '%s': %w", parts[3], err)
	}

	averageTicketValue, err := parseNullableFloatFromString(parts[4])
	if err != nil {
		return Client{}, fmt.Errorf("erro ao converter TICKET MÉDIO '%s': %w", parts[4], err)
	}

	lastPurchaseTicketValue, err := parseNullableFloatFromString(parts[5])
	if err != nil {
		return Client{}, fmt.Errorf("erro ao converter TICKET DA ÚLTIMA COMPRA '%s': %w", parts[5], err)
	}

	mostFrequentStore := strings.TrimSpace(parts[6])
	lastPurchaseStore := strings.TrimSpace(parts[7])

	return Client{
		CPF:                cpf,
		Private:            isPrivateClient,
		Incomplete:         hasIncompleteData,
		LastPurchaseDate:   lastPurchaseDate,
		AvgTicket:          averageTicketValue,
		LastPurchaseTicket: lastPurchaseTicketValue,
		MostFrequentStore:  mostFrequentStore,
		LastPurchaseStore:  lastPurchaseStore,
	}, nil
}
