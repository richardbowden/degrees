//
//  SettingsView.swift
//  P402Demo
//

import SwiftUI

struct SettingsView: View {
    @State private var apiURL: String
    @State private var showSuccess = false
    @Environment(\.dismiss) var dismiss

    init() {
        _apiURL = State(initialValue: APIClient.shared.baseURL)
    }

    var body: some View {
        VStack(spacing: 0) {
            // Header
            HStack {
                Text("Settings")
                    .font(.title2)
                    .fontWeight(.semibold)
                Spacer()
                Button("Done") {
                    dismiss()
                }
                .buttonStyle(.borderedProminent)
            }
            .padding(.horizontal, 20)
            .padding(.vertical, 16)
            .background(Color(nsColor: .controlBackgroundColor))

            Divider()

            // Content
            ScrollView {
                VStack(alignment: .leading, spacing: 20) {
                    VStack(alignment: .leading, spacing: 12) {
                        Text("API Configuration")
                            .font(.headline)

                        VStack(alignment: .leading, spacing: 6) {
                            Text("Base URL")
                                .font(.caption)
                                .foregroundColor(.secondary)

                            TextField("http://localhost:8080/api/v1", text: $apiURL)
                                .textFieldStyle(.roundedBorder)
                                .font(.system(.body, design: .monospaced))
                        }

                        HStack(spacing: 12) {
                            Button("Save") {
                                saveURL()
                            }
                            .buttonStyle(.borderedProminent)

                            Button("Reset to Default") {
                                apiURL = "http://localhost:8080/api/v1"
                            }
                            .buttonStyle(.bordered)
                        }

                        if showSuccess {
                            HStack {
                                Image(systemName: "checkmark.circle.fill")
                                    .foregroundColor(.green)
                                Text("URL saved successfully")
                                    .font(.caption)
                                    .foregroundColor(.green)
                            }
                        }
                    }
                    .padding()
                    .background(
                        RoundedRectangle(cornerRadius: 8)
                            .fill(Color(nsColor: .controlBackgroundColor))
                    )

                    VStack(alignment: .leading, spacing: 8) {
                        Text("Common URLs")
                            .font(.headline)

                        VStack(spacing: 8) {
                            URLPreset(
                                title: "Local Development",
                                url: "http://localhost:8080/api/v1",
                                onSelect: { apiURL = $0 }
                            )

                            URLPreset(
                                title: "Local IP (for testing on devices)",
                                url: "http://192.168.1.100:8080/api/v1",
                                onSelect: { apiURL = $0 }
                            )

                            URLPreset(
                                title: "Staging Server",
                                url: "https://staging.example.com/api/v1",
                                onSelect: { apiURL = $0 }
                            )
                        }
                    }
                    .padding()
                    .background(
                        RoundedRectangle(cornerRadius: 8)
                            .fill(Color(nsColor: .controlBackgroundColor))
                    )

                    VStack(alignment: .leading, spacing: 8) {
                        Text("ℹ️ Note")
                            .font(.caption)
                            .fontWeight(.semibold)
                            .foregroundColor(.blue)

                        Text("Changes take effect immediately. Make sure to include /api/v1 in the URL.")
                            .font(.caption)
                            .foregroundColor(.secondary)

                        Text("For production servers, use https:// for secure connections.")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding()
                }
                .padding(20)
            }
        }
        .frame(width: 500, height: 400)
    }

    private func saveURL() {
        APIClient.shared.baseURL = apiURL
        showSuccess = true

        // Hide success message after 2 seconds
        DispatchQueue.main.asyncAfter(deadline: .now() + 2) {
            showSuccess = false
        }
    }
}

struct URLPreset: View {
    let title: String
    let url: String
    let onSelect: (String) -> Void

    var body: some View {
        Button(action: { onSelect(url) }) {
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text(title)
                        .font(.caption)
                        .fontWeight(.medium)
                    Text(url)
                        .font(.caption2)
                        .foregroundColor(.secondary)
                        .lineLimit(1)
                }
                Spacer()
                Image(systemName: "chevron.right")
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
            .padding(8)
            .background(
                RoundedRectangle(cornerRadius: 6)
                    .fill(Color(nsColor: .textBackgroundColor))
            )
        }
        .buttonStyle(.plain)
    }
}
