CREATE TABLE bluetooth_connections (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    motorbike_id BIGINT NOT NULL,
    connected_at TIMESTAMPTZ NOT NULL,
    disconnected_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    -- Opsiyonel: ilişkiler (foreign key'ler)
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_motorbike FOREIGN KEY (motorbike_id) REFERENCES motorbikes(id)
);

-- 1. Trigger fonksiyonunu oluştur
CREATE OR REPLACE FUNCTION update_bluetooth_connections_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Trigger'ı tabloya ekle
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON bluetooth_connections
    FOR EACH ROW
EXECUTE FUNCTION update_bluetooth_connections_updated_at();