package validation

// getDefaultMessages returns default error messages for locale
func getDefaultMessages(locale string) map[string]string {
	switch locale {
	case "id":
		return getIndonesianMessages()
	default:
		return getEnglishMessages()
	}
}

// getEnglishMessages returns English error messages
func getEnglishMessages() map[string]string {
	return map[string]string{
		"required": "{field} is required",
		"min":      "{field} must be at least {param} characters",
		"max":      "{field} must be at most {param} characters",
		"email":    "{field} must be a valid email address",
		"oneof":    "{field} must be one of: {param}",
		"numeric":  "{field} must be numeric",
		"alpha":    "{field} must contain only letters",
		"alphanum": "{field} must contain only letters and numbers",
		"url":      "{field} must be a valid URL",
		"uuid":     "{field} must be a valid UUID",
		"phone":    "{field} must be a valid phone number",
		"nik":      "{field} must be a valid NIK",
		"phone_id": "{field} must be a valid Indonesian phone number",
	}
}

// getIndonesianMessages returns Indonesian error messages
func getIndonesianMessages() map[string]string {
	return map[string]string{
		"required": "{field} wajib diisi",
		"min":      "{field} minimal {param} karakter",
		"max":      "{field} maksimal {param} karakter",
		"email":    "{field} harus berupa alamat email yang valid",
		"oneof":    "{field} harus salah satu dari: {param}",
		"numeric":  "{field} harus berupa angka",
		"alpha":    "{field} hanya boleh berisi huruf",
		"alphanum": "{field} hanya boleh berisi huruf dan angka",
		"url":      "{field} harus berupa URL yang valid",
		"uuid":     "{field} harus berupa UUID yang valid",
		"phone":    "{field} harus berupa nomor telepon yang valid",
		"nik":      "{field} harus berupa NIK yang valid",
		"phone_id": "{field} harus berupa nomor telepon Indonesia yang valid",
	}
}
