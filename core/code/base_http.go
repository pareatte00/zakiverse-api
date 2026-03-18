package code

const (
	// @Http StatusOK
	// @Locale EN "OK"
	// @Locale TH "สำเร็จ"
	HttpOK Code = "HTP-200"

	// @Http StatusCreated
	// @Locale EN "Created"
	// @Locale TH "สร้างสำเร็จ"
	HttpCreated Code = "HTP-201"

	// @Http StatusNoContent
	// @Locale EN "No content"
	// @Locale TH "ไม่มีเนื้อหา"
	HttpNoContent Code = "HTP-204"

	// @Http StatusBadRequest
	// @Locale EN "Bad request"
	// @Locale TH "คำขอไม่ถูกต้อง"
	HttpBadRequest Code = "HTP-400"

	// @Http StatusUnauthorized
	// @Locale EN "Unauthorized"
	// @Locale TH "ไม่ได้รับอนุญาต"
	HttpUnauthorized Code = "HTP-401"

	// @Http StatusForbidden
	// @Locale EN "Forbidden"
	// @Locale TH "ไม่มีสิทธิ์เข้าถึง"
	HttpForbidden Code = "HTP-403"

	// @Http StatusNotFound
	// @Locale EN "Not found"
	// @Locale TH "ไม่พบข้อมูล"
	HttpNotFound Code = "HTP-404"

	// @Http StatusConflict
	// @Locale EN "Conflict"
	// @Locale TH "ข้อมูลขัดแย้ง"
	HttpConflict Code = "HTP-409"

	// @Http StatusInternalServerError
	// @Locale EN "Internal server error"
	// @Locale TH "เกิดข้อผิดพลาดภายในระบบ"
	HttpInternalServerError Code = "HTP-500"
)
