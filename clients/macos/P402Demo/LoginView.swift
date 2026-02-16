//
//  LoginView.swift
//  P402Demo
//

import SwiftUI

struct LoginView: View {
    @EnvironmentObject var authVM: AuthViewModel
    @State private var email = ""
    @State private var password = ""
    @State private var rememberMe = false

    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                Text("Login")
                    .font(.title2)
                    .fontWeight(.semibold)

                VStack(alignment: .leading, spacing: 16) {
                    VStack(alignment: .leading, spacing: 6) {
                        Text("Email")
                            .font(.caption)
                            .foregroundColor(.secondary)
                        TextField("email@example.com", text: $email)
                            .textFieldStyle(.roundedBorder)
                            .textContentType(.emailAddress)
                    }

                    VStack(alignment: .leading, spacing: 6) {
                        Text("Password")
                            .font(.caption)
                            .foregroundColor(.secondary)
                        SecureField("••••••••", text: $password)
                            .textFieldStyle(.roundedBorder)
                            .textContentType(.password)
                    }

                    Toggle("Remember Me (30 days)", isOn: $rememberMe)
                        .toggleStyle(.checkbox)
                }
                .padding(.horizontal)

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

                Button(action: login) {
                    if authVM.isLoading {
                        ProgressView()
                            .controlSize(.small)
                            .frame(maxWidth: .infinity)
                    } else {
                        Text("Login")
                            .frame(maxWidth: .infinity)
                    }
                }
                .buttonStyle(.borderedProminent)
                .disabled(authVM.isLoading || email.isEmpty || password.isEmpty)
                .padding(.horizontal)
                .padding(.top, 8)

                Text("First time? Register and verify your email first")
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
            .padding(.vertical, 20)
        }
        .frame(maxWidth: 400)
    }

    private func login() {
        authVM.clearMessages()
        Task {
            await authVM.login(email: email, password: password, rememberMe: rememberMe)
        }
    }
}
