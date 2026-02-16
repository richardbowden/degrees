//
//  UserListView.swift
//  P402Demo
//

import SwiftUI

struct UserListView: View {
    @EnvironmentObject var authVM: AuthViewModel
    @State private var users: [User] = []
    @State private var isLoading = false
    @State private var errorMessage: String?

    var body: some View {
        VStack(spacing: 0) {
            // Header
            HStack {
                Text("User List")
                    .font(.title2)
                    .fontWeight(.semibold)

                Spacer()

                Button(action: loadUsers) {
                    Label("Refresh", systemImage: "arrow.clockwise")
                }
                .disabled(isLoading)
            }
            .padding(.horizontal, 20)
            .padding(.vertical, 16)

            Divider()

            // Content
            if isLoading {
                ProgressView("Loading users...")
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if let error = errorMessage {
                VStack(spacing: 16) {
                    Image(systemName: "exclamationmark.triangle")
                        .font(.system(size: 48))
                        .foregroundColor(.orange)

                    Text(error)
                        .font(.caption)
                        .foregroundColor(.red)
                        .multilineTextAlignment(.center)

                    Button("Try Again") {
                        loadUsers()
                    }
                    .buttonStyle(.borderedProminent)
                }
                .padding()
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if users.isEmpty {
                VStack(spacing: 16) {
                    Image(systemName: "person.3")
                        .font(.system(size: 48))
                        .foregroundColor(.secondary)

                    Text("No users found")
                        .font(.headline)
                        .foregroundColor(.secondary)

                    Button("Load Users") {
                        loadUsers()
                    }
                    .buttonStyle(.borderedProminent)
                }
                .padding()
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                ScrollView {
                    LazyVStack(spacing: 12) {
                        ForEach(users) { user in
                            UserRow(user: user)
                        }
                    }
                    .padding()
                }
            }
        }
        .onAppear {
            if users.isEmpty {
                loadUsers()
            }
        }
    }

    private func loadUsers() {
        guard let token = authVM.sessionToken else {
            errorMessage = "Not authenticated"
            return
        }

        isLoading = true
        errorMessage = nil

        Task {
            do {
                let fetchedUsers = try await APIClient.shared.listUsers(authToken: token)
                await MainActor.run {
                    users = fetchedUsers
                    isLoading = false
                }
            } catch {
                await MainActor.run {
                    errorMessage = error.localizedDescription
                    isLoading = false
                }
            }
        }
    }
}

struct UserRow: View {
    let user: User

    var body: some View {
        HStack(spacing: 16) {
            // Avatar
            Circle()
                .fill(LinearGradient(
                    colors: [.blue, .purple],
                    startPoint: .topLeading,
                    endPoint: .bottomTrailing
                ))
                .frame(width: 40, height: 40)
                .overlay(
                    Text(user.firstName.prefix(1).uppercased())
                        .font(.headline)
                        .foregroundColor(.white)
                )

            // User info
            VStack(alignment: .leading, spacing: 4) {
                HStack {
                    Text(user.firstName)
                        .font(.headline)
                    if let middleName = user.middleName, !middleName.isEmpty {
                        Text(middleName)
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                    if let surname = user.surname, !surname.isEmpty {
                        Text(surname)
                            .font(.headline)
                    }
                }

                Text(user.email)
                    .font(.caption)
                    .foregroundColor(.secondary)
            }

            Spacer()

            // Status badge
            HStack(spacing: 8) {
                Circle()
                    .fill(user.enabled ? Color.green : Color.red)
                    .frame(width: 8, height: 8)
                Text(user.enabled ? "Active" : "Disabled")
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
        }
        .padding()
        .background(
            RoundedRectangle(cornerRadius: 8)
                .fill(Color(nsColor: .controlBackgroundColor))
        )
    }
}
