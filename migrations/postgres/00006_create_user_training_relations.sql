-- +goose Up
CREATE TABLE user_trainings (
    training_id BIGINT NOT NULL,
    user_id     BIGINT NOT NULL,
    PRIMARY KEY (user_id, training_id),
    UNIQUE (training_id)
);

-- +goose Down
DROP TABLE user_trainings;
