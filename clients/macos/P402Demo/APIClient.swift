//
//  APIClient.swift
//  P402Demo
//

import Foundation

enum APIError: LocalizedError {
    case invalidURL
    case networkError(Error)
    case invalidResponse
    case decodingError(Error)
    case serverError(String)
    case unauthorized

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid URL"
        case .networkError(let error):
            return "Network error: \(error.localizedDescription)"
        case .invalidResponse:
            return "Invalid response from server"
        case .decodingError(let error):
            return "Failed to decode response: \(error.localizedDescription)"
        case .serverError(let message):
            return message
        case .unauthorized:
            return "Unauthorized - please login"
        }
    }
}

class APIClient {
    static let shared = APIClient()

    // Configurable base URL - defaults to localhost
    var baseURL: String {
        get {
            UserDefaults.standard.string(forKey: "api_base_url") ?? "http://localhost:8080/api/v1"
        }
        set {
            UserDefaults.standard.set(newValue, forKey: "api_base_url")
        }
    }

    private let decoder: JSONDecoder = {
        let decoder = JSONDecoder()
        return decoder
    }()

    private let encoder: JSONEncoder = {
        let encoder = JSONEncoder()
        return encoder
    }()

    private init() {}

    // MARK: - Generic Request Method
    private func request<T: Decodable>(
        path: String,
        method: String = "GET",
        body: Encodable? = nil,
        authToken: String? = nil
    ) async throws -> T {
        guard let url = URL(string: "\(baseURL)\(path)") else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        if let token = authToken {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        if let body = body {
            request.httpBody = try encoder.encode(body)
        }

        print("📡 \(method) \(path)")
        if let token = authToken {
            print("🔑 Auth: Bearer \(token.prefix(20))...")
        }

        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw APIError.invalidResponse
            }

            print("✅ Status: \(httpResponse.statusCode)")

            // Handle error responses
            if httpResponse.statusCode >= 400 {
                if let errorResponse = try? decoder.decode(ErrorResponse.self, from: data) {
                    throw APIError.serverError(errorResponse.message ?? errorResponse.error ?? "Unknown error")
                }
                if httpResponse.statusCode == 401 {
                    throw APIError.unauthorized
                }
                throw APIError.serverError("HTTP \(httpResponse.statusCode)")
            }

            // Decode success response
            do {
                let result = try decoder.decode(T.self, from: data)
                return result
            } catch {
                print("❌ Decoding error: \(error)")
                if let jsonString = String(data: data, encoding: .utf8) {
                    print("📄 Response: \(jsonString)")
                }
                throw APIError.decodingError(error)
            }
        } catch let error as APIError {
            throw error
        } catch {
            throw APIError.networkError(error)
        }
    }

    // MARK: - Auth Endpoints

    func register(_ request: RegisterRequest) async throws -> RegisterResponse {
        return try await self.request(
            path: "/auth/register",
            method: "POST",
            body: request
        )
    }

    func verifyEmail(_ request: VerifyEmailRequest) async throws -> VerifyEmailResponse {
        return try await self.request(
            path: "/auth/verify-email",
            method: "POST",
            body: request
        )
    }

    func login(_ request: LoginRequest) async throws -> LoginResponse {
        return try await self.request(
            path: "/auth/login",
            method: "POST",
            body: request
        )
    }

    func logout(_ request: LogoutRequest, authToken: String) async throws -> LogoutResponse {
        return try await self.request(
            path: "/auth/logout",
            method: "POST",
            body: request,
            authToken: authToken
        )
    }

    func changePassword(_ request: ChangePasswordRequest, authToken: String) async throws -> ChangePasswordResponse {
        return try await self.request(
            path: "/user/change-password",
            method: "POST",
            body: request,
            authToken: authToken
        )
    }

    func resetPassword(_ request: ResetPasswordRequest) async throws -> ResetPasswordResponse {
        return try await self.request(
            path: "/auth/reset-password",
            method: "POST",
            body: request
        )
    }

    func completePasswordReset(_ request: CompletePasswordResetRequest) async throws -> CompletePasswordResetResponse {
        return try await self.request(
            path: "/auth/complete-password-reset",
            method: "POST",
            body: request
        )
    }

    // MARK: - User Endpoints

    func listUsers(authToken: String) async throws -> [User] {
        // The API might return an array directly or wrapped in an object
        // Try direct array first, then fallback to wrapped response
        do {
            let users: [User] = try await self.request(
                path: "/admin/users",
                method: "GET",
                authToken: authToken
            )
            return users
        } catch {
            // If that fails, try wrapped response
            let response: UserListResponse = try await self.request(
                path: "/admin/users",
                method: "GET",
                authToken: authToken
            )
            return response.users ?? []
        }
    }
}
