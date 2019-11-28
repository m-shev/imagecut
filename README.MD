# Сервис "Нарезка картинок"

## Команды
* Скомпилировать проект ```make build```
* Скомпилировать и запустить ```make run```  
проверить ```http://localhost:3000/crop/900/900/?origin=https://www.1zoom.ru/big2/49/179561-melisenta.jpg```
* Запустить тесты ```make test```
* Запустить в docker container для тестирования ```make qa```
* Запустить в docker container для production  ```make prod```

## Переменные окружения
* IMAGECUT_CACHE_SIZE - максимальный размер LRU в байтах:  
```docker run -p=3000:3000 --rm -e IMAGECUT_ENV=PROD -e IMAGECUT_CACHE_SIZE=524288000 -v imagecut-data:/opt/imagecut/app-data -d mshev/imagecut```
* IMAGECUT_CACHE_PATH - путь к файлу (json) с данными LRU кэша, при запуске сервера выполняется попытка прочитать, при остановке перезаписывается  
```IMAGECUT_CACHE_PATH=folder/cache.json ./build/imagecut```
* IMAGECUT_IMAGE_FOLDER - путь к папке для хранения кэшированных изображений  
```IMAGECUT_IMAGE_FOLDER=/folder/images ./build/imagecut```