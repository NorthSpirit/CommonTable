-- +goose Up

ALTER TABLE chat_messages
    ADD COLUMN deleted_at    timestamp,
    ADD COLUMN deleted_by    text,
    ADD COLUMN delete_reason text,
    ADD CONSTRAINT chat_messages_delete_reason_length CHECK (char_length(delete_reason) <= 255),
    ADD CONSTRAINT chat_messages_deleted_fields CHECK (deleted_at IS NOT NULL OR (deleted_by IS NULL AND delete_reason IS NULL));

-- +goose Down

ALTER TABLE chat_messages
    DROP CONSTRAINT chat_messages_deleted_fields,
    DROP CONSTRAINT chat_messages_delete_reason_length,
    DROP COLUMN delete_reason,
    DROP COLUMN deleted_by,
    DROP COLUMN deleted_at;
