# Multiline Comments Test

This test case demonstrates how the protoc-gen-connect-go-mcp plugin handles multiline comments and special characters in proto files.

## Purpose

This test verifies that the code generator correctly escapes:
- Newline characters (`\n`)
- Double quotes (`"`)
- Backslashes (`\`)
- Tab characters (`\t`)
- Carriage returns (`\r`)

## Test Coverage

The `test.proto` file includes:
1. **CRUD operations** - Complete Create, Read, Update, Delete operations for user management
2. **Multiline method comments** - Comments spanning multiple lines that become tool names
3. **Multiline message comments** - Comments that become tool descriptions
4. **Special characters** - Quotes and backslashes in comments
5. **Warning messages** - Comments with special formatting (e.g., "WARNING:")

## Expected Behavior

When running `buf generate`, the plugin should:
1. Parse all comments from the proto file
2. Properly escape special characters for Go string literals
3. Generate valid Go code without syntax errors

## Running the Test

```bash
# Generate code
buf generate

# Verify syntax
gofmt -e gen/testv1mcp/test.mcpserver.go

# Build (requires dependencies)
go build ./gen/testv1mcp/...
```

## Example Output

The generated code should contain properly escaped strings like:

### CRUD Operations
```go
// Create
Name: "CreateUser - Creates a new user account in the system.\nNOTE: This method requires valid email verification before account activation.",

// Read
Name: "GetUser - Retrieves user information by user ID.\n\nPrerequisites:\n- Valid user ID must be provided\n- User must exist in the system",

// Update
Name: "UpdateUser - Updates an existing user's information.\n\nThis operation allows updating user profile data.\nOnly the fields provided in the request will be updated.",

// Delete
Name: "DeleteUser - Deletes a user account from the system.\n\nWARNING: This operation is irreversible!\nAll user data will be permanently deleted.",
```

### Special Characters
```go
Name: "ProcessPayment - Processes a payment transaction.\n\nThis is a test RPC for demonstrating multiline comments.\nIt includes special characters like \"quotes\" and \\backslashes\\.\nMultiple lines should be properly escaped in the generated code.",
```
