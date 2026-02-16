//
//  ResetPasswordView.swift
//  P402Demo
//

import SwiftUI

struct ResetPasswordView: View {
    @EnvironmentObject var authVM: AuthViewModel
    @State private var step = 1 // 1 = request reset, 2 = complete reset
    @State private var email = ""
    @State private var token = ""
    @State private var newPassword = ""
    @State private var newPasswordConfirm = ""

    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                Text("Reset Password")
                    .font(.title2)
                    .fontWeight(.semibold)

                if step == 1 {
                    stepOne
                } else {
                    stepTwo
                }

                if let error = authVM.errorMessage {
                    Text(error)
                        .font(.caption)
                        .foregroundColor(.red)
                        .padding(.horizontal)
                }

                if let success = authVM.successMessage {
                    Text(success)
                        .font(.caption)
                        .foregroundColor(.green)
                        .padding(.horizontal)
                }

                if step == 1 {
                    Button(action: requestReset) {
                        if authVM.isLoading {
                            ProgressView()
                                .controlSize(.small)
                                .frame(maxWidth: .infinity)
                        } else {
                            Text("Request Reset")
                                .frame(maxWidth: .infinity)
                        }
                    }
                    .buttonStyle(.borderedProminent)
                    .disabled(authVM.isLoading || email.isEmpty)
                    .padding(.horizontal)
                } else {
                    Button(action: completeReset) {
                        if authVM.isLoading {
                            ProgressView()
                                .controlSize(.small)
                                .frame(maxWidth: .infinity)
                        } else {
                            Text("Reset Password")
                                .frame(maxWidth: .infinity)
                        }
                    }
                    .buttonStyle(.borderedProminent)
                    .disabled(authVM.isLoading || !isStepTwoValid)
                    .padding(.horizontal)
                }

                if step == 2 {
                    Button("← Back to Request") {
                        step = 1
                        authVM.clearMessages()
                    }
                    .buttonStyle(.link)
                }
            }
            .padding(.vertical, 20)
        }
        .frame(maxWidth: 400)
        .onChange(of: authVM.successMessage) { message in
            if step == 1 && message != nil {
                // Move to step 2 after successful request
                DispatchQueue.main.asyncAfter(deadline: .now() + 2) {
                    step = 2
                    authVM.clearMessages()
                }
            }
        }
    }

    private var stepOne: some View {
        VStack(spacing: 16) {
            Text("Step 1: Request Reset")
                .font(.headline)

            Text("Enter your email to receive a password reset token")
                .font(.caption)
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)

            VStack(alignment: .leading, spacing: 6) {
                Text("Email")
                    .font(.caption)
                    .foregroundColor(.secondary)
                TextField("email@example.com", text: $email)
                    .textFieldStyle(.roundedBorder)
                    .textContentType(.emailAddress)
            }
            .padding(.horizontal)

            Text("💡 Check server logs for reset token")
                .font(.caption)
                .foregroundColor(.orange)
        }
    }

    private var stepTwo: some View {
        VStack(spacing: 16) {
            Text("Step 2: Complete Reset")
                .font(.headline)

            Text("Enter the reset token and your new password")
                .font(.caption)
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)

            VStack(alignment: .leading, spacing: 16) {
                VStack(alignment: .leading, spacing: 6) {
                    Text("Reset Token")
                        .font(.caption)
                        .foregroundColor(.secondary)
                    TextField("Token from email...", text: $token)
                        .textFieldStyle(.roundedBorder)
                        .font(.system(.body, design: .monospaced))
                }

                VStack(alignment: .leading, spacing: 6) {
                    Text("New Password")
                        .font(.caption)
                        .foregroundColor(.secondary)
                    SecureField("••••••••", text: $newPassword)
                        .textFieldStyle(.roundedBorder)
                        .textContentType(.newPassword)
                }

                VStack(alignment: .leading, spacing: 6) {
                    Text("Confirm New Password")
                        .font(.caption)
                        .foregroundColor(.secondary)
                    SecureField("••••••••", text: $newPasswordConfirm)
                        .textFieldStyle(.roundedBorder)
                        .textContentType(.newPassword)
                }
            }
            .padding(.horizontal)
        }
    }

    private var isStepTwoValid: Bool {
        !token.isEmpty && !newPassword.isEmpty && !newPasswordConfirm.isEmpty
    }

    private func requestReset() {
        authVM.clearMessages()
        Task {
            await authVM.initiatePasswordReset(email: email)
        }
    }

    private func completeReset() {
        authVM.clearMessages()
        Task {
            await authVM.completePasswordReset(
                token: token,
                newPassword: newPassword,
                newPasswordConfirm: newPasswordConfirm
            )
        }
    }
}
