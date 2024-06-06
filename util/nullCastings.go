package util

import "database/sql"

func ToNullInt32(v *int32) (r sql.NullInt32) {
	if v != nil {
		r.Int32 = *v
		r.Valid = true
	}
	return
}

func ToInt32ToNullInt32(v int32) (r sql.NullInt32) {
	r.Int32 = v
	r.Valid = true
	return
}
