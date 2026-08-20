-- +goose Up
CREATE TABLE training_item_relations (
    training_id      BIGINT NOT NULL,
    training_item_id BIGINT NOT NULL,
    PRIMARY KEY (training_id, training_item_id),
    UNIQUE (training_item_id)
);

-- +goose Down
DROP TABLE training_item_relations;
