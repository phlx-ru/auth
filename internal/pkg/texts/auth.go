package texts

const (
	UsernameNotFound       = `пользователь с почтой или телефоном '{{ .username }}' не найден`
	WrongPassword          = `неправильный пароль`
	WrongCode              = `неправильный код`
	WrongResetHash         = `неправильный код сброса пароля`
	WrongOldPassword       = `неправильно указан старый пароль`
	LoginTooOften          = `слишком много попыток ввода неправильного пароля`
	LoginByCodeTooOften    = `слишком много попыток ввода неправильного кода`
	ResetPasswordTooOften  = `сброс пароля запрашивается слишком часто`
	NewPasswordTooOften    = `попытка установить новый пароль происходит слишком часто`
	ChangePasswordTooOften = `попытка сменить пароль происходит слишком часто`
	GenerateCodeTooOften   = `код уже был запрошен, следующая попытка возможна в течение минуты`

	AuthCodeEmailSubject = `🗝️ Ваш временный код авторизации`
	AuthCodeEmailBody    = `{{ .code }} — Ваш временный код авторизации, действительный в течение {{ .minutes }} минут.

Сообщение сформировано автоматически. Пожалуйста, не отвечайте на него.
Если Вы не запрашивали код, то просто проигнорируйте это письмо.`

	AuthCodeTelegramBody = `{{ .code }} — Ваш временный код авторизации, действительный в течение {{ .minutes }} минут`

	ResetPasswordEmailSubject = `🗝️ Ваш код для сброса пароля`
	ResetPasswordEmailBody    = `{{ .hash }} — Ваш код для сброса пароля.
Используйте его на странице, которая появится позже // TODO сделать страницу
или перейдите по ссылке, которая появится позже // TODO сделать ссылку

Сообщение сформировано автоматически. Пожалуйста, не отвечайте на него.
Если Вы не запрашивали код, то просто проигнорируйте это письмо.`
	// TODO ^

	ResetPasswordTelegramBody = `{{ .hash }} — Ваш код для сброса пароля.
Для создания нового пароля перейдите по ссылке, которая появится позже // TODO сделать ссылку`
	// TODO ^
)