```
go mod init improve_lives_bot

go mod tidy
```
### Билдим проект и запускаем в фоне
```
GOOS=windows GOARCH=amd64 go build
nohup go run improve_lives* &
ctrl+c
```

## Чтобы посмотреть процессы приложения

```
ps aux | grep "improve"
```
и потом сделать
```
kill 3432_процесса
```

   ```bash
   sudo ln -fs /etc/nginx/sites-available/uprank.conf /etc/nginx/sites-enabled/uprank.conf
   ```

**Проверьте конфигурацию Nginx** на наличие ошибок:

   ```bash
   sudo nginx -t
   ```
**Перезапустите Nginx**, чтобы применить новую конфигурацию:

   ```bash
   sudo systemctl restart nginx
   ```

**Проверить запущенные процессы**
```netstat -tupln```
sudo netstat -antlp

# Миграции через liquibase
запускаем все миграции
```bash
liquibase update
```
на сколько миграций откатываем
```bash
liquibase rollbackCount 1 
```