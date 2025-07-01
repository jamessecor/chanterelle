Chanterelle is a Vermont based band. This is their website. It is currently under construction.

### Login and Authentication

This site uses a two-factor verification system fot logins. If a user enters a valid phone number, the api will send a verification code to that number. The user must then enter the code into the form where it is checked against the existing code for that phone number (5 minute expiration). If the code is a match, the user receives a jwt. All admin routes are protected via the jwt.

&copy; James Secor 2025
