//
//  P402DemoApp.swift
//  P402Demo
//
//  macOS test client for P402 authentication API
//

import SwiftUI

@main
struct P402DemoApp: App {
    @StateObject private var authViewModel = AuthViewModel()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(authViewModel)
                .frame(minWidth: 600, minHeight: 500)
        }
        .windowStyle(.hiddenTitleBar)
        .windowResizability(.contentSize)
    }
}
