package entity

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ProposedAt  time.Time      `gorm:"not null" json:"proposed_at"` // waktu pengajuan jadwal
	StartTime   time.Time      `gorm:"not null" json:"start_time"`  // waktu mulai yang diajukan
	EndTime     time.Time      `gorm:"not null" json:"end_time"`    // waktu selesai yang diajukan
	Status      ScheduleStatus `gorm:"default:pending" json:"status"`
	Description string         `json:"description,omitempty"` // opsional: deskripsi singkat
	Location    string         `json:"location,omitempty"`    // bisa untuk link meet / ruang bimbingan

	// Relasi
	ThesisID uuid.UUID `gorm:"type:uuid;index" json:"thesis_id"`
	Thesis   Thesis    `gorm:"foreignKey:ThesisID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"thesis"`

	CreatedByID uuid.UUID `gorm:"type:uuid;index" json:"created_by_id"` // user (mahasiswa/dosen) yang buat jadwal
	CreatedBy   User      `gorm:"foreignKey:CreatedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"created_by"`

	ApprovedByID *uuid.UUID `gorm:"type:uuid;index" json:"approved_by_id,omitempty"` // dosen yang menyetujui
	ApprovedBy   *User      `gorm:"foreignKey:ApprovedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"approved_by,omitempty"`

	// Jika jadwal sudah disetujui dan sesi sudah berjalan, bisa dihubungkan
	// SessionID *uuid.UUID `gorm:"type:uuid;index" json:"session_id,omitempty"`
	// Session   *Session   `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"session,omitempty"`

	TimeStamp
}
