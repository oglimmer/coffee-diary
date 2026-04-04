import SwiftUI

@main
struct CoffeeDiaryApp: App {
    @State private var authViewModel = AuthViewModel()

    var body: some Scene {
        WindowGroup {
            Group {
                if authViewModel.isLoading {
                    ProgressView("Checking session...")
                } else if authViewModel.isAuthenticated {
                    DiaryListView(authViewModel: authViewModel)
                } else {
                    LoginView(authViewModel: authViewModel)
                }
            }
            .task {
                await authViewModel.checkSession()
            }
        }
    }
}
