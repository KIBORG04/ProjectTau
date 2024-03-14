# SS13 Statistics

* **Tau Ceti:** https://github.com/TauCetiStation/TauCetiClassic/
* **Code:** https://github.com/KIBORG04/ProjectTau
* **Поддержать проект (и посмотреть на кошку):** [Boosty](https://boosty.to/kib.org) 

---

## Описание
SS13 Statistics - это система для агрегации и обработке статистики по билду игры SS13 Tau Ceti. Сейчас реализовано множество разнообразных представлений данных, личная статистика игроков, тепловые карты и другое.

## Тепловая карта смертей
<p align="center">
  <img src="https://github.com/KIBORG04/ProjectTau/blob/master/deaths.png?raw=true" alt="Diagram"/>
</p>

## Развёртывание 
```shell
git clone https://github.com/KIBORG04/ProjectTau.git
```

Свои настройки конфигурации вы можете указать в `./config/config.yaml` и `docker-compose.yml`.

> [!IMPORTANT]
> Для работы конфигов, вам необходимо убрать из названия файла `./config/config-example.yaml` 
> постфикс `-example`. Затем вписать свои настройки.

Запуск через docker-compose:
```shell
docker-compose build
docker-compose up
```

Сервис запустится по адресу `localhost:8080`. 

