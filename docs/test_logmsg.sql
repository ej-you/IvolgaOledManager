DROP TABLE IF EXISTS storage;

CREATE TABLE storage (
    id INT PRIMARY KEY AUTO_INCREMENT,
    level VARCHAR(1) NOT NULL DEFAULT 0,
    header VARCHAR(30) NOT NULL,
    content VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
    storage (level, header, content)
VALUES (
        5,
        "kernel panic",
        "error was occured after using ssc-hmc binary"
    ),
    (
        5,
        "system is malformed",
        "rm -rf / was used"
    ),
    (
        2,
        "new data",
        "humidity sensor data was received"
    ),
    (
        4,
        "ошибка соединения",
        "не удалось установить соединение с датчиком влажности"
    ),
    (
        4,
        "обработка данных",
        "ошибка при расчёте pH"
    ),
    (
        3,
        "Lorem",
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
    ),
    (
        4,
        "интернет соединение",
        "какой-то текст"
    ),
    (
        4,
        "соединение с устройством",
        "Порт занят, невозможно установить соединение с датчиком"
    );