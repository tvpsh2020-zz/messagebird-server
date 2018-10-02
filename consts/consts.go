package consts

// API response result
const (
	APISuccess = "success"
	APIFail    = "fail"
)

// Constants for SMS
const (
	SMSType             = "binary"
	Plain               = "plain"
	PlainSMSLength      = 160
	PlainCSMSLength     = 153
	PlainSMSMaxLength   = 1377
	Unicode             = "unicode"
	UnicodeSMSLength    = 70
	UnicodeCSMSLength   = 67
	UnicodeSMSMaxLength = 603
)

// Constants for regular expression pattern
const (
	GSM0338Regex    = `[^A-Za-z0-9 \\r\\n@£$¥èéùìòÇØøÅåΔ_ΦΓΛ¤ΩΠΨΣΘΞÆæßÉ!\"#$%&amp;'()*+,./:;&lt;=&gt;?¡ÄÖÑÜ§¿äöñüà^{}\[\~\|\]\<\>\-\\]*`
	OriginatorRegex = `^[a-zA-Z0-9]{1,11}$`
	RecipientsRegex = `[0-9]*\,*`
)
