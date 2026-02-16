//
//  VerifyEmailView.swift
//  P402Demo
//

import SwiftUI

struct VerifyEmailView: View {
    @EnvironmentObject var authVM: AuthViewModel
    @State private var token = ""

    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                Text("Verify Email")
                    .font(.title2)
                    .fontWeight(.semibold)

                Text("Enter the verification token from your email\n(or check server logs)")
                    .font(.caption)
                    .foregroundColor(.secondary)
                    .multilineTextAlignment(.center)

                VStack(alignment: .leading, spacing: 16) {
                    VStack(alignment: .leading, spacing: 6) {
                        Text("Verification Token")
                            .font(.caption)
                            .foregroundColor(.secondary)
                        TextField("Token from email...", text: $token)
                            .textFieldStyle(.roundedBorder)
                            .font(.system(.body, design: .monospaced))
                    }
                }
                .padding(.horizontal)

                if let error = authVM.errorMessage {
                    Text(error)
                        .font(.caption)
                        .foregroundColor(.red)
                        .padding(.horizontal)
                }

                if let success = authVM.successMessage {
                    VStack(spacing: 8) {
                        Text(success)
                            .font(.caption)
                            .foregroundColor(.green)

                        Text("You can now login!")
                            .font(.caption)
                            .foregroundColor(.blue)
                    }
                    .padding(.horizontal)
                }

                Button(action: verify) {
                    if authVM.isLoading {
                        ProgressView()
                            .controlSize(.small)
                            .frame(maxWidth: .infinity)
                    } else {
                        Text("Verify Email")
                            .frame(maxWidth: .infinity)
                    }
                }
                .buttonStyle(.borderedProminent)
                .disabled(authVM.isLoading || token.isEmpty)
                .padding(.horizontal)
                .padding(.top, 8)

                VStack(spacing: 4) {
                    Text("💡 Tip: Check server logs for the token")
                        .font(.caption)
                        .foregroundColor(.orange)
                    Text("Look for: \"sending verification email\"")
                        .font(.caption2)
                        .foregroundColor(.secondary)
                }
            }
            .padding(.vertical, 20)
        }
        .frame(maxWidth: 400)
    }

    private func verify() {
        authVM.clearMessages()
        Task {
            await authVM.verifyEmail(token: token)
        }
    }
}
