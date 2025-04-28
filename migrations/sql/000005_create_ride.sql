CREATE TABLE IF NOT EXISTS rides (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    motorbike_id BIGINT NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP WITH TIME ZONE,
    duration VARCHAR(255),
    cost FLOAT8,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_motorbike FOREIGN KEY (motorbike_id) REFERENCES motorbikes(id)
);

-- İndeksler
CREATE INDEX idx_rides_user_id ON rides(user_id);
CREATE INDEX idx_rides_motorbike_id ON rides(motorbike_id);
CREATE INDEX idx_rides_start_time ON rides(start_time);
CREATE INDEX idx_rides_end_time ON rides(end_time);
CREATE INDEX idx_rides_duration ON rides(duration);

-- Trigger için fonksiyon
CREATE OR REPLACE FUNCTION update_ride_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_rides_updated_at
    BEFORE UPDATE ON rides
    FOR EACH ROW
    EXECUTE FUNCTION update_ride_updated_at();
