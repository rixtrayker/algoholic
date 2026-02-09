# Algoholic API - Postman Collection

This directory contains a comprehensive Postman collection for testing the Algoholic API, along with Newman scripts for automated testing.

## Contents

- `algoholic-api.postman_collection.json` - Complete API collection with 30+ endpoints
- `algoholic-local.postman_environment.json` - Local environment configuration
- `run-tests.sh` - Automated test runner script using Newman
- `reports/` - Test execution reports (generated after running tests)

## Features

### Collection Highlights

- **Complete API Coverage**: All 36 endpoints across 7 categories
- **Automated Tests**: Each request includes test scripts that validate:
  - HTTP status codes
  - Response structure
  - Required fields
  - Data types
  - Business logic
- **Variable Management**: Automatically extracts and stores tokens, IDs, and other values
- **Example Responses**: Sample responses for key endpoints
- **Authentication Flow**: Automatic token management after login/registration

### Endpoint Categories

1. **Authentication** (3 endpoints)
   - Register User
   - Login
   - Get Current User

2. **Questions** (4 endpoints)
   - Get Random Question
   - Get Random Question with Filters
   - Submit Answer
   - Request Hint

3. **Problems** (4 endpoints)
   - List Problems
   - Get Problem by ID
   - Get Problem Questions
   - Search Problems

4. **User Stats & Progress** (3 endpoints)
   - Get User Stats
   - Get Progress History
   - Get Attempt History

5. **Training Plans** (5 endpoints)
   - List Available Plans
   - Get Plan Details
   - Enroll in Plan
   - Get My Plans
   - Update Plan Progress

6. **Topics** (2 endpoints)
   - List All Topics
   - Get Topic Performance

7. **Health Check** (1 endpoint)
   - API Health

## Prerequisites

Before running the tests, ensure you have:

1. **Backend Running**: The Go backend must be running on `http://localhost:4000`
   ```bash
   cd ../backend
   go run main.go
   ```

2. **Database Setup**: PostgreSQL database must be configured and running

3. **Newman Installed**: Newman CLI is required for automated testing
   ```bash
   npm install -g newman
   ```

## Usage

### Option 1: Using Postman GUI

1. **Import Collection**:
   - Open Postman
   - Click "Import" button
   - Select `algoholic-api.postman_collection.json`

2. **Import Environment**:
   - Click "Environments" in the sidebar
   - Click "Import"
   - Select `algoholic-local.postman_environment.json`

3. **Select Environment**:
   - Choose "Algoholic - Local" from the environment dropdown

4. **Run Requests**:
   - Start with "Authentication > Register User" to create an account
   - The auth token will be automatically saved
   - Continue with other requests

### Option 2: Using Newman (CLI)

#### Basic Test Run

```bash
cd postman
newman run algoholic-api.postman_collection.json \
  --environment algoholic-local.postman_environment.json
```

#### Using the Test Script

```bash
cd postman
./run-tests.sh
```

The script will:
- Check if the backend is running
- Execute all tests in sequence
- Generate detailed reports
- Exit with appropriate status code

#### Advanced Newman Options

**Run specific folder:**
```bash
newman run algoholic-api.postman_collection.json \
  --environment algoholic-local.postman_environment.json \
  --folder "Authentication"
```

**Generate HTML report:**
```bash
newman run algoholic-api.postman_collection.json \
  --environment algoholic-local.postman_environment.json \
  --reporters cli,htmlextra \
  --reporter-htmlextra-export reports/report.html
```

**Run with custom iterations:**
```bash
newman run algoholic-api.postman_collection.json \
  --environment algoholic-local.postman_environment.json \
  --iteration-count 5
```

**Delay between requests:**
```bash
newman run algoholic-api.postman_collection.json \
  --environment algoholic-local.postman_environment.json \
  --delay-request 1000
```

## Test Coverage

The collection includes comprehensive test assertions for:

### Status Code Validation
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
```

### Response Structure Validation
```javascript
pm.test("Response has user and token", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('user');
    pm.expect(jsonData).to.have.property('token');
});
```

### Data Type Validation
```javascript
pm.test("Points earned is a number", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.points_earned).to.be.a('number');
});
```

### Business Logic Validation
```javascript
pm.test("Question matches filters", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.difficulty_score).to.be.within(40, 60);
});
```

## Environment Variables

The collection uses the following variables:

| Variable | Description | Auto-populated? |
|----------|-------------|-----------------|
| `base_url` | API base URL | ❌ Manual |
| `auth_token` | JWT authentication token | ✅ After login/register |
| `user_id` | Current user ID | ✅ After login/register |
| `question_id` | Last fetched question ID | ✅ After getting question |
| `problem_id` | Last fetched problem ID | ✅ After listing problems |
| `plan_id` | Last fetched plan ID | ✅ After listing plans |

## Workflow Example

Here's a typical test workflow:

1. **Register** → Creates user, saves `auth_token` and `user_id`
2. **Get Random Question** → Fetches question, saves `question_id`
3. **Submit Answer** → Uses saved `question_id` to submit
4. **Get User Stats** → Verifies stats updated after answer
5. **List Problems** → Fetches problems, saves first `problem_id`
6. **Get Problem by ID** → Uses saved `problem_id`
7. **List Training Plans** → Fetches plans, saves first `plan_id`
8. **Enroll in Plan** → Uses saved `plan_id`

## Reports

After running tests with the script, reports are generated in the `reports/` directory:

- `newman-report.json` - Detailed JSON report with all test results
- `newman-report.xml` - JUnit XML format for CI/CD integration

### Interpreting Results

**Successful Run:**
```
┌─────────────────────────┬──────────┬──────────┐
│                         │ executed │   failed │
├─────────────────────────┼──────────┼──────────┤
│              iterations │        1 │        0 │
├─────────────────────────┼──────────┼──────────┤
│                requests │       30 │        0 │
├─────────────────────────┼──────────┼──────────┤
│            test-scripts │       30 │        0 │
├─────────────────────────┼──────────┼──────────┤
│      prerequest-scripts │        0 │        0 │
├─────────────────────────┼──────────┼──────────┤
│              assertions │       85 │        0 │
└─────────────────────────┴──────────┴──────────┘
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: API Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_DB: algoholic
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'

      - name: Install dependencies
        run: cd backend && go mod download

      - name: Start backend
        run: cd backend && go run main.go &

      - name: Wait for backend
        run: |
          for i in {1..30}; do
            curl -s http://localhost:4000/health && break
            sleep 1
          done

      - name: Install Newman
        run: npm install -g newman

      - name: Run API tests
        run: cd postman && ./run-tests.sh
```

## Troubleshooting

### Backend Not Running

**Error:**
```
✗ Backend is not running!
Please start the backend with: cd backend && go run main.go
```

**Solution:**
```bash
cd ../backend
go run main.go
```

### Authentication Failures

If you get 401 errors:
1. Run the "Register User" or "Login" request first
2. Check that `auth_token` is populated in environment variables
3. Verify the token hasn't expired (default: 72 hours)

### Database Connection Errors

If endpoints return database errors:
1. Verify PostgreSQL is running
2. Check database credentials in `backend/config/config.go`
3. Run database migrations if needed

### Port Conflicts

If backend can't start on port 4000:
1. Change `base_url` in environment to use different port
2. Update backend configuration to use different port
3. Kill process using port 4000: `lsof -ti:4000 | xargs kill`

## Contributing

When adding new endpoints to the API:

1. Add the request to appropriate folder in collection
2. Include test scripts following existing patterns
3. Add example responses
4. Update this README with new endpoint details
5. Test with Newman before committing

## Best Practices

1. **Always run Health Check first** to verify backend is responsive
2. **Start with Authentication** to get a valid token
3. **Use Collection Runner** in Postman to run all tests in sequence
4. **Review test results** in reports for failed assertions
5. **Keep environment variables updated** for different environments (dev, staging, prod)

## Resources

- [Newman Documentation](https://learning.postman.com/docs/running-collections/using-newman-cli/command-line-integration-with-newman/)
- [Postman Test Scripts](https://learning.postman.com/docs/writing-scripts/test-scripts/)
- [Postman Variables](https://learning.postman.com/docs/sending-requests/variables/)
- [Algoholic API Documentation](../docs/api-reference.md)
