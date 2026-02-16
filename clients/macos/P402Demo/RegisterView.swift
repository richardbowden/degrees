//
//  RegisterView.swift
//  P402Demo
//

import SwiftUI

struct RegisterView: View {
    @EnvironmentObject var authVM: AuthViewModel
    @State private var email = ""
    @State private var password = ""
    @State private var passwordConfirm = ""
    @State private var firstName = ""
    @State private var middleName = ""
    @State private var surname = ""
    @State private var username = ""

    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                Text("Register")
                    .font(.title2)
                    .fontWeight(.semibold)

                VStack(alignment: .leading, spacing: 16) {
                    // Required fields
                    Group {
                        VStack(alignment: .leading, spacing: 6) {
                            Text("First Name *")
                                .font(.caption)
                                .foregroundColor(.secondary)
                            TextField("John", text: $firstName)
                                .textFieldStyle(.roundedBorder)
                        }

                        VStack(alignment: .leading, spacing: 6) {
                            Text("Email *")
                                .font(.caption)
                                .foregroundColor(.secondary)
                            TextField("email@example.com", text: $email)
                                .textFieldStyle(.roundedBorder)
                                .textContentType(.emailAddress)
                        }

                        VStack(alignment: .leading, spacing: 6) {
                            Text("Password *")
                                .font(.caption)
                                .foregroundColor(.secondary)
                            SecureField("••••••••", text: $password)
                                .textFieldStyle(.roundedBorder)
                                .textContentType(.newPassword)
                        }

                        VStack(alignment: .leading, spacing: 6) {
                            Text("Confirm Password *")
                                .font(.caption)
                                .foregroundColor(.secondary)
                            SecureField("••••••••", text: $passwordConfirm)
                                .textFieldStyle(.roundedBorder)
                                .textContentType(.newPassword)
                        }
                    }

                    Divider()

                    // Optional fields
                    Group {
                        VStack(alignment: .leading, spacing: 6) {
                            Text("Middle Name (optional)")
                                .font(.caption)
                                .foregroundColor(.secondary)
                            TextField("", text: $middleName)
                                .textFieldStyle(.roundedBorder)
                        }

                        VStack(alignment: .leading, spacing: 6) {
                            Text("Surname (optional)")
                                .font(.caption)
                                .foregroundColor(.secondary)
                            TextField("Doe", text: $surname)
                                .textFieldStyle(.roundedBorder)
                        }

                        VStack(alignment: .leading, spacing: 6) {
                            Text("Username (optional)")
                                .font(.caption)
                                .foregroundColor(.secondary)
                            TextField("johndoe", text: $username)
                                .textFieldStyle(.roundedBorder)
                        }
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

                        Text("Check server logs for verification token")
                            .font(.caption)
                            .foregroundColor(.blue)
                    }
                    .padding(.horizontal)
                }

                Button(action: register) {
                    if authVM.isLoading {
                        ProgressView()
                            .controlSize(.small)
                            .frame(maxWidth: .infinity)
                    } else {
                        Text("Register")
                            .frame(maxWidth: .infinity)
                    }
                }
                .buttonStyle(.borderedProminent)
                .disabled(authVM.isLoading || !isValid)
                .padding(.horizontal)
                .padding(.top, 8)

                Text("* Required fields")
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
            .padding(.vertical, 20)
        }
        .frame(maxWidth: 400)
    }

    private var isValid: Bool {
        !email.isEmpty && !password.isEmpty && !passwordConfirm.isEmpty && !firstName.isEmpty
    }

    private func register() {
        authVM.clearMessages()
        Task {
            await authVM.register(
                email: email,
                password: password,
                passwordConfirm: passwordConfirm,
                firstName: firstName,
                middleName: middleName.isEmpty ? nil : middleName,
                surname: surname.isEmpty ? nil : surname,
                username: username.isEmpty ? nil : username
            )
        }
    }
}
