package utils

import (
	"strings"

	commonv1 "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/common"
)

func ParseUserType(s string) *commonv1.UserType {
	s = strings.ToLower(s)
	switch s {
	case "student":
		return commonv1.UserType_USER_TYPE_STUDENT.Enum()
	case "employee":
		return commonv1.UserType_USER_TYPE_EMPLOYEE.Enum()
	case "service":
		return commonv1.UserType_USER_TYPE_SERVICE.Enum()
	case "other":
		return commonv1.UserType_USER_TYPE_OTHER.Enum()
	default:
		return nil
	}
}

func ParseUserStatus(s string) *commonv1.UserStatus {
	s = strings.ToLower(s)
	switch s {
	case "pending":
		return commonv1.UserStatus_USER_STATUS_PENDING.Enum()
	case "active":
		return commonv1.UserStatus_USER_STATUS_ACTIVE.Enum()
	case "graduated":
		return commonv1.UserStatus_USER_STATUS_GRADUATED.Enum()
	case "leave":
		return commonv1.UserStatus_USER_STATUS_LEAVE.Enum()
	case "suspended":
		return commonv1.UserStatus_USER_STATUS_SUSPENDED.Enum()
	case "withdrawn":
		return commonv1.UserStatus_USER_STATUS_WITHDRAWN.Enum()
	case "expelled":
		return commonv1.UserStatus_USER_STATUS_EXPELLED.Enum()
	default:
		return nil
	}
}

// ParseStudentClassStatus 문자열을 StudentClassStatus enum으로 변환
func ParseStudentClassStatus(status string) *commonv1.StudentClassStatus {
	switch status {
	case "active":
		return commonv1.StudentClassStatus_STUDENT_CLASS_STATUS_ACTIVE.Enum()
	case "graduated":
		return commonv1.StudentClassStatus_STUDENT_CLASS_STATUS_GRADUATED.Enum()
	default:
		return nil
	}
}

func ParseUserRole(s string) *commonv1.UserRole {
	s = strings.ToLower(s)
	switch s {
	case "student":
		return commonv1.UserRole_USER_ROLE_STUDENT.Enum()
	case "doorkeeper":
		return commonv1.UserRole_USER_ROLE_DOORKEEPER.Enum()
	case "class_rep":
		return commonv1.UserRole_USER_ROLE_CLASS_REP.Enum()
	case "ta":
		return commonv1.UserRole_USER_ROLE_TA.Enum()
	case "professor":
		return commonv1.UserRole_USER_ROLE_PROFESSOR.Enum()
	case "admin":
		return commonv1.UserRole_USER_ROLE_ADMIN.Enum()
	case "default":
		return commonv1.UserRole_USER_ROLE_DEFAULT.Enum()
	default:
		return nil
	}
}

func StringFromUserType(ut commonv1.UserType) string {
	return strings.ToLower(strings.TrimPrefix(ut.String(), "USER_TYPE_"))
}

func StringFromUserStatus(us commonv1.UserStatus) string {
	return strings.ToLower(strings.TrimPrefix(us.String(), "USER_STATUS_"))
}

func StringFromStudentClassStatus(scs commonv1.StudentClassStatus) string {
	return strings.ToLower(strings.TrimPrefix(scs.String(), "STUDENT_CLASS_STATUS_"))
}

func StringFromUserRole(ur commonv1.UserRole) string {
	return strings.ToLower(strings.TrimPrefix(ur.String(), "USER_ROLE_"))
}
