//
//  AuthViewModel.swift
//  P402Demo
//

import Foundation
import Combine

@MainActor
class AuthViewModel: ObservableObject {
    @Published var isAuthenticated = false
    @Published var currentUser: User?
    @Published var sessionToken: String?
    @Published var errorMessage: String?
    @Published var successMessage: String?
    @Published var isLoading = false
    @Published var isCheckingSession = true

    private let api = APIClient.shared

    // MARK: - Registration
    func register(
        email: String,
        password: String,
        passwordConfirm: String,
        firstName: String,
        middleName: String? = nil,
        surname: String? = nil,
        username: String? = nil
    ) async {
        isLoading = true
        errorMessage = nil
        successMessage = nil

        let request = RegisterRequest(
            email: email,
            password: password,
            passwordConfirm: passwordConfirm,
            firstName: firstName,
            middleName: middleName,
            surname: surname,
            username: username
        )

        do {
            let response = try await api.register(request)
            successMessage = response.message
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    // MARK: - Email Verification
    func verifyEmail(token: String) async {
        isLoading = true
        errorMessage = nil
        successMessage = nil

        let request = VerifyEmailRequest(token: token)

        do {
            let response = try await api.verifyEmail(request)
            successMessage = response.message
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    // MARK: - Login
    func login(email: String, password: String, rememberMe: Bool = false) async {
        isLoading = true
        errorMessage = nil
        successMessage = nil

        let request = LoginRequest(
            email: email,
            password: password,
            rememberMe: rememberMe
        )

        do {
            let response = try await api.login(request)
            sessionToken = response.sessionToken
            currentUser = response.user
            isAuthenticated = true
            successMessage = "Welcome, \(response.user.firstName)!"

            // Store credentials securely (simplified for demo)
            UserDefaults.standard.set(response.sessionToken, forKey: "sessionToken")
            
            // Store user data
            if let userData = try? JSONEncoder().encode(response.user) {
                UserDefaults.standard.set(userData, forKey: "currentUser")
            }
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    // MARK: - Logout
    func logout() async {
        guard let token = sessionToken else {
            self.performLogout()
            return
        }

        isLoading = true
        errorMessage = nil

        let request = LogoutRequest(sessionToken: token)

        do {
            _ = try await api.logout(request, authToken: token)
            performLogout()
            successMessage = "Logged out successfully"
        } catch {
            // Still logout locally even if server call fails
            performLogout()
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    private func performLogout() {
        sessionToken = nil
        currentUser = nil
        isAuthenticated = false
        UserDefaults.standard.removeObject(forKey: "sessionToken")
        UserDefaults.standard.removeObject(forKey: "currentUser")
    }

    // MARK: - Change Password
    func changePassword(oldPassword: String, newPassword: String, newPasswordConfirm: String) async {
        guard let token = sessionToken else {
            errorMessage = "Not authenticated"
            return
        }

        isLoading = true
        errorMessage = nil
        successMessage = nil

        let request = ChangePasswordRequest(
            oldPassword: oldPassword,
            newPassword: newPassword,
            newPasswordConfirm: newPasswordConfirm
        )

        do {
            let response = try await api.changePassword(request, authToken: token)
            successMessage = response.message
            // Logout after password change as all sessions are invalidated
            performLogout()
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    // MARK: - Reset Password
    func initiatePasswordReset(email: String) async {
        isLoading = true
        errorMessage = nil
        successMessage = nil

        let request = ResetPasswordRequest(email: email)

        do {
            let response = try await api.resetPassword(request)
            successMessage = response.message
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    func completePasswordReset(token: String, newPassword: String, newPasswordConfirm: String) async {
        isLoading = true
        errorMessage = nil
        successMessage = nil

        let request = CompletePasswordResetRequest(
            token: token,
            newPassword: newPassword,
            newPasswordConfirm: newPasswordConfirm
        )

        do {
            let response = try await api.completePasswordReset(request)
            successMessage = response.message
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    // MARK: - Restore Session
    func restoreSession() async {
        defer { isCheckingSession = false }
        
        guard let token = UserDefaults.standard.string(forKey: "sessionToken") else {
            return
        }

        // Restore user data
        if let userData = UserDefaults.standard.data(forKey: "currentUser"),
           let user = try? JSONDecoder().decode(User.self, from: userData) {
            currentUser = user
        }
        
        // Try to use the stored token
        sessionToken = token

        // Validate by fetching user list or another authenticated endpoint
        // For simplicity, we'll just set it as authenticated
        // In production, you'd validate the token with a "GET /user/me" endpoint
        isAuthenticated = true
    }

    // MARK: - Clear Messages
    func clearMessages() {
        errorMessage = nil
        successMessage = nil
    }
}
