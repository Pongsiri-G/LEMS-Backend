package exceptions

import "errors"

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrNotAllowAccess = errors.New("you can't access to this")
var ErrUserNotFound = errors.New("user not found")
var ErrInactiveUser = errors.New("user has been created\nbut user is not activated")
var ErrUserPending = errors.New("บัญชีของคุณยังไม่ได้รับการอนุมัติ กรุณารอผู้ดูแลระบบอนุมัติ")
var ErrUserDeactivated = errors.New("บัญชีของคุณถูกระงับการใช้งาน กรุณาติดต่อผู้ดูแลระบบ")
var ErrUserRejected = errors.New("บัญชีของคุณถูกปฏิเสธ กรุณาติดต่อผู้ดูแลระบบ")
var ErrRegistrationSuccess = errors.New("ลงทะเบียนสำเร็จ กรุณารอผู้ดูแลระบบอนุมัติ")
