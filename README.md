# 🚀 Go API Toolkit Documentation

A modular, production-ready toolkit for building consistent REST APIs in Go with automatic response formatting, error handling, pagination, and logging.

## 📦 Features

- ✅ **Consistent Response Format** - Standardized JSON responses across all endpoints
- ✅ **Smart Pagination** - Automatic pagination with bulletproof defaults
- ✅ **Structured Error Handling** - Built-in error types with proper HTTP status codes
- ✅ **Automatic Logging** - Optional request/response logging with Zap integration
- ✅ **Framework Agnostic Core** - Supports Gin (more frameworks coming soon)
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
│   └── gin.go         # Gin framework integration
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

### Basic Setup

```go
// main.go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/fiqrioemry/go-api-toolkit/response"
    "github.com/fiqrioemry/go-api-toolkit/pagination"
    "your-app/utils" // Your existing logger
)

func main() {
    // Initialize your existing logger
    utils.InitLogger()

    // 🚀 ONE LINE SETUP - Initialize response handler
    response.InitGin(response.InitConfig{
        Logger:              utils.GetLogger(), // Your existing Zap logger
        LogSuccessResponses: false,             // Set to true to log all success responses
        LogErrorResponses:   true,              // Always log error responses
    })

    // Pagination module doesn't need initialization - it's standalone!

    // Your existing setup...
    r := gin.Default()
    // ... rest of your application
}
```

### Advanced Configuration

```go
// Custom logger configuration
response.InitGin(response.InitConfig{
    Logger:              customZapLogger,
    LogSuccessResponses: true,  // Log all responses in development
    LogErrorResponses:   true,
})

// Or use without logging (NoOp logger by default)
response.InitGin(response.InitConfig{
    Logger: nil, // Will use NoOpLogger - no logging
})
```

## 🔄 Migration Guide

### Before (Manual Response Handling)

```go
// ❌ BEFORE: Manual, repetitive, error-prone

// handlers/user_handler.go
func (h *UserHandler) GetUsers(c *gin.Context) {
    var req dto.GetUsersRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "Invalid parameters",
            "error":   err.Error(),
        })
        return
    }

    // Manual pagination validation
    if req.Page == 0 {
        req.Page = 1
    }
    if req.Limit == 0 {
        req.Limit = 10
    }
    if req.Limit > 100 {
        req.Limit = 100
    }
    if req.SortBy == "" {
        req.SortBy = "created_at"
    }
    if req.SortOrder == "" {
        req.SortOrder = "desc"
    }

    users, total, err := h.service.GetUsers(req)
    if err != nil {
        // Manual error handling
        if strings.Contains(err.Error(), "not found") {
            c.JSON(404, gin.H{
                "success": false,
                "message": "Users not found",
            })
        } else {
            c.JSON(500, gin.H{
                "success": false,
                "message": "Internal server error",
            })
        }
        return
    }

    // Manual pagination calculation
    totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))
    offset := (req.Page - 1) * req.Limit

    // Manual response formatting
    c.JSON(200, gin.H{
        "success": true,
        "message": "Users retrieved successfully",
        "data":    users,
        "meta": gin.H{
            "pagination": gin.H{
                "page":       req.Page,
                "limit":      req.Limit,
                "totalItems": total,
                "totalPages": totalPages,
                "offset":     offset,
            },
        },
    })
}

// services/user_service.go
func (s *userService) GetUsers(req dto.GetUsersRequest) ([]dto.UserResponse, int, error) {
    if req.MinAge != nil && req.MaxAge != nil && *req.MinAge > *req.MaxAge {
        return nil, 0, errors.New("min age cannot be greater than max age")
    }

    users, total, err := s.repo.GetUsers(req)
    if err != nil {
        log.Printf("Error getting users: %v", err) // Manual logging
        return nil, 0, fmt.Errorf("failed to get users: %w", err)
    }

    return users, total, nil
}
```

### After (With Go API Toolkit)

```go
// ✅ AFTER: Clean, consistent, automatic

// handlers/user_handler.go
import (
    "github.com/fiqrioemry/go-api-toolkit/response"
    "github.com/fiqrioemry/go-api-toolkit/pagination"
)

func (h *UserHandler) GetUsers(c *gin.Context) {
    var req dto.GetUsersRequest

    // 🚀 Smart binding with automatic defaults
    if err := pagination.BindAndSetDefaults(c, &req); err != nil {
        response.Error(c, response.BadRequest(err.Error()))
        return
    }

    users, total, err := h.service.GetUsers(req)
    if err != nil {
        response.Error(c, err) // 🚀 Automatic error handling + logging
        return
    }

    // 🚀 One-liner pagination response
    pag := pagination.Build(req.Page, req.Limit, total)
    response.OKWithPagination(c, "Users retrieved successfully", users, pag)
}

// services/user_service.go
func (s *userService) GetUsers(req dto.GetUsersRequest) ([]dto.UserResponse, int, error) {
    if req.MinAge != nil && req.MaxAge != nil && *req.MinAge > *req.MaxAge {
        return nil, 0, response.BadRequest("Min age cannot be greater than max age") // 🚀 Structured errors
    }

    users, total, err := s.repo.GetUsers(req)
    if err != nil {
        return nil, 0, response.DatabaseError("Failed to get users", err) // 🚀 Automatic logging
    }

    return users, total, nil
}
```

### DTO Updates (Remove Binding Constraints)

```go
// ❌ BEFORE: Binding validation that rejects valid requests
type GetUsersRequest struct {
    Page    int    `form:"page" json:"page" binding:"omitempty,min=1"`        // Rejects page=0
    Limit   int    `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"` // Rejects limit=0
    SortBy  string `form:"sortBy" json:"sortBy"`
    MinAge  *int   `form:"minAge" json:"minAge" binding:"omitempty,min=0"`
    MaxAge  *int   `form:"maxAge" json:"maxAge" binding:"omitempty,min=0"`
}

// ✅ AFTER: Let the package handle validation and defaults
type GetUsersRequest struct {
    Page    int    `form:"page" json:"page" binding:"omitempty"`              // Allow 0, will become 1
    Limit   int    `form:"limit" json:"limit" binding:"omitempty"`            // Allow 0, will become 10
    SortBy  string `form:"sortBy" json:"sortBy"`
    MinAge  *int   `form:"minAge" json:"minAge" binding:"omitempty,min=0"`
    MaxAge  *int   `form:"maxAge" json:"maxAge" binding:"omitempty,min=0"`
}
```

## 📚 Usage Examples

### 1. Basic Response Handling

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequestMsg(c, "Invalid request data")
        return
    }

    user, err := h.service.CreateUser(req)
    if err != nil {
        response.Error(c, err) // Automatic error handling
        return
    }

    response.Created(c, "User created successfully", user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")

    user, err := h.service.GetUserByID(userID)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.OK(c, "User retrieved successfully", user)
}
```

### 2. Pagination Examples

#### Approach 1: Using Your Existing DTO

```go
func (h *ProductHandler) GetProducts(c *gin.Context) {
    var req dto.GetProductsRequest

    // Bind and apply smart defaults automatically
    if err := pagination.BindAndSetDefaults(c, &req); err != nil {
        response.Error(c, response.BadRequest(err.Error()))
        return
    }

    // Business validation
    if req.MinPrice != nil && req.MaxPrice != nil && *req.MinPrice > *req.MaxPrice {
        response.BadRequestMsg(c, "Min price cannot be greater than max price")
        return
    }

    products, total, err := h.service.GetProducts(req)
    if err != nil {
        response.Error(c, err)
        return
    }

    // Build pagination and respond
    pag := pagination.Build(req.Page, req.Limit, total)
    response.OKWithPagination(c, "Products retrieved successfully", products, pag)
}
```

#### Approach 2: Using Pagination Module Types

```go
func (h *ProductHandler) GetProductsFlexible(c *gin.Context) {
    var params pagination.FlexibleQueryParams

    // Smart binding with validation
    if err := pagination.SmartBindFlexible(c, &params); err != nil {
        response.Error(c, response.BadRequest(err.Error()))
        return
    }

    products, total, err := h.service.GetProductsFlexible(params)
    if err != nil {
        response.Error(c, err)
        return
    }

    // Quick pagination
    pag := pagination.QuickFlexible(params, total)
    response.OKWithPagination(c, "Products retrieved successfully", products, pag)
}
```

#### Approach 3: Simple Pagination

```go
func (h *ProductHandler) GetProductsSimple(c *gin.Context) {
    var params pagination.DefaultQueryParams

    if err := pagination.SmartBind(c, &params); err != nil {
        response.Error(c, response.BadRequest(err.Error()))
        return
    }

    products, total, err := h.service.GetProductsSimple(params)
    if err != nil {
        response.Error(c, err)
        return
    }

    // One-liner pagination
    pag := pagination.Quick(params, total)
    response.OKWithPagination(c, "Products retrieved successfully", products, pag)
}
```

### 3. Error Handling in Services

```go
func (s *userService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
    // Check if email already exists
    exists, err := s.repo.EmailExists(req.Email)
    if err != nil {
        return nil, response.DatabaseError("Failed to check email existence", err)
    }
    if exists {
        return nil, response.Conflict("Email already exists")
    }

    // Validate age
    if req.Age < 18 {
        return nil, response.BadRequest("User must be at least 18 years old")
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

### 4. Working with int64 Total Counts

```go
func (s *productService) GetProducts(req dto.GetProductsRequest) ([]dto.ProductResponse, int, error) {
    products, total, err := s.repo.GetProducts(req) // total is int64
    if err != nil {
        return nil, 0, response.DatabaseError("Failed to get products", err)
    }

    var result []dto.ProductResponse
    for _, product := range products {
        result = append(result, s.convertToResponse(&product))
    }

    return result, int(total), nil // Convert int64 to int
}
```

### 5. Advanced Response with Permissions

```go
func (h *AdminHandler) GetUsers(c *gin.Context) {
    userID := c.GetString("user_id")
    var req dto.GetUsersRequest

    if err := pagination.BindAndSetDefaults(c, &req); err != nil {
        response.Error(c, response.BadRequest(err.Error()))
        return
    }

    users, total, err := h.service.GetUsers(req)
    if err != nil {
        response.Error(c, err)
        return
    }

    // Get user permissions
    permissions := map[string]bool{
        "canEdit":   h.authService.CanEdit(userID),
        "canDelete": h.authService.CanDelete(userID),
        "canCreate": h.authService.CanCreate(userID),
    }

    pag := pagination.Build(req.Page, req.Limit, total)
    response.OKWithPaginationAndPermissions(c, "Users retrieved successfully", users, pag, permissions)
}
```

## 📤 Response Examples

### Success Response

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "John Doe",
      "email": "john@example.com",
      "createdAt": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### Success Response with Pagination

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "John Doe",
      "email": "john@example.com",
      "createdAt": "2024-01-15T10:30:00Z"
    }
  ],
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

### Success Response with Pagination and Permissions

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
    },
    "permissions": {
      "canEdit": true,
      "canDelete": false,
      "canCreate": true
    }
  }
}
```

### Error Response

```json
{
  "success": false,
  "message": "User not found",
  "code": "NOT_FOUND"
}
```

### Error Response with Validation Details

```json
{
  "success": false,
  "message": "Validation failed",
  "code": "INVALID_INPUT",
  "errors": {
    "email": "Email is required",
    "age": "Age must be at least 18"
  }
}
```

## 🎯 Available Functions

### Response Functions

```go
// Basic responses
response.OK(c, "Success message", data)
response.Created(c, "Created message", data)
response.Error(c, err)

// Quick error messages
response.BadRequestMsg(c, "Invalid data")
response.NotFoundMsg(c, "Resource not found")
response.UnauthorizedMsg(c, "Access denied")
response.ForbiddenMsg(c, "Forbidden")

// Pagination responses
response.OKWithPagination(c, "Success", data, pagination)
response.OKWithPaginationAndPermissions(c, "Success", data, pagination, permissions)
```

### Error Constructors

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

## 📊 Logging Output

When error logging is enabled, you'll see structured logs like:

```json
{
  "level": "error",
  "ts": "2024-01-15T10:30:00.000Z",
  "msg": "Client error occurred",
  "path": "/api/v1/users",
  "method": "GET",
  "client_ip": "192.168.1.1",
  "user_agent": "Mozilla/5.0...",
  "user_id": "user123",
  "trace_id": "trace456",
  "error_code": "NOT_FOUND",
  "error_message": "User not found"
}
```

## 🧪 Testing

### Unit Testing Handler

```go
func TestGetUsers(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    response.InitGin(response.InitConfig{Logger: nil}) // No logging in tests

    mockService := &MockUserService{}
    handler := NewUserHandler(mockService)

    // Test
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("GET", "/users?page=1&limit=10", nil)

    handler.GetUsers(c)

    // Assert
    assert.Equal(t, 200, w.Code)
    // ... more assertions
}
```

### Testing Core Handler (Framework Agnostic)

```go
func TestHandler(t *testing.T) {
    handler := response.NewHandler()
    mockWriter := &MockJSONWriter{}

    err := response.NotFound("User not found")
    handler.HandleError(mockWriter, nil, err)

    assert.Equal(t, 404, mockWriter.StatusCode)
    // ... more assertions
}
```

## 🚀 Migration Steps

1. **Install Package**

   ```bash
   go get github.com/fiqrioemry/go-api-toolkit
   ```

2. **Initialize in main.go**

   ```go
   response.InitGin(response.InitConfig{
       Logger: utils.GetLogger(),
       LogErrorResponses: true,
   })
   ```

3. **Update Imports**

   ```go
   import (
       "github.com/fiqrioemry/go-api-toolkit/response"
       "github.com/fiqrioemry/go-api-toolkit/pagination"
   )
   ```

4. **Update Handlers**

   - Replace manual error handling with `response.Error(c, err)`
   - Replace manual responses with `response.OK()`, `response.Created()`
   - Use `pagination.BindAndSetDefaults()` for smart binding
   - Use `pagination.Build()` for pagination

5. **Update Services**

   - Replace custom errors with `response.BadRequest()`, `response.NotFound()`, etc.
   - Return `int` total count instead of pagination objects

6. **Update DTOs**
   - Remove `min=1` constraints from Page/Limit binding tags
   - Let the package handle validation and defaults

## 🎉 Benefits Summary

- ✅ **80% Less Boilerplate Code** - Eliminate repetitive response formatting
- ✅ **Consistent API Responses** - Same format across all endpoints
- ✅ **Automatic Error Logging** - Never miss an error again
- ✅ **Bulletproof Pagination** - Handle all edge cases automatically
- ✅ **Type-Safe Errors** - Structured error handling with proper HTTP codes
- ✅ **Framework Flexibility** - Easy to add support for Fiber, Echo, etc.
- ✅ **Zero Breaking Changes** - Migrate gradually without disruption
- ✅ **Production Ready** - Battle-tested patterns and best practices

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Made with ❤️ for the Go community**
