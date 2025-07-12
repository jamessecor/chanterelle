Chanterelle is a Vermont based band. This is their website. It is currently under construction.

### Login and Authentication

This site uses a two-factor verification system fot logins. If a user enters a valid phone number, the api will send a verification code to that number. The user must then enter the code into the form where it is checked against the existing code for that phone number (5 minute expiration). If the code is a match, the user receives a jwt. All admin routes are protected via the jwt.

&copy; James Secor 2025

## Testing

### Running Mailchimp Integration Tests

To run the Mailchimp integration tests:

1. Set up your Mailchimp credentials in the `.env` file:
```bash
MAILCHIMP_API_KEY=your_api_key
MAILCHIMP_LIST_ID=your_list_id
```

2. Run the tests with:
```bash
MAILCHIMP_TEST=true go test ./internal/services/... -v
```

Note: The integration tests will create test contacts in your Mailchimp list using unique email addresses (test+timestamp@example.com). These contacts will remain in your list after the tests complete.

### Running All Tests

To run all tests (excluding Mailchimp integration tests):
```bash
go test ./... -v
```
