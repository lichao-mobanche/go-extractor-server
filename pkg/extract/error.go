package extract

// ContentEmptyError TODO
type ContentEmptyError string

func (e ContentEmptyError) Error() string {
	return fmt.Sprintf("content empty %s", string(e))
}

// ContentTypeError TODO
type ContentTypeError string

func (e ContentTypeError) Error() string {
	return fmt.Sprintf("content type %s illegal", string(e))
}

// DetectorError TODO
type DetectorError string

func (e DetectorError) Error() string {
	return fmt.Sprintf("detector failed %s", string(e))
}

// EncodeError TODO
type EncodeError string

func (e EncodeError) Error() string {
	return fmt.Sprintf("encode failed %s", string(e))
}

// Base64Error TODO
type Base64Error string

func (e Base64Error) Error() string {
	return fmt.Sprintf("base64 decode failed %s", string(e))
}

// DocError TODO
type DocError string

func (e DocError) Error() string {
	return fmt.Sprintf("doc load failed %s", string(e))
} 