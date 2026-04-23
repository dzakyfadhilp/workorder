-- Database Schema for Workorder Updates API
-- PostgreSQL 12+

CREATE TABLE IF NOT EXISTS workorder_updates (
    id BIGSERIAL PRIMARY KEY,
    
    -- Request tracking
    request_id VARCHAR(100) NOT NULL,
    
    -- Required fields (indexed for fast lookup)
    wonum VARCHAR(50) NOT NULL,
    siteid VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    
    -- Request data fields
    task VARCHAR(100),
    memo TEXT,
    wolo1 VARCHAR(100),
    wolo3 VARCHAR(100),
    latitude VARCHAR(50),
    longitude VARCHAR(50),
    cpe_model VARCHAR(100),
    cpe_vendor VARCHAR(100),
    cpe_serial_number VARCHAR(100),
    errorcode VARCHAR(50),
    suberrorcode VARCHAR(50),
    labor_scmt VARCHAR(50),
    statusiface VARCHAR(10),
    urlevidence TEXT,
    engineermemo TEXT,
    np_statusmemo TEXT,
    task_name VARCHAR(200),
    tk_custom_header_03 VARCHAR(100),
    tk_custom_header_04 VARCHAR(100),
    tk_custom_header_09 VARCHAR(100),
    tk_custom_header_10 VARCHAR(100),
    
    -- Response data (for audit trail)
    response_data TEXT,
    response_message TEXT,
    
    -- Audit fields
    raw_payload JSONB NOT NULL,
    received_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT unique_wonum_siteid UNIQUE (wonum, siteid)
);

-- Indexes for performance
CREATE INDEX idx_request_id ON workorder_updates(request_id);
CREATE INDEX idx_wonum ON workorder_updates(wonum);
CREATE INDEX idx_siteid ON workorder_updates(siteid);
CREATE INDEX idx_status ON workorder_updates(status);
CREATE INDEX idx_received_at ON workorder_updates(received_at DESC);
CREATE INDEX idx_raw_payload ON workorder_updates USING GIN (raw_payload);

-- Comments for documentation
COMMENT ON TABLE workorder_updates IS 'Stores workorder status updates with full audit trail';
COMMENT ON COLUMN workorder_updates.request_id IS 'Unique request ID for tracking';
COMMENT ON COLUMN workorder_updates.wonum IS 'Work order number - primary business key';
COMMENT ON COLUMN workorder_updates.siteid IS 'Site identifier - part of composite key';
COMMENT ON COLUMN workorder_updates.raw_payload IS 'Complete JSON payload for audit purposes';
COMMENT ON COLUMN workorder_updates.received_at IS 'Timestamp when request was received';
COMMENT ON COLUMN workorder_updates.updated_at IS 'Timestamp of last update';
