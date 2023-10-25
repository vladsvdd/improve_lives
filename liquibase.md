# Install Liquibase on Linux with Debian/Ubuntu

You can install, upgrade, and uninstall Liquibase on your Linux machine with the Debian/Ubuntu package installer by
following the steps on this page.

## Install

To install Liquibase using the Debian/Ubuntu installer package, follow the steps below:

Open a terminal.
Run the following command to import the Liquibase GPG key and add the Liquibase repository to the apt sources list:

```bash
wget -O- https://repo.liquibase.com/liquibase.asc | gpg --dearmor > liquibase-keyring.gpg && \
cat liquibase-keyring.gpg | sudo tee /usr/share/keyrings/liquibase-keyring.gpg > /dev/null && \
echo 'deb [arch=amd64 signed-by=/usr/share/keyrings/liquibase-keyring.gpg] https://repo.liquibase.com stable main' | sudo tee /etc/apt/sources.list.d/liquibase.list
```

Update the package lists:

```bash
sudo apt-get update
```

Install Liquibase:

```bash
sudo apt-get install liquibase
```

Liquibase is now installed on your system.

### Optional: install a specific version

To install a specific version of Liquibase using the Debian/Ubuntu installer package, you can use the package manager's
version pinning mechanism. Run the following command and replace x.y.z with the desired version number:

```bash
sudo apt-get install liquibase= x.y.z
```

The specified version of Liquibase is now installed on your system.

## Upgrade

To upgrade Liquibase to the latest version, follow these steps:

Open a terminal.
Update the package lists:

```bash
sudo apt-get update
```

Upgrade Liquibase:

```bash
sudo apt-get upgrade liquibase
```

Liquibase has now been upgraded to the latest version available in the repository.

## Uninstall

To uninstall Liquibase from your system, follow these steps:

Open a terminal.
Remove the Liquibase package:

```bash
sudo apt-get remove liquibase
```

Liquibase is now uninstalled from your system.

# После установки

Интеграция Liquibase в ваш проект включает в себя несколько шагов. Liquibase - это инструмент управления версиями и
миграции схемы базы данных, который помогает вам управлять изменениями схемы базы данных со временем. Вот пошаговое
руководство по интеграции Liquibase в ваш проект на Go:

**Шаг 1: Установите Liquibase**

Перед тем как использовать Liquibase, вам нужно его установить. Вы можете скачать Liquibase с официального
веб-сайта: https://www.liquibase.org/download

После скачивания распакуйте архив и убедитесь, что исполняемый файл Liquibase (`liquibase` или `liquibase.bat` в
Windows) находится в переменной PATH вашей системы.

**Шаг 2: Создайте проект Liquibase**

В корневом каталоге вашего проекта создайте новый каталог для хранения конфигурации Liquibase и файлов журнала
изменений. Например:

```
mkdir db
cd db
```

**Шаг 3: Создайте файл конфигурации Liquibase**

Создайте файл конфигурации Liquibase с именем `liquibase.properties` в каталоге `db`. Вот пример базового
файла `liquibase.properties`:

```properties
url=jdbc:mysql://localhost:3306/your_database
username=your_username
password=your_password
driver=com.mysql.cj.jdbc.Driver
classpath=path/to/your/database/driver.jar
changeLogFile=db-changelog.xml
```

Замените `your_database`, `your_username`, `your_password` и `path/to/your/database/driver.jar` на фактические данные
для подключения к вашей базе данных.

Путь к `driver.jar` зависит от базы данных, которую вы используете, и способа, которым вы устанавливаете этот драйвер.
Обычно вы можете скачать необходимый драйвер напрямую с сайта базы данных или воспользоваться менеджером зависимостей,
таким как Maven или Gradle, чтобы автоматически загрузить драйвер.

Вот несколько примеров для различных баз данных:

**MySQL:**

Для MySQL драйвер можно скачать с официального сайта MySQL или использовать инструмент управления зависимостями, такой
как `go-sql-driver/mysql`:

```shell
go get github.com/go-sql-driver/mysql
```

**PostgreSQL:**

Для PostgreSQL драйвер можно скачать с официального сайта PostgreSQL или использовать инструмент управления
зависимостями, такой как `lib/pq`:

```shell
go get github.com/lib/pq
```

**SQLite:**

Для SQLite драйвер обычно включен в стандартную библиотеку Go, поэтому дополнительной установки не требуется.

**Oracle:**

Для Oracle вам, возможно, потребуется загрузить официальный драйвер Oracle Database с официального сайта Oracle и
добавить его в свой проект Go вручную.

После того как вы установили драйвер, путь к `driver.jar` будет зависеть от вашей операционной системы и того, как вы
организовали свой проект. В общем случае, вы можете указать путь к `driver.jar` относительно корневой директории вашего
проекта. Например:

```
path/to/your/database/driver.jar
```

Если вы используете инструмент управления зависимостями, путь к драйверу может быть автоматически установлен в GOPATH
или vendor-директорию вашего проекта, и вы можете указать путь относительно этой директории.

**Шаг 4: Создайте файл журнала изменений Liquibase**

Создайте файл журнала изменений Liquibase с именем `db-changelog.xml` в том же каталоге `db`. Этот XML-файл будет
содержать ваши изменения схемы базы данных в структурированном формате.

Вот пример базового файла `db-changelog.xml` с простым набором изменений:

```php
<?xml version="1.0" encoding="UTF-8"?>
<databaseChangeLog xmlns="http://www.liquibase.org/xml/ns/dbchangelog"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="http://www.liquibase.org/xml/ns/dbchangelog
    http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-4.8.xsd">
    
    <changeSet id="1" author="your_name">
        <createTable tableName="example_table">
            <column name="id" type="BIGINT" autoIncrement="true">
                <constraints primaryKey="true" nullable="false"/>
            </column>
            <column name="name" type="VARCHAR(255)"/>
        </createTable>
    </changeSet>

</databaseChangeLog>
```

**Шаг 5: Запустите миграции Liquibase**

Теперь, когда вы настроили проект Liquibase с файлом конфигурации и файлом журнала изменений, вы можете запустить
Liquibase для применения изменений схемы базы данных. Откройте терминал и перейдите в каталог `db`, где находится ваш
файл `liquibase.properties`. Затем выполните следующую команду:

```bash
liquibase update
```
or with logs
```bash
liquibase update --log-level flag
```

Эта команда выполнит изменения базы данных, указанные в вашем файле `db-changelog.xml`.

**Шаг 6: Автоматизация Liquibase в вашем проекте на Go**

Для автоматизации миграций Liquibase как части вашего проекта на Go вы можете использовать пакет `os/exec` Go для
программного выполнения команд Liquibase. Вот пример того, как это можно сделать:

```
package main

import (
	"os"
	"os/exec"
)

func main() {
	// Перейдите в каталог проекта Liquibase
	os.Chdir("db")

	// Запустите команду Liquibase для обновления
	cmd := exec.Command("liquibase", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// Обработайте ошибку
	}
}
```

Этот код переключает выполнение в каталог `db` и запускает команду обновления Liquibase. Вы можете вызывать этот код как
часть процесса сборки или развертывания вашего проекта.

С этими шагами вы интегрировали Liquibase в свой проект на Go и можете более эффективно управлять изменениями схемы базы
данных с помощью инструментов управления версиями и миграции Liquibase.