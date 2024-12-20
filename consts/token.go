package consts

const (
	TOKEN_EXPIRE_DURATION = 7 * 24 * 60 * 60

	TOKEN_SECRET = "SAMPLE_BLOG_BACKEND_SECRET"

	TOKEN_ISSUER = "org.mehak.blog"

	MAX_TOKENS_PER_USER = 5

	REDIS_AVAILABLE_USER_TOKEN_LIST = "USER:TOKENS"
)
