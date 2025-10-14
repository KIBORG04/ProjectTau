package forecast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"strings"
	"time"
)

// UpdateDailyForecast теперь сохраняет новую запись в историю
func UpdateDailyForecast() {
	log.Println("Updating daily forecast history...")
	// ... (логика получения данных из БД остается такой же) ...
    var dbResult []struct {
		Date    time.Time
		Players float64
	}
	daysToQuery := 180
	dateFrom := time.Now().AddDate(0, 0, -daysToQuery)
	tx := r.Database.Raw(`
        SELECT date::date, round(avg(s.crew_total)) as players
        FROM roots JOIN scores s ON round_id = s.root_id WHERE date >= ?
        GROUP BY date::date ORDER BY date::date;
    `, dateFrom).Scan(&dbResult)
	if tx.Error != nil {
		log.Printf("DB Error in UpdateDailyForecast: %v", tx.Error)
		return
	}
	historicalData := make(map[string]float64, len(dbResult))
	for _, row := range dbResult {
		historicalData[row.Date.Format("2006-01-02")] = row.Players
	}
	if len(historicalData) < 30 {
		log.Println("Not enough data for daily forecast, skipping.")
		return
	}

	pythonServiceURL := "http://forecast-service:5000/forecast/daily"
	forecastJSON, err := getForecastFromPython(pythonServiceURL, historicalData)
	if err != nil {
		log.Printf("Error getting daily forecast from Python: %v", err)
		return
	}

	// Создаем новую запись в истории
	saveNewForecastHistory("daily", forecastJSON)
	log.Println("Successfully added new daily forecast to history.")
}

// UpdateWeeklyForecast теперь сохраняет новую запись в историю
func UpdateWeeklyForecast() {
	log.Println("Updating weekly forecast history...")
    // ... (логика получения данных из БД остается такой же) ...
	var dbResult []struct {
		WeekDate string
		Players  int
	}
	tx := r.Database.Raw(`
        SELECT date_part('isoyear', date) || '-' || date_part('week', date) AS week_date, round(avg(s.crew_total)) AS players
        FROM roots JOIN scores s ON round_id = s.root_id
        GROUP BY week_date
        ORDER BY to_date(date_part('isoyear', date) || '-' || date_part('week', date), 'YYYY-WW');
    `).Scan(&dbResult)
	if tx.Error != nil {
		log.Printf("DB Error in UpdateWeeklyForecast: %v", tx.Error)
		return
	}
	historicalData := make(map[string]float64, len(dbResult))
	for _, onlineDay := range dbResult {
		dateParts := strings.Split(onlineDay.WeekDate, "-")
		if len(dateParts) == 2 {
			if len(dateParts[1]) == 1 { dateParts[1] = "0" + dateParts[1] }
			historicalData[fmt.Sprintf("%s-%s", dateParts[0], dateParts[1])] = float64(onlineDay.Players)
		}
	}
	if len(historicalData) < 20 {
		log.Println("Not enough data for weekly forecast, skipping.")
		return
	}

	pythonServiceURL := "http://forecast-service:5000/forecast"
	forecastJSON, err := getForecastFromPython(pythonServiceURL, historicalData)
	if err != nil {
		log.Printf("Error getting weekly forecast from Python: %v", err)
		return
	}
    
	// Создаем новую запись в истории
	saveNewForecastHistory("weekly", forecastJSON)
	log.Println("Successfully added new weekly forecast to history.")
}

// saveNewForecastHistory всегда создает новую запись
func saveNewForecastHistory(forecastType string, data json.RawMessage) {
	newForecast := domain.ForecastHistory{
		ForecastType: forecastType,
		Data:         data,
	}
	if err := r.Database.Create(&newForecast).Error; err != nil {
		log.Printf("Error creating new '%s' forecast history entry: %v", forecastType, err)
	}
}

// getForecastFromPython (без изменений)
func getForecastFromPython(url string, data map[string]float64) (json.RawMessage, error) {
    // ... код этой функции остается точно таким же
	requestPayload := map[string]interface{}{"data": data}
    payloadBytes, _ := json.Marshal(requestPayload)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil { return nil, fmt.Errorf("error calling Python service: %w", err) }
    defer resp.Body.Close()
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil { return nil, fmt.Errorf("failed to read forecast response: %w", err) }
    if resp.StatusCode != http.StatusOK { return nil, fmt.Errorf("python service returned error. Status: %s, Body: %s", resp.Status, string(responseBody)) }
    return responseBody, nil
}