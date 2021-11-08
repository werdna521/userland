package redis

const (
	userKey                    = "user"
	verificationKey            = "verification"
	forgotPasswordKey          = "forgotPassword"
	tokenKey                   = "token"
	emailChangeVerificationKey = "emailChangeVerification"

	hEmailChangeNewEmailKey = "email"
	hEmailChangeToken       = "token"
)

const (
	sessionKey      = "session"
	accessTokenKey  = "accesstoken"
	refreshTokenKey = "refreshtoken"

	hSessionClientKey    = "client"
	hSessionCreatedAtKey = "created_at"
	hSessionUpdatedAtKey = "updated_at"
)
