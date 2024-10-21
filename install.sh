#!/bin/bash

# Определение переменных
PROJECT_NAME="WednesdayPollTelegramBot"
PROJECT_DIR="$(pwd)"
BUILD_DIR="$PROJECT_DIR/build"
AUTOSTART_DIR="$HOME/.config/autostart"
DESKTOP_FILE="$AUTOSTART_DIR/$PROJECT_NAME.desktop"
CONFIG_FILE="$PROJECT_DIR/app/config/config.yaml"

# Функция для сборки проекта
build_project() {
    echo "Сборка проекта..."
    if [ ! -f "$CONFIG_FILE" ]; then
        echo "Файл конфигурации не найден. Пожалуйста, создайте его."
        exit 1
    else
        echo "Файл конфигурации найден. Продолжение сборки проекта..."
    fi

    cd "$PROJECT_DIR/app"
    go build -o "$BUILD_DIR/$PROJECT_NAME"
    if [ $? -eq 0 ]; then
        echo "Проект успешно собран."
        chmod +x "$BUILD_DIR/$PROJECT_NAME"
        echo "Права на выполнение добавлены."
        #скопировать файл конфигурации в директорию сборки
        mkdir -p "$BUILD_DIR/config"
        cp "$CONFIG_FILE" "$BUILD_DIR/config"
        echo "Файл конфигурации скопирован в директорию сборки."
    else
        echo "Ошибка при сборке проекта."
        exit 1
    fi
}

# Функция для настройки автозапуска после загрузки рабочего стола
setup_desktop_autostart() {
    # проверка наличия исполняемого файла
    if [ ! -f "$BUILD_DIR/$PROJECT_NAME" ]; then
        echo "Ошибка: исполняемый файл не найден."
        exit 1
    fi

    # создание директории автозапуска, если она не существует
    if [ ! -d "$AUTOSTART_DIR" ]; then
        mkdir -p "$AUTOSTART_DIR"
    fi

    # настройка автозапуска
    if [ ! -f "$DESKTOP_FILE" ]; then
        echo "Настройка автозапуска после загрузки рабочего стола..."
        cat << EOF > "$DESKTOP_FILE"
[Desktop Entry]
Type=Application
Exec=$BUILD_DIR/$PROJECT_NAME
Hidden=false
NoDisplay=false
Name=$PROJECT_NAME
Comment=Автозапуск $PROJECT_NAME по средам и четвергам
Icon=$PROJECT_DIR/app/logo.png
X-GNOME-Autostart-enabled=true
EOF
        echo "Автозапуск настроен."
    else
        echo "Автозапуск уже настроен."
    fi
}

# Основная логика скрипта
main() {
    build_project
    setup_desktop_autostart

    echo "Установка завершена."
}

# Запуск основной функции
main
