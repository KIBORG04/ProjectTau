# internal/service/sarima/sarima_service.py (ОБНОВЛЕННАЯ ВЕРСИЯ)

import pandas as pd
from flask import Flask, request, jsonify
from statsmodels.tsa.statespace.sarimax import SARIMAX
import warnings

warnings.filterwarnings("ignore")
app = Flask(__name__)

def make_prediction(data, date_format_string, frequency, seasonal_period, forecast_steps=12):
    """
    Универсальная функция для построения прогноза SARIMA.
    """
    try:
        ts = pd.Series(data, dtype=float)
        if date_format_string == 'weekly':
            # Для недель "YYYY-WW"
            ts.index = pd.to_datetime(
                ts.index.str.split('-').str[0] + '-W' + ts.index.str.split('-').str[1] + '-1', 
                format='%G-W%V-%u'
            )
        else:
            # Для дней "YYYY-MM-DD"
            ts.index = pd.to_datetime(ts.index, format=date_format_string)
        
        ts = ts.sort_index().asfreq(frequency).fillna(method='ffill')
    except Exception as e:
        return {"error": f"Failed to parse time series data: {e}"}, 400

    # Настраиваем и обучаем модель
    order = (1, 0, 1)
    seasonal_order = (1, 1, 0, seasonal_period)
    
    try:
        model = SARIMAX(ts, order=order, seasonal_order=seasonal_order, enforce_stationarity=False, enforce_invertibility=False)
        model_fit = model.fit(disp=False)
    except Exception as e:
         return {"error": f"Failed to fit SARIMA model: {e}"}, 500

    # Генерируем прогноз
    prediction = model_fit.get_forecast(steps=forecast_steps)
    
    # Форматируем результат
    forecast_result = {}
    for date, value in zip(prediction.predicted_mean.index, prediction.predicted_mean.values):
        if date_format_string == 'weekly':
            year, week, _ = date.isocalendar()
            week_key = f"{year}-{week:02d}" # Добавляем ведущий ноль
            forecast_result[week_key] = round(value, 2)
        else:
            date_key = date.strftime(date_format_string)
            forecast_result[date_key] = round(value, 2)

    return forecast_result, 200


@app.route('/forecast', methods=['POST'])
def create_weekly_forecast():
    """ Endpoint для НЕДЕЛЬНОГО прогноза (старый) """
    json_data = request.get_json()
    if not json_data or 'data' not in json_data:
        return jsonify({"error": "Missing 'data' in request body"}), 400
    
    result, status_code = make_prediction(
        data=json_data['data'],
        date_format_string='weekly',
        frequency='W-MON',      # Еженедельная частота
        seasonal_period=52,     # Сезонность - 52 недели (1 год)
        forecast_steps=12       # Прогноз на 12 недель
    )
    return jsonify(result), status_code

@app.route('/forecast/daily', methods=['POST'])
def create_daily_forecast():
    """ НОВЫЙ Endpoint для ДНЕВНОГО прогноза """
    json_data = request.get_json()
    if not json_data or 'data' not in json_data:
        return jsonify({"error": "Missing 'data' in request body"}), 400

    result, status_code = make_prediction(
        data=json_data['data'],
        date_format_string='%Y-%m-%d',
        frequency='D',          # Дневная частота
        seasonal_period=7,      # Сезонность - 7 дней (1 неделя)
        forecast_steps=21       # Прогноз на 21 день вперед
    )
    return jsonify(result), status_code


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=False) # Выключаем debug для продакшена