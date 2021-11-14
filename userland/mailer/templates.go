package mailer

const emailVerificationTemplate = `
	Hi Userlanders,<br/>
	<br/>
	Please verify your Email by clicking <a href="%s">here</a><br/>
	<br/>
	Cheers,<br/>
	Your Userland Team
`

const passwordResetTemplate = `
	Hi Userlanders,<br/>
	<br/>
	Here is your token to reset your password:
	<p style="font-size: 18px; font-weight: 600;">%s</p><br/>
	<br/>
	If you don't request a password reset, please ignore this email.<br/>
	<br/>
	Cheers,<br/>
	Your Userland Team
`
