# P402 Demo - macOS Test Client

A simple macOS app to test the P402 authentication API.

## Features

- User Registration
- Email Verification
- Login
- Password Change
- Password Reset
- View User List (sysop only)

## Setup

1. Make sure the P402 server is running on `localhost:8080`
2. Open `P402Demo.xcodeproj` in Xcode
3. Build and run (⌘R)

## Usage

### First Time Setup
1. Register a new account
2. Check server logs for verification token
3. Verify your email
4. Login with credentials

### Testing Features
- **Login**: Use registered credentials
- **Change Password**: Must be logged in
- **Reset Password**: Enter email, check logs for token
- **User List**: First user is automatically sysop

## API Endpoints

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/verify-email`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `POST /api/v1/user/change-password`
- `POST /api/v1/auth/reset-password`
- `POST /api/v1/auth/complete-password-reset`
- `GET /api/v1/admin/users`

## Architecture

```
P402Demo/
├── P402DemoApp.swift          # Main app entry point
├── Models/
│   ├── User.swift             # User model
│   └── APIModels.swift        # Request/Response models
├── Services/
│   └── APIClient.swift        # HTTP client for API calls
├── Views/
│   ├── ContentView.swift      # Main navigation
│   ├── RegisterView.swift     # Registration form
│   ├── LoginView.swift        # Login form
│   ├── DashboardView.swift    # Authenticated home
│   ├── VerifyEmailView.swift  # Email verification
│   ├── ChangePasswordView.swift
│   ├── ResetPasswordView.swift
│   └── UserListView.swift     # Admin user list
└── ViewModels/
    └── AuthViewModel.swift    # Handles auth state
```
