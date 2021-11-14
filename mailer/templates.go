package mailer

const emailVerificationTemplate = `
	Hi Userlanders,
	<br/>
	Please verify your Email by clicking <a href="%s">here</a>
	<br/>
	Cheers,<br/>
	Your Userland Team
`

const passwordResetTemplate = `
	Hi Userlanders,
	<br/>
	Here is your token to reset your password:
	<p style="font-size: 18px; font-weight: 600;">%s</p>
	<br/>
	If you don't request a password reset, please ignore this email.
	<br/>
	Cheers,<br/>
	Your Userland Team
`
