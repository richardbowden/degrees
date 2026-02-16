//
//  Models.swift
//  P402Demo
//

import Foundation

// MARK: - User Model
struct User: Codable, Identifiable {
    let id: String
    let firstName: String
    let middleName: String?
    let surname: String?
    let email: String
    let enabled: Bool
    let sysop: Bool?
    let signUpStage: Int?
    let createdOn: String?
    let updatedAt: String?

    // No CodingKeys needed - gRPC-Gateway returns camelCase JSON by default
}

// MARK: - Request Models
struct RegisterRequest: Codable {
    let email: String
    let password: String
    let passwordConfirm: String
    let firstName: String
    let middleName: String?
    let surname: String?
    let username: String?
}

struct LoginRequest: Codable {
    let email: String
    let password: String
    let rememberMe: Bool?
}

struct VerifyEmailRequest: Codable {
    let token: String
}

struct LogoutRequest: Codable {
    let sessionToken: String
}

struct ChangePasswordRequest: Codable {
    let oldPassword: String
    let newPassword: String
    let newPasswordConfirm: String
}

struct ResetPasswordRequest: Codable {
    let email: String
}

struct CompletePasswordResetRequest: Codable {
    let token: String
    let newPassword: String
    let newPasswordConfirm: String
}

// MARK: - Response Models
struct RegisterResponse: Codable {
    let message: String
}

struct LoginResponse: Codable {
    let sessionToken: String
    let user: User
    let message: String
}

struct VerifyEmailResponse: Codable {
    let message: String
}

struct LogoutResponse: Codable {
    let message: String
}

struct ChangePasswordResponse: Codable {
    let message: String
}

struct ResetPasswordResponse: Codable {
    let message: String
}

struct CompletePasswordResetResponse: Codable {
    let message: String
}

struct UserListResponse: Codable {
    let users: [User]?
}

// MARK: - Error Response
struct ErrorResponse: Codable {
    let error: String?
    let message: String?
    let code: String?
}
