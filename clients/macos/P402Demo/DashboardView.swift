//
//  DashboardView.swift
//  P402Demo
//

import SwiftUI

struct DashboardView: View {
    @EnvironmentObject var authVM: AuthViewModel
    @State private var selectedTab = 0
    @State private var showSettings = false

    var body: some View {
        VStack(spacing: 0) {
            // Header
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text("P402 Demo")
                        .font(.title2)
                        .fontWeight(.bold)
                    if let user = authVM.currentUser {
                        Text("Welcome, \(user.firstName)!")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                }

                Spacer()

                Button(action: { showSettings = true }) {
                    Image(systemName: "gear")
                }
                .buttonStyle(.borderless)
                .help("Settings")

                Button(action: logout) {
                    Label("Logout", systemImage: "rectangle.portrait.and.arrow.right")
                }
                .buttonStyle(.bordered)
            }
            .padding(.horizontal, 20)
            .padding(.vertical, 16)
            .background(Color(nsColor: .controlBackgroundColor))

            Divider()

            // Tab Selector
            Picker("", selection: $selectedTab) {
                Text("Profile").tag(0)
                Text("Change Password").tag(1)
                Text("Users (Admin)").tag(2)
            }
            .pickerStyle(.segmented)
            .padding(.horizontal, 20)
            .padding(.vertical, 12)

            Divider()

            // Content
            TabView(selection: $selectedTab) {
                ProfileView()
                    .tag(0)
                ChangePasswordView()
                    .tag(1)
                UserListView()
                    .tag(2)
            }
            .tabViewStyle(.automatic)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(Color(nsColor: .windowBackgroundColor))
        .sheet(isPresented: $showSettings) {
            SettingsView()
        }
    }

    private func logout() {
        Task {
            await authVM.logout()
        }
    }
}

struct ProfileView: View {
    @EnvironmentObject var authVM: AuthViewModel

    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                Text("Profile Information")
                    .font(.title2)
                    .fontWeight(.semibold)
                    .padding(.top, 20)

                if let user = authVM.currentUser {
                    VStack(alignment: .leading, spacing: 16) {
                        ProfileField(label: "ID", value: user.id)
                        ProfileField(label: "First Name", value: user.firstName)
                        if let middleName = user.middleName, !middleName.isEmpty {
                            ProfileField(label: "Middle Name", value: middleName)
                        }
                        if let surname = user.surname, !surname.isEmpty {
                            ProfileField(label: "Surname", value: surname)
                        }
                        ProfileField(label: "Email", value: user.email)
                        ProfileField(
                            label: "Status",
                            value: user.enabled ? "Enabled" : "Disabled",
                            valueColor: user.enabled ? .green : .red
                        )
                    }
                    .padding()
                    .background(
                        RoundedRectangle(cornerRadius: 8)
                            .fill(Color(nsColor: .controlBackgroundColor))
                    )
                    .padding(.horizontal, 40)
                }

                if let token = authVM.sessionToken {
                    VStack(alignment: .leading, spacing: 8) {
                        Text("Session Token")
                            .font(.caption)
                            .foregroundColor(.secondary)
                        Text(token)
                            .font(.system(.caption, design: .monospaced))
                            .foregroundColor(.secondary)
                            .textSelection(.enabled)
                            .padding(8)
                            .background(
                                RoundedRectangle(cornerRadius: 4)
                                    .fill(Color(nsColor: .textBackgroundColor))
                            )
                    }
                    .padding(.horizontal, 40)
                }

                Spacer()
            }
        }
    }
}

struct ProfileField: View {
    let label: String
    let value: String
    var valueColor: Color = .primary

    var body: some View {
        HStack {
            Text(label)
                .font(.caption)
                .foregroundColor(.secondary)
                .frame(width: 100, alignment: .leading)
            Text(value)
                .font(.body)
                .foregroundColor(valueColor)
            Spacer()
        }
    }
}
