//
//  ChangePasswordView.swift
//  P402Demo
//

import SwiftUI

struct ChangePasswordView: View {
    @EnvironmentObject var authVM: AuthViewModel
    @State private var oldPassword = ""
    @State private var newPassword = ""
    @State private var newPasswordConfirm = ""

    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                Text("Change Password")
                    .font(.title2)
                    .fontWeight(.semibold)
                    .padding(.top, 20)

                VStack(alignment: .leading, spacing: 16) {
                    VStack(alignment: .leading, spacing: 6) {
                        Text("Current Password")
                            .font(.caption)
                            .foregroundColor(.secondary)
                        SecureField("••••••••", text: $oldPassword)
                            .textFieldStyle(.roundedBorder)
                            .textContentType(.password)
                    }

                    Divider()

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
                .padding(.horizontal, 40)

                if let error = authVM.errorMessage {
                    Text(error)
                        .font(.caption)
                        .foregroundColor(.red)
                        .padding(.horizontal, 40)
                }

                if let success = authVM.successMessage {
                    VStack(spacing: 8) {
                        Text(success)
                            .font(.caption)
                            .foregroundColor(.green)

                        Text("You will be logged out and need to login again")
                            .font(.caption)
                            .foregroundColor(.blue)
                    }
                    .padding(.horizontal, 40)
                }

                Button(action: changePassword) {
                    if authVM.isLoading {
                        ProgressView()
                            .controlSize(.small)
                            .frame(maxWidth: .infinity)
                    } else {
                        Text("Change Password")
                            .frame(maxWidth: .infinity)
                    }
                }
                .buttonStyle(.borderedProminent)
                .disabled(authVM.isLoading || !isValid)
                .padding(.horizontal, 40)
                .padding(.top, 8)

                VStack(spacing: 4) {
                    Text("⚠️ Warning")
                        .font(.caption)
                        .fontWeight(.semibold)
                        .foregroundColor(.orange)
                    Text("Changing password will invalidate all active sessions")
                        .font(.caption)
                        .foregroundColor(.secondary)
                        .multilineTextAlignment(.center)
                }
                .padding(.horizontal, 40)

                Spacer()
            }
        }
    }

    private var isValid: Bool {
        !oldPassword.isEmpty && !newPassword.isEmpty && !newPasswordConfirm.isEmpty
    }

    private func changePassword() {
        authVM.clearMessages()
        Task {
            await authVM.changePassword(
                oldPassword: oldPassword,
                newPassword: newPassword,
                newPasswordConfirm: newPasswordConfirm
            )
        }
    }
}
