//
//  ContentView.swift
//  P402Demo
//

import SwiftUI

struct ContentView: View {
    @EnvironmentObject var authVM: AuthViewModel

    var body: some View {
        Group {
            if authVM.isCheckingSession {
                // Show a loading state while checking for existing session
                VStack(spacing: 16) {
                    ProgressView()
                        .scaleEffect(1.5)
                    Text("Loading...")
                        .foregroundColor(.secondary)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if authVM.isAuthenticated {
                DashboardView()
            } else {
                WelcomeView()
            }
        }
        .task {
            await authVM.restoreSession()
        }
    }
}

struct WelcomeView: View {
    @State private var selectedTab = 0
    @State private var showSettings = false

    var body: some View {
        VStack(spacing: 0) {
            // Header
            HStack {
                Spacer()
                VStack(spacing: 8) {
                    Text("P402 Demo")
                        .font(.system(size: 32, weight: .bold))
                    Text("Authentication Test Client")
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                }
                Spacer()
            }
            .overlay(
                Button(action: { showSettings = true }) {
                    Image(systemName: "gear")
                        .font(.title3)
                }
                .buttonStyle(.borderless)
                .help("Settings")
                .padding(.trailing, 20),
                alignment: .trailing
            )
            .padding(.top, 30)
            .padding(.bottom, 20)

            // Tab Selector
            Picker("", selection: $selectedTab) {
                Text("Login").tag(0)
                Text("Register").tag(1)
                Text("Verify Email").tag(2)
                Text("Reset Password").tag(3)
            }
            .pickerStyle(.segmented)
            .padding(.horizontal, 40)
            .padding(.bottom, 20)

            // Content
            TabView(selection: $selectedTab) {
                LoginView()
                    .tag(0)
                RegisterView()
                    .tag(1)
                VerifyEmailView()
                    .tag(2)
                ResetPasswordView()
                    .tag(3)
            }
            .tabViewStyle(.automatic)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(Color(nsColor: .windowBackgroundColor))
        .sheet(isPresented: $showSettings) {
            SettingsView()
        }
    }
}
