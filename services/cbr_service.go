package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/beevik/etree"
)

// buildSOAPRequest формирует SOAP-запрос для получения ключевой ставки.
// Формируется запрос за последние 30 дней.
func buildSOAPRequest() string {
	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
		<soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
			<soap12:Body>
				<KeyRate xmlns="http://web.cbr.ru/">
					<fromDate>%s</fromDate>
					<ToDate>%s</ToDate>
				</KeyRate>
			</soap12:Body>
		</soap12:Envelope>`, fromDate, toDate)
}

// sendSOAPRequest отправляет SOAP-запрос к ЦБ РФ и возвращает сырой ответ.
func sendSOAPRequest(soapRequest string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx", bytes.NewBuffer([]byte(soapRequest)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://web.cbr.ru/KeyRate")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	return body, nil
}

// parseSOAPResponse парсит XML-ответ и извлекает ключевую ставку.
func parseSOAPResponse(rawBody []byte) (float64, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(rawBody); err != nil {
		return 0, fmt.Errorf("failed to parse XML: %v", err)
	}
	// Ищем элементы KR внутри diffgram.
	krElements := doc.FindElements("//diffgram/KeyRate/KR")
	if len(krElements) == 0 {
		return 0, errors.New("key rate data not found")
	}
	latestKR := krElements[0]
	rateElement := latestKR.FindElement("./Rate")
	if rateElement == nil {
		return 0, errors.New("Rate element not found")
	}
	rateStr := rateElement.Text()
	var rate float64
	if _, err := fmt.Sscanf(rateStr, "%f", &rate); err != nil {
		return 0, fmt.Errorf("failed to convert rate: %v", err)
	}
	return rate, nil
}

// GetCentralBankRate обращается к ЦБ РФ, получает ключевую ставку, и добавляет маржу банка (например, +5%).
func GetCentralBankRate() (float64, error) {
	soapRequest := buildSOAPRequest()
	rawBody, err := sendSOAPRequest(soapRequest)
	if err != nil {
		return 0, err
	}
	rate, err := parseSOAPResponse(rawBody)
	if err != nil {
		return 0, err
	}
	// Добавляем маржу, например +5
	rate += 5
	return rate, nil
}
