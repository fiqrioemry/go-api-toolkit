# 🚀 Go API Toolkit Documentation

A modular, production-ready toolkit for building consistent REST APIs in Go with automatic response formatting, error handling, validation, pagination, and logging.

## 📦 Features

- ✅ **Consistent Response Format** - Standardized JSON responses across all endpoints
- ✅ **Smart Validation** - Comprehensive validation with Indonesian localization support
- ✅ **Smart Pagination** - Automatic pagination with bulletproof defaults
- ✅ **Structured Error Handling** - Built-in error types with proper HTTP status codes
- ✅ **Automatic Logging** - Optional request/response logging with Zap integration
- ✅ **Framework Agnostic Core** - Supports Gin & Native HTTP (more frameworks coming soon)
- ✅ **Zero Breaking Changes** - Backward compatible with existing code
- ✅ **Modular Architecture** - Use only what you need

## 🏗️ Architecture

```
go-api-toolkit/
├── response/           # HTTP response handling & error management
│   ├── types.go       # Core types and structures
│   ├── handler.go     # Framework-agnostic core engine
│   ├── errors.go      # Error constructors and utilities
│   ├── logger.go      # Logger interface and implementations
│   ├── zap_adapter.go # Zap logger integration
│   ├── gin.go         # Gin framework integration
│   └── http.go        # Native HTTP integration
├── validation/         # Validation logic and utilities
│   ├── types.go       # Core validation types and structures
│   ├── handler.go     # Framework-agnostic validation engine
│   ├── rules.go       # Built-in validation rules
│   ├── errors.go      # Validation error handling
│   ├── locales.go     # Localization support (EN/ID)
│   ├── gin.go         # Gin framework integration
│   └── http.go        # Native HTTP integration
├── pagination/         # Pagination logic and utilities
│   ├── types.go       # Pagination types and structures
│   ├── builder.go     # Pagination building logic
│   └── gin.go         # Gin framework pagination helpers
└── go.mod
```

## 📋 Installation

```bash
go get github.com/fiqrioemry/go-api-toolkit
```

## ⚙️ Configuration

### Basic Setup (Gin Framework)

```go
// main.go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/fiqrioemry/go-api-toolkit/response"
    "github.com/fiqrioemry/go-api-toolkit/validation"
    "github.com/fiqrioemry/go-api-toolkit/pagination"
    "your-app/utils" // Your existing logger
)

func main() {
    // Initialize your existing logger
    utils.InitLogger()

    // 🚀 ONE LINE SETUP - Initialize response handler
    response.InitGin(response.InitConfig{
        Logger:              utils.GetLogger(),
        LogSuccessResponses: false,
        LogErrorResponses:   true,
    })

    // 🚀 ONE LINE SETUP - Initialize validation
    validation.InitGin(validation.InitConfig{
        Logger:           utils.GetLogger(),
        CustomMessages:   true,
        Locale:           "en", // or "id" for Indonesian
        CustomRules:      validation.GetBuiltInRules(),
    })

    // Pagination module doesn't need initialization - it's standalone!

    r := gin.Default()
    // ... rest of your application
}
```

### Basic Setup (Native HTTP)

```go
// main.go
package main

import (
    "net/http"
    "github.com/fiqrioemry/go-api-toolkit/response"
    "github.com/fiqrioemry/go-api-toolkit/validation"
    "your-app/utils"
)

func main() {
    utils.InitLogger()

    // 🚀 Initialize for native HTTP
    response.InitHTTP(response.InitConfig{
        Logger:            utils.GetLogger(),
        LogErrorResponses: true,
    })

    validation.InitHTTP(validation.InitConfig{
        Logger:         utils.GetLogger(),
        CustomMessages: true,
        Locale:         "id", // Indonesian support
        CustomRules:    validation.GetBuiltInRules(),
    })

    http.HandleFunc("/users", createUserHandler)
    http.ListenAndServe(":8080", nil)
}
```

### Advanced Configuration

```go
// Custom configuration with Indonesian localization
validation.InitGin(validation.InitConfig{
    Logger:           customZapLogger,
    CustomMessages:   true,
    StopOnFirstError: false, // Validate all fields
    Locale:           "id",  // Indonesian messages
    CustomRules: map[string]validation.Rule{
        "nik":      validation.NewNIKRule(),
        "phone_id": validation.NewIndonesianPhoneRule(),
        "custom":   &MyCustomRule{},
    },
    ErrorMessages: map[string]string{
        "required": "Field ini wajib diisi",
        "email":    "Format email tidak valid",
    },
})
```

## 🔄 Migration Guide

### Before (Manual Validation & Response)

```go
// ❌ BEFORE: Manual, repetitive, error-prone

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "Invalid JSON format",
            "error":   err.Error(),
        })
        return
    }

    // Manual validation
    if req.Name == "" {
        c.JSON(400, gin.H{"success": false, "message": "Name is required"})
        return
    }
    if len(req.Name) < 2 {
        c.JSON(400, gin.H{"success": false, "message": "Name must be at least 2 characters"})
        return
    }
    if req.Email == "" {
        c.JSON(400, gin.H{"success": false, "message": "Email is required"})
        return
    }
    if !strings.Contains(req.Email, "@") {
        c.JSON(400, gin.H{"success": false, "message": "Invalid email format"})
        return
    }
    if req.Age < 18 {
        c.JSON(400, gin.H{"success": false, "message": "Age must be at least 18"})
        return
    }
    // ... 20+ lines of validation

    user, err := h.service.CreateUser(req)
    if err != nil {
        if strings.Contains(err.Error(), "duplicate") {
            c.JSON(409, gin.H{"success": false, "message": "Email already exists"})
        } else {
            c.JSON(500, gin.H{"success": false, "message": "Internal server error"})
        }
        return
    }

    // Manual response formatting
    c.JSON(201, gin.H{
        "success": true,
        "message": "User created successfully",
        "data":    user,
    })
}
```

### After (With Go API Toolkit)

```go
// ✅ AFTER: Clean, consistent, automatic

import (
    "github.com/fiqrioemry/go-api-toolkit/response"
    "github.com/fiqrioemry/go-api-toolkit/validation"
    "github.com/fiqrioemry/go-api-toolkit/pagination"
)

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest

    // 🚀 ONE LINER: bind + validate all rules
    if err := validation.BindAndValidate(c, &req); err != nil {
        response.Error(c, err) // Automatic error handling + logging
        return
    }

    user, err := h.service.CreateUser(req)
    if err != nil {
        response.Error(c, err) // Automatic error handling + logging
        return
    }

    response.Created(c, "User created successfully", user)
}
```

### DTO with Validation Tags

```go
// ✅ Comprehensive validation with single tags
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50,alpha" message:"Name is required and must be 2-50 letters only"`
    Email    string `json:"email" validate:"required,email" message:"Please provide a valid email address"`
    Age      int    `json:"age" validate:"required,min=18,max=120" message:"Age must be between 18 and 120"`
    Phone    string `json:"phone" validate:"required,phone_id" message:"Please provide a valid Indonesian phone number"`
    NIK      string `json:"nik" validate:"omitempty,nik" message:"NIK must be a valid Indonesian identity number"`
    Password string `json:"password" validate:"required,min=8,password=upper,lower,number" message:"Password must be at least 8 characters with uppercase, lowercase, and number"`
    Website  string `json:"website" validate:"omitempty,url" message:"Website must be a valid URL"`
    UserType string `json:"userType" validate:"required,oneof=admin user guest" message:"User type must be admin, user, or guest"`
}

// Remove old binding constraints - let validation package handle it
type GetUsersRequest struct {
    Page      int    `form:"page" json:"page" validate:"omitempty,min=1"`              // Let validation handle defaults
    Limit     int    `form:"limit" json:"limit" validate:"omitempty,min=1,max=100"`   // Let validation handle defaults
    SortBy    string `form:"sortBy" json:"sortBy" validate:"omitempty,oneof=name email age createdAt"`
    SortOrder string `form:"sortOrder" json:"sortOrder" validate:"omitempty,oneof=asc desc"`
    MinAge    *int   `form:"minAge" json:"minAge" validate:"omitempty,min=0,max=120"`
    MaxAge    *int   `form:"maxAge" json:"maxAge" validate:"omitempty,min=0,max=120"`
    Search    string `form:"q" json:"q" validate:"omitempty,min=2,max=100"`
}
```

## 📚 Usage Examples

### 1. Basic Validation & Response (Gin)

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest

    // 🚀 Auto-bind + validate with comprehensive rules
    if err := validation.BindAndValidate(c, &req); err != nil {
        response.Error(c, err)
        return
    }

    user, err := h.service.CreateUser(req)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Created(c, "User created successfully", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
    var req dto.UpdateUserRequest
    userID := c.Param("id")

    // Advanced: validation with context
    if err := validation.BindAndValidate(c, &req,
        validation.WithContext(map[string]interface{}{
            "user_id": userID,
            "action":  "update",
        })); err != nil {
        response.Error(c, err)
        return
    }

    user, err := h.service.UpdateUser(userID, req)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.OK(c, "User updated successfully", user)
}
```

### 2. Basic Validation & Response (Native HTTP)

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateUserRequest

    // 🚀 Same API for native HTTP
    if err := validation.BindAndValidateHTTP(r, &req); err != nil {
        response.ErrorHTTP(w, r, err)
        return
    }

    user, err := userService.CreateUser(req)
    if err != nil {
        response.ErrorHTTP(w, r, err)
        return
    }

    response.CreatedHTTP(w, r, "User created successfully", user)
}
```

### 6. Update Services

```go
// Focus on business logic, not validation
func (s *userService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
    // Structural validation already done in handler

    // Only business validation here
    exists, err := s.repo.EmailExists(req.Email)
    if err != nil {
        return nil, response.DatabaseError("Failed to check email", err)
    }
    if exists {
        return nil, response.Conflict("Email already exists")
    }

    user, err := s.repo.CreateUser(req)
    if err != nil {
        return nil, response.InternalServerError("Failed to create user", err)
    }

    return s.convertToResponse(user), nil
}
```

## 🎉 Benefits Summary

### Before vs After Comparison

| Aspect                 | Before (Manual)               | After (Go API Toolkit)         |
| ---------------------- | ----------------------------- | ------------------------------ |
| **Validation Code**    | 30+ lines per handler         | 1 line with tags               |
| **Error Handling**     | Manual JSON responses         | Automatic structured responses |
| **Localization**       | Not supported                 | Built-in EN/ID support         |
| **Indonesian Context** | Manual implementation         | Built-in NIK, phone validation |
| **Response Format**    | Inconsistent across endpoints | Standardized JSON format       |
| **Error Logging**      | Manual logging                | Automatic structured logging   |
| **Pagination**         | Manual calculation            | One-liner with defaults        |
| **Framework Support**  | Single framework              | Gin + Native HTTP + extensible |
| **Testing**            | Complex setup                 | Simple, focused tests          |
| **Maintenance**        | High (repetitive code)        | Low (centralized logic)        |

### Key Benefits

- ✅ **90% Less Validation Code** - Eliminate repetitive validation logic
- ✅ **Consistent API Responses** - Same format across all endpoints
- ✅ **Indonesian Context Ready** - Built-in NIK, phone, localization
- ✅ **Multi-Framework Support** - Gin, Native HTTP, easily extensible
- ✅ **Automatic Error Logging** - Never miss an error again
- ✅ **Bulletproof Pagination** - Handle all edge cases automatically
- ✅ **Type-Safe Errors** - Structured error handling with proper HTTP codes
- ✅ **Zero Breaking Changes** - Migrate gradually without disruption
- ✅ **Production Ready** - Battle-tested patterns and best practices
- ✅ **Developer Experience** - Focus on business logic, not boilerplate

## 🔮 Advanced Use Cases

### 1. Multi-Language API

```go
// Dynamic locale based on request header
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    locale := c.GetHeader("Accept-Language") // "en" or "id"

    if err := validation.BindAndValidate(c, &req,
        validation.WithLocale(locale)); err != nil {
        response.Error(c, err)
        return
    }

    // Business logic...
}
```

### 2. Role-Based Validation

```go
func (h *AdminHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    userRole := c.GetString("user_role")

    if err := validation.BindAndValidate(c, &req,
        validation.WithContext(map[string]interface{}{
            "user_role": userRole,
            "action":    "create_user",
        })); err != nil {
        response.Error(c, err)
        return
    }

    // Business logic...
}

// Custom validation rule that uses context
type RoleBasedRule struct{}

func (r *RoleBasedRule) Validate(value interface{}, params string, context map[string]interface{}) error {
    userRole := context["user_role"].(string)
    if userRole != "admin" && value.(string) == "admin" {
        return fmt.Errorf("only admins can create admin users")
    }
    return nil
}
```

### 3. File Upload with Validation

```go
type FileUploadRequest struct {
    Title       string `form:"title" validate:"required,min=2,max=100"`
    Description string `form:"description" validate:"omitempty,max=500"`
    Category    string `form:"category" validate:"required,oneof=image document video"`
    IsPublic    bool   `form:"isPublic"`
    Tags        string `form:"tags" validate:"omitempty"`
}

func (h *FileHandler) UploadFile(c *gin.Context) {
    var req FileUploadRequest

    // Force multipart form binding for file uploads
    if err := validation.BindAndValidate(c, &req,
        validation.ForceForm()); err != nil {
        response.Error(c, err)
        return
    }

    // Handle file upload...
}
```

### 4. Complex Business Validation

```go
type TransferRequest struct {
    FromAccount string  `json:"fromAccount" validate:"required,uuid"`
    ToAccount   string  `json:"toAccount" validate:"required,uuid"`
    Amount      float64 `json:"amount" validate:"required,min=0.01"`
    Currency    string  `json:"currency" validate:"required,oneof=IDR USD EUR"`
    Note        string  `json:"note" validate:"omitempty,max=200"`
}

func (s *paymentService) Transfer(req TransferRequest) (*TransferResponse, error) {
    // Structural validation already done in handler

    // Business validation
    if req.FromAccount == req.ToAccount {
        return nil, response.BadRequest("Cannot transfer to the same account")
    }

    balance, err := s.getAccountBalance(req.FromAccount)
    if err != nil {
        return nil, response.DatabaseError("Failed to check balance", err)
    }

    if balance < req.Amount {
        return nil, response.BadRequest("Insufficient balance")
    }

    // Process transfer...
}
```

## 🛠️ Extending the Toolkit

### Adding New Framework Support

The toolkit is designed to be easily extensible. Here's how to add support for a new framework:

#### 1. Create Framework Integration File

```go
// validation/fiber.go
package validation

import "github.com/gofiber/fiber/v2"

type FiberWriter struct {
    ctx *fiber.Ctx
}

func (fw *FiberWriter) WriteJSON(statusCode int, data interface{}) error {
    return fw.ctx.Status(statusCode).JSON(data)
}

func BindAndValidateFiber(ctx *fiber.Ctx, obj interface{}, opts ...ValidationOption) error {
    // Implement Fiber-specific binding logic
    // Similar to BindAndValidate but using Fiber's API
}

func InitFiber(configs ...InitConfig) {
    // Initialize for Fiber framework
}
```

#### 2. Add Response Support

```go
// response/fiber.go
package response

import "github.com/gofiber/fiber/v2"

func ErrorFiber(ctx *fiber.Ctx, err error) {
    // Implement Fiber error handling
}

func OKFiber(ctx *fiber.Ctx, message string, data interface{}) {
    // Implement Fiber success response
}
```

### Creating Domain-Specific Rules

```go
// Custom rule for Indonesian bank account validation
type BankAccountRule struct{}

func (r *BankAccountRule) Validate(value interface{}, params string, context map[string]interface{}) error {
    account, ok := value.(string)
    if !ok {
        return fmt.Errorf("bank account must be a string")
    }

    // Indonesian bank account validation logic
    if len(account) < 10 || len(account) > 16 {
        return fmt.Errorf("bank account must be 10-16 digits")
    }

    // Add specific bank validation based on params
    bank := params // e.g., "bca", "mandiri", "bni"
    if !r.validateBankSpecificFormat(account, bank) {
        return fmt.Errorf("invalid account format for %s", bank)
    }

    return nil
}

func (r *BankAccountRule) validateBankSpecificFormat(account, bank string) bool {
    switch bank {
    case "bca":
        return len(account) >= 10 && strings.HasPrefix(account, "0")
    case "mandiri":
        return len(account) >= 13
    case "bni":
        return len(account) >= 10
    default:
        return true
    }
}

func (r *BankAccountRule) GetMessage() string {
    return "field must be a valid Indonesian bank account"
}

// Usage
type PaymentRequest struct {
    BankAccount string `json:"bankAccount" validate:"required,bank_account=bca"`
}
```

## 📈 Performance Considerations

### Validation Performance

- **Reflection Caching**: Struct reflection is cached for better performance
- **Rule Compilation**: Regex patterns are compiled once during initialization
- **Memory Efficiency**: Minimal allocations during validation
- **Concurrent Safe**: All validators are thread-safe

### Response Performance

- **JSON Marshaling**: Efficient JSON encoding with minimal allocations
- **Logger Buffering**: Structured logging with proper buffering
- **Error Pooling**: Error objects are reused when possible

### Benchmarks

```go
// Example benchmark results (approximate)
BenchmarkValidation-8           1000000    1200 ns/op    240 B/op    3 allocs/op
BenchmarkPagination-8          5000000     280 ns/op     64 B/op     1 allocs/op
BenchmarkResponseJSON-8        2000000     680 ns/op    128 B/op     2 allocs/op
```

## 🔍 Troubleshooting

### Common Issues

#### 1. Validation Not Working

```go
// ❌ Problem: Validation not initialized
func main() {
    r := gin.Default()
    // Missing validation.InitGin()
}

// ✅ Solution: Initialize validation
func main() {
    validation.InitGin(validation.InitConfig{})
    r := gin.Default()
}
```

#### 2. Custom Rules Not Found

```go
// ❌ Problem: Custom rule not registered
type MyRequest struct {
    Field string `validate:"mycustom"` // Rule not found
}

// ✅ Solution: Register custom rules
validation.InitGin(validation.InitConfig{
    CustomRules: map[string]validation.Rule{
        "mycustom": &MyCustomRule{},
    },
})
```

#### 3. Localization Not Working

```go
// ❌ Problem: Wrong locale setting
validation.InitGin(validation.InitConfig{
    Locale: "indonesia", // Should be "id"
})

// ✅ Solution: Use correct locale codes
validation.InitGin(validation.InitConfig{
    Locale: "id", // or "en"
})
```

#### 4. Pagination Defaults Not Applied

```go
// ❌ Problem: Missing pagination defaults
func GetUsers(c *gin.Context) {
    var req GetUsersRequest
    validation.BindAndValidate(c, &req)
    // req.Page might be 0
}

// ✅ Solution: Apply pagination defaults
func GetUsers(c *gin.Context) {
    var req GetUsersRequest
    validation.BindAndValidate(c, &req)
    pagination.ApplyDefaultsToStruct(&req) // Apply defaults
}
```

### Debug Mode

```go
// Enable debug logging to see what's happening
validation.InitGin(validation.InitConfig{
    Logger: logger, // Make sure logger is not nil
})

// Check logs for validation details
```

## 🤝 Contributing

We welcome contributions! Here's how you can help:

### Areas for Contribution

1. **New Framework Support** - Add Fiber, Echo, Chi integrations
2. **Additional Validation Rules** - Country-specific or domain-specific rules
3. **Localization** - Add more language support
4. **Performance Improvements** - Optimize validation and response handling
5. **Documentation** - Improve examples and guides

### Getting Started

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-framework`
3. Make your changes with tests
4. Ensure all tests pass: `go test ./...`
5. Submit a pull request

### Code Guidelines

- Follow Go conventions and best practices
- Add comprehensive tests for new features
- Update documentation for any API changes
- Maintain backward compatibility
- Use structured logging for debugging

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built for the Indonesian Go developer community
- Inspired by best practices from Laravel, Express.js, and other frameworks
- Special thanks to contributors and early adopters

---

**Made with ❤️ for the Go community in Indonesia and beyond**

For more examples, advanced usage, and community discussions, visit our [GitHub repository](https://github.com/fiqrioemry/go-api-toolkit)., err)
return
}

    response.CreatedHTTP(w, r, "User created successfully", user)

}

func getUsersHandler(w http.ResponseWriter, r \*http.Request) {
var req dto.GetUsersRequest

    if err := validation.BindAndValidateHTTP(r, &req); err != nil {
        response.ErrorHTTP(w, r, err)
        return
    }

    // Apply pagination defaults
    pagination.ApplyDefaultsToStruct(&req)

    users, total, err := userService.GetUsers(req)
    if err != nil {
        response.ErrorHTTP(w, r, err)
        return
    }

    pag := pagination.Build(req.Page, req.Limit, total)
    response.OKWithPaginationHTTP(w, r, "Users retrieved successfully", users, pag)

}

````

### 3. Pagination with Validation

```go
func (h *ProductHandler) GetProducts(c *gin.Context) {
    var req dto.GetProductsRequest

    // 1. Bind and validate (includes pagination fields)
    if err := validation.BindAndValidate(c, &req); err != nil {
        response.Error(c, err)
        return
    }

    // 2. Apply pagination defaults
    pagination.ApplyDefaultsToStruct(&req)

    // 3. Business logic
    products, total, err := h.service.GetProducts(req)
    if err != nil {
        response.Error(c, err)
        return
    }

    // 4. Response with pagination
    pag := pagination.Build(req.Page, req.Limit, total)
    response.OKWithPagination(c, "Products retrieved successfully", products, pag)
}
````

### 4. Advanced Validation Options

```go
func (h *FileHandler) UploadFile(c *gin.Context) {
    var req dto.FileUploadRequest

    // Force form-data binding for file uploads
    if err := validation.BindAndValidate(c, &req,
        validation.ForceForm()); err != nil {
        response.Error(c, err)
        return
    }

    file, err := h.service.UploadFile(req)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Created(c, "File uploaded successfully", file)
}

func (h *APIHandler) ProcessData(c *gin.Context) {
    var req dto.ProcessRequest

    // Multiple validation options
    if err := validation.BindAndValidate(c, &req,
        validation.WithLocale("id"),
        validation.WithContext(map[string]interface{}{
            "user_role": c.GetString("user_role"),
        }),
        validation.WithCustomRules(map[string]validation.Rule{
            "custom": &MyCustomRule{},
        })); err != nil {
        response.Error(c, err)
        return
    }

    // Business logic...
}
```

### 5. Service Layer with Structured Errors

```go
func (s *userService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
    // Validation is already done in handler, focus on business logic

    // Check business rules
    exists, err := s.repo.EmailExists(req.Email)
    if err != nil {
        return nil, response.DatabaseError("Failed to check email existence", err)
    }
    if exists {
        return nil, response.Conflict("Email already exists")
    }

    // Additional business validation
    if req.Age < 21 && req.UserType == "admin" {
        return nil, response.BadRequest("Admin users must be at least 21 years old")
    }

    user, err := s.repo.CreateUser(req)
    if err != nil {
        return nil, response.InternalServerError("Failed to create user", err)
    }

    return s.convertToResponse(user), nil
}

func (s *userService) GetUserByID(userID string) (*dto.UserResponse, error) {
    user, err := s.repo.GetByID(userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, response.NotFound("User not found")
        }
        return nil, response.DatabaseError("Failed to get user", err)
    }

    return s.convertToResponse(user), nil
}
```

## 🎯 Validation Rules Reference

### Basic Rules

- `required` - Field must be present and not empty
- `min=5` - Minimum length/value (works for strings, numbers, slices)
- `max=100` - Maximum length/value
- `oneof=admin user guest` - Must be one of specified values

### String Rules

- `email` - Valid email format
- `alpha` - Letters only (a-z, A-Z)
- `alphanum` - Letters and numbers only
- `numeric` - Numbers only
- `url` - Valid URL format
- `uuid` - Valid UUID v4 format

### Indonesian-Specific Rules

- `phone_id` - Valid Indonesian phone number (+62, 62, 08xx)
- `nik` - Valid Indonesian NIK (16 digits identity number)

### Password Rules

- `password=upper,lower,number,special` - Password strength requirements
  - `upper` - At least one uppercase letter
  - `lower` - At least one lowercase letter
  - `number` - At least one number
  - `special` - At least one special character

### Usage Examples

```go
type ExampleRequest struct {
    // Basic validation
    Name     string `validate:"required,min=2,max=50,alpha"`
    Email    string `validate:"required,email"`
    Age      int    `validate:"required,min=18,max=120"`

    // Indonesian specific
    Phone    string `validate:"required,phone_id"`
    NIK      string `validate:"omitempty,nik"`

    // Advanced
    Password string `validate:"required,min=8,password=upper,lower,number"`
    Website  string `validate:"omitempty,url"`
    Role     string `validate:"required,oneof=admin user guest"`
    UserID   string `validate:"required,uuid"`
}
```

## 🌍 Localization Support

### English (Default)

```go
validation.InitGin(validation.InitConfig{
    Locale: "en",
})
```

Error response:

```json
{
  "success": false,
  "message": "Validation failed",
  "code": "VALIDATION_ERROR",
  "errors": {
    "name": "Name is required and must be 2-50 characters",
    "email": "Email must be a valid email address",
    "phone": "Phone must be a valid Indonesian phone number"
  }
}
```

### Indonesian

```go
validation.InitGin(validation.InitConfig{
    Locale: "id",
})
```

Error response:

```json
{
  "success": false,
  "message": "Validasi gagal",
  "code": "VALIDATION_ERROR",
  "errors": {
    "nama": "Nama wajib diisi dan minimal 2 karakter",
    "email": "Email harus berupa alamat email yang valid",
    "telepon": "Telepon harus berupa nomor telepon Indonesia yang valid"
  }
}
```

## 📤 Response Examples

### Success Response

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "John Doe",
    "email": "john@example.com",
    "createdAt": "2024-01-15T10:30:00Z"
  }
}
```

### Success Response with Pagination

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [...],
  "meta": {
    "pagination": {
      "page": 1,
      "limit": 10,
      "totalItems": 25,
      "totalPages": 3,
      "offset": 0
    }
  }
}
```

### Validation Error Response

```json
{
  "success": false,
  "message": "Validation failed",
  "code": "VALIDATION_ERROR",
  "errors": {
    "name": "Name is required and must be 2-50 letters only",
    "email": "Please provide a valid email address",
    "age": "Age must be between 18 and 120",
    "phone": "Please provide a valid Indonesian phone number"
  }
}
```

### Business Logic Error Response

```json
{
  "success": false,
  "message": "Email already exists",
  "code": "CONFLICT"
}
```

## 🎯 Available Functions

### Validation Functions

#### Gin Framework

```go
// Basic validation
validation.BindAndValidate(c, &req)

// With options
validation.BindAndValidate(c, &req,
    validation.WithContext(ctx),
    validation.ForceJSON(),
    validation.WithLocale("id"))

// Standalone validation
validation.Validate(&req)
validation.ValidateWithContext(&req, ctx)
```

#### Native HTTP

```go
// HTTP validation
validation.BindAndValidateHTTP(r, &req)

// With options
validation.BindAndValidateHTTP(r, &req,
    validation.WithContext(ctx),
    validation.ForceForm())
```

### Response Functions

#### Gin Framework

```go
// Basic responses
response.OK(c, "Success message", data)
response.Created(c, "Created message", data)
response.Error(c, err)

// Quick error messages
response.BadRequestMsg(c, "Invalid data")
response.NotFoundMsg(c, "Resource not found")
response.UnauthorizedMsg(c, "Access denied")

// Pagination responses
response.OKWithPagination(c, "Success", data, pagination)
response.OKWithPaginationAndPermissions(c, "Success", data, pagination, permissions)
```

#### Native HTTP

```go
// Basic responses
response.OKHTTP(w, r, "Success message", data)
response.CreatedHTTP(w, r, "Created message", data)
response.ErrorHTTP(w, r, err)

// Quick error messages
response.BadRequestMsgHTTP(w, r, "Invalid data")
response.NotFoundMsgHTTP(w, r, "Resource not found")

// Pagination responses
response.OKWithPaginationHTTP(w, r, "Success", data, pagination)
```

### Error Constructors (Framework Agnostic)

```go
// Client errors (4xx)
response.BadRequest("Invalid input")
response.Unauthorized("Access token required")
response.Forbidden("Insufficient permissions")
response.NotFound("Resource not found")
response.Conflict("Email already exists")

// Server errors (5xx)
response.InternalServerError("Something went wrong", err)
response.DatabaseError("Database operation failed", err)
```

### Pagination Functions

```go
// Building pagination
pagination.Build(page, limit, total)
pagination.Quick(params, total)
pagination.QuickFlexible(params, total)

// Smart binding
pagination.SmartBind(c, &params)
pagination.SmartBindFlexible(c, &params)
pagination.BindAndSetDefaults(c, &anyStruct)

// Apply defaults to any struct
pagination.ApplyDefaultsToStruct(&anyStruct)
```

## 🔧 Custom Validation Rules

### Creating Custom Rules

```go
type MyCustomRule struct{}

func (r *MyCustomRule) Validate(value interface{}, params string, context map[string]interface{}) error {
    str, ok := value.(string)
    if !ok {
        return fmt.Errorf("value must be a string")
    }

    // Your custom validation logic
    if !strings.HasPrefix(str, "CUSTOM_") {
        return fmt.Errorf("value must start with CUSTOM_")
    }

    return nil
}

func (r *MyCustomRule) GetMessage() string {
    return "field must start with CUSTOM_"
}

// Register during initialization
validation.InitGin(validation.InitConfig{
    CustomRules: map[string]validation.Rule{
        "mycustom": &MyCustomRule{},
    },
})
```

### Using Custom Rules

```go
type MyDTO struct {
    Code string `json:"code" validate:"required,mycustom"`
}
```

### Built-in Custom Rules

```go
// Get all built-in rules including Indonesian-specific ones
validation.InitGin(validation.InitConfig{
    CustomRules: validation.GetBuiltInRules(),
})

// Or register individually
validation.InitGin(validation.InitConfig{
    CustomRules: map[string]validation.Rule{
        "nik":      validation.NewNIKRule(),
        "phone_id": validation.NewIndonesianPhoneRule(),
        "url":      validation.NewURLRule(),
        "uuid":     validation.NewUUIDRule(),
        "password": validation.NewPasswordRule(),
    },
})
```

## 📊 Logging Output

When error logging is enabled, you'll see structured logs like:

```json
{
  "level": "error",
  "ts": "2024-01-15T10:30:00.000Z",
  "msg": "Validation error occurred",
  "path": "/api/v1/users",
  "method": "POST",
  "client_ip": "192.168.1.1",
  "user_agent": "Mozilla/5.0...",
  "user_id": "user123",
  "trace_id": "trace456",
  "error_code": "VALIDATION_ERROR",
  "error_message": "Validation failed",
  "validation_errors": {
    "name": "Name is required",
    "email": "Invalid email format"
  }
}
```

## 🧪 Testing

### Unit Testing Handler with Validation

```go
func TestCreateUser(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    response.InitGin(response.InitConfig{Logger: nil})
    validation.InitGin(validation.InitConfig{Logger: nil})

    mockService := &MockUserService{}
    handler := NewUserHandler(mockService)

    // Test valid request
    validUser := `{
        "name": "John Doe",
        "email": "john@example.com",
        "age": 25,
        "phone": "+6281234567890"
    }`

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/users", strings.NewReader(validUser))
    c.Request.Header.Set("Content-Type", "application/json")

    handler.CreateUser(c)

    assert.Equal(t, 201, w.Code)

    // Test invalid request
    invalidUser := `{
        "name": "A",
        "email": "invalid-email",
        "age": 17
    }`

    w2 := httptest.NewRecorder()
    c2, _ := gin.CreateTestContext(w2)
    c2.Request = httptest.NewRequest("POST", "/users", strings.NewReader(invalidUser))
    c2.Request.Header.Set("Content-Type", "application/json")

    handler.CreateUser(c2)

    assert.Equal(t, 400, w2.Code)

    var response map[string]interface{}
    json.Unmarshal(w2.Body.Bytes(), &response)
    assert.Equal(t, "VALIDATION_ERROR", response["code"])
    assert.NotNil(t, response["errors"])
}
```

### Testing Validation Rules

```go
func TestNIKValidation(t *testing.T) {
    rule := validation.NewNIKRule()

    // Valid NIK
    err := rule.Validate("1234567890123456", "", nil)
    assert.NoError(t, err)

    // Invalid NIK (wrong length)
    err = rule.Validate("123456789012345", "", nil)
    assert.Error(t, err)

    // Invalid NIK (non-numeric)
    err = rule.Validate("123456789012345a", "", nil)
    assert.Error(t, err)
}
```

## 🚀 Migration Steps

### 1. Install Package

```bash
go get github.com/fiqrioemry/go-api-toolkit
```

### 2. Initialize Modules

#### For Gin Framework

```go
// main.go
response.InitGin(response.InitConfig{
    Logger:            utils.GetLogger(),
    LogErrorResponses: true,
})

validation.InitGin(validation.InitConfig{
    Logger:         utils.GetLogger(),
    CustomMessages: true,
    Locale:         "id", // or "en"
    CustomRules:    validation.GetBuiltInRules(),
})
```

#### For Native HTTP

```go
// main.go
response.InitHTTP(response.InitConfig{
    Logger:            utils.GetLogger(),
    LogErrorResponses: true,
})

validation.InitHTTP(validation.InitConfig{
    Logger:         utils.GetLogger(),
    CustomMessages: true,
    Locale:         "id",
    CustomRules:    validation.GetBuiltInRules(),
})
```

### 3. Update Imports

```go
import (
    "github.com/fiqrioemry/go-api-toolkit/response"
    "github.com/fiqrioemry/go-api-toolkit/validation"
    "github.com/fiqrioemry/go-api-toolkit/pagination"
)
```

### 4. Update DTOs with Validation Tags

```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50,alpha"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"required,min=18,max=120"`
    Phone    string `json:"phone" validate:"required,phone_id"`
    Password string `json:"password" validate:"required,min=8,password=upper,lower,number"`
}
```

### 5. Update Handlers

#### Gin Framework

```go
// Before: 30+ lines of manual validation
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    // ... 30+ lines of manual validation and error handling
}

// After: 3 lines + business logic
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest

    if err := validation.BindAndValidate(c, &req); err != nil {
        response.Error(c, err)
        return
    }

    // Business logic immediately
    user, err := h.service.CreateUser(req)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Created(c, "User created successfully", user)
}
```

#### Native HTTP

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateUserRequest

    if err := validation.BindAndValidateHTTP(r, &req); err != nil {
        response.ErrorHTTP(w, r, err)
        return
    }

    user, err := userService.CreateUser(req)
    if err != nil {
        response.ErrorHTTP(w, r, err)
        return
    }

    response.CreatedHTTP(w, r, "User created successfully", user)
}
```

### 6. Update Services

```go
// Focus on business logic, not validation
func (s *userService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
    // Structural validation already done in handler

    // Only business validation here
    exists, err := s.repo.EmailExists(req.Email)
    if err != nil {
        return nil, response.DatabaseError("Failed to check email", err)
    }
    if exists {
        return nil, response.Conflict("Email already exists")
    }

    user, err := s.repo.CreateUser(req)
    if err != nil {
        return nil, response.InternalServerError("Failed to create user", err)
    }

    return s.convertToResponse(user), nil
}
```

## 🎉 Benefits Summary

### Before vs After Comparison

| Aspect                 | Before (Manual)               | After (Go API Toolkit)         |
| ---------------------- | ----------------------------- | ------------------------------ |
| **Validation Code**    | 30+ lines per handler         | 1 line with tags               |
| **Error Handling**     | Manual JSON responses         | Automatic structured responses |
| **Localization**       | Not supported                 | Built-in EN/ID support         |
| **Indonesian Context** | Manual implementation         | Built-in NIK, phone validation |
| **Response Format**    | Inconsistent across endpoints | Standardized JSON format       |
| **Error Logging**      | Manual logging                | Automatic structured logging   |
| **Pagination**         | Manual calculation            | One-liner with defaults        |
| **Framework Support**  | Single framework              | Gin + Native HTTP + extensible |
| **Testing**            | Complex setup                 | Simple, focused tests          |
| **Maintenance**        | High (repetitive code)        | Low (centralized logic)        |

### Key Benefits

- ✅ **90% Less Validation Code** - Eliminate repetitive validation logic
- ✅ **Consistent API Responses** - Same format across all endpoints
- ✅ **Indonesian Context Ready** - Built-in NIK, phone, localization
- ✅ **Multi-Framework Support** - Gin, Native HTTP, easily extensible
- ✅ **Automatic Error Logging** - Never miss an error again
- ✅ **Bulletproof Pagination** - Handle all edge cases automatically
- ✅ **Type-Safe Errors** - Structured error handling with proper HTTP codes
- ✅ **Zero Breaking Changes** - Migrate gradually without disruption
- ✅ **Production Ready** - Battle-tested patterns and best practices
- ✅ **Developer Experience** - Focus on business logic, not boilerplate

## 🔮 Advanced Use Cases

### 1. Multi-Language API

```go
// Dynamic locale based on request header
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    locale := c.GetHeader("Accept-Language") // "en" or "id"

    if err := validation.BindAndValidate(c, &req,
        validation.WithLocale(locale)); err != nil {
        response.Error(c, err)
        return
    }

    // Business logic...
}
```

### 2. Role-Based Validation

```go
func (h *AdminHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    userRole := c.GetString("user_role")

    if err := validation.BindAndValidate(c, &req,
        validation.WithContext(map[string]interface{}{
            "user_role": userRole,
            "action":    "create_user",
        })); err != nil {
        response.Error(c, err)
        return
    }

    // Business logic...
}

// Custom validation rule that uses context
type RoleBasedRule struct{}

func (r *RoleBasedRule) Validate(value interface{}, params string, context map[string]interface{}) error {
    userRole := context["user_role"].(string)
    if userRole != "admin" && value.(string) == "admin" {
        return fmt.Errorf("only admins can create admin users")
    }
    return nil
}
```

### 3. File Upload with Validation

```go
type FileUploadRequest struct {
    Title       string `form:"title" validate:"required,min=2,max=100"`
    Description string `form:"description" validate:"omitempty,max=500"`
    Category    string `form:"category" validate:"required,oneof=image document video"`
    IsPublic    bool   `form:"isPublic"`
    Tags        string `form:"tags" validate:"omitempty"`
}

func (h *FileHandler) UploadFile(c *gin.Context) {
    var req FileUploadRequest

    // Force multipart form binding for file uploads
    if err := validation.BindAndValidate(c, &req,
        validation.ForceForm()); err != nil {
        response.Error(c, err)
        return
    }

    // Handle file upload...
}
```

### 4. Complex Business Validation

```go
type TransferRequest struct {
    FromAccount string  `json:"fromAccount" validate:"required,uuid"`
    ToAccount   string  `json:"toAccount" validate:"required,uuid"`
    Amount      float64 `json:"amount" validate:"required,min=0.01"`
    Currency    string  `json:"currency" validate:"required,oneof=IDR USD EUR"`
    Note        string  `json:"note" validate:"omitempty,max=200"`
}

func (s *paymentService) Transfer(req TransferRequest) (*TransferResponse, error) {
    // Structural validation already done in handler

    // Business validation
    if req.FromAccount == req.ToAccount {
        return nil, response.BadRequest("Cannot transfer to the same account")
    }

    balance, err := s.getAccountBalance(req.FromAccount)
    if err != nil {
        return nil, response.DatabaseError("Failed to check balance", err)
    }

    if balance < req.Amount {
        return nil, response.BadRequest("Insufficient balance")
    }

    // Process transfer...
}
```

## 🛠️ Extending the Toolkit

### Adding New Framework Support

The toolkit is designed to be easily extensible. Here's how to add support for a new framework:

#### 1. Create Framework Integration File

```go
// validation/fiber.go
package validation

import "github.com/gofiber/fiber/v2"

type FiberWriter struct {
    ctx *fiber.Ctx
}

func (fw *FiberWriter) WriteJSON(statusCode int, data interface{}) error {
    return fw.ctx.Status(statusCode).JSON(data)
}

func BindAndValidateFiber(ctx *fiber.Ctx, obj interface{}, opts ...ValidationOption) error {
    // Implement Fiber-specific binding logic
    // Similar to BindAndValidate but using Fiber's API
}

func InitFiber(configs ...InitConfig) {
    // Initialize for Fiber framework
}
```

#### 2. Add Response Support

```go
// response/fiber.go
package response

import "github.com/gofiber/fiber/v2"

func ErrorFiber(ctx *fiber.Ctx, err error) {
    // Implement Fiber error handling
}

func OKFiber(ctx *fiber.Ctx, message string, data interface{}) {
    // Implement Fiber success response
}
```

### Creating Domain-Specific Rules

```go
// Custom rule for Indonesian bank account validation
type BankAccountRule struct{}

func (r *BankAccountRule) Validate(value interface{}, params string, context map[string]interface{}) error {
    account, ok := value.(string)
    if !ok {
        return fmt.Errorf("bank account must be a string")
    }

    // Indonesian bank account validation logic
    if len(account) < 10 || len(account) > 16 {
        return fmt.Errorf("bank account must be 10-16 digits")
    }

    // Add specific bank validation based on params
    bank := params // e.g., "bca", "mandiri", "bni"
    if !r.validateBankSpecificFormat(account, bank) {
        return fmt.Errorf("invalid account format for %s", bank)
    }

    return nil
}

func (r *BankAccountRule) validateBankSpecificFormat(account, bank string) bool {
    switch bank {
    case "bca":
        return len(account) >= 10 && strings.HasPrefix(account, "0")
    case "mandiri":
        return len(account) >= 13
    case "bni":
        return len(account) >= 10
    default:
        return true
    }
}

func (r *BankAccountRule) GetMessage() string {
    return "field must be a valid Indonesian bank account"
}

// Usage
type PaymentRequest struct {
    BankAccount string `json:"bankAccount" validate:"required,bank_account=bca"`
}
```

## 📈 Performance Considerations

### Validation Performance

- **Reflection Caching**: Struct reflection is cached for better performance
- **Rule Compilation**: Regex patterns are compiled once during initialization
- **Memory Efficiency**: Minimal allocations during validation
- **Concurrent Safe**: All validators are thread-safe

### Response Performance

- **JSON Marshaling**: Efficient JSON encoding with minimal allocations
- **Logger Buffering**: Structured logging with proper buffering
- **Error Pooling**: Error objects are reused when possible

### Benchmarks

```go
// Example benchmark results (approximate)
BenchmarkValidation-8           1000000    1200 ns/op    240 B/op    3 allocs/op
BenchmarkPagination-8          5000000     280 ns/op     64 B/op     1 allocs/op
BenchmarkResponseJSON-8        2000000     680 ns/op    128 B/op     2 allocs/op
```

## 🔍 Troubleshooting

### Common Issues

#### 1. Validation Not Working

```go
// ❌ Problem: Validation not initialized
func main() {
    r := gin.Default()
    // Missing validation.InitGin()
}

// ✅ Solution: Initialize validation
func main() {
    validation.InitGin(validation.InitConfig{})
    r := gin.Default()
}
```

#### 2. Custom Rules Not Found

```go
// ❌ Problem: Custom rule not registered
type MyRequest struct {
    Field string `validate:"mycustom"` // Rule not found
}

// ✅ Solution: Register custom rules
validation.InitGin(validation.InitConfig{
    CustomRules: map[string]validation.Rule{
        "mycustom": &MyCustomRule{},
    },
})
```

#### 3. Localization Not Working

```go
// ❌ Problem: Wrong locale setting
validation.InitGin(validation.InitConfig{
    Locale: "indonesia", // Should be "id"
})

// ✅ Solution: Use correct locale codes
validation.InitGin(validation.InitConfig{
    Locale: "id", // or "en"
})
```

#### 4. Pagination Defaults Not Applied

```go
// ❌ Problem: Missing pagination defaults
func GetUsers(c *gin.Context) {
    var req GetUsersRequest
    validation.BindAndValidate(c, &req)
    // req.Page might be 0
}

// ✅ Solution: Apply pagination defaults
func GetUsers(c *gin.Context) {
    var req GetUsersRequest
    validation.BindAndValidate(c, &req)
    pagination.ApplyDefaultsToStruct(&req) // Apply defaults
}
```

### Debug Mode

```go
// Enable debug logging to see what's happening
validation.InitGin(validation.InitConfig{
    Logger: logger, // Make sure logger is not nil
})

// Check logs for validation details
```

## 🤝 Contributing

We welcome contributions! Here's how you can help:

### Areas for Contribution

1. **New Framework Support** - Add Fiber, Echo, Chi integrations
2. **Additional Validation Rules** - Country-specific or domain-specific rules
3. **Localization** - Add more language support
4. **Performance Improvements** - Optimize validation and response handling
5. **Documentation** - Improve examples and guides

### Getting Started

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-framework`
3. Make your changes with tests
4. Ensure all tests pass: `go test ./...`
5. Submit a pull request

### Code Guidelines

- Follow Go conventions and best practices
- Add comprehensive tests for new features
- Update documentation for any API changes
- Maintain backward compatibility
- Use structured logging for debugging

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built for the Indonesian Go developer community
- Inspired by best practices from Laravel, Express.js, and other frameworks
- Special thanks to contributors and early adopters

---

**Made with ❤️ for the Go community in Indonesia and beyond**

For more examples, advanced usage, and community discussions, visit our [GitHub repository](https://github.com/fiqrioemry/go-api-toolkit).
