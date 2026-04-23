package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"workorder-api/model"
)

type WorkorderRepository struct {
	db *sql.DB
}

func NewWorkorderRepository(db *sql.DB) *WorkorderRepository {
	return &WorkorderRepository{db: db}
}

func (r *WorkorderRepository) UpsertWorkorder(req *model.WorkorderRequest) error {
	rawPayload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	query := `
		INSERT INTO workorder_updates (
			wonum, siteid, status, task, memo, wolo1, wolo3,
			latitude, longitude, cpe_model, cpe_vendor, cpe_serial_number,
			errorcode, suberrorcode, labor_scmt, statusiface,
			urlevidence, engineermemo, np_statusmemo, task_name,
			tk_custom_header_03, tk_custom_header_04, tk_custom_header_09, tk_custom_header_10,
			response_data, response_message, raw_payload, received_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29
		)
		ON CONFLICT (wonum, siteid) 
		DO UPDATE SET
			status = EXCLUDED.status,
			task = EXCLUDED.task,
			memo = EXCLUDED.memo,
			wolo1 = EXCLUDED.wolo1,
			wolo3 = EXCLUDED.wolo3,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			cpe_model = EXCLUDED.cpe_model,
			cpe_vendor = EXCLUDED.cpe_vendor,
			cpe_serial_number = EXCLUDED.cpe_serial_number,
			errorcode = EXCLUDED.errorcode,
			suberrorcode = EXCLUDED.suberrorcode,
			labor_scmt = EXCLUDED.labor_scmt,
			statusiface = EXCLUDED.statusiface,
			urlevidence = EXCLUDED.urlevidence,
			engineermemo = EXCLUDED.engineermemo,
			np_statusmemo = EXCLUDED.np_statusmemo,
			task_name = EXCLUDED.task_name,
			tk_custom_header_03 = EXCLUDED.tk_custom_header_03,
			tk_custom_header_04 = EXCLUDED.tk_custom_header_04,
			tk_custom_header_09 = EXCLUDED.tk_custom_header_09,
			tk_custom_header_10 = EXCLUDED.tk_custom_header_10,
			response_data = EXCLUDED.response_data,
			response_message = EXCLUDED.response_message,
			raw_payload = EXCLUDED.raw_payload,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.db.Exec(query,
		req.Req.Wonum,
		req.Req.Siteid,
		req.Req.Status,
		req.Req.Task,
		req.Req.Memo,
		req.Req.Wolo1,
		req.Req.Wolo3,
		req.Req.Latitude,
		req.Req.Longitude,
		req.Req.CpeModel,
		req.Req.CpeVendor,
		req.Req.CpeSerialNumber,
		req.Req.Errorcode,
		req.Req.Suberrorcode,
		req.Req.LaborScmt,
		req.Req.Statusiface,
		req.Req.Urlevidence,
		req.Req.Engineermemo,
		req.Req.NpStatusmemo,
		req.Req.TaskName,
		req.Req.TkCustomHeader03,
		req.Req.TkCustomHeader04,
		req.Req.TkCustomHeader09,
		req.Req.TkCustomHeader10,
		req.Res.Data,
		req.Res.Message,
		string(rawPayload),
		req.Date,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to upsert workorder: %w", err)
	}

	return nil
}
