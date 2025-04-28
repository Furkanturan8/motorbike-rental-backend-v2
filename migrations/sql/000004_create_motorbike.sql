CREATE TABLE IF NOT EXISTS motorbikes (
    id BIGSERIAL PRIMARY KEY,
    model VARCHAR(255),
    location_latitude FLOAT8,
    location_longitude FLOAT8,
    status motorbike_status NOT NULL,
    lock_status lock_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS motorbike_photos (
    id BIGSERIAL PRIMARY KEY,
    motorbike_id BIGINT NOT NULL,
    photo_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_motorbike FOREIGN KEY (motorbike_id) REFERENCES motorbikes(id)
);

-- İndeksler
CREATE INDEX idx_motorbikes_status ON motorbikes(status);
CREATE INDEX idx_motorbikes_lock_status ON motorbikes(lock_status);
CREATE INDEX idx_motorbike_photos_motorbike_id ON motorbike_photos(motorbike_id);

-- Trigger için fonksiyon
CREATE OR REPLACE FUNCTION update_motorbike_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_motorbikes_updated_at
    BEFORE UPDATE ON motorbikes
    FOR EACH ROW
    EXECUTE FUNCTION update_motorbike_updated_at();

CREATE OR REPLACE FUNCTION update_motorbike_photo_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_motorbike_photos_updated_at
    BEFORE UPDATE ON motorbike_photos
    FOR EACH ROW
    EXECUTE FUNCTION update_motorbike_photo_updated_at();
