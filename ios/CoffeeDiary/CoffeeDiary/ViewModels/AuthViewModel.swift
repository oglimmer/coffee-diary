import AuthenticationServices
import SwiftUI

@Observable
final class AuthViewModel {
    var user: User?
    var isLoading = true
    var error: String?

    var isAuthenticated: Bool { user != nil }

    private let authService = AuthService()

    func checkSession() async {
        isLoading = true
        user = await authService.checkSession()
        isLoading = false
    }

    func login(anchor: ASPresentationAnchor) async {
        error = nil
        do {
            try await authService.login(anchor: anchor)
            user = await authService.checkSession()
        } catch {
            self.error = error.localizedDescription
        }
    }

    func logout() async {
        await authService.logout()
        user = nil
    }
}
