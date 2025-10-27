package entity

import "github.com/Amierza/chat-service/constants"

type (
	Role           string
	Degree         string
	Progress       string
	SessionStatus  string
	ScheduleStatus string
)

const (
	STUDENT            Role = constants.ENUM_ROLE_STUDENT
	LECTURER           Role = constants.ENUM_ROLE_LECTURER
	PRIMARY_LECTURER   Role = constants.ENUM_ROLE_PRIMARY_LECTURER
	SECONDARY_LECTURER Role = constants.ENUM_ROLE_SECONDARY_LECTURER

	SCHEDULE_PENDING  ScheduleStatus = constants.ENUM_SCHEDULE_STATUS_PENDING
	SCHEDULE_APPROVED ScheduleStatus = constants.ENUM_SCHEDULE_STATUS_APPROVED
	SCHEDULE_REJECTED ScheduleStatus = constants.ENUM_SCHEDULE_STATUS_REJECTED

	S1 Degree = constants.ENUM_DEGREE_S1
	S2 Degree = constants.ENUM_DEGREE_S2
	S3 Degree = constants.ENUM_DEGREE_S3

	BAB1             Progress = constants.ENUM_PROGRESS_BAB1
	BAB2             Progress = constants.ENUM_PROGRESS_BAB2
	BAB3             Progress = constants.ENUM_PROGRESS_BAB3
	BAB4             Progress = constants.ENUM_PROGRESS_BAB4
	BAB5             Progress = constants.ENUM_PROGRESS_BAB5
	SEMINAR_PROPOSAL Progress = constants.ENUM_PROGRESS_SEMINAR_PROPOSAL
	SEMINAR_HASIL    Progress = constants.ENUM_PROGRESS_SEMINAR_HASIL

	WAITING            SessionStatus = constants.ENUM_SESSION_STATUS_WAITING
	ONGOING            SessionStatus = constants.ENUM_SESSION_STATUS_ONGOING
	PROCESSING_SUMMARY SessionStatus = constants.ENUM_SESSION_STATUS_PROCESSING_SUMMARY
	FINISHED           SessionStatus = constants.ENUM_SESSION_STATUS_FINSIHED
)

func IsValidRole(r Role) bool {
	return r == STUDENT || r == PRIMARY_LECTURER || r == SECONDARY_LECTURER
}
func IsValidDegree(d Degree) bool {
	return d == S1 || d == S2 || d == S3
}
func IsValidProgress(p Progress) bool {
	return p == BAB1 || p == BAB2 || p == BAB3 || p == BAB4 || p == BAB5 || p == SEMINAR_PROPOSAL || p == SEMINAR_HASIL
}
func IsValidSessionStatus(ss SessionStatus) bool {
	return ss == WAITING || ss == ONGOING || ss == FINISHED
}
func IsValidScheduleStatus(ss ScheduleStatus) bool {
	return ss == SCHEDULE_PENDING || ss == SCHEDULE_APPROVED || ss == SCHEDULE_REJECTED
}
