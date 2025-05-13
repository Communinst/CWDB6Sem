package entities

type Dump struct {
	ID       int    `db:"id" json:"id"`             // Unique identifier for the dump record
	Filename string `db:"filename" json:"filename"` // The name of the dump file
	Size     int64  `db:"size" json:"size"`         // Size of the dump file in bytes (optional)
}
