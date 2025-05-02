package emailsender

import "fmt"


func GenerateVerificationHTMLBody(code string, message string) string {
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>One-Time Verification Code</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        h1 {
            color: #007bff;
        }
        p {
            color: #555;
        }
        .code {
            font-size: 24px;
            font-weight: bold;
            color: #28a745;
            background-color: #f0fdf4;
            padding: 10px;
            border-radius: 5px;
            margin: 10px 0;
            display: inline-block;
        }
        .expiry {
            font-size: 12px;
            color: #6c757d;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Medscribe: One-Time Verification Code</h1>
        <p>%s</p>
        <p class="code">%s</p>
        <p class="expiry">This code will expire in 120 seconds.</p>
    </div>
</body>
</html>`, message, code)
	return htmlBody
}

func GeneratePasswordResetHTMLBody(resetLink string) string {
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Reset Request</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        h1 {
            color: #007bff;
        }
        p {
            color: #555;
        }
        .link-container {
            margin: 20px 0;
            padding: 15px;
            background-color: #e9ecef;
            border-radius: 5px;
        }
        .reset-link {
            display: inline-block;
            padding: 10px 20px;
            background-color: #28a745;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            font-weight: bold;
        }
        .instructions {
            font-size: 14px;
            color: #6c757d;
            margin-top: 20px;
        }
        .expiry {
            font-size: 12px;
            color: #6c757d;
            margin-top: 10px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Medscribe: Password Reset Request</h1>
        <p>You have requested to reset your password for your Medscribe account.</p>
        <div class="link-container">
            <a class="reset-link" href="%s">Click here to reset your password</a>
        </div>
        <p class="instructions">If the button above doesn't work, you can also copy and paste the following link into your web browser:</p>
        <p><a href="%s">%s</a></p>
        <p class="instructions">If you did not request a password reset, please ignore this email. Your password will remain unchanged. This password reset link will expire shortly for security reasons.</p>
        <p class="expiry">This link will expire within the next 15 minutes.</p>
    </div>
</body>
</html>`, resetLink, resetLink, resetLink)
	return htmlBody
}