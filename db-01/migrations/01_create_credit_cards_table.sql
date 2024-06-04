-- +goose Up
CREATE TABLE credit_cards(
    id    SERIAL NOT NULL,
    number VARCHAR(255) NOT NULL,
    expiration_date VARCHAR(255) NOT NULL,
    cvv INT NOT NULL,
    holder_name VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE credit_cards;